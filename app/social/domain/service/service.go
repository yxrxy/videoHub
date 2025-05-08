package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yxrxy/videoHub/app/social/domain/model"
	"github.com/yxrxy/videoHub/app/social/infrastructure/ws"
)

func (s *SocialService) BroadcastSystemMessage(content string) {
	msg := &ws.Message{
		Type:      ws.MessageTypeSystem,
		Content:   content,
		Timestamp: time.Now().Unix(),
	}
	s.wsService.Broadcast(msg.Content)
}

func (s *SocialService) SavePrivateMessage(ctx context.Context, senderID, receiverID int64, content string) error {
	msg := &model.PrivateMessage{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		IsRead:     false,
	}
	if err := s.db.SendPrivateMessage(ctx, msg); err != nil {
		return err
	}
	return nil
}

func (s *SocialService) SendMessage(ctx context.Context, senderID, receiverID int64, content string) error {
	// 序列化消息
	wsMsg := &ws.Message{
		Type:    ws.MessageTypePrivate,
		From:    senderID,
		To:      receiverID,
		Content: content,
	}

	data, err := json.Marshal(wsMsg)
	if err != nil {
		return err
	}

	// 发送序列化后的消息
	if err := s.wsService.SendPrivateMessage(senderID, receiverID, string(data)); err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}
	return nil
}

// CreateChatRoom 创建聊天室成员
func (s *SocialService) CreateChatRoom(ctx context.Context, name string, creatorID int64, memberIDs []int64) ([]model.ChatRoomMember, error) {
	members := make([]model.ChatRoomMember, 0, len(memberIDs)+1)
	// 添加创建者
	members = append(members, model.ChatRoomMember{
		UserID:   creatorID,
		Role:     model.ChatRoomRoleCreator,
		Nickname: name, // 可以从用户服务获取昵称
	})

	// 添加其他成员
	for _, memberID := range memberIDs {
		if memberID != creatorID {
			members = append(members, model.ChatRoomMember{
				UserID:   memberID,
				Role:     model.ChatRoomRoleMember,
				Nickname: name, // 可以从用户服务获取昵称
			})
		}
	}

	return members, nil
}
