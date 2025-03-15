package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/yxrrxy/videoHub/app/video/controller"
	"github.com/yxrrxy/videoHub/kitex_gen/video/videoservice"
	"github.com/yxrrxy/videoHub/pkg/middleware"
)

func Register(h *server.Hertz, client videoservice.Client) {
	videoController := controller.NewVideoController(client)

	videoGroup := h.Group("/video")
	{
		authedGroup := videoGroup.Group("", middleware.JWT())
		{
			authedGroup.POST("/publish", videoController.Publish)                     // 发布视频
			authedGroup.GET("/list", videoController.GetUserVideoList)                // 获取用户视频列表
			authedGroup.GET("/hot", videoController.GetHotVideoList)                  // 获取热门视频列表
			authedGroup.POST("/increment_visit", videoController.IncrementVisitCount) // 增加播放量
			authedGroup.POST("/increment_like", videoController.IncrementLikeCount)   // 增加点赞数
		}
	}
}
