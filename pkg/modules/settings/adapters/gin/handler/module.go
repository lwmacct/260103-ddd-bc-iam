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

// Handlers 聚合 Settings BC 模块的所有 HTTP 处理器。
//
// 设计说明：使用聚合结构体而非 fx.Out 导出单独 Handler，
// 与 IAM 模块保持风格一致，维护成本更低。
type Handlers struct {
	UserSetting *UserSettingHandler
	OrgSetting  *OrgSettingHandler
	TeamSetting *TeamSettingHandler
}

// NewAllHandlers 创建所有 Handler
func NewAllHandlers(p HandlerParams) *Handlers {
	return &Handlers{
		UserSetting: NewUserSettingHandler(p.UseCases),
		OrgSetting:  NewOrgSettingHandler(p.UseCases.OrgSet, p.UseCases.OrgReset, p.UseCases.OrgGet, p.UseCases.OrgList),
		TeamSetting: NewTeamSettingHandler(p.UseCases.TeamSet, p.UseCases.TeamReset, p.UseCases.TeamGet, p.UseCases.TeamList),
	}
}
