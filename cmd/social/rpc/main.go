package main

import (
	"context"
	"net"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	kitexserver "github.com/cloudwego/kitex/server"
	"github.com/hertz-contrib/websocket"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/yxrrxy/videoHub/app/social/repository"
	"github.com/yxrrxy/videoHub/app/social/service"
	"github.com/yxrrxy/videoHub/app/social/ws"
	"github.com/yxrrxy/videoHub/config"
	"github.com/yxrrxy/videoHub/kitex_gen/social/socialservice"
	"github.com/yxrrxy/videoHub/kitex_gen/user/userservice"
	pkgcontext "github.com/yxrrxy/videoHub/pkg/context"
	"github.com/yxrrxy/videoHub/pkg/errno"
	"github.com/yxrrxy/videoHub/pkg/middleware"
	"github.com/yxrrxy/videoHub/pkg/response"
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

	// 启动WebSocket管理器
	go wsManager.Start(context.Background())
	go wsManager.StartHeartbeat(context.Background(), 30*time.Second)

	db := repository.InitDB()
	socialRepo := repository.NewSocial(db)
	socialService := service.NewSocialService(socialRepo, wsManager, userClient)

	// 启动RPC服务
	go func() {
		addr, err := net.ResolveTCPAddr("tcp", config.Social.RPCAddr)
		if err != nil {
			panic(err)
		}

		svr := socialservice.NewServer(
			socialService,
			kitexserver.WithServiceAddr(addr),
			kitexserver.WithMiddleware(middleware.Auth()),
			kitexserver.WithRegistry(r),
			kitexserver.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
				ServiceName: config.Social.Name,
			}),
		)

		if err := svr.Run(); err != nil {
			panic(err)
		}
	}()

	// 启动HTTP服务（WebSocket服务）
	startHTTPServer(socialService)
}

// startHTTPServer 启动HTTP服务（用于WebSocket）
func startHTTPServer(socialService *service.SocialService) {
	httpAddr := config.Social.HTTPAddr
	if httpAddr == "" {
		httpAddr = ":8083"
	}

	h := server.Default(server.WithHostPorts(httpAddr))

	// 使用JWT中间件
	h.Use(middleware.JWT())

	// WebSocket路由
	h.GET("/ws/connect", func(ctx context.Context, c *app.RequestContext) {
		userID, exists := pkgcontext.GetUserID(ctx)
		if !exists {
			c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
			return
		}

		upgrader := websocket.HertzUpgrader{
			CheckOrigin: func(ctx *app.RequestContext) bool {
				return true // 在生产环境中应该检查来源
			},
		}

		err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
			socialService.RegisterWebSocketClient(userID, conn)
		})

		if err != nil {
			c.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
			return
		}
	})

	// 健康检查
	h.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, map[string]interface{}{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	h.Spin()
}
