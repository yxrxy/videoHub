package model

import (
	"time"
)

type Video struct {
	ID           int64
	UserID       int64
	VideoURL     string
	CoverURL     string
	Title        string
	Description  string
	Duration     int64
	Category     string
	Tags         string
	VisitCount   int64
	LikeCount    int64
	CommentCount int64
	IsPrivate    bool
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
