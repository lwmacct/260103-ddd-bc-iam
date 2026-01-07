package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	handler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/handler"
)

// User 返回用户模块的所有路由
func User(
	userProfileHandler *handler.UserProfileHandler,
	adminUserHandler *handler.AdminUserHandler,
	userOrgHandler *handler.UserOrgHandler,
) []routes.Route {
	var allRoutes []routes.Route

	// ==================== 用户资料路由 ====================
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/user/profile",
			Handlers:    []gin.HandlerFunc{userProfileHandler.GetProfile},
			OperationID: "self:profile:get",
			Tags:        []string{"user-profile"},
			Summary:     "个人资料",
			Description: "获取当前用户资料",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/user/profile",
			Handlers:    []gin.HandlerFunc{userProfileHandler.UpdateProfile},
			OperationID: "self:profile:update",
			Tags:        []string{"user-profile"},
			Summary:     "更新资料",
			Description: "更新当前用户资料",
		},
		{
			Method:      routes.PUT,
			Path:        "/api/user/password",
			Handlers:    []gin.HandlerFunc{userProfileHandler.ChangePassword},
			OperationID: "self:password:change",
			Tags:        []string{"user-profile"},
			Summary:     "修改密码",
			Description: "修改当前用户密码",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/user/account",
			Handlers:    []gin.HandlerFunc{userProfileHandler.DeleteAccount},
			OperationID: "self:account:delete",
			Tags:        []string{"user-profile"},
			Summary:     "删除账户",
			Description: "删除当前用户账户",
		},
	}...)

	// ==================== 用户组织视图路由 ====================
	allRoutes = append(allRoutes, []routes.Route{
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
	}...)

	// ==================== 用户管理路由（管理员）====================
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

	return allRoutes
}
