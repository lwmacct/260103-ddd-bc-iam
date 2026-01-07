package auth

import (
	"go.uber.org/fx"
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
	fx.Provide(
		NewLoginHandler,
		NewLogin2FAHandler,
		NewRegisterHandler,
		NewRefreshTokenHandler,
		newAuthUseCases,
	),
)

// authUseCasesParams 聚合构造参数
type authUseCasesParams struct {
	fx.In

	Login        *LoginHandler
	Login2FA     *Login2FAHandler
	Register     *RegisterHandler
	RefreshToken *RefreshTokenHandler
}

func newAuthUseCases(p authUseCasesParams) *AuthUseCases {
	return &AuthUseCases{
		Login:        p.Login,
		Login2FA:     p.Login2FA,
		Register:     p.Register,
		RefreshToken: p.RefreshToken,
	}
}
