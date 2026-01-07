package captcha

import (
	"go.uber.org/fx"

	captchaInfra "github.com/lwmacct/260103-ddd-shared/pkg/shared/captcha"
)

// CaptchaUseCases 验证码用例处理器聚合
type CaptchaUseCases struct {
	Generate *GenerateHandler
}

// Module 注册 Captcha 子模块依赖
var Module = fx.Module("iam.captcha",
	fx.Provide(newCaptchaUseCases),
)

type captchaUseCasesParams struct {
	fx.In

	CaptchaCommand captchaInfra.CommandRepository
	CaptchaSvc     captchaInfra.Service
}

func newCaptchaUseCases(p captchaUseCasesParams) *CaptchaUseCases {
	return &CaptchaUseCases{
		Generate: NewGenerateHandler(p.CaptchaCommand, p.CaptchaSvc),
	}
}
