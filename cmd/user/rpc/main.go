package main

import (
	"net"

	"github.com/cloudwego/kitex/server"
	"github.com/yxrrxy/videoHub/app/user/repository"
	"github.com/yxrrxy/videoHub/app/user/service"
	"github.com/yxrrxy/videoHub/config"
	"github.com/yxrrxy/videoHub/kitex_gen/user/userservice"
)

func main() {
	config.Init()

	db := repository.InitDB()
	userRepo := repository.NewUser(db)
	userService := service.NewUserService(userRepo)

	addr, err := net.ResolveTCPAddr("tcp", config.User.RPCAddr)
	if err != nil {
		panic(err)
	}

	svr := userservice.NewServer(
		userService,
		server.WithServiceAddr(addr),
	)

	if err := svr.Run(); err != nil {
		panic(err)
	}
}
