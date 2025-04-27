package mw

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// CORS 处理跨域请求的中间件
func CORS() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 设置允许的源，可以是单个域名，也可以是*表示允许所有域名
		ctx.Header("Access-Control-Allow-Origin", "*")
		// 设置允许的HTTP方法
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 设置允许的头部
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		// 设置是否允许携带凭证（如cookies）
		ctx.Header("Access-Control-Allow-Credentials", "true")
		// 设置预检请求的缓存时间（秒）
		ctx.Header("Access-Control-Max-Age", "86400")

		// 对于OPTIONS请求（预检请求），直接返回200
		if string(ctx.Method()) == "OPTIONS" {
			ctx.AbortWithStatus(consts.StatusOK)
			return
		}

		// 继续处理请求
		ctx.Next(c)
	}
}
