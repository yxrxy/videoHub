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
			authedGroup.POST("/publish", videoController.Publish)
			authedGroup.GET("/list", videoController.GetUserVideoList)
		}
	}
}
