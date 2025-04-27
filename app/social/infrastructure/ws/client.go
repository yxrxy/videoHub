package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hertz-contrib/websocket"
	"github.com/yxrxy/videoHub/pkg/constants"
)

const (
	// 写入超时时间
	writeWait = 10 * time.Second

	// 读取超时时间
	pongWait = 60 * time.Second

	// 发送ping的时间间隔，必须小于pongWait
	pingPeriod = time.Duration(float64(pongWait) * constants.WebSocketPingRatio)

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
	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("error setting read deadline: %v", err)
	}
	c.Conn.SetPongHandler(func(string) error {
		if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			log.Printf("error setting read deadline in pong handler: %v", err)
		}
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

		log.Printf("received message: %s", string(message))

		// 添加写入锁
		c.mu.Lock()

		// 检查消息类型
		var content interface{}
		if err := json.Unmarshal(message, &content); err != nil {
			// 如果解析失败，将消息作为纯文本处理
			content = string(message)
		}

		// 获取房间ID（不需要再加锁，因为已经在锁内）
		var roomID int64
		for rid := range c.Rooms {
			roomID = rid
			break // 目前只取第一个房间
		}

		// 如果没有房间ID，记录错误并继续
		if roomID == 0 {
			log.Printf("error: client not in any room")
			if err := c.Conn.WriteJSON(map[string]interface{}{
				"type":    "error",
				"message": "您还未加入任何聊天室",
			}); err != nil {
				log.Printf("error sending error message: %v", err)
			}
			c.mu.Unlock()
			continue
		}

		// 构造消息对象
		msg := Message{
			Type:      MessageTypeGroup,
			From:      c.UserID,
			RoomID:    roomID,
			Content:   fmt.Sprint(content),
			Timestamp: time.Now().Unix(),
			Extra: map[string]any{
				"type": "text",
			},
		}

		// 发送消息到聊天室
		manager.SendToRoom(roomID, &msg)
		log.Printf("message sent to room %d: %+v", roomID, msg)

		c.mu.Unlock()
	}
}

// WritePump 处理WebSocket写入
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for range ticker.C {
		if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
			log.Printf("error setting write deadline: %v", err)
			return
		}
		if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			return
		}
	}
}

// SendMessage 发送消息给客户端
func (c *Client) SendMessage(messageType int, data []byte) error {
	if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		log.Printf("error setting write deadline: %v", err)
		return err
	}
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
