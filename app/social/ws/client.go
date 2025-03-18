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

// ReadPump 处理WebSocket读取
func (c *Client) ReadPump(manager *Manager) {
	defer func() {
		manager.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// 解析消息
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("error parsing message: %v", err)
			continue
		}

		// 设置发送者ID和时间戳
		msg.From = c.UserID
		msg.Timestamp = time.Now().Unix()

		// 根据消息类型处理
		switch msg.Type {
		case MessageTypePrivate:
			if msg.To > 0 {
				data, err := json.Marshal(msg)
				if err != nil {
					log.Printf("error marshaling message: %v", err)
					continue
				}
				manager.SendToUser(msg.To, data)
			}
		case MessageTypeGroup:
			if msg.RoomID > 0 {
				manager.SendToRoom(msg.RoomID, &msg)
			}
		}
	}
}

// WritePump 处理WebSocket写入
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendMessage 发送消息给客户端
func (c *Client) SendMessage(messageType int, data []byte) error {
	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Conn.WriteMessage(messageType, data)
}

// JoinRoom 加入聊天室
func (c *Client) JoinRoom(roomID int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Rooms[roomID] = true
}

// LeaveRoom 离开聊天室
func (c *Client) LeaveRoom(roomID int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.Rooms, roomID)
}

// IsInRoom 检查是否在聊天室中
func (c *Client) IsInRoom(roomID int64) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Rooms[roomID]
}
