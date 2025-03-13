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

func (v *Video) GetVideoList(ctx context.Context, userID int64) ([]*video.Video, int64, error) {
	var videos []*model.Video
	var total int64

	if err := v.db.WithContext(ctx).Where("user_id = ?", userID).Find(&videos).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*video.Video, 0, len(videos))
	for _, v := range videos {
		var deletedAt int64
		if v.DeletedAt != nil {
			deletedAt = v.DeletedAt.Unix()
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
			DeletedAt:    &deletedAt,
		})
	}

	return result, total, nil
}

func (v *Video) GetHotVideos(ctx context.Context, limit int32, category string, cursor int64) ([]*video.Video, int64, error) {
	// 先从 Redis 获取热门视频 ID
	videoIDs, err := cache.GetHotVideos(ctx, category, int(limit), cursor)
	if err != nil {
		return nil, 0, err
	}

	// 如果 Redis 中没有数据，从数据库获取并写入 Redis
	if len(videoIDs) == 0 {
		var videos []*model.Video
		query := v.db.WithContext(ctx).Order("visit_count DESC, like_count DESC")

		if category != "" {
			query = query.Where("category = ?", category)
		}

		if cursor > 0 {
			query = query.Where("id < ?", cursor)
		}

		if limit <= 0 {
			limit = 10
		}

		if err := query.Limit(int(limit)).Find(&videos).Error; err != nil {
			return nil, 0, err
		}

		// 写入 Redis
		for _, video := range videos {
			score := float64(video.VisitCount) + float64(video.LikeCount)*1.5
			if err := cache.UpdateVideoScore(ctx, video.ID, video.VisitCount, video.LikeCount); err != nil {
				// 这里只记录错误，不影响返回结果
				log.Printf("Failed to update video score: %v", err)
			}
		}

		// 转换为 RPC 响应格式
		result := make([]*video.Video, 0, len(videos))
		for _, v := range videos {
			var deletedAt int64
			if v.DeletedAt != nil {
				deletedAt = v.DeletedAt.Unix()
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
				DeletedAt:    &deletedAt,
			})
		}

		return result, int64(len(result)), nil
	}

	// 如果 Redis 中有数据，根据 ID 获取视频详情
	var videos []*model.Video
	if err := v.db.WithContext(ctx).Where("id IN ?", videoIDs).Find(&videos).Error; err != nil {
		return nil, 0, err
	}

	// 转换为 RPC 响应格式
	result := make([]*video.Video, 0, len(videos))
	for _, v := range videos {
		var deletedAt int64
		if v.DeletedAt != nil {
			deletedAt = v.DeletedAt.Unix()
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
			DeletedAt:    &deletedAt,
		})
	}

	return result, int64(len(result)), nil
}

// IncrementVisitCount 增加视频播放量
func (v *Video) IncrementVisitCount(ctx context.Context, videoID int64) error {
	if err := v.db.WithContext(ctx).Model(&model.Video{}).Where("id = ?", videoID).
		UpdateColumn("visit_count", gorm.Expr("visit_count + ?", 1)).Error; err != nil {
		return err
	}

	// 更新 Redis 热度分数
	return cache.UpdateVideoScore(ctx, videoID, 1, 0)
}

// IncrementLikeCount 增加视频点赞数
func (v *Video) IncrementLikeCount(ctx context.Context, videoID int64) error {
	if err := v.db.WithContext(ctx).Model(&model.Video{}).Where("id = ?", videoID).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
		return err
	}

	// 更新 Redis 热度分数
	return cache.UpdateVideoScore(ctx, videoID, 0, 1)
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.GetDSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
