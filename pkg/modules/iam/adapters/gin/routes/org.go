package routes

import (
	"github.com/lwmacct/260101-go-pkg-gin/pkg/routes"
	handler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/handler"
)

// Org 返回 Org 域路由
// 包含：Org 域成员管理、Org 域团队管理、Org 域团队成员管理
func Org(
	orgMemberHandler *handler.OrgMemberHandler,
	teamHandler *handler.TeamHandler,
	teamMemberHandler *handler.TeamMemberHandler,
) []routes.Route {
	return []routes.Route{
		// ==================== Org 域 - 成员管理 ====================
		{
			Method:      routes.GET,
			Path:        "/api/org/:org_id/members",
			Handler:     orgMemberHandler.List,
			Operation:   "org:members:list",
			Tags:        "Org - Members",
			Summary:     "成员列表",
			Description: "获取组织成员列表",
		},
		{
			Method:      routes.POST,
			Path:        "/api/org/:org_id/members",
			Handler:     orgMemberHandler.Add,
			Operation:   "org:members:add",
			Tags:        "Org - Members",
			Summary:     "添加成员",
			Description: "添加组织成员",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/org/:org_id/members/:user_id",
			Handler:     orgMemberHandler.Remove,
			Operation:   "org:members:remove",
			Tags:        "Org - Members",
			Summary:     "移除成员",
			Description: "移除组织成员",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/org/:org_id/members/:user_id/role",
			Handler:     orgMemberHandler.UpdateRole,
			Operation:   "org:members:update:role",
			Tags:        "Org - Members",
			Summary:     "更新成员角色",
			Description: "更新组织成员角色",
		},

		// ==================== Org 域 - 团队管理 ====================
		{
			Method:      routes.POST,
			Path:        "/api/org/:org_id/teams",
			Handler:     teamHandler.Create,
			Operation:   "org:teams:create",
			Tags:        "Org - Teams",
			Summary:     "创建团队",
			Description: "创建团队",
		},
		{
			Method:      routes.GET,
			Path:        "/api/org/:org_id/teams",
			Handler:     teamHandler.List,
			Operation:   "org:teams:list",
			Tags:        "Org - Teams",
			Summary:     "团队列表",
			Description: "获取团队列表",
		},
		{
			Method:      routes.GET,
			Path:        "/api/org/:org_id/teams/:team_id",
			Handler:     teamHandler.Get,
			Operation:   "org:teams:get",
			Tags:        "Org - Teams",
			Summary:     "团队详情",
			Description: "获取团队详情",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/org/:org_id/teams/:team_id",
			Handler:     teamHandler.Update,
			Operation:   "org:teams:update",
			Tags:        "Org - Teams",
			Summary:     "更新团队",
			Description: "更新团队信息",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/org/:org_id/teams/:team_id",
			Handler:     teamHandler.Delete,
			Operation:   "org:teams:delete",
			Tags:        "Org - Teams",
			Summary:     "删除团队",
			Description: "删除团队",
		},

		// ==================== Org 域 - 团队成员管理 ====================
		{
			Method:      routes.GET,
			Path:        "/api/org/:org_id/teams/:team_id/members",
			Handler:     teamMemberHandler.List,
			Operation:   "org:team:members:list",
			Tags:        "Org - Teams",
			Summary:     "团队成员列表",
			Description: "获取团队成员列表",
		},
		{
			Method:      routes.POST,
			Path:        "/api/org/:org_id/teams/:team_id/members",
			Handler:     teamMemberHandler.Add,
			Operation:   "org:team:members:add",
			Tags:        "Org - Teams",
			Summary:     "添加团队成员",
			Description: "添加团队成员",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/org/:org_id/teams/:team_id/members/:user_id",
			Handler:     teamMemberHandler.Remove,
			Operation:   "org:team:members:remove",
			Tags:        "Org - Teams",
			Summary:     "移除团队成员",
			Description: "移除团队成员",
		},
	}
}
