package container

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/internal/bootstrap"
	"github.com/lwmacct/260103-ddd-bc-iam/internal/config"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/health"

	// Application UseCases (only for middleware dependencies)
	iamapplication "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/application"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infrastructure/auth"
	iampersistence "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infrastructure/persistence"

	// Handlers (injected via fx.In from their modules)
	iamhandler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/transport/gin/handler"

	ginHttp "github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin"
)

// HTTPModule 提供 HTTP 路由和服务器。
//
// 注意：所有 HTTP 处理器由各自的业务模块 (app/iam/crm) 提供，
// 本模块只负责路由注册和服务器生命周期管理。
var HTTPModule = fx.Module("http",
	fx.Provide(
		health.NewSystemChecker,
		newRouter,
		newHTTPServer,
	),
	fx.Invoke(startHTTPServer),
)

// newHTTPServer 创建 HTTP 服务器实例。
func newHTTPServer(router *gin.Engine, cfg *config.Config) *ginHttp.Server {
	return ginHttp.NewServer(router, cfg.Server.Addr)
}

// startHTTPServer 注册 HTTP 服务器启动和关闭钩子。
func startHTTPServer(lc fx.Lifecycle, server *ginHttp.Server, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.Info("Starting HTTP server", "addr", cfg.Server.Addr, "env", cfg.Server.Env)

			// 在 goroutine 中启动服务器，避免阻塞 OnStart
			go func() {
				if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					slog.Error("HTTP server error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("Shutting down HTTP server")
			return server.Shutdown(ctx)
		},
	})
}

// routerParams 聚合创建路由所需的依赖。
type routerParams struct {
	fx.In

	Config      *config.Config
	RedisClient *redis.Client

	// Services
	JWTManager      *auth.JWTManager
	PATService      *auth.PATService
	PermissionCache *auth.PermissionCacheService

	// UseCases
	Audit *iamapplication.AuditUseCases

	// Repositories (for middleware)
	MemberRepos     iampersistence.OrgMemberRepositories
	TeamRepos       iampersistence.TeamRepositories
	TeamMemberRepos iampersistence.TeamMemberRepositories

	// Handlers
	Auth        *iamhandler.AuthHandler
	Captcha     *iamhandler.CaptchaHandler
	AdminUser   *iamhandler.AdminUserHandler
	UserProfile *iamhandler.UserProfileHandler
	Role        *iamhandler.RoleHandler
	PAT         *iamhandler.PATHandler
	AuditH      *iamhandler.AuditHandler
	TwoFA       *iamhandler.TwoFAHandler
	Org         *iamhandler.OrgHandler
	OrgMember   *iamhandler.OrgMemberHandler
	Team        *iamhandler.TeamHandler
	TeamMember  *iamhandler.TeamMemberHandler
	UserOrg     *iamhandler.UserOrgHandler
}

func newRouter(p routerParams) *gin.Engine {
	// Create Gin Engine using bootstrap
	engine := bootstrap.NewEngine()

	// Get all routes from modules using the new routes function
	allRoutes := AllRoutes(
		// IAM Handlers
		p.Auth,
		p.TwoFA,
		p.UserProfile,
		p.UserOrg,
		p.PAT,
		p.AdminUser,
		p.Role,
		p.Captcha,
		p.AuditH,
		p.Org,
		p.OrgMember,
		p.Team,
		p.TeamMember,
	)

	// Create MiddlewareInjector with all dependencies
	injector := NewMiddlewareInjector(RouterDepsParams{
		Config:             p.Config,
		RedisClient:        p.RedisClient,
		JWTManager:         p.JWTManager,
		PATService:         p.PATService,
		PermissionCache:    p.PermissionCache,
		AuditCreateHandler: p.Audit.CreateLog,
		MemberRepos:        p.MemberRepos,
		TeamRepos:          p.TeamRepos,
		TeamMemberRepos:    p.TeamMemberRepos,
	})

	// Register routes to engine with middleware injection
	RegisterRoutes(engine, allRoutes, injector)

	return engine
}
