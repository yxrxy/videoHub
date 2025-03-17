package model

import "time"

// PrivateMessage 私信模型
type PrivateMessage struct {
	ID         int64      `gorm:"primarykey" json:"id"`
	SenderID   int64      `gorm:"index" json:"sender_id"`            // 发送者ID
	ReceiverID int64      `gorm:"index" json:"receiver_id"`          // 接收者ID
	Content    string     `gorm:"type:text" json:"content"`          // 消息内容
	IsRead     bool       `gorm:"default:false" json:"is_read"`      // 是否已读
	CreatedAt  time.Time  `json:"created_at"`                        // 创建时间
	UpdatedAt  time.Time  `json:"updated_at"`                        // 更新时间
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"` // 删除时间
}

// ChatRoom 聊天室模型
type ChatRoom struct {
	ID        int64            `gorm:"primarykey" json:"id"`
	Name      string           `gorm:"type:varchar(64)" json:"name"`      // 聊天室名称
	CreatorID int64            `gorm:"index" json:"creator_id"`           // 创建者ID
	Type      int8             `gorm:"type:tinyint;index" json:"type"`    // 类型：1=私聊,2=群聊
	CreatedAt time.Time        `json:"created_at"`                        // 创建时间
	UpdatedAt time.Time        `json:"updated_at"`                        // 更新时间
	DeletedAt *time.Time       `gorm:"index" json:"deleted_at,omitempty"` // 删除时间
	Members   []ChatRoomMember `gorm:"foreignKey:RoomID" json:"members"`  // 聊天室成员
}

// ChatRoomMember 聊天室成员模型
type ChatRoomMember struct {
	ID        int64      `gorm:"primarykey" json:"id"`
	RoomID    int64      `gorm:"uniqueIndex:idx_room_user" json:"room_id"` // 聊天室ID
	UserID    int64      `gorm:"uniqueIndex:idx_room_user" json:"user_id"` // 用户ID
	Nickname  string     `gorm:"type:varchar(32)" json:"nickname"`         // 在群里的昵称
	Role      int8       `gorm:"type:tinyint;default:0" json:"role"`       // 角色：0=普通成员,1=管理员,2=群主
	CreatedAt time.Time  `json:"created_at"`                               // 创建时间
	UpdatedAt time.Time  `json:"updated_at"`                               // 更新时间
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`        // 删除时间
}

// ChatMessage 聊天消息模型
type ChatMessage struct {
	ID        int64      `gorm:"primarykey" json:"id"`
	RoomID    int64      `gorm:"index:idx_room_created" json:"room_id"`    // 聊天室ID
	SenderID  int64      `gorm:"index" json:"sender_id"`                   // 发送者ID
	Content   string     `gorm:"type:text" json:"content"`                 // 消息内容
	Type      int8       `gorm:"type:tinyint;default:0" json:"type"`       // 消息类型：0=文本,1=图片,2=视频,3=文件
	CreatedAt time.Time  `gorm:"index:idx_room_created" json:"created_at"` // 创建时间
	UpdatedAt time.Time  `json:"updated_at"`                               // 更新时间
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`        // 删除时间
}

// Friendship 好友关系模型
type Friendship struct {
	ID        int64      `gorm:"primarykey" json:"id"`
	UserID    int64      `gorm:"uniqueIndex:idx_user_friend" json:"user_id"`   // 用户ID
	FriendID  int64      `gorm:"uniqueIndex:idx_user_friend" json:"friend_id"` // 好友ID
	Status    int8       `gorm:"type:tinyint;default:0;index" json:"status"`   // 状态：0=待确认,1=已接受,2=已拒绝,3=已拉黑
	Remark    string     `gorm:"type:varchar(32)" json:"remark"`               // 好友备注
	CreatedAt time.Time  `json:"created_at"`                                   // 创建时间
	UpdatedAt time.Time  `json:"updated_at"`                                   // 更新时间
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`            // 删除时间
}

// FriendRequest 好友申请模型
type FriendRequest struct {
	ID         int64      `gorm:"primarykey" json:"id"`
	SenderID   int64      `gorm:"index" json:"sender_id"`                     // 发送者ID
	ReceiverID int64      `gorm:"index" json:"receiver_id"`                   // 接收者ID
	Message    string     `gorm:"type:varchar(255)" json:"message"`           // 申请消息
	Status     int8       `gorm:"type:tinyint;default:0;index" json:"status"` // 状态：0=待处理,1=已接受,2=已拒绝
	CreatedAt  time.Time  `json:"created_at"`                                 // 创建时间
	UpdatedAt  time.Time  `json:"updated_at"`                                 // 更新时间
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`          // 删除时间
}

// MessageReadStatus 消息已读状态模型
type MessageReadStatus struct {
	ID        int64      `gorm:"primarykey" json:"id"`
	MessageID int64      `gorm:"uniqueIndex:idx_message_user" json:"message_id"` // 消息ID
	UserID    int64      `gorm:"uniqueIndex:idx_message_user" json:"user_id"`    // 用户ID
	IsRead    bool       `gorm:"default:false" json:"is_read"`                   // 是否已读
	ReadAt    *time.Time `json:"read_at"`                                        // 读取时间
	CreatedAt time.Time  `json:"created_at"`                                     // 创建时间
}

// 定义常量
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

// TableName 指定表名
func (PrivateMessage) TableName() string {
	return "private_messages"
}

func (ChatRoom) TableName() string {
	return "chat_rooms"
}

func (ChatRoomMember) TableName() string {
	return "chat_room_members"
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}

func (Friendship) TableName() string {
	return "friendships"
}

func (FriendRequest) TableName() string {
	return "friend_requests"
}

func (MessageReadStatus) TableName() string {
	return "message_read_status"
}
