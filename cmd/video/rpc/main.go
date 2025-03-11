package main

import (
	"net"

	"github.com/cloudwego/kitex/server"
	"github.com/yxrrxy/videoHub/app/video/repository"
	"github.com/yxrrxy/videoHub/app/video/service"
	"github.com/yxrrxy/videoHub/config"
	"github.com/yxrrxy/videoHub/kitex_gen/video/videoservice"
	"github.com/yxrrxy/videoHub/pkg/middleware"
)

func main() {
	config.Init()

	db := repository.InitDB()
	videoRepo := repository.NewVideo(db)
	videoService := service.NewVideoService(videoRepo)

	addr, err := net.ResolveTCPAddr("tcp", config.Video.RPCAddr)
	if err != nil {
		panic(err)
	}

	svr := videoservice.NewServer(
		videoService,
		server.WithServiceAddr(addr),
		server.WithMiddleware(middleware.Auth()),
	)

	if err := svr.Run(); err != nil {
		panic(err)
	}
}
