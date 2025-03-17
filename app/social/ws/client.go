package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hertz-contrib/websocket"
)

const (
	// 写入超时时间
	writeWait = 10 * time.Second

	// 读取超时时间
	pongWait = 60 * time.Second

	// 发送ping的时间间隔，必须小于pongWait
	pingPeriod = (pongWait * 9) / 10

	// 最大消息大小
	maxMessageSize = 512 * 1024 // 512KB
)

type ExtendedClient struct {
	Client
	Send chan []byte // 发送消息的通道
}

// NewExtendedClient 创建新的客户端连接
func NewExtendedClient(userID int64, conn *websocket.Conn) *ExtendedClient {
	return &ExtendedClient{
		Client: Client{
			UserID: userID,
			Conn:   conn,
			Rooms:  make(map[int64]bool),
		},
		Send: make(chan []byte, 256),
	}
}

// ReadPump 处理WebSocket读取
func (c *ExtendedClient) ReadPump(manager *Manager) {
	defer func() {
		manager.unregister <- &c.Client
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		messageType, data, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		if messageType != websocket.TextMessage && messageType != websocket.BinaryMessage {
			continue
		}

		// 解析消息
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		// 设置发送者ID和时间戳
		msg.From = c.UserID
		msg.Timestamp = time.Now().Unix()

		// 根据消息类型处理
		switch msg.Type {
		case MessageTypePrivate:
			// 发送私信
			if msg.To > 0 {
				data, err := json.Marshal(msg)
				if err != nil {
					log.Printf("error marshaling message: %v", err)
					continue
				}
				manager.SendToUser(msg.To, data)
			}
		case MessageTypeGroup:
			// 发送群聊消息
			if msg.RoomID > 0 {
				manager.SendToRoom(msg.RoomID, &msg)
			}
		}
	}
}

// WritePump 处理WebSocket写入
func (c *ExtendedClient) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 通道已关闭
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			// 发送ping消息
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// JoinRoom 加入聊天室
func (c *ExtendedClient) JoinRoom(roomID int64) {
	c.Client.mu.Lock()
	c.Client.Rooms[roomID] = true
	c.Client.mu.Unlock()
}

// LeaveRoom 离开聊天室
func (c *ExtendedClient) LeaveRoom(roomID int64) {
	c.Client.mu.Lock()
	delete(c.Client.Rooms, roomID)
	c.Client.mu.Unlock()
}

// IsInRoom 检查是否在聊天室中
func (c *ExtendedClient) IsInRoom(roomID int64) bool {
	c.Client.mu.Lock()
	defer c.Client.mu.Unlock()
	return c.Client.Rooms[roomID]
}
