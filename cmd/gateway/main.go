package main

import (
	"log"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/yxrrxy/videoHub/app/gateway/router"
	userController "github.com/yxrrxy/videoHub/app/user/controller"
	videoController "github.com/yxrrxy/videoHub/app/video/controller"
	"github.com/yxrrxy/videoHub/config"
	"github.com/yxrrxy/videoHub/kitex_gen/user/userservice"
	"github.com/yxrrxy/videoHub/kitex_gen/video/videoservice"
)

func main() {
	// 初始化配置
	config.Init()

	// 创建etcd解析器
	etcdResolver, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		log.Fatalf("创建etcd解析器失败: %v", err)
	}

	// 创建RPC客户端
	userClient, err := userservice.NewClient(
		"user",
		client.WithResolver(etcdResolver), // 使用etcd服务发现
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()), // 负载均衡
	)
	if err != nil {
		log.Fatalf("创建用户服务客户端失败: %v", err)
	}

	videoClient, err := videoservice.NewClient(
		"video",
		client.WithResolver(etcdResolver), // 使用etcd服务发现
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()), // 负载均衡
	)
	if err != nil {
		log.Fatalf("创建视频服务客户端失败: %v", err)
	}

	// 创建控制器
	userCtrl := userController.NewUserController(userClient)
	videoCtrl := videoController.NewVideoController(videoClient)

	// 创建路由
	router := router.NewRouter(userCtrl, videoCtrl)

	// 创建HTTP服务器
	h := server.Default(server.WithHostPorts(config.Gateway.Addr))

	// 注册路由
	router.Register(h)

	if err := h.Run(); err != nil {
		log.Fatal(err)
	}
}
