package main

import (
	"context"
	"log"
	"net"
	"strconv"
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
	"github.com/yxrrxy/videoHub/pkg/errno"
	"github.com/yxrrxy/videoHub/pkg/jwt"
	"github.com/yxrrxy/videoHub/pkg/middleware"
	"github.com/yxrrxy/videoHub/pkg/response"
)

func main() {
	config.Init()

	// 创建etcd注册中心
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		log.Fatalf("创建etcd注册中心失败: %v", err)
	}

	// 创建etcd解析器
	etcdResolver, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		log.Fatalf("创建etcd解析器失败: %v", err)
	}

	userClient, err := userservice.NewClient(
		config.User.Name,
		client.WithResolver(etcdResolver),
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()),
	)
	if err != nil {
		log.Fatalf("创建用户服务客户端失败: %v", err)
	}

	wsManager := ws.NewManager()

	ctx := context.Background()
	go wsManager.Start(ctx)
	go wsManager.StartHeartbeat(ctx, 30*time.Second)

	db := repository.InitDB()
	socialRepo := repository.NewSocial(db)
	socialService := service.NewSocialService(socialRepo, wsManager, userClient)

	go func() {
		addr, err := net.ResolveTCPAddr("tcp", config.Social.RPCAddr)
		if err != nil {
			log.Fatalf("解析RPC地址失败: %v", err)
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
			log.Fatalf("RPC服务运行失败: %v", err)
		}
	}()

	startHTTPServer(socialService)
}

// startHTTPServer 启动HTTP服务（用于WebSocket）
func startHTTPServer(socialService *service.SocialService) {
	h := server.Default(server.WithHostPorts(config.Social.HTTPAddr))
	h.NoHijackConnPool = true

	// WebSocket路由
	h.GET("/ws/connect", func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("WebSocket处理器恢复自panic: %v", r)
			}
		}()

		token := c.Query("token")
		if token == "" {
			c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidToken.ErrCode, errno.ErrInvalidToken.ErrMsg))
			return
		}

		// 获取聊天室ID
		roomIDStr := c.Query("room_id")
		if roomIDStr == "" {
			c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, "缺少room_id参数"))
			return
		}
		roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
		if err != nil {
			c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, "room_id参数无效"))
			return
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
			return
		}

		userID := claims.UserID

		upgrader := websocket.HertzUpgrader{
			CheckOrigin: func(ctx *app.RequestContext) bool {
				return true
			},
			HandshakeTimeout: 30 * time.Second,
			ReadBufferSize:   1024,
			WriteBufferSize:  1024,
		}

		if err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
			defer conn.Close()
			defer func() {
				// 离开聊天室
				if err := socialService.LeaveChatRoom(userID, roomID); err != nil {
					log.Printf("离开聊天室失败: %v", err)
				}
			}()

			// 验证聊天室是否存在
			if _, err := socialService.GetChatRoom(context.Background(), roomID); err != nil {
				log.Printf("聊天室不存在: %v", err)
				conn.WriteJSON(map[string]interface{}{
					"type":    "error",
					"message": "聊天室不存在",
				})
				return
			}

			socialService.RegisterWebSocketClient(userID, conn)

			// 加入聊天室
			if err := socialService.JoinChatRoom(userID, roomID); err != nil {
				log.Printf("加入聊天室失败: %v", err)
				conn.WriteJSON(map[string]interface{}{
					"type":    "error",
					"message": "加入聊天室失败: " + err.Error(),
				})
				return
			}

			// 获取最近的消息历史
			messages, err := socialService.GetChatMessages(context.Background(), roomID, 1, 20)
			if err != nil {
				log.Printf("获取历史消息失败: %v", err)
			} else {
				conn.WriteJSON(map[string]interface{}{
					"type":     "history",
					"messages": messages.Messages,
				})
			}

			// 发送欢迎消息
			conn.WriteJSON(map[string]interface{}{
				"type":    "system",
				"message": "欢迎加入聊天室",
				"room_id": roomID,
			})

			// 持续读取消息
			for {
				messageType, message, err := conn.ReadMessage()
				if err != nil {
					log.Printf("读取消息错误: %v", err)
					return
				}

				if messageType != websocket.TextMessage {
					continue
				}

				// 发送消息到聊天室
				if err := socialService.SendChatMessage(context.Background(), roomID, userID, string(message), 0); err != nil {
					log.Printf("发送消息失败: %v", err)
					conn.WriteJSON(map[string]interface{}{
						"type":    "error",
						"message": "发送消息失败: " + err.Error(),
					})
				}
			}
		}); err != nil {
			log.Printf("升级WebSocket连接失败: %v", err)
		}
	})

	// 健康检查路由
	h.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, map[string]interface{}{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	h.Spin()
}
