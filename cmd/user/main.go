package main

import (
	"github.com/cloudwego/kitex/server"
	"github.com/yxrrxy/videoHub/app/user/service"
	"github.com/yxrrxy/videoHub/kitex_gen/user/userservice"
	"github.com/yxrrxy/videoHub/pkg/middleware"
)

func main() {
	svr := userservice.NewServer(
		new(service.UserService),
		server.WithMiddleware(middleware.Auth()),
		server.WithServiceAddr(addr),

	svr.Run()
}
