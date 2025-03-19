package model

// Like 点赞模型
type Like struct {
	ID        int64  `gorm:"primarykey;column:id;comment:点赞ID"`
	UserID    int64  `gorm:"index:idx_user_video;not null;column:user_id;comment:用户ID"`
	VideoID   int64  `gorm:"index:idx_user_video;not null;column:video_id;comment:视频ID"`
	DeletedAt *int64 `gorm:"column:deleted_at;comment:删除时间"`
}

// TableName 指定表名
func (Like) TableName() string {
	return "likes"
}

// Comment 评论模型
type Comment struct {
	ID        int64  `gorm:"primarykey;column:id;comment:评论ID"`
	UserID    int64  `gorm:"index:idx_video;not null;column:user_id;comment:用户ID"`
	VideoID   int64  `gorm:"index:idx_video;not null;column:video_id;comment:视频ID"`
	Content   string `gorm:"type:text;not null;column:content;comment:评论内容"`
	ParentID  *int64 `gorm:"column:parent_id;comment:父评论ID"`
	LikeCount int32  `gorm:"default:0;column:like_count;comment:点赞数"`
	DeletedAt *int64 `gorm:"column:deleted_at;comment:删除时间"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}

// CommentLike 评论点赞模型
type CommentLike struct {
	UserID    int64 `gorm:"index:idx_user_comment;not null;comment:用户ID"`
	CommentID int64 `gorm:"index:idx_user_comment;not null;comment:评论ID"`
}

// TableName 指定表名
func (CommentLike) TableName() string {
	return "comment_likes"
}
