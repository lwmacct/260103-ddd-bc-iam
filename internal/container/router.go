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
	// IAM Handlers（聚合）
	iam *handler.Handlers,

	// Settings BC Handlers（聚合）
	settingsBC *userSettingsHandler.Handlers,

	// Settings Handlers (external dependency)
	settingHandler *settingsHandler.SettingHandler,
	settingsCfg settingsconfig.Config,
) []ginroutes.Route {
	// IAM 域路由
	iamRoutes := routes.All(
		iam.Auth,
		iam.TwoFA,
		iam.Captcha,
		iam.UserProfile,
		iam.UserOrg,
		iam.AdminUser,
		iam.Role,
		iam.PAT,
		iam.Audit,
		iam.Org,
		iam.OrgMember,
		iam.Team,
		iam.TeamMember,
	)

	// Settings BC 路由（User + Org + Team）
	settingsRouteList := settingsRoutes.All(
		settingsBC.UserSetting,
		settingsBC.OrgSetting,
		settingsBC.TeamSetting,
	)

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

		// 转换 OpenAPI 风格路径参数 {param} 为 Gin 风格 :param
		ginPath := ginroutes.ToGinPath(route.Path)

		// 注册到 Gin Engine
		switch route.Method {
		case ginroutes.GET:
			engine.GET(ginPath, handlers...)
		case ginroutes.POST:
			engine.POST(ginPath, handlers...)
		case ginroutes.PUT:
			engine.PUT(ginPath, handlers...)
		case ginroutes.DELETE:
			engine.DELETE(ginPath, handlers...)
		case ginroutes.PATCH:
			engine.PATCH(ginPath, handlers...)
		case ginroutes.HEAD:
			engine.HEAD(ginPath, handlers...)
		case ginroutes.OPTIONS:
			engine.OPTIONS(ginPath, handlers...)
		}
	}
}
