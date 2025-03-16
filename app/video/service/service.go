package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yxrrxy/videoHub/app/video/model"
	"github.com/yxrrxy/videoHub/app/video/repository"
	"github.com/yxrrxy/videoHub/config"
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

	// 检查视频大小
	if len(req.VideoData) > config.Upload.Video.MaxSize {
		return nil, errno.ErrInvalidParam
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
	case "video/quicktime":
		ext = ".mov"
	}

	// 生成唯一文件名
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%d%s", userID, timestamp, ext)

	// 确保目录存在
	if err := os.MkdirAll(config.Upload.Video.UploadDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %v", err)
	}

	// 保存视频文件
	videoPath := filepath.Join(config.Upload.Video.UploadDir, filename)
	if err := os.WriteFile(videoPath, req.VideoData, 0644); err != nil {
		return nil, fmt.Errorf("保存视频失败: %v", err)
	}

	// 使用 ffmpeg 获取视频时长和生成封面
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", videoPath)
	output, err := cmd.Output()
	if err != nil {
		os.Remove(videoPath)
		return nil, fmt.Errorf("获取视频时长失败: %v", err)
	}
	duration, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		os.Remove(videoPath)
		return nil, fmt.Errorf("解析视频时长失败: %v", err)
	}

	// 生成封面
	coverFilename := fmt.Sprintf("%d_%d.jpg", userID, timestamp)
	coverDir := filepath.Join(config.Upload.Video.UploadDir, "../covers")
	if err := os.MkdirAll(coverDir, 0755); err != nil {
		os.Remove(videoPath)
		return nil, fmt.Errorf("创建封面目录失败: %v", err)
	}
	coverPath := filepath.Join(coverDir, coverFilename)

	cmd = exec.Command("ffmpeg", "-i", videoPath, "-ss", "00:00:01", "-vframes", "1", coverPath)
	if err := cmd.Run(); err != nil {
		os.Remove(videoPath)
		return nil, fmt.Errorf("生成封面失败: %v", err)
	}

	// 构建URL
	videoURL := fmt.Sprintf("%s/%s", config.Upload.Video.BaseURL, filename)
	coverURL := fmt.Sprintf("%s/%s", strings.Replace(config.Upload.Video.BaseURL, "/videos", "/covers", 1), coverFilename)

	// 创建 Video 记录
	videoRecord := &model.Video{
		UserID:       userID,
		VideoURL:     videoURL,
		CoverURL:     coverURL,
		Title:        req.Title,
		Description:  *req.Description,
		Duration:     int64(duration),
		Category:     *req.Category,
		Tags:         strings.Join(req.Tags, ","),
		VisitCount:   0,
		LikeCount:    0,
		CommentCount: 0,
		IsPrivate:    *req.IsPrivate,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 存入数据库
	if err := s.repo.CreateVideo(ctx, videoRecord); err != nil {
		os.Remove(videoPath)
		os.Remove(coverPath)
		return nil, fmt.Errorf("保存视频记录失败: %v", err)
	}

	// 返回结果
	return &videoproto.PublishResponse{
		VideoUrl: videoRecord.VideoURL,
		CoverUrl: videoRecord.CoverURL,
	}, nil
}

func isValidVideoType(contentType string) bool {
	validTypes := map[string]bool{
		"video/mp4":       true,
		"video/quicktime": true,
	}
	return validTypes[contentType]
}

func (s *VideoService) GetVideoList(ctx context.Context, req *video.VideoListRequest) (r *video.VideoListResponse, err error) {
	// 设置默认值
	page := int64(1)
	size := int32(10)
	var category string

	if req.Page != nil {
		page = *req.Page
	}
	if req.Size != nil {
		size = *req.Size
	}
	if req.Category != nil {
		category = *req.Category
	}

	videos, total, err := s.repo.GetVideoList(ctx, req.UserId, page, size, &category)
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
	// 设置默认值
	limit := int32(10)
	category := ""
	lastVisit := int64(0)
	lastLike := int64(0)
	lastID := int64(0)

	// 检查可选参数
	if req.Limit != nil {
		limit = *req.Limit
	}
	if req.Category != nil {
		category = *req.Category
	}
	if req.LastVisit != nil {
		lastVisit = *req.LastVisit
	}
	if req.LastLike != nil {
		lastLike = *req.LastLike
	}
	if req.LastId != nil {
		lastID = *req.LastId
	}

	videos, total, nextVisit, nextLike, nextID, err := s.repo.GetHotVideos(
		ctx, limit, category, lastVisit, lastLike, lastID)
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
