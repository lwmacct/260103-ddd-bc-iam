// Package routes 定义 App 模块的 HTTP 路由。
//
// App 模块包含：健康检查、系统管理（任务、缓存、概览）
package routes

import (
	"github.com/lwmacct/260101-go-pkg-gin/pkg/routes"

	corehandler "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/transport/gin/handler"
)

// All 返回 App 域的所有路由
func All(
	healthHandler *corehandler.HealthHandler,
	settingHandler *corehandler.SettingHandler,
	userSettingHandler *corehandler.UserSettingHandler,
	cacheHandler *corehandler.CacheHandler,
	overviewHandler *corehandler.OverviewHandler,
) []routes.Route {
	var allRoutes []routes.Route

	// Public 路由（健康检查 + 用户配置）
	allRoutes = append(allRoutes, Public(
		healthHandler,
		userSettingHandler,
	)...)

	// Admin 路由（系统管理）
	allRoutes = append(allRoutes, Admin(
		settingHandler,
		cacheHandler,
		overviewHandler,
	)...)

	return allRoutes
}
