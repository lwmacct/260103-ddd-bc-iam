// Package routes 定义 Settings 模块的 HTTP 路由。
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

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
			Handlers:    []gin.HandlerFunc{h.List},
			OperationID: "self:settings:list",
			Tags:        []string{"user-setting"},
			Summary:     "配置列表",
			Description: "获取当前用户的配置列表（系统默认值+用户自定义值合并视图）",
		},
		{
			Method:      routes.GET,
			Path:        "/api/user/settings/categories",
			Handlers:    []gin.HandlerFunc{h.ListCategories},
			OperationID: "self:settings:categories",
			Tags:        []string{"user-setting"},
			Summary:     "分类列表",
			Description: "获取配置分类列表",
		},
		{
			Method:      routes.GET,
			Path:        "/api/user/settings/:key",
			Handlers:    []gin.HandlerFunc{h.Get},
			OperationID: "self:settings:get",
			Tags:        []string{"user-setting"},
			Summary:     "获取配置",
			Description: "获取指定配置项的值",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/user/settings/:key",
			Handlers:    []gin.HandlerFunc{h.Set},
			OperationID: "self:settings:update",
			Tags:        []string{"user-setting"},
			Summary:     "设置配置",
			Description: "设置指定配置项的值（用户自定义覆盖）",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/user/settings/:key",
			Handlers:    []gin.HandlerFunc{h.Reset},
			OperationID: "self:settings:delete",
			Tags:        []string{"user-setting"},
			Summary:     "重置配置",
			Description: "重置指定配置项（恢复系统默认值）",
		},
		{
			Method:      routes.POST,
			Path:        "/api/user/settings/batch",
			Handlers:    []gin.HandlerFunc{h.BatchSet},
			OperationID: "self:settings:batch",
			Tags:        []string{"user-setting"},
			Summary:     "批量设置",
			Description: "批量设置多个配置项",
		},
		{
			Method:      routes.POST,
			Path:        "/api/user/settings/reset-all",
			Handlers:    []gin.HandlerFunc{h.ResetAll},
			OperationID: "self:settings:reset-all",
			Tags:        []string{"user-setting"},
			Summary:     "重置所有",
			Description: "重置所有用户自定义配置",
		},
	}
}
