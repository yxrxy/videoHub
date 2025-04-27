package mysql

import (
	"context"

	"github.com/yxrxy/videoHub/app/interaction/domain/model"
	model2 "github.com/yxrxy/videoHub/app/video/domain/model"
	"gorm.io/gorm"
)

type Interaction struct {
	db *gorm.DB
}

func NewInteraction(db *gorm.DB) *Interaction {
	return &Interaction{db: db}
}

// GetLikeList 获取点赞列表
func (i *Interaction) GetLikeList(ctx context.Context, videoID int64, offset, limit int) ([]*model.Like, error) {
	var likes []*Like
	err := i.db.WithContext(ctx).
		Where("video_id = ? AND deleted_at IS NULL", videoID).
		Offset(offset).
		Limit(limit).
		Find(&likes).Error
	var result []*model.Like
	for _, like := range likes {
		result = append(result, &model.Like{
			UserID:  like.UserID,
			VideoID: like.VideoID,
		})
	}
	return result, err
}

// CreateLike 创建点赞
func (i *Interaction) CreateLike(ctx context.Context, userID, videoID int64) error {
	return i.db.WithContext(ctx).Create(&Like{
		UserID:  userID,
		VideoID: videoID,
	}).Error
}

// GetLike 获取点赞记录
func (i *Interaction) GetLike(ctx context.Context, userID, videoID int64) (*model.Like, error) {
	var like Like
	err := i.db.WithContext(ctx).
		Where("video_id = ? AND deleted_at IS NULL", videoID).
		First(&like).Error
	if err != nil {
		return nil, err
	}
	return &model.Like{
		UserID:  like.UserID,
		VideoID: like.VideoID,
	}, nil
}

// DeleteLike 删除点赞
func (i *Interaction) DeleteLike(ctx context.Context, like *model.Like) error {
	return i.db.WithContext(ctx).Delete(&Like{
		UserID:  like.UserID,
		VideoID: like.VideoID,
	}).Error
}

// CreateComment 创建评论
func (i *Interaction) CreateComment(ctx context.Context, userID, videoID int64, content string, parentID *int64) error {
	return i.db.WithContext(ctx).Create(&Comment{
		UserID:   userID,
		VideoID:  videoID,
		Content:  content,
		ParentID: parentID,
	}).Error
}

// GetComment 获取评论
func (i *Interaction) GetComment(ctx context.Context, commentID int64) (*model.Comment, error) {
	var comment Comment
	err := i.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", commentID).
		First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &model.Comment{
		ID:       comment.ID,
		UserID:   comment.UserID,
		VideoID:  comment.VideoID,
		Content:  comment.Content,
		ParentID: comment.ParentID,
	}, nil
}

// IsVideoExist 检查视频是否存在
func (i *Interaction) IsVideoExist(ctx context.Context, videoID int64) (bool, error) {
	var count int64
	err := i.db.WithContext(ctx).Model(&model2.Video{}).
		Where("id = ?", videoID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetCommentList 获取评论列表
func (i *Interaction) GetCommentList(ctx context.Context, videoID int64, offset, limit int) ([]*model.Comment, int64, error) {
	var comments []*Comment
	var total int64

	// 获取总数
	err := i.db.WithContext(ctx).Model(&Comment{}).
		Where("video_id = ? AND deleted_at IS NULL", videoID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取评论列表
	err = i.db.WithContext(ctx).
		Where("video_id = ? AND deleted_at IS NULL", videoID).
		Offset(offset).
		Limit(limit).
		Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}
	var result []*model.Comment
	for _, comment := range comments {
		result = append(result, &model.Comment{
			ID:        comment.ID,
			UserID:    comment.UserID,
			VideoID:   comment.VideoID,
			Content:   comment.Content,
			ParentID:  comment.ParentID,
			LikeCount: comment.LikeCount,
		})
	}
	return result, total, nil
}

// DeleteComment 删除评论
func (i *Interaction) DeleteComment(ctx context.Context, comment *model.Comment) error {
	return i.db.WithContext(ctx).Delete(&Comment{
		ID:      comment.ID,
		UserID:  comment.UserID,
		VideoID: comment.VideoID,
	}).Error
}

// LikeComment 点赞评论
func (i *Interaction) LikeComment(ctx context.Context, userID, commentID int64) error {
	return i.db.WithContext(ctx).Create(&CommentLike{
		UserID:    userID,
		CommentID: commentID,
	}).Error
}

// UnlikeComment 取消评论点赞
func (i *Interaction) UnlikeComment(ctx context.Context, userID, commentID int64) error {
	return i.db.WithContext(ctx).Where("user_id = ? AND comment_id = ?", userID, commentID).
		Delete(&CommentLike{}).Error
}

// IsCommentLiked 检查用户是否点赞了评论
func (i *Interaction) IsCommentLiked(ctx context.Context, userID, commentID int64) (bool, error) {
	var count int64
	err := i.db.WithContext(ctx).Model(&CommentLike{}).
		Where("user_id = ? AND comment_id = ?", userID, commentID).
		Count(&count).Error
	return count > 0, err
}

// GetCommentLikeCount 获取评论点赞数
func (i *Interaction) GetCommentLikeCount(ctx context.Context, commentID int64) (int64, error) {
	var count int64
	err := i.db.WithContext(ctx).Model(&CommentLike{}).
		Where("comment_id = ?", commentID).
		Count(&count).Error
	return count, err
}
