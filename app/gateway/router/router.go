package router

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	socialController "github.com/yxrrxy/videoHub/app/social/controller"
	"github.com/yxrrxy/videoHub/app/user/controller"
	videoController "github.com/yxrrxy/videoHub/app/video/controller"
	interactionController "github.com/yxrrxy/videoHub/app/videoInteractions/controller"
	"github.com/yxrrxy/videoHub/pkg/middleware"
)

type Router struct {
	userCtrl        *controller.UserController
	videoCtrl       *videoController.VideoController
	socialCtrl      *socialController.SocialHandler
	interactionCtrl *interactionController.InteractionController
}

func NewRouter(
	userCtrl *controller.UserController,
	videoCtrl *videoController.VideoController,
	socialCtrl *socialController.SocialHandler,
	interactionCtrl *interactionController.InteractionController,
) *Router {
	return &Router{
		userCtrl:        userCtrl,
		videoCtrl:       videoCtrl,
		socialCtrl:      socialCtrl,
		interactionCtrl: interactionCtrl,
	}
}

func (r *Router) Register(h *server.Hertz) {
	h.Use(middleware.CORS())

	// 静态文件服务
	h.StaticFS("/", &app.FS{
		Root: "src/storage",
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
		// 添加ping路由用于测试连接
		api.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
			c.JSON(200, map[string]interface{}{
				"status":  "ok",
				"message": "pong",
				"time":    time.Now().Format(time.RFC3339),
			})
		})

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

		// 互动服务路由组
		interaction := api.Group("/interaction")
		{
			// 公开接口
			interaction.GET("/likes", r.interactionCtrl.GetLikes)
			interaction.GET("/comments", r.interactionCtrl.GetComments)

			// 需要认证的接口
			authed := interaction.Group("", middleware.JWT())
			{
				authed.POST("/like", r.interactionCtrl.Like)
				authed.POST("/comment", r.interactionCtrl.Comment)
				authed.DELETE("/comment", r.interactionCtrl.DeleteComment)
				authed.POST("/comment/like", r.interactionCtrl.LikeComment)
			}
		}

		// 社交功能路由组
		social := api.Group("/social", middleware.JWT())
		{
			// 私信相关
			social.POST("/messages", r.socialCtrl.SendPrivateMessage)
			social.GET("/messages", r.socialCtrl.GetPrivateMessages)

			// 聊天室相关
			social.POST("/chat/rooms", r.socialCtrl.CreateChatRoom)
			social.GET("/chat/rooms/:id", r.socialCtrl.GetChatRoom)
			social.GET("/chat/rooms", r.socialCtrl.GetUserChatRooms)
			social.POST("/chat/rooms/:id/messages", r.socialCtrl.SendChatMessage)
			social.GET("/chat/rooms/:id/messages", r.socialCtrl.GetChatMessages)

			// 好友相关
			social.POST("/friends", r.socialCtrl.AddFriend)
			social.GET("/friends", r.socialCtrl.GetUserFriends)
			social.GET("/friends/:id", r.socialCtrl.GetFriendship)

			// 好友申请相关
			social.POST("/friend-requests", r.socialCtrl.CreateFriendRequest)
			social.GET("/friend-requests", r.socialCtrl.GetFriendRequests)
			social.PUT("/friend-requests/:id", r.socialCtrl.HandleFriendRequest)

			// 消息状态相关
			social.PUT("/messages/:id/read", r.socialCtrl.MarkMessageRead)
			social.GET("/messages/unread/count", r.socialCtrl.GetUnreadMessageCount)
		}
	}
}
