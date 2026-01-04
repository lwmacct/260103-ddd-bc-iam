// Package routes 定义 Settings 模块的 Team Settings HTTP 路由。
package routes

import (
	"github.com/lwmacct/260101-go-pkg-gin/pkg/routes"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/adapters/gin/handler"
)

// TeamSettings 返回团队配置的所有路由
func TeamSettings(h *handler.TeamSettingHandler) []routes.Route {
	return []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/org/:org_id/teams/:team_id/settings",
			Handler:     h.List,
			Operation:   "org:team:settings:list",
			Tags:        "Team - Settings",
			Summary:     "团队配置列表",
			Description: "获取当前团队的配置列表（支持三级继承：团队>组织>系统默认值）",
		},
		{
			Method:      routes.GET,
			Path:        "/api/org/:org_id/teams/:team_id/settings/:key",
			Handler:     h.Get,
			Operation:   "org:team:settings:get",
			Tags:        "Team - Settings",
			Summary:     "获取团队配置",
			Description: "获取指定配置项的值（支持三级继承：团队>组织>系统默认值）",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/org/:org_id/teams/:team_id/settings/:key",
			Handler:     h.Set,
			Operation:   "org:team:settings:update",
			Tags:        "Team - Settings",
			Summary:     "设置团队配置",
			Description: "设置指定配置项的值（团队自定义覆盖）",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/org/:org_id/teams/:team_id/settings/:key",
			Handler:     h.Reset,
			Operation:   "org:team:settings:reset",
			Tags:        "Team - Settings",
			Summary:     "重置团队配置",
			Description: "重置指定配置项（删除团队自定义值，恢复组织配置或系统默认值）",
		},
	}
}
