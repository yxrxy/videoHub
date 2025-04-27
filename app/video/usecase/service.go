package usecase

import (
	"context"

	"github.com/yxrxy/videoHub/app/video/domain/model"
	"github.com/yxrxy/videoHub/pkg/errno"
)

func (s *useCase) Publish(ctx context.Context, userID int64, videoData []byte, contentType,
	title string, description, category *string, tags []string, isPrivate bool) (string, error) {
	if !s.svc.CheckVideo(ctx, videoData, contentType) {
		return "", errno.NewErrNo(errno.ServiceUserExist, "video format error")
	}
	videoPath, err := s.svc.SaveVideo(ctx, userID, videoData, contentType, title, description, category, tags, isPrivate)
	if err != nil {
		return "", err
	}

	return videoPath, nil
}

func (s *useCase) GetVideoList(ctx context.Context, userID, page int64, size int32, category *string) ([]*model.Video, int64, error) {
	return s.svc.GetVideoList(ctx, userID, page, size, category)
}

func (s *useCase) SearchVideo(ctx context.Context, keywords string, pageSize, pageNum int32, fromDate, toDate *int64, username *string) ([]*model.Video, int64, error) {
	return nil, 0, nil
}

func (s *useCase) GetVideoDetail(ctx context.Context, videoID, userID int64) (*model.Video, error) {
	return s.svc.GetVideoDetail(ctx, videoID, userID)
}

// 从重构前的三层架构那边复制的，就不改了
func (s *useCase) GetHotVideos(ctx context.Context, limit int32, category string, lastVisit, lastLike, lastID int64) ([]*model.Video, int64, int64, int64, int64, error) {
	return s.svc.GetHotVideos(ctx, limit, category, lastVisit, lastLike, lastID)
}

func (s *useCase) DeleteVideo(ctx context.Context, videoID, userID int64) error {
	video, err := s.db.GetVideoByID(ctx, videoID)
	if err != nil {
		return err
	}
	if video.UserID != userID {
		return errno.NewErrNo(errno.InternalServiceErrorCode, "无权限删除该视频")
	}

	return s.db.DeleteVideo(ctx, videoID)
}

func (s *useCase) IncrementVisitCount(ctx context.Context, videoID int64) error {
	return s.db.IncrementVisitCount(ctx, videoID)
}

func (s *useCase) IncrementLikeCount(ctx context.Context, videoID int64) error {
	return s.db.IncrementLikeCount(ctx, videoID)
}
