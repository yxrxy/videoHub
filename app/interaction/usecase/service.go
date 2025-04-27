package usecase

import (
	"context"
	"errors"

	"github.com/yxrxy/videoHub/app/interaction/domain/model"
	"gorm.io/gorm"
)

// Like 点赞操作
func (s *useCase) Like(ctx context.Context, userID int64, videoID int64) (bool, error) {
	return s.svc.Like(ctx, userID, videoID)
}

// GetLikes 获取点赞列表
func (s *useCase) GetLikes(ctx context.Context, videoID int64, page int32, size int32) ([]*model.Like, error) {
	return s.svc.GetLikes(ctx, videoID, page, size)
}

// Comment 发表评论
func (s *useCase) Comment(ctx context.Context, userID int64, videoID int64, content string, parentID int64) (bool, error) {
	return s.svc.Comment(ctx, userID, videoID, content, parentID)
}

// GetComments 获取评论列表
func (s *useCase) GetComments(ctx context.Context, videoID int64, page int32, size int32) ([]*model.Comment, error) {
	return s.svc.GetComments(ctx, videoID, page, size)
}

// DeleteComment 删除评论
func (s *useCase) DeleteComment(ctx context.Context, userID int64, commentID int64) (bool, error) {
	// 获取评论
	comment, err := s.db.GetComment(ctx, commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
		return false, err
	}

	// 检查是否是评论作者
	if comment.UserID != userID {
		return false, err
	}

	// 删除评论
	if err := s.db.DeleteComment(ctx, comment); err != nil {
		return false, err
	}
	return true, nil
}

// LikeComment 点赞评论
func (s *useCase) LikeComment(ctx context.Context, userID int64, commentID int64) (bool, error) {
	return s.svc.LikeComment(ctx, userID, commentID)
}
