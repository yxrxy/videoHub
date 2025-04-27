package main

import (
	"context"
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/yxrxy/videoHub/app/interaction"
	"github.com/yxrxy/videoHub/config"
	"github.com/yxrxy/videoHub/kitex_gen/interaction/interactionservice"
	"github.com/yxrxy/videoHub/pkg/base"
	"github.com/yxrxy/videoHub/pkg/middleware"
)

func init() {
	config.Init("interaction")
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		log.Fatalf("User: new etcd registry failed, err: %v", err)
	}
	listenAddr := config.VideoInteractions.RPCAddr
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		log.Fatalf("VideoInteractions: resolve tcp addr failed, err: %v", err)
	}
	p := base.TelemetryProvider(config.Server.Name, config.Otel.CollectorAddr)
	defer func() {
		if err := p.Shutdown(context.Background()); err != nil {
			log.Fatalf("VideoInteractions: shutdown telemetry provider failed, err: %v", err)
		}
	}()

	svr := interactionservice.NewServer(
		// 注入依赖
		interaction.InjectInteractionHandler(),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "InteractionService",
		}),
		server.WithMuxTransport(),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithLimit(&limit.Option{
			MaxConnections: int(config.Server.MaxConnections),
			MaxQPS:         int(config.Server.MaxQPS),
		}),

		server.WithMiddleware(middleware.ErrorLog()),
		server.WithMiddleware(middleware.Respond()),
	)

	if err := svr.Run(); err != nil {
		log.Printf("VideoInteractions: run server failed, err: %v", err)
		return
	}
}
