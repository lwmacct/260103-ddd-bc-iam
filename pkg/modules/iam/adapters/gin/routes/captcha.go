package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	handler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/handler"
)

// Captcha 返回验证码模块的所有路由
func Captcha(captchaHandler *handler.CaptchaHandler) []routes.Route {
	return []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/auth/captcha",
			Handlers:    []gin.HandlerFunc{captchaHandler.GetCaptcha},
			OperationID: "public:auth:captcha",
			Tags:        []string{"auth"},
			Summary:     "获取验证码",
			Description: "生成图形验证码用于登录",
			Public:      true,
		},
	}
}
