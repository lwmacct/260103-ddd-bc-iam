// Package routes 定义 Settings 模块的 Org Settings HTTP 路由。
package routes

import (
	"github.com/lwmacct/260101-go-pkg-gin/pkg/routes"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/adapters/gin/handler"
)

// OrgSettings 返回组织配置的所有路由
func OrgSettings(h *handler.OrgSettingHandler) []routes.Route {
	return []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/org/:org_id/settings",
			Handler:     h.List,
			Operation:   "org:settings:list",
			Tags:        "Org - Settings",
			Summary:     "组织配置列表",
			Description: "获取当前组织的配置列表（系统默认值+组织自定义值合并视图）",
		},
		{
			Method:      routes.GET,
			Path:        "/api/org/:org_id/settings/:key",
			Handler:     h.Get,
			Operation:   "org:settings:get",
			Tags:        "Org - Settings",
			Summary:     "获取组织配置",
			Description: "获取指定配置项的值（系统默认值或组织自定义值）",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/org/:org_id/settings/:key",
			Handler:     h.Set,
			Operation:   "org:settings:update",
			Tags:        "Org - Settings",
			Summary:     "设置组织配置",
			Description: "设置指定配置项的值（组织自定义覆盖）",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/org/:org_id/settings/:key",
			Handler:     h.Reset,
			Operation:   "org:settings:reset",
			Tags:        "Org - Settings",
			Summary:     "重置组织配置",
			Description: "重置指定配置项（删除组织自定义值，恢复系统默认值）",
		},
	}
}
