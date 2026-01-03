package routes

import (
	"github.com/lwmacct/260101-go-pkg-gin/pkg/routes"

	iamhandler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/transport/gin/handler"
)

// Public 公开路由（无需认证）
//
// 中间件由应用层注入：RequestID, OperationID, Logger
func Public(
	authHandler *iamhandler.AuthHandler,
	captchaHandler *iamhandler.CaptchaHandler,
) []routes.Route {
	return []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/auth/captcha",
			Handler:     captchaHandler.GetCaptcha,
			Operation:   "public:auth:captcha",
			Tags:        "Authentication",
			Summary:     "获取验证码",
			Description: "生成图形验证码用于登录",
		},
		{
			Method:      routes.POST,
			Path:        "/api/auth/register",
			Handler:     authHandler.Register,
			Operation:   "public:auth:register",
			Tags:        "Authentication",
			Summary:     "注册",
			Description: "用户注册",
		},
		{
			Method:      routes.POST,
			Path:        "/api/auth/login",
			Handler:     authHandler.Login,
			Operation:   "public:auth:login",
			Tags:        "Authentication",
			Summary:     "登录",
			Description: "用户登录",
		},
		{
			Method:      routes.POST,
			Path:        "/api/auth/login/2fa",
			Handler:     authHandler.Login2FA,
			Operation:   "public:auth:login2fa",
			Tags:        "Authentication",
			Summary:     "2FA 登录",
			Description: "两步验证登录",
		},
		{
			Method:      routes.POST,
			Path:        "/api/auth/refresh",
			Handler:     authHandler.RefreshToken,
			Operation:   "public:auth:refresh",
			Tags:        "Authentication",
			Summary:     "刷新令牌",
			Description: "使用 refresh token 获取新的 access token",
		},
	}
}
