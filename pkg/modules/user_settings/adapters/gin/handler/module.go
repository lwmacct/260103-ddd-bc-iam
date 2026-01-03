package handler

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/user_settings/app"
)

// HandlerModule Handler Fx 模块
var HandlerModule = fx.Module("user_settings.handler",
	fx.Provide(NewUserSettingHandler),
)

// HandlerParams Handler 构造参数（供外部 Fx 注入使用）
type HandlerParams struct {
	fx.In

	UseCases *app.UseCases
}
