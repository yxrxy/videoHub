package main

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/yxrrxy/videoHub/app/video/router"
	"github.com/yxrrxy/videoHub/config"
	"github.com/yxrrxy/videoHub/kitex_gen/video/videoservice"
)

func main() {
	config.Init()

	h := server.Default(server.WithHostPorts(config.Video.HTTPAddr))

	c, err := videoservice.NewClient(
		"video",
		client.WithHostPorts("127.0.0.1"+config.Video.RPCAddr),
	)
	if err != nil {
		panic(err)
	}

	// 静态文件服务配置
	h.StaticFS("/videos", &app.FS{
		Root: "src/storage/videos",
	})

	router.Register(h, c)

	h.Spin()
}
