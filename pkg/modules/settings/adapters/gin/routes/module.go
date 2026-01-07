package routes

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/adapters/gin/handler"
)

// RoutesModule 注册 Settings BC 模块的所有路由
var RoutesModule = fx.Module("settings.routes",
	fx.Provide(
		fx.Annotate(
			NewAllRoutes,
			fx.ResultTags(`name:"settings"`),
		),
	),
)

// NewAllRoutes 聚合所有 Settings BC 路由
// 直接注入 Handlers 聚合结构体，消除参数传递
func NewAllRoutes(h *handler.Handlers) []routes.Route {
	userRoutes := AllUser(h.UserSetting)
	orgRoutes := OrgSettings(h.OrgSetting)
	teamRoutes := TeamSettings(h.TeamSetting)

	return append(userRoutes, append(orgRoutes, teamRoutes...)...)
}
