package pack

import (
	"github.com/yxrxy/videoHub/app/gateway/model/model"
	rpcmodel "github.com/yxrxy/videoHub/kitex_gen/model"
	"github.com/yxrxy/videoHub/pkg/base"
)

// BuildLikeInfo 将 RPC 交流实体转换成 http 返回的实体
func BuildLikeInfo(l *model.LikeInfo) *rpcmodel.LikeInfo {
	if l == nil {
		return nil
	}
	return &rpcmodel.LikeInfo{
		Id:        l.ID,
		UserId:    l.UserID,
		VideoId:   l.VideoID,
		CreatedAt: l.CreatedAt,
		DeletedAt: l.DeletedAt,
	}
}

// BuildLikeInfoList 构建点赞列表
func BuildLikeInfoList(ls []*model.LikeInfo) []*rpcmodel.LikeInfo {
	return base.BuildTypeList(ls, BuildLikeInfo)
}

// BuildCommentInfo 将 RPC 交流实体转换成 http 返回的实体
func BuildCommentInfo(c *model.CommentInfo) *rpcmodel.CommentInfo {
	if c == nil {
		return nil
	}
	return &rpcmodel.CommentInfo{
		Id:        c.ID,
		UserId:    c.UserID,
		VideoId:   c.VideoID,
		Content:   c.Content,
		ParentId:  c.ParentID,
		CreatedAt: c.CreatedAt,
		DeletedAt: c.DeletedAt,
		LikeCount: c.LikeCount,
		IsLiked:   c.IsLiked,
	}
}

// BuildCommentInfoList 构建评论列表
func BuildCommentInfoList(cs []*model.CommentInfo) []*rpcmodel.CommentInfo {
	return base.BuildTypeList(cs, BuildCommentInfo)
}
