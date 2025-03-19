package main

import (
	"log"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/yxrrxy/videoHub/app/gateway/router"
	socialController "github.com/yxrrxy/videoHub/app/social/controller"
	userController "github.com/yxrrxy/videoHub/app/user/controller"
	videoController "github.com/yxrrxy/videoHub/app/video/controller"
	interactionController "github.com/yxrrxy/videoHub/app/videoInteractions/controller"
	"github.com/yxrrxy/videoHub/config"
	"github.com/yxrrxy/videoHub/kitex_gen/social/socialservice"
	"github.com/yxrrxy/videoHub/kitex_gen/user/userservice"
	"github.com/yxrrxy/videoHub/kitex_gen/video/videoservice"
	"github.com/yxrrxy/videoHub/kitex_gen/videoInteractions/interactionservice"
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
		config.User.Name,
		client.WithResolver(etcdResolver), // 使用etcd服务发现
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()), // 负载均衡
	)
	if err != nil {
		log.Fatalf("创建用户服务客户端失败: %v", err)
	}

	videoClient, err := videoservice.NewClient(
		config.Video.Name,
		client.WithResolver(etcdResolver),                          // 使用etcd服务发现
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()), // 负载均衡
	)
	if err != nil {
		log.Fatalf("创建视频服务客户端失败: %v", err)
	}

	// 创建社交服务客户端
	socialClient, err := socialservice.NewClient(
		config.Social.Name,
		client.WithResolver(etcdResolver),                          // 使用etcd服务发现
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()), // 负载均衡
	)
	if err != nil {
		log.Fatalf("创建社交服务客户端失败: %v", err)
	}

	// 创建互动服务客户端
	interactionClient, err := interactionservice.NewClient(
		config.VideoInteractions.Name,
		client.WithResolver(etcdResolver),                          // 使用etcd服务发现
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()), // 负载均衡
	)
	if err != nil {
		log.Fatalf("创建互动服务客户端失败: %v", err)
	}

	// 创建控制器
	userCtrl := userController.NewUserController(userClient)
	videoCtrl := videoController.NewVideoController(videoClient)
	socialCtrl := socialController.NewSocialHandler(socialClient)
	interactionCtrl := interactionController.NewInteractionController(interactionClient)

	// 创建路由
	router := router.NewRouter(userCtrl, videoCtrl, socialCtrl, interactionCtrl)

	// 创建HTTP服务器
	h := server.Default(server.WithHostPorts(config.Gateway.Addr))

	// 注册路由
	router.Register(h)

	// 启动服务器
	h.Spin()
}
