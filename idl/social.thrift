namespace go social

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

// 基础消息结构
struct Message {
    1: required i64 id                   // 消息ID
    2: required i64 sender_id            // 发送者ID
    3: required string content           // 消息内容
    4: required i64 created_at           // 创建时间
    5: required bool is_read             // 是否已读
}

// 私信消息
struct PrivateMessage {
    1: required i64 id                   // 消息ID
    2: required i64 sender_id            // 发送者ID
    3: required i64 receiver_id          // 接收者ID
    4: required string content           // 消息内容
    5: required bool is_read             // 是否已读
    6: required i64 created_at           // 创建时间
    7: required i64 updated_at           // 更新时间
    8: optional i64 deleted_at           // 删除时间
}

// 聊天室
struct ChatRoom {
    1: required i64 id                   // 聊天室ID
    2: required string name              // 聊天室名称
    3: required i64 creator_id           // 创建者ID
    4: required i8 type                  // 类型：1=私聊,2=群聊
    5: required i64 created_at           // 创建时间
    6: required i64 updated_at           // 更新时间
    7: optional i64 deleted_at           // 删除时间
    8: optional list<ChatRoomMember> members // 成员列表
}

// 聊天室成员
struct ChatRoomMember {
    1: required i64 id                   // 成员ID
    2: required i64 room_id              // 聊天室ID
    3: required i64 user_id              // 用户ID
    4: optional string nickname          // 在群里的昵称
    5: required i8 role                  // 角色：0=普通成员,1=管理员,2=群主
    6: required i64 created_at           // 创建时间
    7: required i64 updated_at           // 更新时间
    8: optional i64 deleted_at           // 删除时间
}

// 聊天消息
struct ChatMessage {
    1: required i64 id                   // 消息ID
    2: required i64 room_id              // 聊天室ID
    3: required i64 sender_id            // 发送者ID
    4: required string content           // 消息内容
    5: required i8 type                  // 消息类型：0=文本,1=图片,2=视频,3=文件
    6: required i64 created_at           // 创建时间
    7: required i64 updated_at           // 更新时间
    8: optional i64 deleted_at           // 删除时间
}

// 好友关系
struct Friendship {
    1: required i64 id                   // 关系ID
    2: required i64 user_id              // 用户ID
    3: required i64 friend_id            // 好友ID
    4: required i8 status                // 状态：0=待确认,1=已接受,2=已拒绝,3=已拉黑
    5: optional string remark            // 好友备注
    6: required i64 created_at           // 创建时间
    7: required i64 updated_at           // 更新时间
    8: optional i64 deleted_at           // 删除时间
}

// 好友申请
struct FriendRequest {
    1: required i64 id                   // 申请ID
    2: required i64 sender_id            // 发送者ID
    3: required i64 receiver_id          // 接收者ID
    4: optional string message           // 申请消息
    5: required i8 status                // 状态：0=待处理,1=已接受,2=已拒绝
    6: required i64 created_at           // 创建时间
    7: required i64 updated_at           // 更新时间
    8: optional i64 deleted_at           // 删除时间
}

// 发送私信请求
struct SendPrivateMessageRequest {
    1: required i64 sender_id            // 发送者ID
    2: required i64 receiver_id          // 接收者ID
    3: required string content           // 消息内容
}

// 发送私信响应
struct SendPrivateMessageResponse {
    1: required PrivateMessage message   // 创建的消息
}

// 获取私信列表请求
struct GetPrivateMessagesRequest {
    1: required i64 user_id              // 用户ID
    2: required i64 other_user_id        // 对方用户ID
    3: optional i64 last_id              // 上一页最后一条消息ID
    4: optional i32 limit                // 每页数量
}

// 获取私信列表响应
struct PrivateMessagesResponse {
    1: required list<PrivateMessage> messages // 消息列表
    2: required i64 total                     // 总数
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
    1: required ChatRoom room            // 创建的聊天室
}

// 获取聊天室请求
struct GetChatRoomRequest {
    1: required i64 room_id              // 聊天室ID
    2: required i64 user_id              // 请求用户ID
}

// 获取聊天室响应
struct GetChatRoomResponse {
    1: required ChatRoom room            // 聊天室信息
}

// 获取用户聊天室列表请求
struct GetUserChatRoomsRequest {
    1: required i64 user_id              // 用户ID
    2: optional i32 page                 // 页码
    3: optional i32 size                 // 每页数量
}

// 获取用户聊天室列表响应
struct GetUserChatRoomsResponse {
    1: required list<ChatRoom> rooms     // 聊天室列表
    2: required i64 total                // 总数
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
    1: required ChatMessage message      // 创建的消息
}

// 获取聊天消息列表请求
struct GetChatMessagesRequest {
    1: required i64 room_id              // 聊天室ID
    2: required i64 user_id              // 用户ID
    3: optional i64 last_id              // 上一页最后一条消息ID
    4: optional i32 limit                // 每页数量
}

// 获取聊天消息列表响应
struct ChatMessagesResponse {
    1: required list<ChatMessage> messages // 消息列表
    2: required i64 total                  // 总数
}

// 添加好友请求
struct AddFriendRequest {
    1: required i64 user_id              // 用户ID
    2: required i64 friend_id            // 好友ID
    3: optional string remark            // 好友备注
}

// 添加好友响应
struct AddFriendResponse {
    1: required Friendship friendship    // 创建的好友关系
}

// 获取用户好友列表请求
struct GetUserFriendsRequest {
    1: required i64 user_id              // 用户ID
    2: optional i32 page                 // 页码
    3: optional i32 size                 // 每页数量
}

// 获取用户好友列表响应
struct GetUserFriendsResponse {
    1: required list<Friendship> friendships // 好友关系列表
    2: required i64 total                    // 总数
}

// 获取好友关系请求
struct GetFriendshipRequest {
    1: required i64 user_id              // 用户ID
    2: required i64 friend_id            // 好友ID
}

// 获取好友关系响应
struct GetFriendshipResponse {
    1: required Friendship friendship    // 好友关系
}

// 创建好友申请请求
struct CreateFriendRequestRequest {
    1: required i64 sender_id            // 发送者ID
    2: required i64 receiver_id          // 接收者ID
    3: optional string message           // 申请消息
}

// 创建好友申请响应
struct CreateFriendRequestResponse {
    1: required FriendRequest request    // 创建的好友申请
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
    1: required list<FriendRequest> requests // 好友申请列表
    2: required i64 total                    // 总数
}

// 处理好友申请请求
struct HandleFriendRequestRequest {
    1: required i64 request_id           // 申请ID
    2: required i64 user_id              // 处理用户ID
    3: required i8 action                // 操作：1=接受,2=拒绝
}

// 处理好友申请响应
struct HandleFriendRequestResponse {
    1: required FriendRequest request    // 更新后的好友申请
    2: optional Friendship friendship    // 如果接受，创建的好友关系
}

// 标记消息已读请求
struct MarkMessageReadRequest {
    1: required i64 message_id           // 消息ID
    2: required i64 user_id              // 用户ID
}

// 标记消息已读响应
struct MarkMessageReadResponse {
    1: required bool success             // 是否成功
}

// 获取未读消息数请求
struct GetUnreadMessageCountRequest {
    1: required i64 user_id              // 用户ID
}

// 获取未读消息数响应
struct GetUnreadMessageCountResponse {
    1: required i64 count                // 未读消息数
}

// WebSocket消息
struct WebSocketMessage {
    1: required string type              // 消息类型
    2: required string data              // 消息数据
}

// 通用异常
exception SocialException {
    1: required string message
}

// 社交服务
service SocialService {
    // 私信相关
    void SendPrivateMessage(1: i64 sender_id, 2: i64 receiver_id, 3: string content) throws (1: SocialException e)
    PrivateMessagesResponse GetPrivateMessages(1: i64 sender_id, 2: i64 receiver_id, 3: i32 page, 4: i32 size) throws (1: SocialException e)
    
    // 聊天室相关
    ChatRoom CreateChatRoom(1: string name, 2: i64 creator_id, 3: i8 room_type, 4: list<i64> member_ids) throws (1: SocialException e)
    ChatRoom GetChatRoom(1: i64 room_id) throws (1: SocialException e)
    list<ChatRoom> GetUserChatRooms(1: i64 user_id) throws (1: SocialException e)
    void SendChatMessage(1: i64 room_id, 2: i64 sender_id, 3: string content, 4: i8 msg_type) throws (1: SocialException e)
    ChatMessagesResponse GetChatMessages(1: i64 room_id, 2: i32 page, 3: i32 size) throws (1: SocialException e)
    
    // 好友相关
    void AddFriend(1: i64 user_id, 2: i64 friend_id) throws (1: SocialException e)
    Friendship GetFriendship(1: i64 user_id, 2: i64 friend_id) throws (1: SocialException e)
    list<Friendship> GetUserFriends(1: i64 user_id) throws (1: SocialException e)
    
    // 好友申请相关
    void CreateFriendRequest(1: i64 sender_id, 2: i64 receiver_id, 3: string message) throws (1: SocialException e)
    list<FriendRequest> GetFriendRequests(1: i64 user_id, 2: i8 status) throws (1: SocialException e)
    void HandleFriendRequest(1: i64 request_id, 2: i8 status) throws (1: SocialException e)
    
    // 消息状态相关
    void MarkMessageRead(1: i64 message_id, 2: i64 user_id) throws (1: SocialException e)
    i64 GetUnreadMessageCount(1: i64 user_id) throws (1: SocialException e)
}