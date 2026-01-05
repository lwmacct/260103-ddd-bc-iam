package container

import (
	"github.com/gin-gonic/gin"
	ginroutes "github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	// IAM
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/handler"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/routes"

	// Settings BC
	userSettingsHandler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/adapters/gin/handler"
	settingsRoutes "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/adapters/gin/routes"

	// Settings (external dependency)
	settingsHandler "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/adapters/gin/handler"
	settingsBCRoutes "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/adapters/gin/routes"
	settingsconfig "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/config"
)

// AllRoutes 聚合所有模块的路由定义。
//
// 职责：
//  1. 从 IAM、Settings BC 模块收集路由定义
//  2. 返回统一的路由列表供注册使用
func AllRoutes(
	// IAM Handlers
	authHandler *handler.AuthHandler,
	twoFAHandler *handler.TwoFAHandler,
	userProfileHandler *handler.UserProfileHandler,
	userOrgHandler *handler.UserOrgHandler,
	patHandler *handler.PATHandler,
	adminUserHandler *handler.AdminUserHandler,
	roleHandler *handler.RoleHandler,
	captchaHandler *handler.CaptchaHandler,
	auditHandler *handler.AuditHandler,
	orgHandler *handler.OrgHandler,
	orgMemberHandler *handler.OrgMemberHandler,
	teamHandler *handler.TeamHandler,
	teamMemberHandler *handler.TeamMemberHandler,

	// Settings BC Handlers
	userSettingHandler *userSettingsHandler.UserSettingHandler,
	orgSettingHandler *userSettingsHandler.OrgSettingHandler,
	teamSettingHandler *userSettingsHandler.TeamSettingHandler,

	// Settings Handlers (external dependency)
	settingHandler *settingsHandler.SettingHandler,
	settingsCfg settingsconfig.Config,
) []ginroutes.Route {
	// IAM 域路由
	iamRoutes := routes.All(
		authHandler,
		twoFAHandler,
		captchaHandler,
		userProfileHandler,
		userOrgHandler,
		adminUserHandler,
		roleHandler,
		patHandler,
		auditHandler,
		orgHandler,
		orgMemberHandler,
		teamHandler,
		teamMemberHandler,
	)

	// Settings BC 路由（User + Org + Team）
	settingsRouteList := settingsRoutes.All(userSettingHandler, orgSettingHandler, teamSettingHandler)

	// Settings 路由 (external dependency)
	settingsBCRouteList := settingsBCRoutes.Admin(settingHandler, &settingsCfg)

	// 合并所有路由
	allRoutes := iamRoutes
	allRoutes = append(allRoutes, settingsRouteList...)
	return append(allRoutes, settingsBCRouteList...)
}

// RegisterRoutes 注册路由到 Gin Engine。
//
// 架构特点：
//  1. BC 层定义路由元数据（OperationID、Path、Method、Handlers）
//  2. 应用层注入中间件（认证、鉴权、审计、日志等）
//  3. 完美解耦 - BC 不依赖具体中间件实现
func RegisterRoutes(engine *gin.Engine, allRoutes []ginroutes.Route, injector *MiddlewareInjector) {
	for _, route := range allRoutes {
		// 应用层注入中间件
		middlewares := injector.InjectMiddlewares(&route)

		// 追加路由定义的 Handlers（BC 层定义的处理函数）
		handlers := middlewares
		handlers = append(handlers, route.Handlers...)

		// 注册到 Gin Engine
		switch route.Method {
		case ginroutes.GET:
			engine.GET(route.Path, handlers...)
		case ginroutes.POST:
			engine.POST(route.Path, handlers...)
		case ginroutes.PUT:
			engine.PUT(route.Path, handlers...)
		case ginroutes.DELETE:
			engine.DELETE(route.Path, handlers...)
		case ginroutes.PATCH:
			engine.PATCH(route.Path, handlers...)
		case ginroutes.HEAD:
			engine.HEAD(route.Path, handlers...)
		case ginroutes.OPTIONS:
			engine.OPTIONS(route.Path, handlers...)
		}
	}
}
