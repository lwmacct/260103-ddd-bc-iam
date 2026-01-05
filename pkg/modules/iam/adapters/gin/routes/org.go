package routes

import (
	"github.com/gin-gonic/gin"
	handler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/handler"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"
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
			Handlers:    []gin.HandlerFunc{orgMemberHandler.List},
			OperationID: "org:members:list",
			Tags:        []string{"org-member"},
			Summary:     "成员列表",
			Description: "获取组织成员列表",
		},
		{
			Method:      routes.POST,
			Path:        "/api/org/:org_id/members",
			Handlers:    []gin.HandlerFunc{orgMemberHandler.Add},
			OperationID: "org:members:add",
			Tags:        []string{"org-member"},
			Summary:     "添加成员",
			Description: "添加组织成员",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/org/:org_id/members/:user_id",
			Handlers:    []gin.HandlerFunc{orgMemberHandler.Remove},
			OperationID: "org:members:remove",
			Tags:        []string{"org-member"},
			Summary:     "移除成员",
			Description: "移除组织成员",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/org/:org_id/members/:user_id/role",
			Handlers:    []gin.HandlerFunc{orgMemberHandler.UpdateRole},
			OperationID: "org:members:update:role",
			Tags:        []string{"org-member"},
			Summary:     "更新成员角色",
			Description: "更新组织成员角色",
		},

		// ==================== Org 域 - 团队管理 ====================
		{
			Method:      routes.POST,
			Path:        "/api/org/:org_id/teams",
			Handlers:    []gin.HandlerFunc{teamHandler.Create},
			OperationID: "org:teams:create",
			Tags:        []string{"org-team"},
			Summary:     "创建团队",
			Description: "创建团队",
		},
		{
			Method:      routes.GET,
			Path:        "/api/org/:org_id/teams",
			Handlers:    []gin.HandlerFunc{teamHandler.List},
			OperationID: "org:teams:list",
			Tags:        []string{"org-team"},
			Summary:     "团队列表",
			Description: "获取团队列表",
		},
		{
			Method:      routes.GET,
			Path:        "/api/org/:org_id/teams/:team_id",
			Handlers:    []gin.HandlerFunc{teamHandler.Get},
			OperationID: "org:teams:get",
			Tags:        []string{"org-team"},
			Summary:     "团队详情",
			Description: "获取团队详情",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/org/:org_id/teams/:team_id",
			Handlers:    []gin.HandlerFunc{teamHandler.Update},
			OperationID: "org:teams:update",
			Tags:        []string{"org-team"},
			Summary:     "更新团队",
			Description: "更新团队信息",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/org/:org_id/teams/:team_id",
			Handlers:    []gin.HandlerFunc{teamHandler.Delete},
			OperationID: "org:teams:delete",
			Tags:        []string{"org-team"},
			Summary:     "删除团队",
			Description: "删除团队",
		},

		// ==================== Org 域 - 团队成员管理 ====================
		{
			Method:      routes.GET,
			Path:        "/api/org/:org_id/teams/:team_id/members",
			Handlers:    []gin.HandlerFunc{teamMemberHandler.List},
			OperationID: "org:team:members:list",
			Tags:        []string{"org-team"},
			Summary:     "团队成员列表",
			Description: "获取团队成员列表",
		},
		{
			Method:      routes.POST,
			Path:        "/api/org/:org_id/teams/:team_id/members",
			Handlers:    []gin.HandlerFunc{teamMemberHandler.Add},
			OperationID: "org:team:members:add",
			Tags:        []string{"org-team"},
			Summary:     "添加团队成员",
			Description: "添加团队成员",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/org/:org_id/teams/:team_id/members/:user_id",
			Handlers:    []gin.HandlerFunc{teamMemberHandler.Remove},
			OperationID: "org:team:members:remove",
			Tags:        []string{"org-team"},
			Summary:     "移除团队成员",
			Description: "移除团队成员",
		},
	}
}
