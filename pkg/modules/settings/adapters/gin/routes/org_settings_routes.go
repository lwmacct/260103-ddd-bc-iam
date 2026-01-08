// Package routes 定义 Settings 模块的 Org Settings HTTP 路由。
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/adapters/gin/handler"
)

// OrgSettings 返回组织配置的所有路由
func OrgSettings(h *handler.OrgSettingHandler) []routes.Route {
	return []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/org/{org_id}/settings",
			Handlers:    []gin.HandlerFunc{h.List},
			OperationID: "org:settings:list",
			Tags:        []string{"org-setting"},
			Summary:     "组织配置列表",
			Description: "获取当前组织的配置列表（系统默认值+组织自定义值合并视图）",
		},
		{
			Method:      routes.GET,
			Path:        "/api/org/{org_id}/settings/{key}",
			Handlers:    []gin.HandlerFunc{h.Get},
			OperationID: "org:settings:get",
			Tags:        []string{"org-setting"},
			Summary:     "获取组织配置",
			Description: "获取指定配置项的值（系统默认值或组织自定义值）",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/org/{org_id}/settings/{key}",
			Handlers:    []gin.HandlerFunc{h.Set},
			OperationID: "org:settings:update",
			Tags:        []string{"org-setting"},
			Summary:     "设置组织配置",
			Description: "设置指定配置项的值（组织自定义覆盖）",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/org/{org_id}/settings/{key}",
			Handlers:    []gin.HandlerFunc{h.Reset},
			OperationID: "org:settings:reset",
			Tags:        []string{"org-setting"},
			Summary:     "重置组织配置",
			Description: "重置指定配置项（删除组织自定义值，恢复系统默认值）",
		},
	}
}
