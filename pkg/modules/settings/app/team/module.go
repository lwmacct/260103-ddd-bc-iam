package team

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/infra/persistence"
	settingpersistence "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/infra/persistence"
)

// TeamUseCases 团队配置用例处理器聚合
type TeamUseCases struct {
	Set   *SetHandler
	Reset *ResetHandler
	Get   *GetHandler
	List  *ListHandler
}

// UseCaseModule 团队配置用例 Fx 模块
var UseCaseModule = fx.Module("settings.team.usecase",
	fx.Provide(newTeamUseCases),
)

type teamUseCasesParams struct {
	fx.In

	TeamRepos    persistence.TeamRepositories           // Team Settings BC 仓储
	OrgRepos     persistence.OrgRepositories            // Org Settings BC 仓储（用于验证团队所属组织）
	SettingRepos settingpersistence.SettingRepositories // 跨 BC 依赖：Settings BC
}

func newTeamUseCases(p teamUseCasesParams) *TeamUseCases {
	// 从 Settings BC 获取 QueryRepository
	settingQueryRepo := p.SettingRepos.Query
	categoryQueryRepo := p.SettingRepos.CategoryQuery

	return &TeamUseCases{
		Set:   NewSetHandler(settingQueryRepo, categoryQueryRepo, p.TeamRepos.Command),
		Reset: NewResetHandler(p.TeamRepos.Command),
		Get:   NewGetHandler(settingQueryRepo, categoryQueryRepo, p.TeamRepos.Query, p.OrgRepos.Query),
		List:  NewListHandler(settingQueryRepo, categoryQueryRepo, p.TeamRepos.Query, p.OrgRepos.Query),
	}
}
