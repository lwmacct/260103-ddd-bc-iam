package container

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/config"
	appauth "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/application/auth"
	appuser "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/application/user"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/auth"
	domain_twofa "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/twofa"
	iampersistence "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infrastructure/persistence"
	infra_twofa "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infrastructure/twofa"
	infracaptcha "github.com/lwmacct/260103-ddd-bc-iam/pkg/shared/captcha/infrastructure"

	infra_auth "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infrastructure/auth"
)

// ServiceModule 提供所有领域服务和基础设施服务。
//
// 服务处理业务逻辑和技术关注点：
//   - Auth: 密码哈希、JWT Token 生成
//   - PermissionCache: 用户权限缓存（Cache-Aside 模式）
//   - TwoFA: 基于 TOTP 的双因素认证
//   - Captcha: 图形验证码生成和存储
var ServiceModule = fx.Module("service",
	fx.Provide(
		// 基础设施服务
		newJWTManager,
		infra_auth.NewTokenGenerator,
		infra_auth.NewLoginSessionService,
		newAuthPermissionCacheService,
		newPATService,
		infracaptcha.NewService,
		infracaptcha.NewMemoryRepository,
		newTwoFAService,

		// 领域服务
		newAuthService,
	),
)

func newJWTManager(cfg *config.Config) *infra_auth.JWTManager {
	return infra_auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry)
}

func newAuthPermissionCacheService(
	permissionCache appauth.PermissionCacheService,
	userWithRolesCache appuser.UserWithRolesCacheService,
	userRepos iampersistence.UserRepositories,
	roleRepos iampersistence.RoleRepositories,
) *infra_auth.PermissionCacheService {
	return infra_auth.NewPermissionCacheService(
		permissionCache,
		userWithRolesCache,
		userRepos.Query,
		roleRepos.Query,
	)
}

func newAuthService(jwt *infra_auth.JWTManager, tokenGen *infra_auth.TokenGenerator) auth.Service {
	passwordPolicy := auth.DefaultPasswordPolicy()
	return infra_auth.NewAuthService(jwt, tokenGen, passwordPolicy)
}

func newPATService(
	patRepos iampersistence.PATRepositories,
	tokenGen *infra_auth.TokenGenerator,
) *infra_auth.PATService {
	return infra_auth.NewPATService(patRepos.Command, patRepos.Query, tokenGen)
}

// twofaServiceParams 聚合 TwoFA 服务所需的依赖。
type twofaServiceParams struct {
	fx.In

	Config    *config.Config
	TwoFA     iampersistence.TwoFARepositories
	UserRepos iampersistence.UserRepositories
}

func newTwoFAService(p twofaServiceParams) domain_twofa.Service {
	return infra_twofa.NewService(p.TwoFA.Command, p.TwoFA.Query, p.UserRepos.Query, p.Config.Auth.TwoFAIssuer)
}
