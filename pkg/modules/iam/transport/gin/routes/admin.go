package routes

import (
	"github.com/lwmacct/260101-go-pkg-gin/pkg/routes"

	iamhandler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/transport/gin/handler"
)

// Admin 管理员路由
func Admin(
	adminUserHandler *iamhandler.AdminUserHandler,
	roleHandler *iamhandler.RoleHandler,
	auditHandler *iamhandler.AuditHandler,
	orgHandler *iamhandler.OrgHandler,
) []routes.Route {
	var allRoutes []routes.Route

	// User management routes
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/admin/users",
			Handler:     adminUserHandler.ListUsers,
			Operation:   "admin:users:list",
			Tags:        "Admin - User Management",
			Summary:     "用户列表",
			Description: "分页获取用户列表",
		},
		{
			Method:      routes.GET,
			Path:        "/api/admin/users/:id",
			Handler:     adminUserHandler.GetUser,
			Operation:   "admin:users:get",
			Tags:        "Admin - User Management",
			Summary:     "用户详情",
			Description: "获取用户详细信息",
		},
		{
			Method:      routes.POST,
			Path:        "/api/admin/users",
			Handler:     adminUserHandler.CreateUser,
			Operation:   "admin:users:create",
			Tags:        "Admin - User Management",
			Summary:     "创建用户",
			Description: "创建新用户",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/admin/users/:id",
			Handler:     adminUserHandler.UpdateUser,
			Operation:   "admin:users:update",
			Tags:        "Admin - User Management",
			Summary:     "更新用户",
			Description: "更新用户信息",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/admin/users/:id",
			Handler:     adminUserHandler.DeleteUser,
			Operation:   "admin:users:delete",
			Tags:        "Admin - User Management",
			Summary:     "删除用户",
			Description: "删除用户",
		},
		{
			Method:      routes.POST,
			Path:        "/api/admin/users/batch",
			Handler:     adminUserHandler.BatchCreateUsers,
			Operation:   "admin:users:batch_create",
			Tags:        "Admin - User Management",
			Summary:     "批量创建",
			Description: "批量创建用户",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/admin/users/:id/roles",
			Handler:     adminUserHandler.AssignRoles,
			Operation:   "admin:users:assign_roles",
			Tags:        "Admin - User Management",
			Summary:     "分配角色",
			Description: "为用户分配角色",
		},
	}...)

	// Role management routes
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/admin/roles",
			Handler:     roleHandler.ListRoles,
			Operation:   "admin:roles:list",
			Tags:        "Admin - Role Management",
			Summary:     "角色列表",
			Description: "分页获取角色列表",
		},
		{
			Method:      routes.GET,
			Path:        "/api/admin/roles/:id",
			Handler:     roleHandler.GetRole,
			Operation:   "admin:roles:get",
			Tags:        "Admin - Role Management",
			Summary:     "角色详情",
			Description: "获取角色详细信息",
		},
		{
			Method:      routes.POST,
			Path:        "/api/admin/roles",
			Handler:     roleHandler.CreateRole,
			Operation:   "admin:roles:create",
			Tags:        "Admin - Role Management",
			Summary:     "创建角色",
			Description: "创建新角色",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/admin/roles/:id",
			Handler:     roleHandler.UpdateRole,
			Operation:   "admin:roles:update",
			Tags:        "Admin - Role Management",
			Summary:     "更新角色",
			Description: "更新角色信息",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/admin/roles/:id",
			Handler:     roleHandler.DeleteRole,
			Operation:   "admin:roles:delete",
			Tags:        "Admin - Role Management",
			Summary:     "删除角色",
			Description: "删除角色",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/admin/roles/:id/permissions",
			Handler:     roleHandler.SetPermissions,
			Operation:   "admin:roles:set_permissions",
			Tags:        "Admin - Role Management",
			Summary:     "设置权限",
			Description: "为角色设置权限",
		},
	}...)

	// Audit log routes
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/admin/audit",
			Handler:     auditHandler.ListLogs,
			Operation:   "admin:audit:list",
			Tags:        "Admin - Audit",
			Summary:     "审计日志列表",
			Description: "获取审计日志列表",
		},
		{
			Method:      routes.GET,
			Path:        "/api/admin/audit/:id",
			Handler:     auditHandler.GetLog,
			Operation:   "admin:audit:get",
			Tags:        "Admin - Audit",
			Summary:     "审计日志详情",
			Description: "获取审计日志详情",
		},
	}...)

	// Organization management routes (admin level)
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/admin/orgs",
			Handler:     orgHandler.List,
			Operation:   "admin:orgs:list",
			Tags:        "Admin - Organizations",
			Summary:     "组织列表",
			Description: "获取组织列表",
		},
		{
			Method:      routes.POST,
			Path:        "/api/admin/orgs",
			Handler:     orgHandler.Create,
			Operation:   "admin:orgs:create",
			Tags:        "Admin - Organizations",
			Summary:     "创建组织",
			Description: "创建新组织",
		},
		{
			Method:      routes.GET,
			Path:        "/api/admin/orgs/:id",
			Handler:     orgHandler.Get,
			Operation:   "admin:orgs:get",
			Tags:        "Admin - Organizations",
			Summary:     "组织详情",
			Description: "获取组织详细信息",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/admin/orgs/:id",
			Handler:     orgHandler.Update,
			Operation:   "admin:orgs:update",
			Tags:        "Admin - Organizations",
			Summary:     "更新组织",
			Description: "更新组织信息",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/admin/orgs/:id",
			Handler:     orgHandler.Delete,
			Operation:   "admin:orgs:delete",
			Tags:        "Admin - Organizations",
			Summary:     "删除组织",
			Description: "删除组织",
		},
	}...)

	return allRoutes
}

// UserOrg 用户组织视图路由
func UserOrg(userOrgHandler *iamhandler.UserOrgHandler) []routes.Route {
	return []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/user/orgs",
			Handler:     userOrgHandler.ListMyOrganizations,
			Operation:   "self:orgs:list",
			Tags:        "User - Organizations",
			Summary:     "我的组织",
			Description: "获取用户所属组织列表",
		},
		{
			Method:      routes.GET,
			Path:        "/api/user/teams",
			Handler:     userOrgHandler.ListMyTeams,
			Operation:   "self:teams:list",
			Tags:        "User - Organizations",
			Summary:     "我的团队",
			Description: "获取用户所属团队列表",
		},
	}
}
