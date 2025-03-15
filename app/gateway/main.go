package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/yxrrxy/videoHub/config"
)

func RegisterRoutes(h *server.Hertz) {
	h := server.New(server.WithHostPorts(config.GatewayConfig))
}
