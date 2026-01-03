package routes

import (
	"github.com/lwmacct/260101-go-pkg-gin/pkg/routes"

	corehandler "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/transport/gin/handler"
)

// Public 公开路由（健康检查 + 用户配置）
func Public(
	healthHandler *corehandler.HealthHandler,
	userSettingHandler *corehandler.UserSettingHandler,
) []routes.Route {
	var allRoutes []routes.Route

	// Health check
	allRoutes = append(allRoutes, routes.Route{
		Method:      routes.GET,
		Path:        "/health",
		Handler:     healthHandler.Check,
		Operation:   "public:health:check",
		Tags:        "System",
		Summary:     "健康检查",
		Description: "系统健康状态检查",
	})

	// ==================== 用户配置 ====================
	// 注意：categories 和 batch 路由必须在 :key 路由之前
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:    routes.GET,
			Path:      "/api/user/settings/categories",
			Handler:   userSettingHandler.ListUserSettingCategories,
			Operation: "self:settings:categories:list",
			Tags:      "User - Settings",
			Summary:   "配置分类列表",
		},
		{
			Method:    routes.POST,
			Path:      "/api/user/settings/batch",
			Handler:   userSettingHandler.BatchSetUserSettings,
			Operation: "self:settings:batch-set",
			Tags:      "User - Settings",
			Summary:   "批量设置配置",
		},
		{
			Method:    routes.GET,
			Path:      "/api/user/settings",
			Handler:   userSettingHandler.GetUserSettings,
			Operation: "self:settings:list",
			Tags:      "User - Settings",
			Summary:   "配置列表",
		},
		{
			Method:    routes.GET,
			Path:      "/api/user/settings/:key",
			Handler:   userSettingHandler.GetUserSetting,
			Operation: "self:settings:get",
			Tags:      "User - Settings",
			Summary:   "配置详情",
		},
		{
			Method:    routes.PUT,
			Path:      "/api/user/settings/:key",
			Handler:   userSettingHandler.SetUserSetting,
			Operation: "self:settings:set",
			Tags:      "User - Settings",
			Summary:   "设置配置",
		},
		{
			Method:    routes.DELETE,
			Path:      "/api/user/settings/:key",
			Handler:   userSettingHandler.ResetUserSetting,
			Operation: "self:settings:reset",
			Tags:      "User - Settings",
			Summary:   "重置配置",
		},
	}...)

	return allRoutes
}
