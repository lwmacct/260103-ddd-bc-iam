package captcha

import (
	"go.uber.org/fx"
)

// CaptchaUseCases 验证码用例处理器聚合
type CaptchaUseCases struct {
	Generate *GenerateHandler
}

// Module 注册 Captcha 子模块依赖
var Module = fx.Module("iam.captcha",
	fx.Provide(
		NewGenerateHandler,
		newCaptchaUseCases,
	),
)

type captchaUseCasesParams struct {
	fx.In

	Generate *GenerateHandler
}

func newCaptchaUseCases(p captchaUseCasesParams) *CaptchaUseCases {
	return &CaptchaUseCases{
		Generate: p.Generate,
	}
}
