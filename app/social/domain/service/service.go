package service

import (
	"context"

	"github.com/yxrxy/videoHub/app/social/domain/model"
)

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
