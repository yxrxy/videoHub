package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/yxrrxy/videoHub/app/user/controller"
	"github.com/yxrrxy/videoHub/kitex_gen/user/userservice"
	"github.com/yxrrxy/videoHub/pkg/middleware"
)

func Register(h *server.Hertz, client userservice.Client) {
	userController := controller.NewUserController(client)

	userGroup := h.Group("/user")
	{
		userGroup.POST("/register", userController.Register) // 注册
		userGroup.POST("/login", userController.Login)       // 登录

		// 需要认证的路由组
		authedGroup := userGroup.Group("", middleware.JWT()) // HTTP 中间件
		{
			authedGroup.GET("/info", userController.GetUserInfo)     // 获取用户信息
			authedGroup.POST("/avatar", userController.UploadAvatar) // 上传头像
		}
	}
}
