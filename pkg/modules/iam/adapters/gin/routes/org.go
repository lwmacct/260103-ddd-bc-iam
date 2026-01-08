package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	handler "github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/adapters/gin/handler"
)

// Org 返回组织模块的所有路由
func Org(
	orgHandler *handler.OrgHandler,
	orgMemberHandler *handler.OrgMemberHandler,
	teamHandler *handler.TeamHandler,
	teamMemberHandler *handler.TeamMemberHandler,
) []routes.Route {
	var allRoutes []routes.Route

	// ==================== 组织管理（管理员）====================
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/admin/orgs",
			Handlers:    []gin.HandlerFunc{orgHandler.List},
			OperationID: "admin:orgs:list",
			Tags:        []string{"admin-org"},
			Summary:     "组织列表",
			Description: "获取组织列表",
		},
		{
			Method:      routes.POST,
			Path:        "/api/admin/orgs",
			Handlers:    []gin.HandlerFunc{orgHandler.Create},
			OperationID: "admin:orgs:create",
			Tags:        []string{"admin-org"},
			Summary:     "创建组织",
			Description: "创建新组织",
		},
		{
			Method:      routes.GET,
			Path:        "/api/admin/orgs/{id}",
			Handlers:    []gin.HandlerFunc{orgHandler.Get},
			OperationID: "admin:orgs:get",
			Tags:        []string{"admin-org"},
			Summary:     "组织详情",
			Description: "获取组织详细信息",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/admin/orgs/{id}",
			Handlers:    []gin.HandlerFunc{orgHandler.Update},
			OperationID: "admin:orgs:update",
			Tags:        []string{"admin-org"},
			Summary:     "更新组织",
			Description: "更新组织信息",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/admin/orgs/{id}",
			Handlers:    []gin.HandlerFunc{orgHandler.Delete},
			OperationID: "admin:orgs:delete",
			Tags:        []string{"admin-org"},
			Summary:     "删除组织",
			Description: "删除组织",
		},
	}...)

	// ==================== 组织成员管理 ====================
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/org/{org_id}/members",
			Handlers:    []gin.HandlerFunc{orgMemberHandler.List},
			OperationID: "org:members:list",
			Tags:        []string{"org-member"},
			Summary:     "成员列表",
			Description: "获取组织成员列表",
		},
		{
			Method:      routes.POST,
			Path:        "/api/org/{org_id}/members",
			Handlers:    []gin.HandlerFunc{orgMemberHandler.Add},
			OperationID: "org:members:add",
			Tags:        []string{"org-member"},
			Summary:     "添加成员",
			Description: "添加组织成员",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/org/{org_id}/members/{user_id}",
			Handlers:    []gin.HandlerFunc{orgMemberHandler.Remove},
			OperationID: "org:members:remove",
			Tags:        []string{"org-member"},
			Summary:     "移除成员",
			Description: "移除组织成员",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/org/{org_id}/members/{user_id}/role",
			Handlers:    []gin.HandlerFunc{orgMemberHandler.UpdateRole},
			OperationID: "org:member-role:update",
			Tags:        []string{"org-member"},
			Summary:     "更新成员角色",
			Description: "更新组织成员角色",
		},
	}...)

	// ==================== 团队管理 ====================
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:      routes.POST,
			Path:        "/api/org/{org_id}/teams",
			Handlers:    []gin.HandlerFunc{teamHandler.Create},
			OperationID: "org:teams:create",
			Tags:        []string{"org-team"},
			Summary:     "创建团队",
			Description: "创建团队",
		},
		{
			Method:      routes.GET,
			Path:        "/api/org/{org_id}/teams",
			Handlers:    []gin.HandlerFunc{teamHandler.List},
			OperationID: "org:teams:list",
			Tags:        []string{"org-team"},
			Summary:     "团队列表",
			Description: "获取团队列表",
		},
		{
			Method:      routes.GET,
			Path:        "/api/org/{org_id}/teams/{team_id}",
			Handlers:    []gin.HandlerFunc{teamHandler.Get},
			OperationID: "org:teams:get",
			Tags:        []string{"org-team"},
			Summary:     "团队详情",
			Description: "获取团队详情",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/org/{org_id}/teams/{team_id}",
			Handlers:    []gin.HandlerFunc{teamHandler.Update},
			OperationID: "org:teams:update",
			Tags:        []string{"org-team"},
			Summary:     "更新团队",
			Description: "更新团队信息",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/org/{org_id}/teams/{team_id}",
			Handlers:    []gin.HandlerFunc{teamHandler.Delete},
			OperationID: "org:teams:delete",
			Tags:        []string{"org-team"},
			Summary:     "删除团队",
			Description: "删除团队",
		},
	}...)

	// ==================== 团队成员管理 ====================
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/org/{org_id}/teams/{team_id}/members",
			Handlers:    []gin.HandlerFunc{teamMemberHandler.List},
			OperationID: "org:team-members:list",
			Tags:        []string{"org-team"},
			Summary:     "团队成员列表",
			Description: "获取团队成员列表",
		},
		{
			Method:      routes.POST,
			Path:        "/api/org/{org_id}/teams/{team_id}/members",
			Handlers:    []gin.HandlerFunc{teamMemberHandler.Add},
			OperationID: "org:team-members:add",
			Tags:        []string{"org-team"},
			Summary:     "添加团队成员",
			Description: "添加团队成员",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/org/{org_id}/teams/{team_id}/members/{user_id}",
			Handlers:    []gin.HandlerFunc{teamMemberHandler.Remove},
			OperationID: "org:team-members:remove",
			Tags:        []string{"org-team"},
			Summary:     "移除团队成员",
			Description: "移除团队成员",
		},
	}...)

	return allRoutes
}
