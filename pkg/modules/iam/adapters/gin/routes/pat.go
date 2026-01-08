package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	handler "github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/adapters/gin/handler"
)

// PAT 返回个人访问令牌模块的所有路由
func PAT(patHandler *handler.PATHandler) []routes.Route {
	return []routes.Route{
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
			Path:        "/api/user/tokens/{id}",
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
			Path:        "/api/user/tokens/{id}",
			Handlers:    []gin.HandlerFunc{patHandler.DeleteToken},
			OperationID: "self:tokens:delete",
			Tags:        []string{"user-pat"},
			Summary:     "删除令牌",
			Description: "删除个人访问令牌",
		},
		{
			Method:      routes.PATCH,
			Path:        "/api/user/tokens/{id}/disable",
			Handlers:    []gin.HandlerFunc{patHandler.DisableToken},
			OperationID: "self:tokens:disable",
			Tags:        []string{"user-pat"},
			Summary:     "禁用令牌",
			Description: "禁用个人访问令牌",
		},
		{
			Method:      routes.PATCH,
			Path:        "/api/user/tokens/{id}/enable",
			Handlers:    []gin.HandlerFunc{patHandler.EnableToken},
			OperationID: "self:tokens:enable",
			Tags:        []string{"user-pat"},
			Summary:     "启用令牌",
			Description: "启用禁用的个人访问令牌",
		},
	}
}
