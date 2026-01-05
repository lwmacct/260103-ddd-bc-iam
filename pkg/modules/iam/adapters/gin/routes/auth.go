package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	handler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/handler"
)

// Auth
func Auth(twoFAHandler *handler.TwoFAHandler) []routes.Route {
	return []routes.Route{
		{
			Method:      routes.POST,
			Path:        "/api/auth/2fa/setup",
			Handlers:    []gin.HandlerFunc{twoFAHandler.Setup},
			OperationID: "self:2fa:setup",
			Tags:        []string{"auth-2fa"},
			Summary:     "设置 2FA",
			Description: "设置两步验证",
		},
		{
			Method:      routes.POST,
			Path:        "/api/auth/2fa/verify",
			Handlers:    []gin.HandlerFunc{twoFAHandler.VerifyAndEnable},
			OperationID: "self:2fa:enable",
			Tags:        []string{"auth-2fa"},
			Summary:     "启用 2FA",
			Description: "验证并启用两步验证",
		},
		{
			Method:      routes.POST,
			Path:        "/api/auth/2fa/disable",
			Handlers:    []gin.HandlerFunc{twoFAHandler.Disable},
			OperationID: "self:2fa:disable",
			Tags:        []string{"auth-2fa"},
			Summary:     "禁用 2FA",
			Description: "禁用两步验证",
		},
		{
			Method:      routes.GET,
			Path:        "/api/auth/2fa/status",
			Handlers:    []gin.HandlerFunc{twoFAHandler.GetStatus},
			OperationID: "self:2fa:status",
			Tags:        []string{"auth-2fa"},
			Summary:     "2FA 状态",
			Description: "获取两步验证状态",
		},
	}
}
