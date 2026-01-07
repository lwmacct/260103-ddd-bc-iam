package org

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/infra/persistence"
	settingpersistence "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/infra/persistence"
)

// OrgUseCases 组织配置用例处理器聚合
type OrgUseCases struct {
	Set   *SetHandler
	Reset *ResetHandler
	Get   *GetHandler
	List  *ListHandler
}

// UseCaseModule 组织配置用例 Fx 模块
var UseCaseModule = fx.Module("settings.org.usecase",
	fx.Provide(newOrgUseCases),
)

type orgUseCasesParams struct {
	fx.In

	OrgRepos     persistence.OrgRepositories            // Org Settings BC 仓储
	SettingRepos settingpersistence.SettingRepositories // 跨 BC 依赖：Settings BC
}

func newOrgUseCases(p orgUseCasesParams) *OrgUseCases {
	// 从 Settings BC 获取 QueryRepository
	settingQueryRepo := p.SettingRepos.Query
	categoryQueryRepo := p.SettingRepos.CategoryQuery

	return &OrgUseCases{
		Set:   NewSetHandler(settingQueryRepo, categoryQueryRepo, p.OrgRepos.Command),
		Reset: NewResetHandler(p.OrgRepos.Command),
		Get:   NewGetHandler(settingQueryRepo, categoryQueryRepo, p.OrgRepos.Query),
		List:  NewListHandler(settingQueryRepo, categoryQueryRepo, p.OrgRepos.Query),
	}
}
