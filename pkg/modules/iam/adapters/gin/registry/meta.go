package registry

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/permission"
)

// Registry 路由注册表（用于查询函数）。
//
// 注意：Registry 在应用启动时由 BuildRegistryFromRoutes() 填充，
// 数据来源于 http.AllRoutes() 中的声明式路由定义。
//
//nolint:gochecknoglobals // 注册表是只读全局配置
var Registry = make(map[permission.Operation]routeMeta)

// HTTPMethod HTTP 请求方法。
type HTTPMethod string

// HTTP 方法常量。
const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
	PATCH  HTTPMethod = "PATCH"
)

// MiddlewareType 中间件类型标识符。
type MiddlewareType string

// 中间件类型常量。
const (
	MiddlewareRequestID   MiddlewareType = "request_id"
	MiddlewareOperationID MiddlewareType = "operation_id"
	MiddlewareAuth        MiddlewareType = "auth"
	MiddlewareOrgContext  MiddlewareType = "org_context"
	MiddlewareTeamContext MiddlewareType = "team_context"
	MiddlewareRBAC        MiddlewareType = "rbac"
	MiddlewareAudit       MiddlewareType = "audit"
	MiddlewareCORS        MiddlewareType = "cors"
	MiddlewareLogger      MiddlewareType = "logger"
)

// MiddlewareConfig 中间件配置。
type MiddlewareConfig struct {
	Name    MiddlewareType
	Options map[string]any // 中间件特定参数，如 {"optional": true}
}

// routeMeta 路由元数据（用于 Registry 查询）。
type routeMeta struct {
	// HTTP 路由
	Method HTTPMethod // HTTP 方法
	Path   string     // 路由路径（Gin 格式），如 /api/admin/users/:id

	// 中间件配置
	ReadOnly bool // 只读操作（对于团队操作，使用 TeamContextOptional 而非 TeamContext）

	// 审计配置
	Audit bool // 是否启用审计（审计详情从 Operation 派生）

	// Swagger 注解字段
	Tags        string // @Tags，如 Admin - Users
	Summary     string // @Summary，如 创建用户
	Description string // @Description（可选）
}

// Route 路由定义（整合元数据与处理器，用于路由注册）。
type Route struct {
	// HTTP 路由
	Method HTTPMethod // HTTP 方法
	Path   string     // 路由路径（Gin 格式），如 /api/admin/users/:id

	// 处理器与权限
	Handler gin.HandlerFunc      // 请求处理器
	Op      permission.Operation // 权限操作标识（URN 格式）

	// 声明式中间件配置
	Middlewares []MiddlewareConfig

	// Swagger 注解字段
	Tags        string // @Tags，如 Admin - Users
	Summary     string // @Summary，如 创建用户
	Description string // @Description（可选）
}

// ToMeta 将 Route 转换为 routeMeta（用于 Registry）。
func (r Route) ToMeta() routeMeta {
	// 从 Middlewares 配置推导 Audit
	hasAudit := false
	for _, mw := range r.Middlewares {
		if mw.Name == MiddlewareAudit {
			hasAudit = true
			break
		}
	}

	return routeMeta{
		Method:      r.Method,
		Path:        r.Path,
		ReadOnly:    false, // 不再使用，设为 false
		Audit:       hasAudit,
		Tags:        r.Tags,
		Summary:     r.Summary,
		Description: r.Description,
	}
}

// ==================== Route 辅助方法 ====================
// 注意：推导方法（NeedsOrgContext、NeedsTeamContext、IsReadOnly 等）已删除
// 因为中间件配置已显式声明在 Middlewares 字段中，无需推导。

// BuildRegistryFromRoutes 从声明式路由列表构建 Registry。
// 在应用启动时调用一次，之后 Registry 可用于所有查询函数。
func BuildRegistryFromRoutes(routes []Route) {
	for _, r := range routes {
		Registry[r.Op] = r.ToMeta()
	}
}
