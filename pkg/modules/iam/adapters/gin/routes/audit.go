package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	handler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/handler"
)

// Audit 返回审计日志模块的所有路由
func Audit(auditHandler *handler.AuditHandler) []routes.Route {
	return []routes.Route{
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
			Path:        "/api/admin/audit/actions",
			Handlers:    []gin.HandlerFunc{auditHandler.GetActions},
			OperationID: "admin:audit:actions",
			Tags:        []string{"admin-audit"},
			Summary:     "审计操作定义",
			Description: "获取所有审计操作的定义、分类和操作类型，供前端筛选器使用",
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
	}
}
