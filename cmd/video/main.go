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
	"github.com/yxrxy/videoHub/app/video"
	"github.com/yxrxy/videoHub/config"
	"github.com/yxrxy/videoHub/kitex_gen/video/videoservice"
	"github.com/yxrxy/videoHub/pkg/base"
	"github.com/yxrxy/videoHub/pkg/middleware"
)

func init() {
	config.Init("video")
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		log.Fatalf("Video: new etcd registry failed, err: %v", err)
	}
	listenAddr := config.Video.RPCAddr
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Video: resolve tcp addr failed, err: %v", err)
	}
	p := base.TelemetryProvider(config.Server.Name, config.Otel.CollectorAddr)
	defer func() {
		if err := p.Shutdown(context.Background()); err != nil {
			log.Fatalf("Video: shutdown telemetry provider failed, err: %v", err)
		}
	}()
	svr := videoservice.NewServer(
		// 注入依赖
		video.InjectVideoHandler(),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "VideoService",
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
		log.Fatalf("Video: run server failed, err: %v", err)
	}
}
