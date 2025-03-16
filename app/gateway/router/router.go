package router

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/yxrrxy/videoHub/app/user/controller"
	videoController "github.com/yxrrxy/videoHub/app/video/controller"
	"github.com/yxrrxy/videoHub/pkg/middleware"
)

type Router struct {
	userCtrl  *controller.UserController
	videoCtrl *videoController.VideoController
}

func NewRouter(userCtrl *controller.UserController, videoCtrl *videoController.VideoController) *Router {
	return &Router{
		userCtrl:  userCtrl,
		videoCtrl: videoCtrl,
	}
}

func (r *Router) Register(h *server.Hertz) {
	// 静态文件服务
	h.StaticFS("/static/avatars", &app.FS{
		Root: "src/storage/avatars",
	})
	h.StaticFS("/static/videos", &app.FS{
		Root: "src/storage/videos",
	})

	// 健康检查
	h.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API 路由组
	api := h.Group("/api")
	{
		// 用户服务路由
		user := api.Group("/user")
		{
			// 无需认证的接口
			user.POST("/register", r.userCtrl.Register)
			user.POST("/login", r.userCtrl.Login)

			// 需要认证的接口
			authed := user.Group("", middleware.JWT())
			{
				authed.GET("/info", r.userCtrl.GetUserInfo)
				authed.POST("/avatar", r.userCtrl.UploadAvatar)
			}
		}

		// 视频服务路由
		video := api.Group("/video")
		{
			// 公开接口
			video.GET("/hot", r.videoCtrl.GetHotVideoList)
			//video.GET("/detail/:id", r.videoCtrl.GetVideoDetail)

			// 需要认证的接口
			authed := video.Group("", middleware.JWT())
			{
				authed.POST("/publish", r.videoCtrl.Publish)
				authed.GET("/list", r.videoCtrl.GetUserVideoList)
				authed.POST("/visit", r.videoCtrl.IncrementVisitCount)
				authed.POST("/like", r.videoCtrl.IncrementLikeCount)
			}
		}
	}
}
