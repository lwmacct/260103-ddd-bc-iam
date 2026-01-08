package routes

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/adapters/gin/handler"
)

// RoutesModule 注册 IAM 模块的所有路由
var RoutesModule = fx.Module("iam.routes",
	fx.Provide(
		fx.Annotate(
			NewAllRoutes,
			fx.ResultTags(`name:"iam"`),
		),
	),
)

// NewAllRoutes 聚合所有 IAM 路由
// 直接注入 Handlers 聚合结构体，消除参数传递
func NewAllRoutes(h *handler.Handlers) []routes.Route {
	var allRoutes []routes.Route

	// Auth 模块
	allRoutes = append(allRoutes, Auth(h.Auth)...)
	allRoutes = append(allRoutes, Captcha(h.Captcha)...)
	allRoutes = append(allRoutes, TwoFA(h.TwoFA)...)

	// User 模块（复合路由）
	allRoutes = append(allRoutes, User(h.UserProfile, h.AdminUser, h.UserOrg)...)

	// PAT 模块
	allRoutes = append(allRoutes, PAT(h.PAT)...)

	// Role 模块
	allRoutes = append(allRoutes, Role(h.Role)...)

	// Audit 模块
	allRoutes = append(allRoutes, Audit(h.Audit)...)

	// Org 模块（复合路由：Org + OrgMember + Team + TeamMember）
	allRoutes = append(allRoutes, Org(h.Org, h.OrgMember, h.Team, h.TeamMember)...)

	return allRoutes
}
