package rpc

import (
	"context"

	"github.com/yxrxy/videoHub/app/video/controllers/rpc/pack"
	"github.com/yxrxy/videoHub/app/video/domain/model"
	"github.com/yxrxy/videoHub/app/video/usecase"
	"github.com/yxrxy/videoHub/kitex_gen/video"
	"github.com/yxrxy/videoHub/pkg/base"
	pkgcontext "github.com/yxrxy/videoHub/pkg/base/context"
)

type VideoHandler struct {
	useCase usecase.VideoUseCase
}

func NewVideoHandler(useCase usecase.VideoUseCase) *VideoHandler {
	return &VideoHandler{useCase: useCase}
}

func (h *VideoHandler) Publish(ctx context.Context, req *video.PublishRequest) (r *video.PublishResponse, err error) {
	r = new(video.PublishResponse)
	userId, err := pkgcontext.GetUserID(ctx)
	if err != nil {
		return r, err
	}
	description := ""
	if req.Description != nil {
		description = *req.Description
	}

	category := "all"
	if req.Category != nil {
		category = *req.Category
	}

	tags := make([]string, 0)
	if req.Tags != nil {
		tags = req.Tags
	}

	isPrivate := false
	if req.IsPrivate != nil {
		isPrivate = *req.IsPrivate
	}
	videoURL, err := h.useCase.Publish(ctx, userId, req.VideoData, req.ContentType, req.Title, &description, &category, tags, isPrivate)
	if err != nil {
		return r, err
	}
	r.Base = base.BuildBaseResp(err)
	r.VideoUrl = videoURL
	return r, err
}

func (h *VideoHandler) List(ctx context.Context, req *video.VideoListRequest) (r *video.VideoListResponse, err error) {
	r = new(video.VideoListResponse)
	userId, err := pkgcontext.GetUserID(ctx)
	if err != nil {
		return
	}
	var videoList []*model.Video
	var total int64
	if videoList, total, err = h.useCase.GetVideoList(ctx, userId, req.Page, req.Size, req.Category); err != nil {
		return
	}
	r.VideoList = pack.Videos(videoList)
	r.Total = total
	r.Base = base.BuildBaseResp(err)
	return
}

func (h *VideoHandler) Detail(ctx context.Context, req *video.DetailRequest) (r *video.DetailResponse, err error) {
	r = new(video.DetailResponse)
	userId, err := pkgcontext.GetUserID(ctx)
	if err != nil {
		return
	}

	var video *model.Video
	if video, err = h.useCase.GetVideoDetail(ctx, req.VideoId, userId); err != nil {
		return
	}
	r.Video = pack.Video(video)
	r.Base = base.BuildBaseResp(err)
	return
}

func (h *VideoHandler) GetHotVideos(ctx context.Context, req *video.HotVideoRequest) (r *video.HotVideoResponse, err error) {
	r = new(video.HotVideoResponse)

	var videoList []*model.Video
	var total int64
	var lastVisit, lastLike, lastID int64
	limit := int32(10)
	category := "all"
	ulastVisit := int64(0)
	ulastLike := int64(0)
	ulastID := int64(0)
	if req.Limit == nil {
		req.Limit = &limit
	}
	if req.Category == nil {
		req.Category = &category
	}
	if req.LastVisit == nil {
		req.LastVisit = &ulastVisit
	}
	if req.LastLike == nil {
		req.LastLike = &ulastLike
	}
	if req.LastId == nil {
		req.LastId = &ulastID
	}

	if videoList, total, lastVisit, lastLike, lastID, err = h.useCase.GetHotVideos(ctx, int32(*req.Limit), *req.Category, *req.LastVisit, *req.LastLike, *req.LastId); err != nil {
		return r, err
	}
	r.Videos = pack.Videos(videoList)
	r.Total = total
	r.NextVisit = &lastVisit
	r.NextLike = &lastLike
	r.NextId = &lastID
	r.Base = base.BuildBaseResp(err)
	return r, err
}

func (h *VideoHandler) Delete(ctx context.Context, req *video.DeleteRequest) (r *video.DeleteResponse, err error) {
	r = new(video.DeleteResponse)

	userId, err := pkgcontext.GetUserID(ctx)
	if err != nil {
		return
	}
	if err = h.useCase.DeleteVideo(ctx, req.VideoId, userId); err != nil {
		return
	}
	r.Base = base.BuildBaseResp(err)
	return
}

func (h *VideoHandler) IncrementVisitCount(ctx context.Context, req *video.IncrementVisitCountRequest) (r *video.IncrementVisitCountResponse, err error) {
	r = new(video.IncrementVisitCountResponse)

	if err = h.useCase.IncrementVisitCount(ctx, req.VideoId); err != nil {
		return
	}
	r.Base = base.BuildBaseResp(err)
	return
}

func (h *VideoHandler) IncrementLikeCount(ctx context.Context, req *video.IncrementLikeCountRequest) (r *video.IncrementLikeCountResponse, err error) {
	r = new(video.IncrementLikeCountResponse)

	if err = h.useCase.IncrementLikeCount(ctx, req.VideoId); err != nil {
		return
	}
	r.Base = base.BuildBaseResp(err)
	return
}

func (h *VideoHandler) Search(ctx context.Context, req *video.SearchRequest) (r *video.SearchResponse, err error) {
	r = new(video.SearchResponse)

	var videoList []*model.Video
	var total int64
	if videoList, total, err = h.useCase.SearchVideo(ctx, req.Keywords, req.PageSize, req.PageNum, req.FromDate, req.ToDate, req.Username); err != nil {
		return
	}
	r.Videos = pack.Videos(videoList)
	r.Total = total
	r.Base = base.BuildBaseResp(err)
	return
}
