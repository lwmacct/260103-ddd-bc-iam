package container

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-iam-bc/internal/config"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/audit"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/infra/auth"
	persistence "github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/infra/persistence"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/adapters/gin/middleware"
	ginmiddleware "github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/middleware"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/permission"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"
)

// RouterDepsParams 聚合创建中间件所需的依赖。
type RouterDepsParams struct {
	fx.In

	Config      *config.Config
	RedisClient *redis.Client

	// Infrastructure Services
	JWTManager      *auth.JWTManager
	PATService      *auth.PATService
	PermissionCache *auth.PermissionCacheService

	// Application Handlers (for middleware)
	AuditCreateHandler *audit.CreateHandler

	// Domain Repositories (for middleware)
	MemberRepos     persistence.OrgMemberRepositories
	TeamRepos       persistence.TeamRepositories
	TeamMemberRepos persistence.TeamMemberRepositories
}

// MiddlewareInjector 中间件注入器。
//
// 职责：根据路由的 Operation 和 Path 决定需要哪些中间件
//
// 优势：
//  1. BC 层完全解耦 - 不依赖具体中间件实现
//  2. 灵活配置 - 应用层可针对不同环境注入不同中间件
//  3. 支持可选依赖 - 不需要认证的应用可以不注入 Auth 中间件
type MiddlewareInjector struct {
	// 中间件工厂函数
	authMiddleware gin.HandlerFunc
	rbacFactory    func(operation string) gin.HandlerFunc
	auditFactory   func(operation string) gin.HandlerFunc
	orgContext     gin.HandlerFunc
	teamContext    gin.HandlerFunc
	teamContextOpt gin.HandlerFunc
}

// NewMiddlewareInjector 创建中间件注入器。
//
// 参数按需传入 - 不需要的功能可以传 nil
func NewMiddlewareInjector(p RouterDepsParams) *MiddlewareInjector {
	injector := &MiddlewareInjector{}

	// 认证中间件（可选 - 支持无认证的应用）
	if p.JWTManager != nil && p.PATService != nil && p.PermissionCache != nil {
		injector.authMiddleware = middleware.Auth(
			p.JWTManager,
			p.PATService,
			p.PermissionCache,
		)
	}

	// RBAC 中间件工厂
	injector.rbacFactory = func(operation string) gin.HandlerFunc {
		// 转换为 permission.Operation 类型
		return middleware.RequireOperation(permission.Operation(operation))
	}

	// Audit 中间件工厂（可选）
	if p.AuditCreateHandler != nil {
		injector.auditFactory = func(operation string) gin.HandlerFunc {
			return nil // TODO: 实现 AuditMiddleware
		}
	}

	// Org/Team 上下文中间件（可选 - 支持无多租户的应用）
	if p.MemberRepos.Query != nil && p.TeamRepos.Query != nil {
		injector.orgContext = middleware.OrgContext(p.MemberRepos.Query)
		injector.teamContext = middleware.TeamContext(p.TeamRepos.Query, p.TeamMemberRepos.Query)
	}

	return injector
}

// InjectMiddlewares 为路由注入中间件。
//
// 规则：
//  1. 所有路由：RequestID + OperationID + Logger
//  2. public:* - 无需认证
//  3. self:*, admin:*, org:* - 需要 Auth + RBAC
//  4. 路径包含 :org_id - 需要 OrgContext
//  5. 路径包含 :team_id - 需要 TeamContext
//  6. 审计操作 - 额外添加 Audit 中间件
func (inj *MiddlewareInjector) InjectMiddlewares(route *routes.Route) []gin.HandlerFunc {
	var m []gin.HandlerFunc

	// 1. 基础中间件（所有路由）
	m = append(m,
		ginmiddleware.RequestID(),
		ginmiddleware.SetOperationID(route.OperationID),
	)

	// 2. 认证中间件（非 public 路由）
	if !strings.HasPrefix(route.OperationID, "public:") {
		if inj.authMiddleware != nil {
			m = append(m, inj.authMiddleware)
		}
	}

	// 3. Org 上下文（路径包含 {org_id}）
	if strings.Contains(route.Path, "{org_id}") && inj.orgContext != nil {
		m = append(m, inj.orgContext)
	}

	// 4. Team 上下文（路径包含 {team_id}）
	if strings.Contains(route.Path, "{team_id}") {
		if inj.teamContextOpt != nil && route.Method == routes.GET {
			// GET 操作使用可选的 Team 上下文
			m = append(m, inj.teamContextOpt)
		} else if inj.teamContext != nil {
			m = append(m, inj.teamContext)
		}
	}

	// 5. RBAC 权限检查（非 public 路由）
	if !strings.HasPrefix(route.OperationID, "public:") {
		if inj.rbacFactory != nil {
			m = append(m, inj.rbacFactory(route.OperationID))
		}
	}

	// 6. 审计中间件（写操作 + 非 public）
	if route.Method != routes.GET &&
		!strings.HasPrefix(route.OperationID, "public:") &&
		inj.auditFactory != nil {
		m = append(m, inj.auditFactory(route.OperationID))
	}

	// 7. Logger 中间件（最后，记录完整请求）
	m = append(m, ginmiddleware.LoggerSkipPaths("/health"))

	return m
}
