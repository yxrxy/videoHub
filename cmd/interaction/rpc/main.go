package main

import (
	"net"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/yxrrxy/videoHub/app/videoInteractions/repository"
	"github.com/yxrrxy/videoHub/app/videoInteractions/service"
	"github.com/yxrrxy/videoHub/config"
	"github.com/yxrrxy/videoHub/kitex_gen/videoInteractions/interactionservice"
	"github.com/yxrrxy/videoHub/pkg/middleware"
)

func main() {
	config.Init()

	db := repository.InitDB()
	interactionRepo := repository.NewInteraction(db)
	interactionService := service.NewInteractionService(interactionRepo)

	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		panic(err)
	}

	addr, err := net.ResolveTCPAddr("tcp", config.VideoInteractions.RPCAddr)
	if err != nil {
		panic(err)
	}

	svr := interactionservice.NewServer(
		interactionService,
		server.WithServiceAddr(addr),
		server.WithMiddleware(middleware.Auth()),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: config.VideoInteractions.Name,
		}),
	)

	if err := svr.Run(); err != nil {
		panic(err)
	}
}
