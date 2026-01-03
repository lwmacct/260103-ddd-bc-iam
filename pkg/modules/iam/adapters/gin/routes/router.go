// Package routes 定义 IAM 模块的所有 HTTP 路由。
//
// 本包遵循 DDD 架构原则，BC 层只负责定义路由结构和 Handler，
// 中间件由应用层（internal/container）注入。
//
// # 路由分组
//
//   - [Public]: 公开路由（登录、注册、验证码）
//   - [Auth]: 认证路由（双因素认证）
//   - [Self]: 用户自服务路由（个人资料、PAT）
//   - [Admin]: 管理员路由（用户管理、角色管理、审计日志、组织管理）
//   - [Org]: 组织管理路由（成员、团队）
//   - [UserOrg]: 用户组织视图路由
//
// # 使用方式
//
//	routes := routes.All(authHandler, twoFAHandler, ...)
//	injector := di.NewMiddlewareInjector(deps)
//	di.RegisterRoutesV2(engine, routes, injector)
package routes

import (
	"github.com/lwmacct/260101-go-pkg-gin/pkg/routes"

	handler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/handler"
)

// All 返回 IAM 域的所有路由
//
// 参数：各模块的 handler（按需传递，避免依赖 god object）
func All(
	// Auth handlers
	authHandler *handler.AuthHandler,
	twoFAHandler *handler.TwoFAHandler,
	captchaHandler *handler.CaptchaHandler,

	// User handlers
	userProfileHandler *handler.UserProfileHandler,
	userOrgHandler *handler.UserOrgHandler,
	userSettingHandler *handler.UserSettingHandler,

	// Admin handlers
	adminUserHandler *handler.AdminUserHandler,
	roleHandler *handler.RoleHandler,
	patHandler *handler.PATHandler,
	auditHandler *handler.AuditHandler,
	orgHandler *handler.OrgHandler,

	// Org handlers (organization/team management)
	orgMemberHandler *handler.OrgMemberHandler,
	teamHandler *handler.TeamHandler,
	teamMemberHandler *handler.TeamMemberHandler,
) []routes.Route {
	var allRoutes []routes.Route

	// Public routes (auth)
	allRoutes = append(allRoutes, Public(
		authHandler,
		captchaHandler,
	)...)

	// Auth routes (2FA)
	allRoutes = append(allRoutes, Auth(twoFAHandler)...)

	// Self routes (user profile, PAT, UserSettings)
	allRoutes = append(allRoutes, Self(
		userProfileHandler,
		patHandler,
		userSettingHandler,
	)...)

	// Admin routes (user/role/audit/org management)
	allRoutes = append(allRoutes, Admin(
		adminUserHandler,
		roleHandler,
		auditHandler,
		orgHandler,
	)...)

	// Org routes (organization/team management)
	allRoutes = append(allRoutes, Org(
		orgMemberHandler,
		teamHandler,
		teamMemberHandler,
	)...)

	// Org routes (user's org view)
	allRoutes = append(allRoutes, UserOrg(userOrgHandler)...)

	return allRoutes
}
