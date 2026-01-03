package routes

import (
	corehandler "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/transport/gin/handler"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/routes"
)

// Admin App 域管理员路由（配置、缓存管理、系统概览）
func Admin(
	settingHandler *corehandler.SettingHandler,
	cacheHandler *corehandler.CacheHandler,
	overviewHandler *corehandler.OverviewHandler,
) []routes.Route {
	var allRoutes []routes.Route

	// ==================== 配置分类 ====================
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:    routes.GET,
			Path:      "/api/admin/settings/categories",
			Handler:   settingHandler.GetCategories,
			Operation: "admin:setting:categories:list",
			Tags:      "Admin - Settings",
			Summary:   "配置分类列表",
		},
		{
			Method:    routes.GET,
			Path:      "/api/admin/settings/categories/:id",
			Handler:   settingHandler.GetCategory,
			Operation: "admin:setting:categories:get",
			Tags:      "Admin - Settings",
			Summary:   "配置分类详情",
		},
		{
			Method:    routes.POST,
			Path:      "/api/admin/settings/categories",
			Handler:   settingHandler.CreateCategory,
			Operation: "admin:setting:categories:create",
			Tags:      "Admin - Settings",
			Summary:   "创建配置分类",
		},
		{
			Method:    routes.PUT,
			Path:      "/api/admin/settings/categories/:id",
			Handler:   settingHandler.UpdateCategory,
			Operation: "admin:setting:categories:update",
			Tags:      "Admin - Settings",
			Summary:   "更新配置分类",
		},
		{
			Method:    routes.DELETE,
			Path:      "/api/admin/settings/categories/:id",
			Handler:   settingHandler.DeleteCategory,
			Operation: "admin:setting:categories:delete",
			Tags:      "Admin - Settings",
			Summary:   "删除配置分类",
		},
	}...)

	// ==================== 系统配置 ====================
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:    routes.POST,
			Path:      "/api/admin/settings/batch",
			Handler:   settingHandler.BatchUpdateSettings,
			Operation: "admin:settings:batch:update",
			Tags:      "Admin - Settings",
			Summary:   "批量更新配置",
		},
		{
			Method:    routes.POST,
			Path:      "/api/admin/settings",
			Handler:   settingHandler.CreateSetting,
			Operation: "admin:settings:create",
			Tags:      "Admin - Settings",
			Summary:   "创建配置",
		},
		{
			Method:    routes.GET,
			Path:      "/api/admin/settings",
			Handler:   settingHandler.GetSettings,
			Operation: "admin:settings:list",
			Tags:      "Admin - Settings",
			Summary:   "配置列表",
		},
		{
			Method:    routes.GET,
			Path:      "/api/admin/settings/:key",
			Handler:   settingHandler.GetSetting,
			Operation: "admin:settings:get",
			Tags:      "Admin - Settings",
			Summary:   "配置详情",
		},
		{
			Method:    routes.PUT,
			Path:      "/api/admin/settings/:key",
			Handler:   settingHandler.UpdateSetting,
			Operation: "admin:settings:update",
			Tags:      "Admin - Settings",
			Summary:   "更新配置",
		},
		{
			Method:    routes.DELETE,
			Path:      "/api/admin/settings/:key",
			Handler:   settingHandler.DeleteSetting,
			Operation: "admin:settings:delete",
			Tags:      "Admin - Settings",
			Summary:   "删除配置",
		},
	}...)

	// Cache management routes
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/admin/cache/info",
			Handler:     cacheHandler.Info,
			Operation:   "admin:cache:info",
			Tags:        "Admin - Cache",
			Summary:     "缓存信息",
			Description: "获取 Redis 缓存信息",
		},
		{
			Method:      routes.GET,
			Path:        "/api/admin/cache/keys",
			Handler:     cacheHandler.ScanKeys,
			Operation:   "admin:cache:keys",
			Tags:        "Admin - Cache",
			Summary:     "扫描缓存键",
			Description: "扫描 Redis 缓存键",
		},
		{
			Method:      routes.GET,
			Path:        "/api/admin/cache/key",
			Handler:     cacheHandler.GetKey,
			Operation:   "admin:cache:key:get",
			Tags:        "Admin - Cache",
			Summary:     "获取缓存键值",
			Description: "获取指定 Redis 键的值",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/admin/cache/key",
			Handler:     cacheHandler.DeleteKey,
			Operation:   "admin:cache:key:delete",
			Tags:        "Admin - Cache",
			Summary:     "删除缓存键",
			Description: "删除指定 Redis 键",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/admin/cache/keys",
			Handler:     cacheHandler.DeleteByPattern,
			Operation:   "admin:cache:keys:delete",
			Tags:        "Admin - Cache",
			Summary:     "批量删除缓存键",
			Description: "按模式批量删除 Redis 键",
		},
	}...)

	// Overview routes
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/admin/overview/stats",
			Handler:     overviewHandler.GetStats,
			Operation:   "admin:overview:stats",
			Tags:        "Admin - Overview",
			Summary:     "系统统计",
			Description: "获取系统统计数据",
		},
	}...)

	return allRoutes
}
