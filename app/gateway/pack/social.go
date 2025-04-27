package pack

import (
	"github.com/yxrxy/videoHub/app/gateway/model/model"
	rpcmodel "github.com/yxrxy/videoHub/kitex_gen/model"
	"github.com/yxrxy/videoHub/pkg/base"
)

// BuildPrivateMessage 将 RPC 交流实体转换成 http 返回的实体
func BuildPrivateMessage(m *model.PrivateMessage) *rpcmodel.PrivateMessage {
	if m == nil {
		return nil
	}
	return &rpcmodel.PrivateMessage{
		Id:         m.ID,
		SenderId:   m.SenderID,
		ReceiverId: m.ReceiverID,
		Content:    m.Content,
		IsRead:     m.IsRead,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		DeletedAt:  m.DeletedAt,
	}
}

// BuildPrivateMessageList 构建私信列表
func BuildPrivateMessageList(ms []*model.PrivateMessage) []*rpcmodel.PrivateMessage {
	return base.BuildTypeList(ms, BuildPrivateMessage)
}

// BuildChatRoom 将 RPC 交流实体转换成 http 返回的实体
func BuildChatRoom(r *model.ChatRoom) *rpcmodel.ChatRoom {
	if r == nil {
		return nil
	}
	return &rpcmodel.ChatRoom{
		Id:        r.ID,
		Name:      r.Name,
		CreatorId: r.CreatorID,
		Type:      r.Type,
		Members:   BuildChatRoomMemberList(r.Members),
	}
}

// BuildChatRoomList 构建聊天室列表
func BuildChatRoomList(rs []*model.ChatRoom) []*rpcmodel.ChatRoom {
	return base.BuildTypeList(rs, BuildChatRoom)
}

// BuildChatRoomMember 将 RPC 交流实体转换成 http 返回的实体
func BuildChatRoomMember(m *model.ChatRoomMember) *rpcmodel.ChatRoomMember {
	if m == nil {
		return nil
	}
	return &rpcmodel.ChatRoomMember{
		Id:        m.ID,
		RoomId:    m.RoomID,
		UserId:    m.UserID,
		Nickname:  m.Nickname,
		Role:      m.Role,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}

// BuildChatRoomMemberList 构建聊天室成员列表
func BuildChatRoomMemberList(ms []*model.ChatRoomMember) []*rpcmodel.ChatRoomMember {
	return base.BuildTypeList(ms, BuildChatRoomMember)
}

// BuildChatMessage 将 RPC 交流实体转换成 http 返回的实体
func BuildChatMessage(m *model.ChatMessage) *rpcmodel.ChatMessage {
	if m == nil {
		return nil
	}
	return &rpcmodel.ChatMessage{
		Id:        m.ID,
		RoomId:    m.RoomID,
		SenderId:  m.SenderID,
		Content:   m.Content,
		Type:      m.Type,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}

// BuildChatMessageList 构建聊天消息列表
func BuildChatMessageList(ms []*model.ChatMessage) []*rpcmodel.ChatMessage {
	return base.BuildTypeList(ms, BuildChatMessage)
}

// BuildFriendship 将 RPC 交流实体转换成 http 返回的实体
func BuildFriendship(f *model.Friendship) *rpcmodel.Friendship {
	if f == nil {
		return nil
	}
	return &rpcmodel.Friendship{
		Id:        f.ID,
		UserId:    f.UserID,
		FriendId:  f.FriendID,
		Status:    f.Status,
		Remark:    f.Remark,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
		DeletedAt: f.DeletedAt,
	}
}

// BuildFriendshipList 构建好友关系列表
func BuildFriendshipList(fs []*model.Friendship) []*rpcmodel.Friendship {
	return base.BuildTypeList(fs, BuildFriendship)
}

// BuildFriendRequest 将 RPC 交流实体转换成 http 返回的实体
func BuildFriendRequest(r *model.FriendRequest) *rpcmodel.FriendRequest {
	if r == nil {
		return nil
	}
	return &rpcmodel.FriendRequest{
		Id:         r.ID,
		SenderId:   r.SenderID,
		ReceiverId: r.ReceiverID,
		Message:    r.Message,
		Status:     r.Status,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
		DeletedAt:  r.DeletedAt,
	}
}

// BuildFriendRequestList 构建好友请求列表
func BuildFriendRequestList(rs []*model.FriendRequest) []*rpcmodel.FriendRequest {
	return base.BuildTypeList(rs, BuildFriendRequest)
}
