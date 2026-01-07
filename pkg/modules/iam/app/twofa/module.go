package twofa

import (
	"go.uber.org/fx"
)

// TwoFAUseCases 双因素认证用例处理器聚合
type TwoFAUseCases struct {
	Setup        *SetupHandler
	VerifyEnable *VerifyEnableHandler
	Disable      *DisableHandler
	GetStatus    *GetStatusHandler
}

// Module 注册 TwoFA 子模块依赖
var Module = fx.Module("iam.twofa",
	fx.Provide(
		NewSetupHandler,
		NewVerifyEnableHandler,
		NewDisableHandler,
		NewGetStatusHandler,
		newTwoFAUseCases,
	),
)

type twofaUseCasesParams struct {
	fx.In

	Setup        *SetupHandler
	VerifyEnable *VerifyEnableHandler
	Disable      *DisableHandler
	GetStatus    *GetStatusHandler
}

func newTwoFAUseCases(p twofaUseCasesParams) *TwoFAUseCases {
	return &TwoFAUseCases{
		Setup:        p.Setup,
		VerifyEnable: p.VerifyEnable,
		Disable:      p.Disable,
		GetStatus:    p.GetStatus,
	}
}
