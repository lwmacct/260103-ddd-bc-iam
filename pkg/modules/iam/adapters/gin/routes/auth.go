package routes

import (
	"github.com/lwmacct/260101-go-pkg-gin/pkg/routes"

	handler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/handler"
)

// Auth
func Auth(twoFAHandler *handler.TwoFAHandler) []routes.Route {
	return []routes.Route{
		{
			Method:      routes.POST,
			Path:        "/api/auth/2fa/setup",
			Handler:     twoFAHandler.Setup,
			Operation:   "self:2fa:setup",
			Tags:        "Authentication - 2FA",
			Summary:     "设置 2FA",
			Description: "设置两步验证",
		},
		{
			Method:      routes.POST,
			Path:        "/api/auth/2fa/verify",
			Handler:     twoFAHandler.VerifyAndEnable,
			Operation:   "self:2fa:enable",
			Tags:        "Authentication - 2FA",
			Summary:     "启用 2FA",
			Description: "验证并启用两步验证",
		},
		{
			Method:      routes.POST,
			Path:        "/api/auth/2fa/disable",
			Handler:     twoFAHandler.Disable,
			Operation:   "self:2fa:disable",
			Tags:        "Authentication - 2FA",
			Summary:     "禁用 2FA",
			Description: "禁用两步验证",
		},
		{
			Method:      routes.GET,
			Path:        "/api/auth/2fa/status",
			Handler:     twoFAHandler.GetStatus,
			Operation:   "self:2fa:status",
			Tags:        "Authentication - 2FA",
			Summary:     "2FA 状态",
			Description: "获取两步验证状态",
		},
	}
}
