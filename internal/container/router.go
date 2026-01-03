package container

import (
	"github.com/gin-gonic/gin"
	ginroutes "github.com/lwmacct/260101-go-pkg-gin/pkg/routes"

	// IAM
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/handler"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/routes"

	// User Settings BC
	userSettingsHandler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/adapters/gin/handler"
	userSettingsRoutes "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/adapters/gin/routes"

	// Settings
	settingsHandler "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/adapters/gin/handler"
	settingsRoutes "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/adapters/gin/routes"
)

// AllRoutes 聚合所有模块的路由定义。
//
// 职责：
//  1. 从 IAM、User Settings 和 Settings 模块收集路由定义
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

	// User Settings BC Handler
	userSettingHandler *userSettingsHandler.UserSettingHandler,

	// Settings Handlers
	settingHandler *settingsHandler.SettingHandler,
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

	// User Settings BC 路由
	userSettingsRouteList := userSettingsRoutes.All(userSettingHandler)

	// Settings 路由
	settingsRouteList := settingsRoutes.Admin(settingHandler)

	// 合并所有路由
	allRoutes := iamRoutes
	allRoutes = append(allRoutes, userSettingsRouteList...)
	return append(allRoutes, settingsRouteList...)
}

// RegisterRoutes 注册路由到 Gin Engine。

// RegisterRoutes 注册路由到 Gin Engine。
//
// 架构特点：
//  1. BC 层定义路由元数据（Operation、Path、Method、Handler）
//  2. 应用层注入中间件（认证、鉴权、审计、日志等）
//  3. 完美解耦 - BC 不依赖具体中间件实现
func RegisterRoutes(engine *gin.Engine, allRoutes []ginroutes.Route, injector *MiddlewareInjector) {
	for _, route := range allRoutes {
		// 应用层注入中间件
		middlewares := injector.InjectMiddlewares(&route)

		// 添加 Handler 作为最后一个处理函数
		middlewares = append(middlewares, route.Handler)

		// 注册到 Gin Engine
		switch route.Method {
		case ginroutes.GET:
			engine.GET(route.Path, middlewares...)
		case ginroutes.POST:
			engine.POST(route.Path, middlewares...)
		case ginroutes.PUT:
			engine.PUT(route.Path, middlewares...)
		case ginroutes.DELETE:
			engine.DELETE(route.Path, middlewares...)
		case ginroutes.PATCH:
			engine.PATCH(route.Path, middlewares...)
		case ginroutes.HEAD:
			engine.HEAD(route.Path, middlewares...)
		case ginroutes.OPTIONS:
			engine.OPTIONS(route.Path, middlewares...)
		}
	}
}
