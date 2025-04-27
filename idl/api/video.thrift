namespace go api.video

include "../video.thrift"

// API 服务
service VideoAPI {
    // 视频上传和管理接口
    video.PublishResponse Publish(1: video.PublishRequest request) (api.post="/api/v1/video/publish")
    video.VideoListResponse GetVideoList(1: video.VideoListRequest request) (api.get="/api/v1/video/list")
    video.DetailResponse GetVideoDetail(1: video.DetailRequest request) (api.get="/api/v1/video/:video_id")
    video.HotVideoResponse GetHotVideos(1: video.HotVideoRequest request) (api.get="/api/v1/video/hot")
    video.DeleteResponse DeleteVideo(1: video.DeleteRequest request) (api.delete="/api/v1/video/:video_id")
    video.SearchResponse SearchVideo(1: video.SearchRequest request) (api.post="/api/v1/video/search")
    

    // 视频互动接口
    video.IncrementVisitCountResponse IncrementVisitCount(1: video.IncrementVisitCountRequest request) (api.post="/api/v1/video/:video_id/visit")
    video.IncrementLikeCountResponse IncrementLikeCount(1: video.IncrementLikeCountRequest request) (api.post="/api/v1/video/:video_id/like")
} 