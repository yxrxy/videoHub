package repository

import (
	"context"

	"github.com/yxrrxy/videoHub/app/videoInteractions/model"
	"github.com/yxrrxy/videoHub/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Interaction struct {
	db *gorm.DB
}

func NewInteraction(db *gorm.DB) *Interaction {
	return &Interaction{db: db}
}

// GetLikeList 获取点赞列表
func (i *Interaction) GetLikeList(ctx context.Context, videoID int64, offset, limit int) ([]*model.Like, error) {
	var likes []*model.Like
	err := i.db.WithContext(ctx).
		Where("video_id = ? AND deleted_at IS NULL", videoID).
		Offset(offset).
		Limit(limit).
		Find(&likes).Error
	return likes, err
}

// CreateLike 创建点赞
func (i *Interaction) CreateLike(ctx context.Context, userID, videoID int64) error {
	return i.db.WithContext(ctx).Create(&model.Like{
		UserID:  userID,
		VideoID: videoID,
	}).Error
}

// GetLike 获取点赞记录
func (i *Interaction) GetLike(ctx context.Context, userID, videoID int64) (*model.Like, error) {
	var like model.Like
	err := i.db.WithContext(ctx).
		Where("user_id = ? AND video_id = ? AND deleted_at IS NULL", userID, videoID).
		First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

// DeleteLike 删除点赞
func (i *Interaction) DeleteLike(ctx context.Context, like *model.Like) error {
	return i.db.WithContext(ctx).Delete(like).Error
}

// CreateComment 创建评论
func (i *Interaction) CreateComment(ctx context.Context, userID, videoID int64, content string, parentID *int64) error {
	return i.db.WithContext(ctx).Create(&model.Comment{
		UserID:   userID,
		VideoID:  videoID,
		Content:  content,
		ParentID: parentID,
	}).Error
}

// GetComment 获取评论
func (i *Interaction) GetComment(ctx context.Context, commentID int64) (*model.Comment, error) {
	var comment model.Comment
	err := i.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", commentID).
		First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetCommentList 获取评论列表
func (i *Interaction) GetCommentList(ctx context.Context, videoID int64, offset, limit int) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	// 获取总数
	err := i.db.WithContext(ctx).Model(&model.Comment{}).
		Where("video_id = ? AND deleted_at IS NULL", videoID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取评论列表
	err = i.db.WithContext(ctx).
		Where("video_id = ? AND deleted_at IS NULL", videoID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// DeleteComment 删除评论
func (i *Interaction) DeleteComment(ctx context.Context, comment *model.Comment) error {
	return i.db.WithContext(ctx).Delete(comment).Error
}

// LikeComment 点赞评论
func (i *Interaction) LikeComment(ctx context.Context, userID, commentID int64) error {
	return i.db.WithContext(ctx).Create(&model.CommentLike{
		UserID:    userID,
		CommentID: commentID,
	}).Error
}

// UnlikeComment 取消评论点赞
func (i *Interaction) UnlikeComment(ctx context.Context, userID, commentID int64) error {
	return i.db.WithContext(ctx).Where("user_id = ? AND comment_id = ?", userID, commentID).
		Delete(&model.CommentLike{}).Error
}

// IsCommentLiked 检查用户是否点赞了评论
func (i *Interaction) IsCommentLiked(ctx context.Context, userID, commentID int64) (bool, error) {
	var count int64
	err := i.db.WithContext(ctx).Model(&model.CommentLike{}).
		Where("user_id = ? AND comment_id = ?", userID, commentID).
		Count(&count).Error
	return count > 0, err
}

// GetCommentLikeCount 获取评论点赞数
func (i *Interaction) GetCommentLikeCount(ctx context.Context, commentID int64) (int64, error) {
	var count int64
	err := i.db.WithContext(ctx).Model(&model.CommentLike{}).
		Where("comment_id = ?", commentID).
		Count(&count).Error
	return count, err
}

// InitDB 初始化数据库连接
func InitDB() *gorm.DB {
	dsn := config.GetDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(&model.Like{}, &model.Comment{}, &model.CommentLike{}); err != nil {
		panic(err)
	}

	return db
}
