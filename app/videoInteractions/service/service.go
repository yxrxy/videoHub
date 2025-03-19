package service

import (
	"context"

	"github.com/yxrrxy/videoHub/app/videoInteractions/repository"
	interactions "github.com/yxrrxy/videoHub/kitex_gen/videoInteractions"
	"gorm.io/gorm"
)

type InteractionService struct {
	repo *repository.Interaction
}

func NewInteractionService(repo *repository.Interaction) *InteractionService {
	return &InteractionService{repo: repo}
}

// Like 点赞操作
func (s *InteractionService) Like(ctx context.Context, userID int64, videoID int64) (bool, error) {
	// 检查是否已经点赞
	like, err := s.repo.GetLike(ctx, userID, videoID)
	if err == nil {
		// 已经点赞，取消点赞
		if err := s.repo.DeleteLike(ctx, like); err != nil {
			return false, err
		}
		return true, nil
	} else if err != gorm.ErrRecordNotFound {
		return false, err
	}

	// 创建新的点赞记录
	if err := s.repo.CreateLike(ctx, userID, videoID); err != nil {
		return false, err
	}
	return true, nil
}

// GetLikes 获取点赞列表
func (s *InteractionService) GetLikes(ctx context.Context, videoID int64, page int32, size int32) ([]*interactions.LikeInfo, error) {
	offset := (page - 1) * size
	likes, err := s.repo.GetLikeList(ctx, videoID, int(offset), int(size))
	if err != nil {
		return nil, err
	}

	result := make([]*interactions.LikeInfo, len(likes))
	for i, like := range likes {
		result[i] = &interactions.LikeInfo{
			UserId:  like.UserID,
			VideoId: like.VideoID,
		}
	}
	return result, nil
}

// Comment 发表评论
func (s *InteractionService) Comment(ctx context.Context, userID int64, videoID int64, content string, parentID int64) (bool, error) {
	var pID *int64
	if parentID != 0 {
		// 检查父评论是否存在
		_, err := s.repo.GetComment(ctx, parentID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return false, err
			}
			return false, err
		}
		pID = &parentID
	}

	// 创建评论
	if err := s.repo.CreateComment(ctx, userID, videoID, content, pID); err != nil {
		return false, err
	}
	return true, nil
}

// GetComments 获取评论列表
func (s *InteractionService) GetComments(ctx context.Context, videoID int64, page int32, size int32) (*interactions.CommentListResponse, error) {
	offset := (page - 1) * size
	comments, total, err := s.repo.GetCommentList(ctx, videoID, int(offset), int(size))
	if err != nil {
		return nil, err
	}

	result := make([]*interactions.CommentInfo, len(comments))
	for i, comment := range comments {
		result[i] = &interactions.CommentInfo{
			Id:        comment.ID,
			UserId:    comment.UserID,
			VideoId:   comment.VideoID,
			Content:   comment.Content,
			ParentId:  comment.ParentID,
			LikeCount: &comment.LikeCount,
		}
	}

	return &interactions.CommentListResponse{
		Comments: result,
		Total:    total,
	}, nil
}

// DeleteComment 删除评论
func (s *InteractionService) DeleteComment(ctx context.Context, userID int64, commentID int64) (bool, error) {
	// 获取评论
	comment, err := s.repo.GetComment(ctx, commentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, err
		}
		return false, err
	}

	// 检查是否是评论作者
	if comment.UserID != userID {
		return false, err
	}

	// 删除评论
	if err := s.repo.DeleteComment(ctx, comment); err != nil {
		return false, err
	}
	return true, nil
}

// LikeComment 点赞评论
func (s *InteractionService) LikeComment(ctx context.Context, userID int64, commentID int64) (bool, error) {
	// 检查评论是否存在
	comment, err := s.repo.GetComment(ctx, commentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, err
		}
		return false, err
	}

	// 检查是否已经点赞
	liked, err := s.repo.IsCommentLiked(ctx, userID, commentID)
	if err != nil {
		return false, err
	}

	if liked {
		// 已点赞，取消点赞
		if err := s.repo.UnlikeComment(ctx, userID, commentID); err != nil {
			return false, err
		}
	} else {
		// 未点赞，添加点赞
		if err := s.repo.LikeComment(ctx, userID, commentID); err != nil {
			return false, err
		}
	}

	// 更新评论的点赞数
	likeCount, err := s.repo.GetCommentLikeCount(ctx, commentID)
	if err != nil {
		return false, err
	}
	comment.LikeCount = int32(likeCount)

	return true, nil
}
