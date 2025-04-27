namespace go api.social

include "../social.thrift"

// API 服务
service SocialAPI {
    // 私信相关接口
    social.SendPrivateMessageResponse SendPrivateMessage(1: social.SendPrivateMessageRequest request) (api.post="/api/v1/social/private/message")
    social.GetPrivateMessagesResponse GetPrivateMessages(1: social.GetPrivateMessagesRequest request) (api.get="/api/v1/social/private/messages")
    
    // 聊天室相关接口
    social.CreateChatRoomResponse CreateChatRoom(1: social.CreateChatRoomRequest request) (api.post="/api/v1/social/chatroom")
    social.GetChatRoomResponse GetChatRoom(1: social.GetChatRoomRequest request) (api.get="/api/v1/social/chatroom/:room_id")
    social.GetUserChatRoomsResponse GetUserChatRooms(1: social.GetUserChatRoomsRequest request) (api.get="/api/v1/social/chatrooms")
    social.SendChatMessageResponse SendChatMessage(1: social.SendChatMessageRequest request) (api.post="/api/v1/social/chatroom/message")
    social.GetChatMessagesResponse GetChatMessages(1: social.GetChatMessagesRequest request) (api.get="/api/v1/social/chatroom/:room_id/messages")
    
    // 好友相关接口
    social.AddFriendResponse AddFriend(1: social.AddFriendRequest request) (api.post="/api/v1/social/friend")
    social.GetFriendshipResponse GetFriendship(1: social.GetFriendshipRequest request) (api.get="/api/v1/social/friendship/:friend_id")
    social.GetUserFriendsResponse GetUserFriends(1: social.GetUserFriendsRequest request) (api.get="/api/v1/social/friends")
    
    // 好友申请相关接口
    social.CreateFriendRequestResponse CreateFriendRequest(1: social.CreateFriendRequestRequest request) (api.post="/api/v1/social/friend/request")
    social.GetFriendRequestsResponse GetFriendRequests(1: social.GetFriendRequestsRequest request) (api.get="/api/v1/social/friend/requests")
    social.HandleFriendRequestResponse HandleFriendRequest(1: social.HandleFriendRequestRequest request) (api.put="/api/v1/social/friend/request/:request_id")
    
    // 消息状态相关接口
    social.MarkMessageReadResponse MarkMessageRead(1: social.MarkMessageReadRequest request) (api.put="/api/v1/social/message/:message_id/read")
    social.GetUnreadMessageCountResponse GetUnreadMessageCount(1: social.GetUnreadMessageCountRequest request) (api.get="/api/v1/social/messages/unread/count")
}
