package model

import (
	"time"
)

type Video struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	VideoURL     string `json:"video_url"`
	CoverURL     string `json:"cover_url"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Duration     int64  `json:"duration"`
	Category     string `json:"category"`
	Tags         string `json:"tags"`
	VisitCount   int64  `json:"visit_count"`
	LikeCount    int64  `json:"like_count"`
	CommentCount int64  `json:"comment_count"`
	IsPrivate    bool   `json:"is_private"`
}

// 语义搜索结果项
type SemanticSearchResultItem struct {
	Videos         []*Video `json:"videos"`
	Summary        string   `json:"summary"`
	RelatedQueries []string `json:"related_queries"`
	FromCache      bool     `json:"from_cache"`
}

type VideoES struct {
	ID          int64     `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Category    string    `json:"category,omitempty"`
	AuthorID    int64     `json:"author_id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	ViewCount   int64     `json:"view_count,omitempty"`
	IsDeleted   bool      `json:"is_deleted,omitempty"`
	SearchText  string    `json:"search_text,omitempty"`

	Keywords string  `json:"keywords,omitempty"`
	FromDate *int64  `json:"from_date,omitempty"`
	ToDate   *int64  `json:"to_date,omitempty"`
	Username *string `json:"username,omitempty"`
}
