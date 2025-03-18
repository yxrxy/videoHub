package ws

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/hertz-contrib/websocket"
)

// 消息类型
const (
	// 系统消息类型
	MessageTypeSystem = "system"
	// 私信消息类型
	MessageTypePrivate = "private"
	// 群聊消息类型
	MessageTypeGroup = "group"
	// 好友请求消息类型
	MessageTypeFriendRequest = "friend_request"
)

// Message WebSocket消息结构
type Message struct {
	Type      string         `json:"type"`              // 消息类型
	From      int64          `json:"from"`              // 发送者ID
	To        int64          `json:"to,omitempty"`      // 接收者ID（私信时使用）
	RoomID    int64          `json:"room_id,omitempty"` // 聊天室ID（群聊时使用）
	Content   string         `json:"content"`           // 消息内容
	Extra     map[string]any `json:"extra,omitempty"`   // 额外数据
	Timestamp int64          `json:"timestamp"`         // 时间戳
}

// Client 表示一个WebSocket客户端连接
type Client struct {
	UserID int64           // 用户ID
	Conn   *websocket.Conn // WebSocket连接
	Rooms  map[int64]bool  // 加入的聊天室
	mu     sync.Mutex      // 保护 Rooms 的互斥锁
}

// Manager 管理WebSocket连接
type Manager struct {
	clients    map[int64]*Client          // 用户ID -> 客户端连接
	roomMap    map[int64]map[*Client]bool // 聊天室ID -> 客户端集合
	broadcast  chan []byte                // 广播消息通道
	register   chan *Client               // 注册客户端通道
	unregister chan *Client               // 注销客户端通道
	mutex      sync.RWMutex               // 保护 clients 和 roomMap 的互斥锁
}

// NewManager 创建一个新的WebSocket管理器
func NewManager() *Manager {
	return &Manager{
		clients:    make(map[int64]*Client),
		roomMap:    make(map[int64]map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// RegisterClient 注册新的WebSocket客户端
func (m *Manager) RegisterClient(userID int64, conn *websocket.Conn) {
	client := &Client{
		UserID: userID,
		Conn:   conn,
		Rooms:  make(map[int64]bool),
	}
	m.register <- client

	go NewExtendedClient(userID, conn).ReadPump(m)
}

// Start 启动WebSocket管理器
func (m *Manager) Start(ctx context.Context) {
	for {
		select {
		case client := <-m.register:
			m.mutex.Lock()
			m.clients[client.UserID] = client
			m.mutex.Unlock()

		case client := <-m.unregister:
			m.mutex.Lock()
			if _, ok := m.clients[client.UserID]; ok {
				delete(m.clients, client.UserID)
				client.Conn.Close()
			}
			m.mutex.Unlock()

		case message := <-m.broadcast:
			m.mutex.RLock()
			for _, client := range m.clients {
				go func(client *Client) {
					if err := client.Conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
						m.unregister <- client
					}
				}(client)
			}
			m.mutex.RUnlock()

		case <-ctx.Done():
			return
		}
	}
}

// SendToUser 发送消息给指定用户
func (m *Manager) SendToUser(userID int64, message []byte) error {
	m.mutex.RLock()
	client, exists := m.clients[userID]
	m.mutex.RUnlock()

	if !exists {
		return nil // 用户不在线，忽略消息
	}

	return client.Conn.WriteMessage(websocket.BinaryMessage, message)
}

// BroadcastMessage 广播消息给所有连接的客户端
func (m *Manager) BroadcastMessage(message []byte) {
	m.broadcast <- message
}

// GetOnlineUsers 获取在线用户ID列表
func (m *Manager) GetOnlineUsers() []int64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	users := make([]int64, 0, len(m.clients))
	for userID := range m.clients {
		users = append(users, userID)
	}
	return users
}

// IsUserOnline 检查用户是否在线
func (m *Manager) IsUserOnline(userID int64) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	_, exists := m.clients[userID]
	return exists
}

// StartHeartbeat 启动心跳检测
func (m *Manager) StartHeartbeat(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.mutex.RLock()
			for _, client := range m.clients {
				go func(client *Client) {
					if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
						m.unregister <- client
					}
				}(client)
			}
			m.mutex.RUnlock()
		case <-ctx.Done():
			return
		}
	}
}

// JoinRoom 加入聊天室
func (m *Manager) JoinRoom(client *Client, roomID int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 将房间ID添加到客户端的房间列表
	client.mu.Lock()
	client.Rooms[roomID] = true
	client.mu.Unlock()

	// 将客户端添加到房间的客户端列表
	if _, exists := m.roomMap[roomID]; !exists {
		m.roomMap[roomID] = make(map[*Client]bool)
	}
	m.roomMap[roomID][client] = true
}

// LeaveRoom 离开聊天室
func (m *Manager) LeaveRoom(client *Client, roomID int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 从客户端的房间列表中移除
	client.mu.Lock()
	delete(client.Rooms, roomID)
	client.mu.Unlock()

	// 从房间的客户端列表中移除
	if room, exists := m.roomMap[roomID]; exists {
		delete(room, client)
		if len(room) == 0 {
			delete(m.roomMap, roomID)
		}
	}
}

// SendToRoom 发送消息给聊天室
func (m *Manager) SendToRoom(roomID int64, message *Message) {
	m.mutex.RLock()
	room, exists := m.roomMap[roomID]
	m.mutex.RUnlock()

	if !exists {
		return
	}

	message.Timestamp = time.Now().Unix()
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	for client := range room {
		client.Conn.WriteMessage(websocket.BinaryMessage, data)
	}
}

// Broadcast 广播消息给所有在线用户
func (m *Manager) Broadcast(message *Message) {
	message.Timestamp = time.Now().Unix()
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	m.broadcast <- data
}

// GetClient 根据用户ID获取客户端
func (m *Manager) GetClient(userID int64) (*Client, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	client, exists := m.clients[userID]
	return client, exists
}
