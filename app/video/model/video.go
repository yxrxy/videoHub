package model

import "time"

// Video 视频模型
type Video struct {
	ID           int64      `gorm:"primarykey" json:"id"`                   // 视频ID
	UserID       int64      `gorm:"index" json:"user_id"`                   // 作者ID
	VideoURL     string     `gorm:"type:varchar(255)" json:"video_url"`     // 视频URL
	CoverURL     string     `gorm:"type:varchar(255)" json:"cover_url"`     // 封面URL
	Title        string     `gorm:"type:varchar(128)" json:"title"`         // 视频标题
	Description  string     `gorm:"type:varchar(512)" json:"description"`   // 视频描述
	Duration     int64      `json:"duration"`                               // 视频时长（秒）
	Category     string     `gorm:"type:varchar(32);index" json:"category"` // 视频分类
	Tags         string     `gorm:"type:varchar(255)" json:"tags"`          // 视频标签，以逗号分隔
	VisitCount   int64      `gorm:"default:0" json:"visit_count"`           // 播放量
	LikeCount    int64      `gorm:"default:0" json:"like_count"`            // 点赞数
	CommentCount int64      `gorm:"default:0" json:"comment_count"`         // 评论数
	IsPrivate    bool       `gorm:"default:false" json:"is_private"`        // 是否私有
	CreatedAt    time.Time  `json:"created_at"`                             // 创建时间
	UpdatedAt    time.Time  `json:"updated_at"`                             // 更新时间
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at"`                // 删除时间
}

// TableName 指定表名
func (Video) TableName() string {
	return "videos"
}
