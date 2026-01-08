package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	handler "github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/adapters/gin/handler"
)

// Auth 返回认证模块的所有路由
func Auth(authHandler *handler.AuthHandler) []routes.Route {
	return []routes.Route{
		{
			Method:      routes.POST,
			Path:        "/api/auth/register",
			Handlers:    []gin.HandlerFunc{authHandler.Register},
			OperationID: "public:auth:register",
			Tags:        []string{"auth"},
			Summary:     "注册",
			Description: "用户注册",
			Public:      true,
		},
		{
			Method:      routes.POST,
			Path:        "/api/auth/login",
			Handlers:    []gin.HandlerFunc{authHandler.Login},
			OperationID: "public:auth:login",
			Tags:        []string{"auth"},
			Summary:     "登录",
			Description: "用户登录",
			Public:      true,
		},
		{
			Method:      routes.POST,
			Path:        "/api/auth/login/2fa",
			Handlers:    []gin.HandlerFunc{authHandler.Login2FA},
			OperationID: "public:auth:login2fa",
			Tags:        []string{"auth"},
			Summary:     "2FA 登录",
			Description: "两步验证登录",
			Public:      true,
		},
		{
			Method:      routes.POST,
			Path:        "/api/auth/refresh",
			Handlers:    []gin.HandlerFunc{authHandler.RefreshToken},
			OperationID: "public:auth:refresh",
			Tags:        []string{"auth"},
			Summary:     "刷新令牌",
			Description: "使用 refresh token 获取新的 access token",
			Public:      true,
		},
	}
}
