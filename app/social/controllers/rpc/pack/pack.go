package pack

import (
	"github.com/yxrxy/videoHub/app/social/domain/model"
	rpcmodel "github.com/yxrxy/videoHub/kitex_gen/model"
)

func PackPrivateMessage(msg *model.PrivateMessage) *rpcmodel.PrivateMessage {
	return &rpcmodel.PrivateMessage{
		Id:         msg.ID,
		SenderId:   msg.SenderID,
		ReceiverId: msg.ReceiverID,
		Content:    msg.Content,
		IsRead:     msg.IsRead,
	}
}

func PackPrivateMessageList(msgs []*model.PrivateMessage) []*rpcmodel.PrivateMessage {
	messageList := make([]*rpcmodel.PrivateMessage, 0, len(msgs))
	for _, msg := range msgs {
		messageList = append(messageList, PackPrivateMessage(msg))
	}
	return messageList
}

func PackChatRoom(room *model.ChatRoom) *rpcmodel.ChatRoom {
	return &rpcmodel.ChatRoom{
		Id:        room.ID,
		Name:      room.Name,
		Type:      room.Type,
		CreatorId: room.CreatorID,
	}
}

func PackChatRoomList(rooms []*model.ChatRoom) []*rpcmodel.ChatRoom {
	roomList := make([]*rpcmodel.ChatRoom, 0, len(rooms))
	for _, room := range rooms {
		roomList = append(roomList, PackChatRoom(room))
	}
	return roomList
}

func PackChatMessage(msg *model.ChatMessage) *rpcmodel.ChatMessage {
	return &rpcmodel.ChatMessage{
		Id:       msg.ID,
		RoomId:   msg.RoomID,
		SenderId: msg.SenderID,
		Content:  msg.Content,
	}
}

func PackChatMessageList(msgs []*model.ChatMessage) []*rpcmodel.ChatMessage {
	messageList := make([]*rpcmodel.ChatMessage, 0, len(msgs))
	for _, msg := range msgs {
		messageList = append(messageList, PackChatMessage(msg))
	}
	return messageList
}

func PackFriend(friend *model.Friendship) *rpcmodel.Friendship {
	return &rpcmodel.Friendship{
		Id:       friend.ID,
		UserId:   friend.UserID,
		FriendId: friend.FriendID,
		Status:   int8(friend.Status),
	}
}

func PackFriendList(friends []*model.Friendship) []*rpcmodel.Friendship {
	friendList := make([]*rpcmodel.Friendship, 0, len(friends))
	for _, friend := range friends {
		friendList = append(friendList, PackFriend(friend))
	}
	return friendList
}

func PackFriendRequest(request *model.FriendRequest) *rpcmodel.FriendRequest {
	return &rpcmodel.FriendRequest{
		Id:         request.ID,
		SenderId:   request.SenderID,
		ReceiverId: request.ReceiverID,
		Message:    &request.Message,
		Status:     int8(request.Status),
	}
}

func PackFriendRequestList(requests []*model.FriendRequest) []*rpcmodel.FriendRequest {
	requestList := make([]*rpcmodel.FriendRequest, 0, len(requests))
	for _, request := range requests {
		requestList = append(requestList, PackFriendRequest(request))
	}
	return requestList
}

func PackFriendship(friendship *model.Friendship) *rpcmodel.Friendship {
	return &rpcmodel.Friendship{
		Id:       friendship.ID,
		UserId:   friendship.UserID,
		FriendId: friendship.FriendID,
		Status:   int8(friendship.Status),
	}
}

func PackFriendshipList(friendships []*model.Friendship) []*rpcmodel.Friendship {
	friendshipList := make([]*rpcmodel.Friendship, 0, len(friendships))
	for _, friendship := range friendships {
		friendshipList = append(friendshipList, PackFriendship(friendship))
	}
	return friendshipList
}
