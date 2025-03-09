package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/endpoint"
	pkgcontext "github.com/yxrrxy/videoHub/pkg/context"
	"github.com/yxrrxy/videoHub/pkg/errno"
	"github.com/yxrrxy/videoHub/pkg/jwt"
)

// Auth 中间件用于验证token并存储用户ID (用于 RPC 服务)
func Auth() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			if c, ok := req.(*app.RequestContext); ok {
				token := c.GetHeader("Authorization")
				ctx = context.WithValue(ctx, "token", token)
			}

			if token, ok := ctx.Value("token").(string); ok {
				claims, err := jwt.ParseToken(token)
				if err == nil {
					ctx = pkgcontext.WithUserID(ctx, claims.UserID)
				}
			}
			return next(ctx, req, resp)
		}
	}
}

// JWT 中间件用于验证token并存储用户ID (用于 HTTP 服务)
func JWT() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := string(c.GetHeader("Authorization"))
		if token == "" {
			c.JSON(401, errno.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(401, errno.ErrInvalidToken)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next(ctx)
	}
}
