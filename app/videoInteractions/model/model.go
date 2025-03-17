package model

import "time"

// VideoComment 视频评论表
type VideoComment struct {
	ID        int64  `gorm:"primarykey"`
	VideoID   int64  `gorm:"index"`
	UserID    int64  `gorm:"index"`
	Content   string `gorm:"type:varchar(512)"`
	LikeCount int64  `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

// TableName 指定表名
func (VideoComment) TableName() string {
	return "video_comments"
}

// VideoVisit 视频访问记录表
type VideoVisit struct {
	ID        int64  `gorm:"primarykey"`
	VideoID   int64  `gorm:"index"`
	UserID    int64  `gorm:"index"`
	IP        string `gorm:"type:varchar(64)"`
	CreatedAt time.Time
}

// TableName 指定表名
func (VideoVisit) TableName() string {
	return "video_visits"
}

// VideoLike 视频点赞表
type VideoLike struct {
	ID        int64 `gorm:"primarykey"`
	VideoID   int64 `gorm:"uniqueIndex:idx_video_user"`
	UserID    int64 `gorm:"uniqueIndex:idx_video_user"`
	CreatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

// TableName 指定表名
func (VideoLike) TableName() string {
	return "video_likes"
}

// VideoTag 视频标签关联表
type VideoTag struct {
	ID        int64  `gorm:"primarykey"`
	VideoID   int64  `gorm:"index:idx_video_tag"`
	Tag       string `gorm:"type:varchar(32);index:idx_video_tag"`
	CreatedAt time.Time
}

// TableName 指定表名
func (VideoTag) TableName() string {
	return "video_tags"
}
