// Package routes 定义 Settings 模块的 Team Settings HTTP 路由。
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/adapters/gin/handler"
)

// TeamSettings 返回团队配置的所有路由
func TeamSettings(h *handler.TeamSettingHandler) []routes.Route {
	return []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/org/{org_id}/teams/{team_id}/settings",
			Handlers:    []gin.HandlerFunc{h.List},
			OperationID: "org:team:settings:list",
			Tags:        []string{"team-setting"},
			Summary:     "团队配置列表",
			Description: "获取当前团队的配置列表（支持三级继承：团队>组织>系统默认值）",
		},
		{
			Method:      routes.GET,
			Path:        "/api/org/{org_id}/teams/{team_id}/settings/{key}",
			Handlers:    []gin.HandlerFunc{h.Get},
			OperationID: "org:team:settings:get",
			Tags:        []string{"team-setting"},
			Summary:     "获取团队配置",
			Description: "获取指定配置项的值（支持三级继承：团队>组织>系统默认值）",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/org/{org_id}/teams/{team_id}/settings/{key}",
			Handlers:    []gin.HandlerFunc{h.Set},
			OperationID: "org:team:settings:update",
			Tags:        []string{"team-setting"},
			Summary:     "设置团队配置",
			Description: "设置指定配置项的值（团队自定义覆盖）",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/org/{org_id}/teams/{team_id}/settings/{key}",
			Handlers:    []gin.HandlerFunc{h.Reset},
			OperationID: "org:team:settings:reset",
			Tags:        []string{"team-setting"},
			Summary:     "重置团队配置",
			Description: "重置指定配置项（删除团队自定义值，恢复组织配置或系统默认值）",
		},
	}
}
