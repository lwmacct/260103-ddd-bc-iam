// Package routes 定义 Settings 模块的 HTTP 路由。
package routes

import (
	"github.com/lwmacct/260101-go-pkg-gin/pkg/routes"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/adapters/gin/handler"
)

// All 返回 Settings 模块的所有路由
func All(h *handler.UserSettingHandler, orgH *handler.OrgSettingHandler, teamH *handler.TeamSettingHandler) []routes.Route {
	userRoutes := AllUser(h)
	orgRoutes := OrgSettings(orgH)
	teamRoutes := TeamSettings(teamH)

	return append(userRoutes, append(orgRoutes, teamRoutes...)...)
}

// AllUser 返回用户配置的所有路由
func AllUser(h *handler.UserSettingHandler) []routes.Route {
	return []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/user/settings",
			Handler:     h.List,
			Operation:   "self:settings:list",
			Tags:        "User - Settings",
			Summary:     "配置列表",
			Description: "获取当前用户的配置列表（系统默认值+用户自定义值合并视图）",
		},
		{
			Method:      routes.GET,
			Path:        "/api/user/settings/categories",
			Handler:     h.ListCategories,
			Operation:   "self:settings:categories",
			Tags:        "User - Settings",
			Summary:     "分类列表",
			Description: "获取配置分类列表",
		},
		{
			Method:      routes.GET,
			Path:        "/api/user/settings/:key",
			Handler:     h.Get,
			Operation:   "self:settings:get",
			Tags:        "User - Settings",
			Summary:     "获取配置",
			Description: "获取指定配置项的值",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/user/settings/:key",
			Handler:     h.Set,
			Operation:   "self:settings:update",
			Tags:        "User - Settings",
			Summary:     "设置配置",
			Description: "设置指定配置项的值（用户自定义覆盖）",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/user/settings/:key",
			Handler:     h.Reset,
			Operation:   "self:settings:delete",
			Tags:        "User - Settings",
			Summary:     "重置配置",
			Description: "重置指定配置项（恢复系统默认值）",
		},
		{
			Method:      routes.POST,
			Path:        "/api/user/settings/batch",
			Handler:     h.BatchSet,
			Operation:   "self:settings:batch",
			Tags:        "User - Settings",
			Summary:     "批量设置",
			Description: "批量设置多个配置项",
		},
		{
			Method:      routes.POST,
			Path:        "/api/user/settings/reset-all",
			Handler:     h.ResetAll,
			Operation:   "self:settings:reset-all",
			Tags:        "User - Settings",
			Summary:     "重置所有",
			Description: "重置所有用户自定义配置",
		},
	}
}
