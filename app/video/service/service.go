package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yxrrxy/videoHub/app/video/model"
	"github.com/yxrrxy/videoHub/app/video/repository"
	"github.com/yxrrxy/videoHub/kitex_gen/video"
	videoproto "github.com/yxrrxy/videoHub/kitex_gen/video"
	"github.com/yxrrxy/videoHub/pkg/errno"
)

type VideoService struct {
	repo *repository.Video
}

func NewVideoService(repo *repository.Video) *VideoService {
	return &VideoService{
		repo: repo,
	}
}

func (s *VideoService) Publish(ctx context.Context, req *videoproto.PublishRequest) (*videoproto.PublishResponse, error) {
	userID := req.UserId
	if userID == 0 {
		return nil, errno.ErrUnauthorized
	}

	contentType := req.ContentType
	if !isValidVideoType(contentType) {
		return nil, errno.ErrInvalidParam
	}

	ext := ".mp4"
	switch contentType {
	case "video/mp4":
		ext = ".mp4"
	case "video/mpeg":
		ext = ".mpeg"
	case "video/webm":
		ext = ".webm"
	}

	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%d%s", userID, timestamp, ext)

	// 保存视频文件
	videoPath := filepath.Join("static/videos", filename)
	if err := os.MkdirAll(filepath.Dir(videoPath), 0755); err != nil {
		return nil, err
	}

	if err := os.WriteFile(videoPath, req.VideoData, 0644); err != nil {
		return nil, err
	}

	video := &model.Video{
		UserID:      userID,
		Title:       req.Title,
		VideoURL:    "/videos/" + filename,
		Description: *req.Description,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.CreateVideo(ctx, video); err != nil {
		os.Remove(videoPath)
		return nil, err
	}

	return &videoproto.PublishResponse{
		VideoUrl: video.VideoURL,
	}, nil
}

func isValidVideoType(contentType string) bool {
	validTypes := map[string]bool{
		"video/mp4":  true,
		"video/mpeg": true,
		"video/webm": true,
	}
	return validTypes[contentType]
}

func (s *VideoService) GetVideoList(ctx context.Context, req *video.VideoListRequest) (r *video.VideoListResponse, err error) {
	videos, total, err := s.repo.GetVideoList(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &video.VideoListResponse{
		Videos: videos,
		Total:  total,
	}, nil
}
