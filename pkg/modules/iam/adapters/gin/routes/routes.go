// Package routes 定义 IAM 模块的所有 HTTP 路由。
//
// 本包遵循 DDD 架构原则，路由按业务模块组织。
// 中间件由应用层（internal/container）注入。
//
// # 路由组织
//
//   - [Auth]: 认证模块（注册、登录、令牌刷新）
//   - [Captcha]: 验证码模块
//   - [TwoFA]: 双因素认证模块
//   - [User]: 用户模块（用户资料、用户管理、用户组织视图）
//   - [PAT]: 个人访问令牌模块
//   - [Role]: 角色管理模块
//   - [Audit]: 审计日志模块
//   - [Org]: 组织模块（组织管理、组织成员、团队、团队成员）
//
// # 使用方式
//
//	routes := routes.All(authHandler, captchaHandler, ...)
//	injector := di.NewMiddlewareInjector(deps)
//	di.RegisterRoutesV2(engine, routes, injector)
package routes

import (
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	handler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/handler"
)

// All 返回 IAM 域的所有路由
//
// 参数：各模块的 handler（按模块传递）
func All(
	// Auth module
	authHandler *handler.AuthHandler,
	captchaHandler *handler.CaptchaHandler,
	twoFAHandler *handler.TwoFAHandler,

	// User module
	userProfileHandler *handler.UserProfileHandler,
	adminUserHandler *handler.AdminUserHandler,
	userOrgHandler *handler.UserOrgHandler,

	// PAT module
	patHandler *handler.PATHandler,

	// Role module
	roleHandler *handler.RoleHandler,

	// Audit module
	auditHandler *handler.AuditHandler,

	// Org module
	orgHandler *handler.OrgHandler,
	orgMemberHandler *handler.OrgMemberHandler,
	teamHandler *handler.TeamHandler,
	teamMemberHandler *handler.TeamMemberHandler,
) []routes.Route {
	var allRoutes []routes.Route

	// Auth module routes
	allRoutes = append(allRoutes, Auth(authHandler)...)
	allRoutes = append(allRoutes, Captcha(captchaHandler)...)
	allRoutes = append(allRoutes, TwoFA(twoFAHandler)...)

	// User module routes
	allRoutes = append(allRoutes, User(userProfileHandler, adminUserHandler, userOrgHandler)...)

	// PAT module routes
	allRoutes = append(allRoutes, PAT(patHandler)...)

	// Role module routes
	allRoutes = append(allRoutes, Role(roleHandler)...)

	// Audit module routes
	allRoutes = append(allRoutes, Audit(auditHandler)...)

	// Org module routes
	allRoutes = append(allRoutes, Org(orgHandler, orgMemberHandler, teamHandler, teamMemberHandler)...)

	return allRoutes
}
