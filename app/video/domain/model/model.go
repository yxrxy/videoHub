package model

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
