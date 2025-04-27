package mysql

import "time"

// PrivateMessage 私信模型
type PrivateMessage struct {
	ID         int64      `json:"id"                   gorm:"primarykey"`
	SenderID   int64      `json:"sender_id"            gorm:"index"`         // 发送者ID
	ReceiverID int64      `json:"receiver_id"          gorm:"index"`         // 接收者ID
	Content    string     `json:"content"              gorm:"type:text"`     // 消息内容
	IsRead     bool       `json:"is_read"              gorm:"default:false"` // 是否已读
	CreatedAt  time.Time  `json:"created_at"`                                // 创建时间
	UpdatedAt  time.Time  `json:"updated_at"`                                // 更新时间
	DeletedAt  *time.Time `json:"deleted_at,omitempty" gorm:"index"`         // 删除时间
}

// ChatRoom 聊天室模型
type ChatRoom struct {
	ID        int64            `json:"id"                   gorm:"primarykey"`
	Name      string           `json:"name"                 gorm:"type:varchar(64)"`   // 聊天室名称
	CreatorID int64            `json:"creator_id"           gorm:"index"`              // 创建者ID
	Type      int8             `json:"type"                 gorm:"type:tinyint;index"` // 类型：1=私聊,2=群聊
	CreatedAt time.Time        `json:"created_at"`                                     // 创建时间
	UpdatedAt time.Time        `json:"updated_at"`                                     // 更新时间
	DeletedAt *time.Time       `json:"deleted_at,omitempty" gorm:"index"`              // 删除时间
	Members   []ChatRoomMember `json:"members"              gorm:"foreignKey:RoomID"`  // 聊天室成员
}

// ChatRoomMember 聊天室成员模型
type ChatRoomMember struct {
	ID        int64      `json:"id"                   gorm:"primarykey"`
	RoomID    int64      `json:"room_id"              gorm:"uniqueIndex:idx_room_user"` // 聊天室ID
	UserID    int64      `json:"user_id"              gorm:"uniqueIndex:idx_room_user"` // 用户ID
	Nickname  string     `json:"nickname"             gorm:"type:varchar(32)"`          // 在群里的昵称
	Role      int8       `json:"role"                 gorm:"type:tinyint;default:0"`    // 角色：0=普通成员,1=管理员,2=群主
	CreatedAt time.Time  `json:"created_at"`                                            // 创建时间
	UpdatedAt time.Time  `json:"updated_at"`                                            // 更新时间
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`                     // 删除时间
}

// ChatMessage 聊天消息模型
type ChatMessage struct {
	ID        int64      `json:"id"                   gorm:"primarykey"`
	RoomID    int64      `json:"room_id"              gorm:"index:idx_room_created"` // 聊天室ID
	SenderID  int64      `json:"sender_id"            gorm:"index"`                  // 发送者ID
	Content   string     `json:"content"              gorm:"type:text"`              // 消息内容
	Type      int8       `json:"type"                 gorm:"type:tinyint;default:0"` // 消息类型：0=文本,1=图片,2=视频,3=文件
	CreatedAt time.Time  `json:"created_at"           gorm:"index:idx_room_created"` // 创建时间
	UpdatedAt time.Time  `json:"updated_at"`                                         // 更新时间
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`                  // 删除时间
}

// Friendship 好友关系模型
type Friendship struct {
	ID        int64      `json:"id"                   gorm:"primarykey"`
	UserID    int64      `json:"user_id"              gorm:"uniqueIndex:idx_user_friend"`  // 用户ID
	FriendID  int64      `json:"friend_id"            gorm:"uniqueIndex:idx_user_friend"`  // 好友ID
	Status    int8       `json:"status"               gorm:"type:tinyint;default:0;index"` // 状态：0=待确认,1=已接受,2=已拒绝,3=已拉黑
	Remark    string     `json:"remark"               gorm:"type:varchar(32)"`             // 好友备注
	CreatedAt time.Time  `json:"created_at"`                                               // 创建时间
	UpdatedAt time.Time  `json:"updated_at"`                                               // 更新时间
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`                        // 删除时间
}

// FriendRequest 好友申请模型
type FriendRequest struct {
	ID         int64      `json:"id"                   gorm:"primarykey"`
	SenderID   int64      `json:"sender_id"            gorm:"index"`                        // 发送者ID
	ReceiverID int64      `json:"receiver_id"          gorm:"index"`                        // 接收者ID
	Message    string     `json:"message"              gorm:"type:varchar(255)"`            // 申请消息
	Status     int8       `json:"status"               gorm:"type:tinyint;default:0;index"` // 状态：0=待处理,1=已接受,2=已拒绝
	CreatedAt  time.Time  `json:"created_at"`                                               // 创建时间
	UpdatedAt  time.Time  `json:"updated_at"`                                               // 更新时间
	DeletedAt  *time.Time `json:"deleted_at,omitempty" gorm:"index"`                        // 删除时间
}

// MessageReadStatus 消息已读状态模型
type MessageReadStatus struct {
	ID        int64      `json:"id"         gorm:"primarykey"`
	MessageID int64      `json:"message_id" gorm:"uniqueIndex:idx_message_user"` // 消息ID
	UserID    int64      `json:"user_id"    gorm:"uniqueIndex:idx_message_user"` // 用户ID
	IsRead    bool       `json:"is_read"    gorm:"default:false"`                // 是否已读
	ReadAt    *time.Time `json:"read_at"`                                        // 读取时间
	CreatedAt time.Time  `json:"created_at"`                                     // 创建时间
}

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
