package rpc

import (
	"context"
	"log"

	"github.com/yxrxy/videoHub/kitex_gen/interaction"
	"github.com/yxrxy/videoHub/kitex_gen/model"
	"github.com/yxrxy/videoHub/pkg/base/client"
	"github.com/yxrxy/videoHub/pkg/errno"
)

// InitInteractionRPC 初始化互动服务客户端
func InitInteractionRPC() {
	c, err := client.InitInteractionRPC()
	if err != nil {
		log.Fatalf("初始化互动服务客户端失败: %v", err)
	}
	interactionClient = *c
}

// LikeVideoRPC 点赞视频
func LikeVideoRPC(ctx context.Context, req *interaction.LikeRequest) error {
	success, err := interactionClient.Like(ctx, req)
	if err != nil {
		log.Printf("点赞RPC调用失败: %v", err)
		return errno.InternalServiceError.WithError(err)
	}
	if success.Base.Code != errno.SuccessCode {
		return errno.InternalServiceError.WithMessage(success.Base.Msg)
	}
	return nil
}

// GetLikesRPC 获取点赞列表
func GetLikesRPC(ctx context.Context, req *interaction.GetLikesRequest) ([]*model.LikeInfo, int64, error) {
	resp, err := interactionClient.GetLikes(ctx, req)
	if err != nil {
		log.Printf("获取点赞列表RPC调用失败: %v", err)
		return nil, 0, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, 0, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.LikeList, resp.Total, nil
}

// AddCommentRPC 添加评论
func CommentVideoRPC(ctx context.Context, req *interaction.CommentRequest) error {
	success, err := interactionClient.Comment(ctx, req)
	if err != nil {
		log.Printf("发表评论RPC调用失败: %v", err)
		return errno.InternalServiceError.WithError(err)
	}
	if success.Base.Code != errno.SuccessCode {
		return errno.InternalServiceError.WithMessage(success.Base.Msg)
	}
	return nil
}

// GetCommentsRPC 获取评论列表
func GetCommentsRPC(ctx context.Context, req *interaction.GetCommentsRequest) ([]*model.Comment, int64, error) {
	resp, err := interactionClient.GetComments(ctx, req)
	if err != nil {
		log.Printf("获取评论列表RPC调用失败: %v", err)
		return nil, 0, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, 0, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.CommentList, resp.Total, nil
}

func DeleteCommentRPC(ctx context.Context, req *interaction.DeleteCommentRequest) error {
	success, err := interactionClient.DeleteComment(ctx, req)
	if err != nil {
		log.Printf("删除评论RPC调用失败: %v", err)
		return errno.InternalServiceError.WithError(err)
	}
	if success.Base.Code != errno.SuccessCode {
		return errno.InternalServiceError.WithMessage(success.Base.Msg)
	}
	return nil
}

func LikeCommentRPC(ctx context.Context, req *interaction.LikeCommentRequest) error {
	success, err := interactionClient.LikeComment(ctx, req)
	if err != nil {
		log.Printf("点赞评论RPC调用失败: %v", err)
		return errno.InternalServiceError.WithError(err)
	}
	if success.Base.Code != errno.SuccessCode {
		return errno.InternalServiceError.WithMessage(success.Base.Msg)
	}
	return nil
}
