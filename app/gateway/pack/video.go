package pack

import (
	"github.com/yxrxy/videoHub/app/gateway/model/model"
	rpcmodel "github.com/yxrxy/videoHub/kitex_gen/model"
	"github.com/yxrxy/videoHub/pkg/base"
)

// BuildVideo 将 RPC 交流实体转换成 http 返回的实体
func BuildVideo(v *model.Video) *rpcmodel.Video {
	if v == nil {
		return nil
	}
	return &rpcmodel.Video{
		Id:            v.ID,
		AuthorId:      v.AuthorId,
		PlayUrl:       v.PlayUrl,
		CoverUrl:      v.CoverUrl,
		Title:         v.Title,
		Description:   v.Description,
		FavoriteCount: v.FavoriteCount,
		CommentCount:  v.CommentCount,
		CreatedAt:     v.CreatedAt,
		UpdatedAt:     v.UpdatedAt,
	}
}

// BuildVideoList 构建视频列表
func BuildVideoList(vs []*model.Video) []*rpcmodel.Video {
	return base.BuildTypeList(vs, BuildVideo)
}
