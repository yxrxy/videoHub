package main

import (
	"net"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/yxrrxy/videoHub/app/user/service"
	"github.com/yxrrxy/videoHub/config"
	"github.com/yxrrxy/videoHub/kitex_gen/user/userservice"
	"github.com/yxrrxy/videoHub/pkg/middleware"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", config.User.Addr)
	if err != nil {
		panic(err)
	}

	// 服务器选项
	opts := []server.Option{
		server.WithServiceAddr(addr),                                       // 服务地址
		server.WithMiddleware(middleware.Auth()),                           // 认证中间件
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // 限流配置
		server.WithMuxTransport(),                                          // 多路复用
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: config.User.Name,
		}),
	}

	svr := userservice.NewServer(new(service.UserService), opts...)

	// 启动 server
	if err := svr.Run(); err != nil {
		panic(err)
	}
}
