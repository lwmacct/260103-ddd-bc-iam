package auth

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/audit"
	authDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/auth"
	twofaDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/twofa"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infra/persistence"
	captchaInfra "github.com/lwmacct/260103-ddd-shared/pkg/shared/captcha"
)

// AuthUseCases 认证用例处理器聚合
type AuthUseCases struct {
	Login        *LoginHandler
	Login2FA     *Login2FAHandler
	Register     *RegisterHandler
	RefreshToken *RefreshTokenHandler
}

// Module 注册 Auth 子模块依赖
var Module = fx.Module("iam.auth",
	fx.Provide(newAuthUseCases),
)

type authUseCasesParams struct {
	fx.In

	UserRepos      persistence.UserRepositories
	CaptchaCommand captchaInfra.CommandRepository
	TwoFARepos     persistence.TwoFARepositories
	AuthSvc        authDomain.Service
	LoginSession   authDomain.SessionService
	TwoFASvc       twofaDomain.Service
	Audit          *audit.AuditUseCases
}

func newAuthUseCases(p authUseCasesParams) *AuthUseCases {
	return &AuthUseCases{
		Login:        NewLoginHandler(p.UserRepos.Query, p.CaptchaCommand, p.TwoFARepos.Query, p.AuthSvc, p.LoginSession, p.Audit.CreateLog),
		Login2FA:     NewLogin2FAHandler(p.UserRepos.Query, p.AuthSvc, p.LoginSession, p.TwoFASvc, p.Audit.CreateLog),
		Register:     NewRegisterHandler(p.UserRepos.Command, p.UserRepos.Query, p.AuthSvc),
		RefreshToken: NewRefreshTokenHandler(p.UserRepos.Query, p.AuthSvc, p.Audit.CreateLog),
	}
}
