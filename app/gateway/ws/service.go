package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/hertz-contrib/websocket"
)

// WsService WebSocket服务
type WsService struct {
	manager *Manager
	mu      sync.RWMutex
}

// NewWsService 创建新的WebSocket服务
func NewWsService(manager *Manager) *WsService {
	return &WsService{
		manager: manager,
		mu:      sync.RWMutex{},
	}
}

// Start 启动WebSocket服务
func (s *WsService) Start(ctx context.Context) {
	// 启动WebSocket管理器
	go s.manager.Start(ctx)
	// 启动心跳检测
	go s.manager.StartHeartbeat(ctx, pingPeriod)
}

// RegisterClient 注册WebSocket客户端
func (s *WsService) RegisterClient(userID int64, conn *websocket.Conn) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查用户是否已经在线
	if s.manager.IsUserOnline(userID) {
		return fmt.Errorf("用户 %d 已经在线", userID)
	}

	// 创建新的客户端
	client := &Client{
		UserID: userID,
		Rooms:  make(map[int64]bool),
		mu:     sync.Mutex{},
		Conn:   conn,
	}

	// 直接注册到管理器
	s.manager.RegisterClient(userID, client.Conn)

	// 同时发送到注册通道，以便其他功能使用
	go func() {
		s.manager.register <- client
	}()
	return nil
}

// JoinChatRoom 加入聊天室
func (s *WsService) JoinChatRoom(userID, roomID int64) error {
	client, exists := s.manager.GetClient(userID)
	if !exists {
		return fmt.Errorf("用户 %d 不在线", userID)
	}

	s.manager.JoinRoom(client, roomID)

	// 发送系统消息通知聊天室
	joinMsg := &Message{
		Type:      MessageTypeSystem,
		Content:   fmt.Sprintf("用户 %d 加入了聊天室", userID),
		RoomID:    roomID,
		Timestamp: time.Now().Unix(),
	}
	s.manager.SendToRoom(userID, roomID, joinMsg)

	return nil
}

// LeaveChatRoom 离开聊天室
func (s *WsService) LeaveChatRoom(userID, roomID int64) error {
	client, exists := s.manager.GetClient(userID)
	if !exists {
		return fmt.Errorf("用户 %d 不在线", userID)
	}

	s.manager.LeaveRoom(client, roomID)

	// 发送系统消息通知聊天室
	leaveMsg := &Message{
		Type:      MessageTypeSystem,
		Content:   fmt.Sprintf("用户 %d 离开了聊天室", userID),
		RoomID:    roomID,
		Timestamp: time.Now().Unix(),
	}
	s.manager.SendToRoom(userID, roomID, leaveMsg)

	return nil
}

// SendChatMessage 发送聊天消息
func (s *WsService) SendChatMessage(userID, roomID int64, content string, msgType int8) error {
	/*client, exists := s.manager.GetClient(userID)
	if !exists {
		return fmt.Errorf("用户 %d 不在线", userID)
	}

	if !client.IsInRoom(roomID) {
		return fmt.Errorf("用户 %d 不在聊天室 %d 中", userID, roomID)
	}*/

	msg := &Message{
		Type:      MessageTypeGroup,
		From:      userID,
		RoomID:    roomID,
		Content:   content,
		Timestamp: time.Now().Unix(),
		Extra: map[string]any{
			"type": msgType,
		},
	}

	s.manager.SendToRoom(userID, roomID, msg)
	return nil
}

// SendPrivateMessage 发送私信
func (s *WsService) SendPrivateMessage(fromID, toID int64, content string) error {
	msg := &Message{
		Type:      MessageTypePrivate,
		From:      fromID,
		To:        toID,
		Content:   content,
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("消息序列化失败: %w", err)
	}

	if err := s.manager.SendToUser(toID, data); err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}

	return nil
}

// Broadcast 广播消息给所有在线用户
func (s *WsService) Broadcast(content string) {
	msg := &Message{
		Type:      MessageTypeSystem,
		Content:   content,
		Timestamp: time.Now().Unix(),
	}
	s.manager.Broadcast(msg)
}

// GetOnlineUsers 获取在线用户列表
func (s *WsService) GetOnlineUsers() []int64 {
	return s.manager.GetOnlineUsers()
}

// IsUserOnline 检查用户是否在线
func (s *WsService) IsUserOnline(userID int64) bool {
	return s.manager.IsUserOnline(userID)
}

// GetManager 获取WebSocket管理器
func (s *WsService) GetManager() *Manager {
	return s.manager
}
