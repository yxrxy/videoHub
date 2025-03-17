package main

import (
	"net"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/yxrrxy/videoHub/app/social/repository"
	"github.com/yxrrxy/videoHub/app/social/service"
	"github.com/yxrrxy/videoHub/app/social/ws"
	"github.com/yxrrxy/videoHub/config"
	"github.com/yxrrxy/videoHub/kitex_gen/social/socialservice"
	"github.com/yxrrxy/videoHub/kitex_gen/user/userservice"
	"github.com/yxrrxy/videoHub/pkg/middleware"
)

func main() {
	config.Init()

	// 创建etcd注册中心
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		panic(err)
	}

	// 创建etcd解析器
	etcdResolver, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		panic(err)
	}

	// 创建用户服务客户端
	userClient, err := userservice.NewClient(
		config.User.Name,
		client.WithResolver(etcdResolver),
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()),
	)
	if err != nil {
		panic(err)
	}

	// 创建WebSocket管理器
	wsManager := ws.NewManager()

	db := repository.InitDB()
	socialRepo := repository.NewSocial(db)
	socialService := service.NewSocialService(socialRepo, wsManager, userClient)

	addr, err := net.ResolveTCPAddr("tcp", config.Social.RPCAddr)
	if err != nil {
		panic(err)
	}

	svr := socialservice.NewServer(
		socialService,
		server.WithServiceAddr(addr),
		server.WithMiddleware(middleware.Auth()),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: config.Social.Name,
		}),
	)

	if err := svr.Run(); err != nil {
		panic(err)
	}
}
