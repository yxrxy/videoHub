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

		authedGroup := userGroup.Group("", middleware.JWT()) // 添加 JWT 中间件
		{
			//authedGroup.GET("/info", userController.GetUserInfo)         // 获取用户信息
			authedGroup.POST("/avatar", userController.UploadAvatar) // 上传头像
			//authedGroup.GET("/mfa/qrcode", userController.GetMFAQRCode) // 获取 MFA 二维码
			//authedGroup.POST("/mfa/bind", userController.BindMFA)       // 绑定 MFA
		}
	}
}
