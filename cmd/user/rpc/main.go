package main

import (
	"net"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/yxrrxy/videoHub/app/user/repository"
	"github.com/yxrrxy/videoHub/app/user/service"
	"github.com/yxrrxy/videoHub/config"
	"github.com/yxrrxy/videoHub/kitex_gen/user/userservice"
	"github.com/yxrrxy/videoHub/pkg/middleware"
)

func main() {
	config.Init()

	db := repository.InitDB()
	userRepo := repository.NewUser(db)
	userService := service.NewUserService(userRepo)

	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		panic(err)
	}

	addr, err := net.ResolveTCPAddr("tcp", config.User.RPCAddr)
	if err != nil {
		panic(err)
	}

	svr := userservice.NewServer(
		userService,
		server.WithServiceAddr(addr),
		server.WithMiddleware(middleware.Auth()),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: config.User.Name,
		}),
	)

	if err := svr.Run(); err != nil {
		panic(err)
	}
}
