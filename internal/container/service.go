package container

import (
	"go.uber.org/fx"

	appauth "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/auth"
	appuser "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/user"
	config "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/config"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/auth"
	twofaDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/twofa"
	persistence "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infra/persistence"
	twofaInfra "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infra/twofa"
	infracaptcha "github.com/lwmacct/260103-ddd-shared/pkg/shared/captcha/infrastructure"

	authInfra "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infra/auth"
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
		authInfra.NewTokenGenerator,
		authInfra.NewLoginSessionService,
		newAuthPermissionCacheService,
		newPATService,
		infracaptcha.NewService,
		infracaptcha.NewMemoryRepository,
		newTwoFAService,

		// 领域服务
		newAuthService,
	),
)

func newJWTManager(iamCfg *config.Config) *authInfra.JWTManager {
	return authInfra.NewJWTManager(iamCfg.JWT.Secret, iamCfg.JWT.AccessTokenExpiry, iamCfg.JWT.RefreshTokenExpiry)
}

func newAuthPermissionCacheService(
	permissionCache appauth.PermissionCacheService,
	userWithRolesCache appuser.UserWithRolesCacheService,
	userRepos persistence.UserRepositories,
	roleRepos persistence.RoleRepositories,
) *authInfra.PermissionCacheService {
	return authInfra.NewPermissionCacheService(
		permissionCache,
		userWithRolesCache,
		userRepos.Query,
		roleRepos.Query,
	)
}

func newAuthService(jwt *authInfra.JWTManager, tokenGen *authInfra.TokenGenerator) auth.Service {
	passwordPolicy := auth.DefaultPasswordPolicy()
	return authInfra.NewAuthService(jwt, tokenGen, passwordPolicy)
}

func newPATService(
	patRepos persistence.PATRepositories,
	tokenGen *authInfra.TokenGenerator,
) *authInfra.PATService {
	return authInfra.NewPATService(patRepos.Command, patRepos.Query, tokenGen)
}

// twofaServiceParams 聚合 TwoFA 服务所需的依赖。
type twofaServiceParams struct {
	fx.In

	IAMConfig *config.Config
	TwoFA     persistence.TwoFARepositories
	UserRepos persistence.UserRepositories
}

func newTwoFAService(p twofaServiceParams) twofaDomain.Service {
	return twofaInfra.NewService(p.TwoFA.Command, p.TwoFA.Query, p.UserRepos.Query, p.IAMConfig.Auth.TwoFAIssuer)
}
