package model

// MessageType 消息类型
type MessageType int8

const (
	// 聊天室类型
	ChatRoomTypePrivate = 1 // 私聊
	ChatRoomTypeGroup   = 2 // 群聊

	// 聊天室成员角色
	ChatRoomRoleMember  = 0 // 普通成员
	ChatRoomRoleAdmin   = 1 // 管理员
	ChatRoomRoleCreator = 2 // 群主

	// 消息类型
	MessageTypeText  = 0 // 文本
	MessageTypeImage = 1 // 图片
	MessageTypeVideo = 2 // 视频
	MessageTypeFile  = 3 // 文件

	// 好友状态
	FriendStatusPending  = 0 // 待确认
	FriendStatusAccepted = 1 // 已接受
	FriendStatusRejected = 2 // 已拒绝
	FriendStatusBlocked  = 3 // 已拉黑

	// 好友申请状态
	FriendRequestStatusPending  = 0 // 待处理
	FriendRequestStatusAccepted = 1 // 已接受
	FriendRequestStatusRejected = 2 // 已拒绝
)

// PrivateMessage 私信模型
type PrivateMessage struct {
	ID         int64  // 主键ID
	SenderID   int64  // 发送者ID
	ReceiverID int64  // 接收者ID
	Content    string // 消息内容
	IsRead     bool   // 是否已读
}

// ChatRoom 聊天室模型
type ChatRoom struct {
	ID        int64            // 主键ID
	Name      string           // 聊天室名称
	CreatorID int64            // 创建者ID
	Type      int8             // 类型：1=私聊,2=群聊
	Members   []ChatRoomMember // 聊天室成员
}

// ChatRoomMember 聊天室成员模型
type ChatRoomMember struct {
	ID       int64  // 主键ID
	RoomID   int64  // 聊天室ID
	UserID   int64  // 用户ID
	Nickname string // 在群里的昵称
	Role     int8   // 角色：0=普通成员,1=管理员,2=群主
}

// ChatMessage 聊天消息模型
type ChatMessage struct {
	ID       int64  // 主键ID
	RoomID   int64  // 聊天室ID
	SenderID int64  // 发送者ID
	Content  string // 消息内容
	Type     int8   // 消息类型：0=文本,1=图片,2=视频,3=文件
}

// Friendship 好友关系模型
type Friendship struct {
	ID       int64  // 主键ID
	UserID   int64  // 用户ID
	FriendID int64  // 好友ID
	Status   int8   // 状态：0=待确认,1=已接受,2=已拒绝,3=已拉黑
	Remark   string // 好友备注
}

// FriendRequest 好友申请模型
type FriendRequest struct {
	ID         int64  // 主键ID
	SenderID   int64  // 发送者ID
	ReceiverID int64  // 接收者ID
	Message    string // 申请消息
	Status     int8   // 状态：0=待处理,1=已接受,2=已拒绝
}

// MessageReadStatus 消息已读状态模型
type MessageReadStatus struct {
	ID        int64 // 主键ID
	MessageID int64 // 消息ID
	UserID    int64 // 用户ID
	IsRead    bool  // 是否已读
}
