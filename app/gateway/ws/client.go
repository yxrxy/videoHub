package ws

import (
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
)

// ReadPump 处理WebSocket读取
func (c *Client) ReadPump(manager *Manager) {
	defer func() {
		manager.unregister <- c
		c.Conn.Close()
	}()
}

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
			log.Printf("error writing ping message: %v", err)
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
