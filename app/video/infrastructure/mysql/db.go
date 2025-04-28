package mysql

import (
	"context"
	"log"

	model2 "github.com/yxrxy/videoHub/app/user/domain/model"
	"github.com/yxrxy/videoHub/app/video/domain/model"
	"github.com/yxrxy/videoHub/app/video/domain/repository"
	"github.com/yxrxy/videoHub/pkg/constants"
	"gorm.io/gorm"
)

type VideoDB struct {
	db    *gorm.DB
	cache repository.VideoCache
}

func NewVideoDB(db *gorm.DB, cache repository.VideoCache) repository.VideoDB {
	return &VideoDB{db: db, cache: cache}
}

func (v *VideoDB) CreateVideo(ctx context.Context, video *model.Video) error {
	return v.db.WithContext(ctx).Create(video).Error
}

func (v *VideoDB) GetVideoByID(ctx context.Context, videoID int64) (*model.Video, error) {
	var video Video
	if err := v.db.WithContext(ctx).Where("id = ?", videoID).First(&video).Error; err != nil {
		return nil, err
	}
	result := &model.Video{
		ID:           video.ID,
		UserID:       video.UserID,
		VideoURL:     video.VideoURL,
		CoverURL:     video.CoverURL,
		Title:        video.Title,
		Description:  video.Description,
		Duration:     video.Duration,
		Category:     video.Category,
		Tags:         video.Tags,
		VisitCount:   video.VisitCount,
		LikeCount:    video.LikeCount,
		CommentCount: video.CommentCount,
	}
	return result, nil
}

func (v *VideoDB) UpdateVideo(ctx context.Context, video *model.Video) error {
	return v.db.WithContext(ctx).Model(&model.Video{}).Where("id = ?", video.ID).Updates(video).Error
}

func (v *VideoDB) GetUsernameByID(ctx context.Context, userID int64) (string, error) {
	var user model2.User
	if err := v.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return "", err
	}
	return user.Username, nil
}

func (v *VideoDB) GetVideoList(ctx context.Context, userID, page int64, size int32, category *string) ([]*model.Video, int64, error) {
	var videos []Video
	var total int64

	offset := (page - 1) * int64(size)

	query := v.db.WithContext(ctx).Model(&model.Video{}).Where("user_id = ?", userID)

	if category != nil && *category != "" {
		query = query.Where("category = ?", *category)
	} else {
		query = query.Where("category = ?", "all")
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

	result := v.convertFormat(videos)

	return result, total, nil
}

func (v *VideoDB) DeleteVideo(ctx context.Context, videoID int64) error {
	return v.db.WithContext(ctx).Delete(&Video{}, videoID).Error
}

func (v *VideoDB) GetHotVideos(ctx context.Context, limit int32, category string, lastVisitCount, lastLikeCount int64,
	lastID int64,
) ([]*model.Video, int64, int64, int64, int64, error) {
	// 从 Redis 获取热门视频 ID
	videoIDs, err := v.cache.GetHotVideos(ctx, category, int(limit), lastVisitCount, lastLikeCount, lastID)
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
func (v *VideoDB) fetchVideosFromDB(
	ctx context.Context,
	limit int32,
	category string,
	lastVisitCount, lastLikeCount, lastID int64,
) ([]*model.Video, int64, int64, int64, int64, error) {
	var videos []Video
	query := v.db.WithContext(ctx).Order("visit_count DESC, like_count DESC, id DESC")

	// 根据分类过滤
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 计算 `score`
	score := float64(lastVisitCount) + float64(lastLikeCount)*constants.VideoScoreWeight

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
		if err := v.cache.UpdateVideoScore(ctx, videoData.ID, videoData.VisitCount, videoData.LikeCount, videoData.Category); err != nil {
			log.Printf("Failed to update video score: %v", err)
		}
	}

	// 转换为 RPC 响应格式
	result := v.convertFormat(videos)
	nextVisitCount, nextLikeCount, nextID := v.extractLastVideoData(videos)

	return result, int64(len(result)), nextVisitCount, nextLikeCount, nextID, nil
}

// getVideosByIDs 根据视频 ID 获取视频详细信息
func (v *VideoDB) getVideosByIDs(ctx context.Context, videoIDs []string) ([]*model.Video, int64, int64, int64, int64, error) {
	var videos []Video
	if err := v.db.WithContext(ctx).Where("id IN ?", videoIDs).Find(&videos).Error; err != nil {
		return nil, 0, 0, 0, 0, err
	}

	result := v.convertFormat(videos)
	nextVisitCount, nextLikeCount, nextID := v.extractLastVideoData(videos)

	return result, int64(len(result)), nextVisitCount, nextLikeCount, nextID, nil
}

// convertFormat 将视频数据转换
func (v *VideoDB) convertFormat(videos []Video) []*model.Video {
	result := make([]*model.Video, len(videos))
	for i, v := range videos {
		result[i] = &model.Video{
			ID:           v.ID,
			UserID:       v.UserID,
			VideoURL:     v.VideoURL,
			CoverURL:     v.CoverURL,
			Title:        v.Title,
			Description:  v.Description,
			Duration:     v.Duration,
			Category:     v.Category,
			Tags:         v.Tags,
			VisitCount:   v.VisitCount,
			LikeCount:    v.LikeCount,
			CommentCount: v.CommentCount,
			IsPrivate:    v.IsPrivate,
		}
	}
	return result
}

// extractLastVideoData 提取最后一个视频的数据
func (v *VideoDB) extractLastVideoData(videos []Video) (int64, int64, int64) {
	if len(videos) == 0 {
		return 0, 0, 0
	}

	lastVideo := videos[len(videos)-1]
	return lastVideo.VisitCount, lastVideo.LikeCount, lastVideo.ID
}

// IncrementVisitCount 增加视频播放量
func (v *VideoDB) IncrementVisitCount(ctx context.Context, videoID int64) error {
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
	return v.cache.UpdateVideoScore(ctx, videoID, 1, 0, category)
}

// IncrementLikeCount 增加视频点赞数
func (v *VideoDB) IncrementLikeCount(ctx context.Context, videoID int64) error {
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
	return v.cache.UpdateVideoScore(ctx, videoID, 0, 1, category)
}
