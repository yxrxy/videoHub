package mysql

import "time"

// Video 视频模型
type Video struct {
	ID           int64      `json:"id"            gorm:"primarykey"`             // 视频ID
	UserID       int64      `json:"user_id"       gorm:"index"`                  // 作者ID
	VideoURL     string     `json:"video_url"     gorm:"type:varchar(255)"`      // 视频URL
	CoverURL     string     `json:"cover_url"     gorm:"type:varchar(255)"`      // 封面URL
	Title        string     `json:"title"         gorm:"type:varchar(128)"`      // 视频标题
	Description  string     `json:"description"   gorm:"type:varchar(512)"`      // 视频描述
	Duration     int64      `json:"duration"`                                    // 视频时长（秒）
	Category     string     `json:"category"      gorm:"type:varchar(32);index"` // 视频分类
	Tags         string     `json:"tags"          gorm:"type:varchar(255)"`      // 视频标签，以逗号分隔
	VisitCount   int64      `json:"visit_count"   gorm:"default:0"`              // 播放量
	LikeCount    int64      `json:"like_count"    gorm:"default:0"`              // 点赞数
	CommentCount int64      `json:"comment_count" gorm:"default:0"`              // 评论数
	IsPrivate    bool       `json:"is_private"    gorm:"default:false"`          // 是否私有
	CreatedAt    time.Time  `json:"created_at"`                                  // 创建时间
	UpdatedAt    time.Time  `json:"updated_at"`                                  // 更新时间
	DeletedAt    *time.Time `json:"deleted_at"    gorm:"index"`                  // 删除时间
}

// TableName 指定表名
func (Video) TableName() string {
	return "video"
}
