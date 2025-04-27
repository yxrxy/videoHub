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
	"github.com/yxrxy/videoHub/app/social"
	"github.com/yxrxy/videoHub/config"
	"github.com/yxrxy/videoHub/kitex_gen/social/socialservice"
	"github.com/yxrxy/videoHub/pkg/base"
	"github.com/yxrxy/videoHub/pkg/middleware"
)

func init() {
	config.Init("social")
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		log.Fatalf("Social: new etcd registry failed, err: %v", err)
	}
	listenAddr := config.Social.RPCAddr
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Social: resolve tcp addr failed, err: %v", err)
	}
	p := base.TelemetryProvider(config.Server.Name, config.Otel.CollectorAddr)
	defer func() {
		if err := p.Shutdown(context.Background()); err != nil {
			log.Fatalf("Social: shutdown telemetry provider failed, err: %v", err)
		}
	}()

	svr := socialservice.NewServer(
		// 注入依赖
		social.InjectSocialHandler(),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "SocialService",
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
	if err = svr.Run(); err != nil {
		log.Printf("Social: run server failed, err: %v", err)
		return
	}
}
