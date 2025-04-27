package model

type Like struct {
	ID      int64
	UserID  int64
	VideoID int64
}

// Comment 评论模型
type Comment struct {
	ID        int64
	UserID    int64
	VideoID   int64
	Content   string
	ParentID  *int64
	LikeCount int32
}

// CommentLike 评论点赞模型
type CommentLike struct {
	UserID    int64
	CommentID int64
}
