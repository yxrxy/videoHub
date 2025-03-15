package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	// 校验视频类型
	contentType := req.ContentType
	if !isValidVideoType(contentType) {
		return nil, errno.ErrInvalidParam
	}

	// 确定文件扩展名
	ext := ".mp4"
	switch contentType {
	case "video/mp4":
		ext = ".mp4"
	case "video/mpeg":
		ext = ".mpeg"
	case "video/webm":
		ext = ".webm"
	}

	// 生成唯一文件名
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%d%s", userID, timestamp, ext)
	videoPath := filepath.Join("static/videos", filename)

	// 创建文件夹并写入视频数据
	if err := os.MkdirAll(filepath.Dir(videoPath), 0755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(videoPath, req.VideoData, 0644); err != nil {
		return nil, err
	}

	// 生成封面（封面图路径）
	coverFilename := fmt.Sprintf("%d_%d.jpg", userID, timestamp)
	coverPath := filepath.Join("static/covers", coverFilename)

	// 创建 Video 记录
	videoRecord := &model.Video{
		UserID:       userID,
		VideoURL:     "/videos/" + filename,
		CoverURL:     "/covers/" + coverFilename,
		Title:        req.Title,
		Description:  *req.Description,
		Duration:     0,                           //duration,
		Category:     *req.Category,               // 分类
		Tags:         strings.Join(req.Tags, ","), // 标签
		VisitCount:   0,
		LikeCount:    0,
		CommentCount: 0,
		IsPrivate:    *req.IsPrivate,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 存入数据库
	if err := s.repo.CreateVideo(ctx, videoRecord); err != nil {
		err := os.Remove(videoPath)
		if err != nil {
			return nil, err
		}
		err = os.Remove(coverPath)
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	// 返回结果
	return &videoproto.PublishResponse{
		VideoUrl: videoRecord.VideoURL,
		CoverUrl: videoRecord.CoverURL,
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
	videos, total, err := s.repo.GetVideoList(ctx, req.UserId, *req.Page, *req.Size, req.Category)
	if err != nil {
		return nil, err
	}

	result := make([]*video.Video, 0, len(videos))

	for _, v := range videos {
		var deletedAt *int64
		if v.DeletedAt != nil {
			timestamp := v.DeletedAt.Unix()
			deletedAt = &timestamp
		}

		result = append(result, &video.Video{
			Id:           v.ID,
			UserId:       v.UserID,
			VideoUrl:     v.VideoURL,
			CoverUrl:     v.CoverURL,
			Title:        v.Title,
			Description:  &v.Description,
			Duration:     v.Duration,
			Category:     v.Category,
			Tags:         strings.Split(v.Tags, ","),
			VisitCount:   v.VisitCount,
			LikeCount:    v.LikeCount,
			CommentCount: v.CommentCount,
			IsPrivate:    v.IsPrivate,
			CreatedAt:    v.CreatedAt.Unix(),
			UpdatedAt:    v.UpdatedAt.Unix(),
			DeletedAt:    deletedAt,
		})
	}
	return &video.VideoListResponse{
		Videos: result,
		Total:  total,
	}, nil
}

func (s *VideoService) GetHotVideos(ctx context.Context, req *video.HotVideoRequest) (r *video.HotVideoResponse, err error) {
	videos, total, nextVisit, nextLike, nextID, err := s.repo.GetHotVideos(
		ctx, *req.Limit, *req.Category, *req.LastVisit, *req.LastLike, *req.LastId)
	if err != nil {
		return nil, err
	}

	return &video.HotVideoResponse{
		Videos:    videos,
		Total:     total,
		NextVisit: &nextVisit,
		NextLike:  &nextLike,
		NextId:    &nextID,
	}, nil
}

func (s *VideoService) IncrementVisitCount(ctx context.Context, req *video.IncrementVisitCountRequest) (
	r *video.IncrementVisitCountResponse, err error) {
	if err := s.repo.IncrementVisitCount(ctx, req.VideoId); err != nil {
		return nil, err
	}

	return &video.IncrementVisitCountResponse{
		Success: true,
	}, nil
}

func (s *VideoService) IncrementLikeCount(ctx context.Context, req *video.IncrementLikeCountRequest) (
	r *video.IncrementLikeCountResponse, err error) {
	if err := s.repo.IncrementLikeCount(ctx, req.VideoId); err != nil {
		return nil, err
	}

	return &video.IncrementLikeCountResponse{
		Success: true,
	}, nil
}
