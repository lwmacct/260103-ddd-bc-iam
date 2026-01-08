package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	handler "github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/adapters/gin/handler"
)

// Role 返回角色管理模块的所有路由
func Role(roleHandler *handler.RoleHandler) []routes.Route {
	return []routes.Route{
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
	}
}
