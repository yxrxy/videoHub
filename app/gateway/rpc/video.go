package rpc

import (
	"context"
	"log"

	"github.com/yxrxy/videoHub/kitex_gen/model"
	"github.com/yxrxy/videoHub/kitex_gen/video"
	"github.com/yxrxy/videoHub/pkg/base/client"
	"github.com/yxrxy/videoHub/pkg/errno"
)

// InitVideoRPC 初始化视频服务客户端
func InitVideoRPC() {
	c, err := client.InitVideoRPC()
	if err != nil {
		log.Fatalf("初始化视频服务客户端失败: %v", err)
	}
	videoClient = *c
}

// PublishVideoRPC 发布视频
func PublishVideoRPC(ctx context.Context, req *video.PublishRequest) error {
	resp, err := videoClient.Publish(ctx, req)
	if err != nil {
		log.Printf("发布视频RPC调用失败: %v", err)
		return errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return nil
}

// GetVideoListRPC 获取视频列表
func GetVideoListRPC(ctx context.Context, req *video.VideoListRequest) ([]*model.Video, error) {
	resp, err := videoClient.List(ctx, req)
	if err != nil {
		log.Printf("获取视频列表RPC调用失败: %v", err)
		return nil, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.VideoList, nil
}

// GetVideoDetailRPC 获取视频详情
func GetVideoDetailRPC(ctx context.Context, req *video.DetailRequest) (*model.Video, error) {
	resp, err := videoClient.Detail(ctx, req)
	if err != nil {
		log.Printf("获取视频详情RPC调用失败: %v", err)
		return nil, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.Video, nil
}

// DeleteVideoRPC 删除视频
func DeleteVideoRPC(ctx context.Context, req *video.DeleteRequest) error {
	resp, err := videoClient.Delete(ctx, req)
	if err != nil {
		log.Printf("删除视频RPC调用失败: %v", err)
		return errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return nil
}

// GetHotVideosRPC 获取热门视频
func GetHotVideosRPC(ctx context.Context, req *video.HotVideoRequest) ([]*model.Video, error) {
	resp, err := videoClient.GetHotVideos(ctx, req)
	if err != nil {
		log.Printf("获取热门视频RPC调用失败: %v", err)
		return nil, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.Videos, nil
}

// IncrementVisitCountRPC 增加视频访问量
func IncrementVisitCountRPC(ctx context.Context, req *video.IncrementVisitCountRequest) error {
	resp, err := videoClient.IncrementVisitCount(ctx, req)
	if err != nil {
		log.Printf("增加视频访问量RPC调用失败: %v", err)
		return errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return nil
}

// IncrementLikeCountRPC 增加视频点赞数
func IncrementLikeCountRPC(ctx context.Context, req *video.IncrementLikeCountRequest) error {
	resp, err := videoClient.IncrementLikeCount(ctx, req)
	if err != nil {
		log.Printf("增加视频点赞数RPC调用失败: %v", err)
		return errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return nil
}
