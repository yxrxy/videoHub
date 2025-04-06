package service

import (
	"context"
	"errors"

	"github.com/yxrxy/videoHub/app/interaction/domain/model"
	"gorm.io/gorm"
)

func (s *InteractionService) Like(ctx context.Context, userID int64, videoID int64) (bool, error) {
	like, err := s.db.GetLike(ctx, userID, videoID)
	if err == nil {
		// 已经点赞，取消点赞
		if err := s.db.DeleteLike(ctx, like); err != nil {
			return false, err
		}
		return true, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	// 创建新的点赞记录
	if err := s.db.CreateLike(ctx, userID, videoID); err != nil {
		return false, err
	}
	return true, nil
}

func (s *InteractionService) GetLikes(ctx context.Context, videoID int64, page int32, size int32) ([]*model.Like, error) {
	offset := (page - 1) * size
	likes, err := s.db.GetLikeList(ctx, videoID, int(offset), int(size))
	if err != nil {
		return nil, err
	}

	result := make([]*model.Like, len(likes))
	for i, like := range likes {
		result[i] = &model.Like{
			UserID:  like.UserID,
			VideoID: like.VideoID,
		}
	}
	return result, nil
}

func (s *InteractionService) Comment(ctx context.Context, userID int64, videoID int64, content string, parentID int64) (bool, error) {
	pa := int64(-1)
	pID := &pa
	if parentID != -1 {
		// 检查父评论是否存在
		_, err := s.db.GetComment(ctx, parentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, err
			}
			return false, err
		}
		pID = &parentID
	}
	// 检查视频是否存在
	exist, err := s.db.IsVideoExist(ctx, videoID)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, errors.New("video not found")
	}
	// 创建评论
	if err := s.db.CreateComment(ctx, userID, videoID, content, pID); err != nil {
		return false, err
	}
	return true, nil
}

func (s *InteractionService) GetComments(ctx context.Context, videoID int64, page int32, size int32) ([]*model.Comment, error) {
	offset := (page - 1) * size
	comments, _, err := s.db.GetCommentList(ctx, videoID, int(offset), int(size))
	if err != nil {
		return nil, err
	}

	result := make([]*model.Comment, len(comments))
	for i, comment := range comments {
		result[i] = &model.Comment{
			ID:        comment.ID,
			UserID:    comment.UserID,
			VideoID:   comment.VideoID,
			Content:   comment.Content,
			ParentID:  comment.ParentID,
			LikeCount: comment.LikeCount,
		}
	}

	return result, nil
}

func (s *InteractionService) LikeComment(ctx context.Context, userID int64, commentID int64) (bool, error) {
	// 检查评论是否存在
	comment, err := s.db.GetComment(ctx, commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
		return false, err
	}

	// 检查是否已经点赞
	liked, err := s.db.IsCommentLiked(ctx, userID, commentID)
	if err != nil {
		return false, err
	}

	if liked {
		// 已点赞，取消点赞
		if err := s.db.UnlikeComment(ctx, userID, commentID); err != nil {
			return false, err
		}
	} else {
		// 未点赞，添加点赞
		if err := s.db.LikeComment(ctx, userID, commentID); err != nil {
			return false, err
		}
	}

	// 更新评论的点赞数
	likeCount, err := s.db.GetCommentLikeCount(ctx, commentID)
	if err != nil {
		return false, err
	}
	comment.LikeCount = int32(likeCount)

	return true, nil
}
