package repository

import (
	"context"
	"log"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/yxrrxy/videoHub/app/video/cache"
	"github.com/yxrrxy/videoHub/app/video/model"
	"github.com/yxrrxy/videoHub/config"
	"github.com/yxrrxy/videoHub/kitex_gen/video"
)

type Video struct {
	db *gorm.DB
}

func NewVideo(db *gorm.DB) *Video {
	return &Video{db: db}
}

func (v *Video) CreateVideo(ctx context.Context, video *model.Video) error {
	return v.db.WithContext(ctx).Create(video).Error
}

func (v *Video) GetVideoList(ctx context.Context, userID, page int64, size int32, category *string) ([]model.Video, int64, error) {
	var videos []model.Video
	var total int64

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	offset := (page - 1) * int64(size)

	query := v.db.WithContext(ctx).Model(&model.Video{}).Where("user_id = ?", userID)

	if category != nil && *category != "" {
		query = query.Where("category = ?", *category)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(size)).
		Find(&videos).Error; err != nil {
		return nil, 0, err
	}

	return videos, total, nil
}

func (v *Video) GetHotVideos(ctx context.Context, limit int32, category string, lastVisitCount, lastLikeCount int64,
	lastID int64) ([]*video.Video, int64, int64, int64, int64, error) {
	// 从 Redis 获取热门视频 ID
	videoIDs, err := cache.GetHotVideos(ctx, category, int(limit), lastVisitCount, lastLikeCount, lastID)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	// 如果 Redis 中没有数据，从数据库获取
	if len(videoIDs) == 0 {
		return v.fetchVideosFromDB(ctx, limit, category, lastVisitCount, lastLikeCount, lastID)
	}

	// 从数据库获取视频详细信息
	return v.getVideosByIDs(ctx, videoIDs)
}

// fetchVideosFromDB 从数据库获取热门视频并更新 Redis
func (v *Video) fetchVideosFromDB(ctx context.Context, limit int32, category string, lastVisitCount, lastLikeCount, lastID int64) ([]*video.Video, int64, int64, int64, int64, error) {
	var videos []*model.Video
	query := v.db.WithContext(ctx).Order("visit_count DESC, like_count DESC, id DESC")

	// 根据分类过滤
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 计算 `score`
	score := float64(lastVisitCount) + float64(lastLikeCount)*1.5

	// 添加游标条件
	if lastID > 0 {
		query = query.Where(`
        (visit_count + like_count * 1.5 < ?) 
        OR (visit_count + like_count * 1.5 = ? AND id < ?)`,
			score, score, lastID,
		)
	}

	// 执行查询
	if err := query.Order("visit_count + like_count * 1.5 DESC, id DESC").
		Limit(int(limit)).
		Find(&videos).Error; err != nil {
		return nil, 0, 0, 0, 0, err
	}

	// 更新 Redis 热度分数
	for _, videoData := range videos {
		if err := cache.UpdateVideoScore(ctx, videoData.ID, videoData.VisitCount, videoData.LikeCount, videoData.Category); err != nil {
			log.Printf("Failed to update video score: %v", err)
		}
	}

	// 转换为 RPC 响应格式
	result := v.convertToRPCFormat(videos)
	nextVisitCount, nextLikeCount, nextID := v.extractLastVideoData(videos)

	return result, int64(len(result)), nextVisitCount, nextLikeCount, nextID, nil
}

// getVideosByIDs 根据视频 ID 获取视频详细信息
func (v *Video) getVideosByIDs(ctx context.Context, videoIDs []string) ([]*video.Video, int64, int64, int64, int64, error) {
	var videos []*model.Video
	if err := v.db.WithContext(ctx).Where("id IN ?", videoIDs).Find(&videos).Error; err != nil {
		return nil, 0, 0, 0, 0, err
	}

	// 转换为 RPC 响应格式
	result := v.convertToRPCFormat(videos)
	nextVisitCount, nextLikeCount, nextID := v.extractLastVideoData(videos)

	return result, int64(len(result)), nextVisitCount, nextLikeCount, nextID, nil
}

// convertToRPCFormat 将视频数据转换为 RPC 响应格式
func (v *Video) convertToRPCFormat(videos []*model.Video) []*video.Video {
	result := make([]*video.Video, len(videos))
	for i, v := range videos {
		var deletedAt int64
		if v.DeletedAt != nil {
			deletedAt = v.DeletedAt.Unix()
		}

		result[i] = &video.Video{
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
			DeletedAt:    &deletedAt,
		}
	}
	return result
}

// extractLastVideoData 提取最后一个视频的数据
func (v *Video) extractLastVideoData(videos []*model.Video) (int64, int64, int64) {
	if len(videos) == 0 {
		return 0, 0, 0
	}

	lastVideo := videos[len(videos)-1]
	return lastVideo.VisitCount, lastVideo.LikeCount, lastVideo.ID
}

// IncrementVisitCount 增加视频播放量
func (v *Video) IncrementVisitCount(ctx context.Context, videoID int64) error {
	var category string
	if err := v.db.WithContext(ctx).Model(&model.Video{}).Select("category").
		Where("id = ?", videoID).Scan(&category).Error; err != nil {
		return err
	}

	if err := v.db.WithContext(ctx).Model(&model.Video{}).Where("id = ?", videoID).
		UpdateColumn("visit_count", gorm.Expr("visit_count + ?", 1)).Error; err != nil {
		return err
	}

	// 更新 Redis 热度分数
	return cache.UpdateVideoScore(ctx, videoID, 1, 0, category)
}

// IncrementLikeCount 增加视频点赞数
func (v *Video) IncrementLikeCount(ctx context.Context, videoID int64) error {
	var category string
	if err := v.db.WithContext(ctx).Model(&model.Video{}).Select("category").
		Where("id = ?", videoID).Scan(&category).Error; err != nil {
		return err
	}

	if err := v.db.WithContext(ctx).Model(&model.Video{}).Where("id = ?", videoID).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
		return err
	}

	// 更新 Redis 热度分数
	return cache.UpdateVideoScore(ctx, videoID, 0, 1, category)
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.GetDSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
