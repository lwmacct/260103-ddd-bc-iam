package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	handler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/handler"
)

// Self 用户自服务路由（资料、令牌）
// 注意：User Settings 已迁移到独立的 User Settings BC
func Self(
	userProfileHandler *handler.UserProfileHandler,
	patHandler *handler.PATHandler,
) []routes.Route {
	var allRoutes []routes.Route

	// User Profile routes
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

	// PAT routes
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:      routes.GET,
			Path:        "/api/user/tokens",
			Handlers:    []gin.HandlerFunc{patHandler.ListTokens},
			OperationID: "self:tokens:list",
			Tags:        []string{"user-pat"},
			Summary:     "令牌列表",
			Description: "获取当前用户的 PAT 列表",
		},
		{
			Method:      routes.POST,
			Path:        "/api/user/tokens",
			Handlers:    []gin.HandlerFunc{patHandler.CreateToken},
			OperationID: "self:tokens:create",
			Tags:        []string{"user-pat"},
			Summary:     "创建令牌",
			Description: "创建个人访问令牌",
		},
		{
			Method:      routes.GET,
			Path:        "/api/user/tokens/:id",
			Handlers:    []gin.HandlerFunc{patHandler.GetToken},
			OperationID: "self:tokens:get",
			Tags:        []string{"user-pat"},
			Summary:     "令牌详情",
			Description: "获取个人访问令牌详情",
		},
		{
			Method:      routes.GET,
			Path:        "/api/user/tokens/scopes",
			Handlers:    []gin.HandlerFunc{patHandler.ListScopes},
			OperationID: "self:tokens:scopes",
			Tags:        []string{"user-pat"},
			Summary:     "令牌作用域",
			Description: "获取令牌作用域列表",
		},
		{
			Method:      routes.DELETE,
			Path:        "/api/user/tokens/:id",
			Handlers:    []gin.HandlerFunc{patHandler.DeleteToken},
			OperationID: "self:tokens:delete",
			Tags:        []string{"user-pat"},
			Summary:     "删除令牌",
			Description: "删除个人访问令牌",
		},
		{
			Method:      routes.PATCH,
			Path:        "/api/user/tokens/:id/disable",
			Handlers:    []gin.HandlerFunc{patHandler.DisableToken},
			OperationID: "self:tokens:disable",
			Tags:        []string{"user-pat"},
			Summary:     "禁用令牌",
			Description: "禁用个人访问令牌",
		},
		{
			Method:      routes.PATCH,
			Path:        "/api/user/tokens/:id/enable",
			Handlers:    []gin.HandlerFunc{patHandler.EnableToken},
			OperationID: "self:tokens:enable",
			Tags:        []string{"user-pat"},
			Summary:     "启用令牌",
			Description: "启用禁用的个人访问令牌",
		},
	}...)

	return allRoutes
}
