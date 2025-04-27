package ws

import (
	"encoding/json"
	"fmt"
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

		// 解析消息
		var content interface{}
		if err := json.Unmarshal(message, &content); err != nil {
			content = string(message)
		}

		msg := Message{
			Type:      MessageTypeGroup,
			From:      c.UserID,
			Content:   fmt.Sprint(content),
			Timestamp: time.Now().Unix(),
			Extra: map[string]any{
				"type": "text",
			},
		}

		if contentMap, ok := content.(map[string]interface{}); ok {
			if roomID, ok := contentMap["room_id"].(float64); ok {
				msg.RoomID = int64(roomID)
			}
			if msgType, ok := contentMap["type"].(string); ok {
				msg.Type = msgType
			}
			if to, ok := contentMap["to"].(float64); ok {
				msg.To = int64(to)
			}
			if contentVal, ok := contentMap["content"]; ok {
				msg.Content = fmt.Sprint(contentVal)
			}
		}

		switch msg.Type {
		case MessageTypePrivate:
			if msg.To > 0 {
				data, err := json.Marshal(msg)
				if err != nil {
					log.Printf("error marshaling message: %v", err)
					continue
				}
				if err := manager.SendToUser(msg.To, data); err != nil {
					log.Printf("error sending message to user: %v", err)
				}
			}
		case MessageTypeGroup:
			if msg.RoomID > 0 {
				manager.SendToRoom(c.UserID, msg.RoomID, &msg)
			}
		}
	}
}

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
