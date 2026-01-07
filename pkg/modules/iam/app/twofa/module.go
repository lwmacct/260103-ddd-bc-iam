package twofa

import (
	"go.uber.org/fx"

	twofaDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/twofa"
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
	fx.Provide(newTwoFAUseCases),
)

type twofaUseCasesParams struct {
	fx.In

	TwoFASvc twofaDomain.Service
}

func newTwoFAUseCases(p twofaUseCasesParams) *TwoFAUseCases {
	return &TwoFAUseCases{
		Setup:        NewSetupHandler(p.TwoFASvc),
		VerifyEnable: NewVerifyEnableHandler(p.TwoFASvc),
		Disable:      NewDisableHandler(p.TwoFASvc),
		GetStatus:    NewGetStatusHandler(p.TwoFASvc),
	}
}
