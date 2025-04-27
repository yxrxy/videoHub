namespace go social

include "model.thrift"

// 消息类型
enum MessageType {
    TEXT = 0,      // 文本消息
    IMAGE = 1,     // 图片消息
    VIDEO = 2,     // 视频消息
    FILE = 3,      // 文件消息
    SYSTEM = 4,    // 系统消息
    PRIVATE = 5,   // 私信
    GROUP = 6      // 群聊消息
}

// 聊天室类型
enum RoomType {
    PRIVATE = 1,   // 私聊
    GROUP = 2      // 群聊
}

// 好友关系状态
enum FriendshipStatus {
    PENDING = 0,   // 待确认
    ACCEPTED = 1,  // 已接受
    REJECTED = 2,  // 已拒绝
    BLOCKED = 3    // 已拉黑
}

// 发送私信请求
struct SendPrivateMessageRequest {
    1: required i64 sender_id            // 发送者ID
    2: required i64 receiver_id          // 接收者ID
    3: required string content           // 消息内容
}

// 发送私信响应
struct SendPrivateMessageResponse {
    1: required model.BaseResp Base      // 基本响应信息
    2: required model.PrivateMessage Message   // 创建的消息
}

// 获取私信列表请求
struct GetPrivateMessagesRequest {
    1: required i64 user_id              // 用户ID
    2: required i64 other_user_id        // 对方用户ID
    3: optional i32 page                 // 页码
    4: optional i32 size                 // 每页数量
}

// 获取私信列表响应
struct GetPrivateMessagesResponse {
    1: required model.BaseResp Base              // 基本响应信息
    2: required list<model.PrivateMessage> MessageList // 消息列表
    3: required i64 Total                        // 总数
}

// 创建聊天室请求
struct CreateChatRoomRequest {
    1: required i64 creator_id           // 创建者ID
    2: required string name              // 聊天室名称
    3: required i8 type                  // 类型：1=私聊,2=群聊
    4: optional list<i64> member_ids     // 初始成员ID列表
}

// 创建聊天室响应
struct CreateChatRoomResponse {
    1: required model.BaseResp Base      // 基本响应信息
    2: required model.ChatRoom Room      // 创建的聊天室
}

// 获取聊天室请求
struct GetChatRoomRequest {
    1: required i64 room_id              // 聊天室ID
    2: required i64 user_id              // 请求用户ID
}

// 获取聊天室响应
struct GetChatRoomResponse {
    1: required model.BaseResp Base      // 基本响应信息
    2: required model.ChatRoom Room      // 聊天室信息
}

// 获取用户聊天室列表请求
struct GetUserChatRoomsRequest {
    1: required i64 user_id              // 用户ID
    2: optional i32 page                 // 页码
    3: optional i32 size                 // 每页数量
}

// 获取用户聊天室列表响应
struct GetUserChatRoomsResponse {
    1: required model.BaseResp Base      // 基本响应信息
    2: required list<model.ChatRoom> RoomList  // 聊天室列表
    3: required i64 Total                // 总数
}

// 发送聊天消息请求
struct SendChatMessageRequest {
    1: required i64 room_id              // 聊天室ID
    2: required i64 sender_id            // 发送者ID
    3: required string content           // 消息内容
    4: optional i8 type                  // 消息类型：0=文本,1=图片,2=视频,3=文件
}

// 发送聊天消息响应
struct SendChatMessageResponse {
    1: required model.BaseResp Base      // 基本响应信息
    2: required model.ChatMessage Message      // 创建的消息
}

// 获取聊天消息列表请求
struct GetChatMessagesRequest {
    1: required i64 room_id              // 聊天室ID
    2: required i64 user_id              // 用户ID
    3: optional i64 page           
    4: optional i32 size              
}

// 获取聊天消息列表响应
struct GetChatMessagesResponse {
    1: required model.BaseResp Base              // 基本响应信息
    2: required list<model.ChatMessage> MessageList    // 消息列表
    3: required i64 Total                        // 总数
}

// 添加好友请求
struct AddFriendRequest {
    1: required i64 user_id              // 用户ID
    2: required i64 friend_id            // 好友ID
    3: optional string remark            // 好友备注
}

// 添加好友响应
struct AddFriendResponse {
    1: required model.BaseResp Base      // 基本响应信息
    2: required model.Friendship Friend        // 创建的好友关系
}

// 获取用户好友列表请求
struct GetUserFriendsRequest {
    1: required i64 user_id              // 用户ID
    2: optional i32 page                 // 页码
    3: optional i32 size                 // 每页数量
}

// 获取用户好友列表响应
struct GetUserFriendsResponse {
    1: required model.BaseResp Base              // 基本响应信息
    2: required list<model.Friendship> FriendList      // 好友关系列表
    3: required i64 Total                        // 总数
}

// 获取好友关系请求
struct GetFriendshipRequest {
    1: required i64 user_id              // 用户ID
    2: required i64 friend_id            // 好友ID
}

// 获取好友关系响应
struct GetFriendshipResponse {
    1: required model.BaseResp Base      // 基本响应信息
    2: required model.Friendship Friend        // 好友关系
}

// 创建好友申请请求
struct CreateFriendRequestRequest {
    1: required i64 sender_id            // 发送者ID
    2: required i64 receiver_id          // 接收者ID
    3: optional string message           // 申请消息
}

// 创建好友申请响应
struct CreateFriendRequestResponse {
    1: required model.BaseResp Base      // 基本响应信息
    2: required model.FriendRequest Request    // 创建的好友申请
}

// 获取好友申请列表请求
struct GetFriendRequestsRequest {
    1: required i64 user_id              // 用户ID
    2: optional i8 type                  // 类型：0=收到的,1=发出的
    3: optional i32 page                 // 页码
    4: optional i32 size                 // 每页数量
}

// 获取好友申请列表响应
struct GetFriendRequestsResponse {
    1: required model.BaseResp Base              // 基本响应信息
    2: required list<model.FriendRequest> RequestList  // 好友申请列表
    3: required i64 Total                        // 总数
}

// 处理好友申请请求
struct HandleFriendRequestRequest {
    1: required i64 request_id           // 申请ID
    2: required i64 user_id              // 处理用户ID
    3: required i8 action                // 操作：1=接受,2=拒绝
}

// 处理好友申请响应
struct HandleFriendRequestResponse {
    1: required model.BaseResp Base      // 基本响应信息
    2: required model.FriendRequest Request    // 更新后的好友申请
    3: optional model.Friendship Friend        // 如果接受，创建的好友关系
}

// 标记消息已读请求
struct MarkMessageReadRequest {
    1: required i64 message_id           // 消息ID
    2: required i64 user_id              // 用户ID
}

// 标记消息已读响应
struct MarkMessageReadResponse {
    1: required model.BaseResp Base      // 基本响应信息
}

// 获取未读消息数请求
struct GetUnreadMessageCountRequest {
    1: required i64 user_id              // 用户ID
}

// 获取未读消息数响应
struct GetUnreadMessageCountResponse {
    1: required model.BaseResp Base      // 基本响应信息
    2: required i64 Count                // 未读消息数
}

// 加入聊天室请求
struct JoinChatRoomRequest {
    1: required i64 user_id              // 用户ID
    2: required i64 room_id              // 聊天室ID
}

// 加入聊天室响应
struct JoinChatRoomResponse {
    1: required model.BaseResp Base      // 基本响应信息
}

// 离开聊天室请求
struct LeaveChatRoomRequest {
    1: required i64 user_id              // 用户ID
    2: required i64 room_id              // 聊天室ID
}

// 离开聊天室响应
struct LeaveChatRoomResponse {
    1: required model.BaseResp Base      // 基本响应信息
}

// 注册WebSocket客户端请求
struct RegisterWebSocketClientRequest {
    1: required i64 user_id              // 用户ID
}

// 注册WebSocket客户端响应
struct RegisterWebSocketClientResponse {
    1: required model.BaseResp Base      // 基本响应信息
}

// 社交服务
service SocialService {
    // 私信相关
    SendPrivateMessageResponse SendPrivateMessage(1: SendPrivateMessageRequest req)
    GetPrivateMessagesResponse GetPrivateMessages(1: GetPrivateMessagesRequest req)
    
    // 聊天室相关
    CreateChatRoomResponse CreateChatRoom(1: CreateChatRoomRequest req)
    GetChatRoomResponse GetChatRoom(1: GetChatRoomRequest req)
    GetUserChatRoomsResponse GetUserChatRooms(1: GetUserChatRoomsRequest req)
    SendChatMessageResponse SendChatMessage(1: SendChatMessageRequest req)
    GetChatMessagesResponse GetChatMessages(1: GetChatMessagesRequest req)
    
    // 好友相关
    AddFriendResponse AddFriend(1: AddFriendRequest req)
    GetFriendshipResponse GetFriendship(1: GetFriendshipRequest req)
    GetUserFriendsResponse GetUserFriends(1: GetUserFriendsRequest req)
    
    // 好友申请相关
    CreateFriendRequestResponse CreateFriendRequest(1: CreateFriendRequestRequest req)
    GetFriendRequestsResponse GetFriendRequests(1: GetFriendRequestsRequest req)
    HandleFriendRequestResponse HandleFriendRequest(1: HandleFriendRequestRequest req)
    
    // 消息状态相关
    MarkMessageReadResponse MarkMessageRead(1: MarkMessageReadRequest req)
    GetUnreadMessageCountResponse GetUnreadMessageCount(1: GetUnreadMessageCountRequest req)

    // 聊天室WebSocket相关
    JoinChatRoomResponse JoinChatRoom(1: JoinChatRoomRequest request)
    LeaveChatRoomResponse LeaveChatRoom(1: LeaveChatRoomRequest request)
    RegisterWebSocketClientResponse RegisterWebSocketClient(1: RegisterWebSocketClientRequest request)
}