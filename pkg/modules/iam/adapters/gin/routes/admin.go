package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	handler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/handler"
)

// Admin 管理员路由
func Admin(
	adminUserHandler *handler.AdminUserHandler,
	roleHandler *handler.RoleHandler,
	auditHandler *handler.AuditHandler,
	orgHandler *handler.OrgHandler,
) []routes.Route {
	var allRoutes []routes.Route

	// User management routes
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/admin/users",
			Handlers:    []gin.HandlerFunc{adminUserHandler.ListUsers},
			OperationID: "admin:users:list",
			Tags:        []string{"admin-user"},
			Summary:     "用户列表",
			Description: "分页获取用户列表",
		},
		{
			Method:      routes.GET,
			Path:        "/api/admin/users/{id}",
			Handlers:    []gin.HandlerFunc{adminUserHandler.GetUser},
			OperationID: "admin:users:get",
			Tags:        []string{"admin-user"},
			Summary:     "用户详情",
			Description: "获取用户详细信息",
		},
		{
			Method:      routes.POST,
			Path:        "/api/admin/users",
			Handlers:    []gin.HandlerFunc{adminUserHandler.CreateUser},
			OperationID: "admin:users:create",
			Tags:        []string{"admin-user"},
			Summary:     "创建用户",
			Description: "创建新用户",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/admin/users/{id}",
			Handlers:    []gin.HandlerFunc{adminUserHandler.UpdateUser},
			OperationID: "admin:users:update",
			Tags:        []string{"admin-user"},
			Summary:     "更新用户",
			Description: "更新用户信息",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/admin/users/{id}",
			Handlers:    []gin.HandlerFunc{adminUserHandler.DeleteUser},
			OperationID: "admin:users:delete",
			Tags:        []string{"admin-user"},
			Summary:     "删除用户",
			Description: "删除用户",
		},
		{
			Method:      routes.POST,
			Path:        "/api/admin/users/batch",
			Handlers:    []gin.HandlerFunc{adminUserHandler.BatchCreateUsers},
			OperationID: "admin:users:batch_create",
			Tags:        []string{"admin-user"},
			Summary:     "批量创建",
			Description: "批量创建用户",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/admin/users/{id}/roles",
			Handlers:    []gin.HandlerFunc{adminUserHandler.AssignRoles},
			OperationID: "admin:users:assign_roles",
			Tags:        []string{"admin-user"},
			Summary:     "分配角色",
			Description: "为用户分配角色",
		},
	}...)

	// Role management routes
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/admin/roles",
			Handlers:    []gin.HandlerFunc{roleHandler.ListRoles},
			OperationID: "admin:roles:list",
			Tags:        []string{"admin-role"},
			Summary:     "角色列表",
			Description: "分页获取角色列表",
		},
		{
			Method:      routes.GET,
			Path:        "/api/admin/roles/{id}",
			Handlers:    []gin.HandlerFunc{roleHandler.GetRole},
			OperationID: "admin:roles:get",
			Tags:        []string{"admin-role"},
			Summary:     "角色详情",
			Description: "获取角色详细信息",
		},
		{
			Method:      routes.POST,
			Path:        "/api/admin/roles",
			Handlers:    []gin.HandlerFunc{roleHandler.CreateRole},
			OperationID: "admin:roles:create",
			Tags:        []string{"admin-role"},
			Summary:     "创建角色",
			Description: "创建新角色",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/admin/roles/{id}",
			Handlers:    []gin.HandlerFunc{roleHandler.UpdateRole},
			OperationID: "admin:roles:update",
			Tags:        []string{"admin-role"},
			Summary:     "更新角色",
			Description: "更新角色信息",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/admin/roles/{id}",
			Handlers:    []gin.HandlerFunc{roleHandler.DeleteRole},
			OperationID: "admin:roles:delete",
			Tags:        []string{"admin-role"},
			Summary:     "删除角色",
			Description: "删除角色",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/admin/roles/{id}/permissions",
			Handlers:    []gin.HandlerFunc{roleHandler.SetPermissions},
			OperationID: "admin:roles:set_permissions",
			Tags:        []string{"admin-role"},
			Summary:     "设置权限",
			Description: "为角色设置权限",
		},
	}...)

	// Audit log routes
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/admin/audit",
			Handlers:    []gin.HandlerFunc{auditHandler.ListLogs},
			OperationID: "admin:audit:list",
			Tags:        []string{"admin-audit"},
			Summary:     "审计日志列表",
			Description: "获取审计日志列表",
		},
		{
			Method:      routes.GET,
			Path:        "/api/admin/audit/{id}",
			Handlers:    []gin.HandlerFunc{auditHandler.GetLog},
			OperationID: "admin:audit:get",
			Tags:        []string{"admin-audit"},
			Summary:     "审计日志详情",
			Description: "获取审计日志详情",
		},
	}...)

	// Organization management routes (admin level)
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

	return allRoutes
}

// UserOrg 用户组织视图路由
func UserOrg(userOrgHandler *handler.UserOrgHandler) []routes.Route {
	return []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/user/orgs",
			Handlers:    []gin.HandlerFunc{userOrgHandler.ListMyOrganizations},
			OperationID: "self:orgs:list",
			Tags:        []string{"user-org"},
			Summary:     "我的组织",
			Description: "获取用户所属组织列表",
		},
		{
			Method:      routes.GET,
			Path:        "/api/user/teams",
			Handlers:    []gin.HandlerFunc{userOrgHandler.ListMyTeams},
			OperationID: "self:teams:list",
			Tags:        []string{"user-org"},
			Summary:     "我的团队",
			Description: "获取用户所属团队列表",
		},
	}
}
