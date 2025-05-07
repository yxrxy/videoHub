package pack

import (
	"github.com/yxrxy/videoHub/app/video/domain/model"
	rpcmodel "github.com/yxrxy/videoHub/kitex_gen/model"
)

func Videos(v []*model.Video) []*rpcmodel.Video {
	rpcVideos := make([]*rpcmodel.Video, 0)
	for _, video := range v {
		rpcVideos = append(rpcVideos, &rpcmodel.Video{
			Id:            video.ID,
			AuthorId:      video.UserID,
			Title:         video.Title,
			PlayUrl:       video.VideoURL,
			CoverUrl:      video.CoverURL,
			FavoriteCount: &video.LikeCount,
			CommentCount:  &video.CommentCount,
			Description:   &video.Description,
		})
	}
	return rpcVideos
}

func Video(v *model.Video) *rpcmodel.Video {
	return &rpcmodel.Video{
		Id:            v.ID,
		AuthorId:      v.UserID,
		Title:         v.Title,
		PlayUrl:       v.VideoURL,
		CoverUrl:      v.CoverURL,
		FavoriteCount: &v.LikeCount,
		CommentCount:  &v.CommentCount,
		Description:   &v.Description,
	}
}

func SemanticSearchResultItems(v []*model.SemanticSearchResultItem) []*rpcmodel.SemanticSearchResultItem {
	rpcResultItems := make([]*rpcmodel.SemanticSearchResultItem, 0)
	for _, resultItem := range v {
		rpcResultItems = append(rpcResultItems, &rpcmodel.SemanticSearchResultItem{
			Videos:         Videos(resultItem.Videos),
			Summary:        &resultItem.Summary,
			RelatedQueries: resultItem.RelatedQueries,
			FromCache:      &resultItem.FromCache,
		})
	}
	return rpcResultItems
}
