package main

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/yxrrxy/videoHub/app/video/cache"
	"github.com/yxrrxy/videoHub/app/video/repository"
	"github.com/yxrrxy/videoHub/app/video/router"
	"github.com/yxrrxy/videoHub/config"
	"github.com/yxrrxy/videoHub/kitex_gen/video/videoservice"
)

func main() {
	config.Init()

	// 初始化数据库
	repository.InitDB()

	// 初始化 Redis
	cache.RedisInit()

	// 初始化 RPC 客户端
	client := videoservice.MustNewClient("video")

	// 初始化 HTTP 服务器
	h := server.Default(server.WithHostPorts(config.Video.HTTPAddr))

	// 静态文件服务配置
	h.StaticFS("/videos", &app.FS{
		Root: "src/storage/videos",
	})

	// 注册路由
	router.Register(h, client)

	// 启动服务器
	h.Spin()
}
