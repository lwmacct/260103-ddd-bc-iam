package handler

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/app"
)

// HandlerModule Handler Fx 模块
var HandlerModule = fx.Module("settings.handler",
	fx.Provide(
		NewAllHandlers,
	),
)

// HandlerParams Handler 构造参数（供外部 Fx 注入使用）
type HandlerParams struct {
	fx.In

	UseCases *app.UseCases
}

// HandlerResult Handler 导出结果
type HandlerResult struct {
	fx.Out

	UserSetting *UserSettingHandler
	OrgSetting  *OrgSettingHandler
	TeamSetting *TeamSettingHandler
}

// NewAllHandlers 创建所有 Handler
func NewAllHandlers(p HandlerParams) HandlerResult {
	return HandlerResult{
		UserSetting: NewUserSettingHandler(p.UseCases),
		OrgSetting:  NewOrgSettingHandler(p.UseCases.OrgSet, p.UseCases.OrgReset, p.UseCases.OrgGet, p.UseCases.OrgList),
		TeamSetting: NewTeamSettingHandler(p.UseCases.TeamSet, p.UseCases.TeamReset, p.UseCases.TeamGet, p.UseCases.TeamList),
	}
}
