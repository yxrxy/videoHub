package middleware

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/yxrrxy/videoHub/kitex_gen/user"
	pkgcontext "github.com/yxrrxy/videoHub/pkg/context"
	"github.com/yxrrxy/videoHub/pkg/errno"
	"github.com/yxrrxy/videoHub/pkg/jwt"
)

// TODO：为什么不起作用，不管了
func Auth() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {

			// 尝试从请求中获取用户 ID
			switch r := req.(type) {
			case *user.UploadAvatarRequest:
				userID := r.UserId
				if userID > 0 {
					// 将用户 ID 存储到上下文
					ctx = pkgcontext.WithUserID(ctx, userID)
				}
			}

			// 验证是否成功设置
			//if userID, exists := pkgcontext.GetUserID(ctx); exists {
			//	fmt.Printf("RPC 中间件验证用户 ID: %d\n", userID)
			//} else {
			//	fmt.Println("RPC 中间件未设置用户 ID")
			//}

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

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			fmt.Printf("解析 token 失败: %v\n", err)
			c.JSON(401, errno.ErrInvalidToken)
			c.Abort()
			return
		}

		newCtx := pkgcontext.WithUserID(ctx, claims.UserID)
		c.Set("user_id", claims.UserID)
		c.Next(newCtx)
	}
}
