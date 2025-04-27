package router

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/yxrxy/videoHub/app/gateway/handler/api/social"
)

// RegisterWebSocket registers all websocket routes
func RegisterWebSocket(r *server.Hertz) {
	// 使用标准路由注册方法
	r.GET("/api/v1/social/ws/connect", func(ctx context.Context, c *app.RequestContext) {
		social.ConnectWebSocket(ctx, c)
	})
}
