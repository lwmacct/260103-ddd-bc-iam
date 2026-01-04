package app

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/app/org"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/app/team"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/app/user"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/infra/persistence"
	settingpersistence "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/infra/persistence"
)

// UseCaseModule 用例 Fx 模块
var UseCaseModule = fx.Module("settings.usecase",
	user.UseCaseModule,
	org.UseCaseModule,
	team.UseCaseModule,
	fx.Provide(newUseCases),
)

// useCasesParams Fx 注入参数
type useCasesParams struct {
	fx.In

	Repos        persistence.Repositories               // User Settings BC 仓储
	OrgRepos     persistence.OrgRepositories            // Org Settings BC 仓储
	TeamRepos    persistence.TeamRepositories           // Team Settings BC 仓储
	SettingRepos settingpersistence.SettingRepositories // 跨 BC 依赖：Settings BC
}

// UseCases 聚合所有用例处理器
type UseCases struct {
	// User Settings
	Set            *user.SetHandler
	BatchSet       *user.BatchSetHandler
	Reset          *user.ResetHandler
	ResetAll       *user.ResetAllHandler
	Get            *user.GetHandler
	List           *user.ListHandler
	ListCategories *user.ListCategoriesHandler

	// Org Settings
	OrgSet   *org.SetHandler
	OrgReset *org.ResetHandler
	OrgGet   *org.GetHandler
	OrgList  *org.ListHandler

	// Team Settings
	TeamSet   *team.SetHandler
	TeamReset *team.ResetHandler
	TeamGet   *team.GetHandler
	TeamList  *team.ListHandler
}

// newUseCases 创建用例聚合实例
func newUseCases(p useCasesParams) *UseCases {
	// 从 Settings BC 获取 QueryRepository
	settingQueryRepo := p.SettingRepos.Query
	categoryQueryRepo := p.SettingRepos.CategoryQuery

	return &UseCases{
		// User Settings
		Set:            user.NewSetHandler(settingQueryRepo, p.Repos.Command),
		BatchSet:       user.NewBatchSetHandler(settingQueryRepo, p.Repos.Command),
		Reset:          user.NewResetHandler(p.Repos.Command),
		ResetAll:       user.NewResetAllHandler(p.Repos.Command),
		Get:            user.NewGetHandler(settingQueryRepo, p.Repos.Query),
		List:           user.NewListHandler(settingQueryRepo, p.Repos.Query),
		ListCategories: user.NewListCategoriesHandler(categoryQueryRepo),

		// Org Settings
		OrgSet:   org.NewSetHandler(settingQueryRepo, p.OrgRepos.Command),
		OrgReset: org.NewResetHandler(p.OrgRepos.Command),
		OrgGet:   org.NewGetHandler(settingQueryRepo, p.OrgRepos.Query),
		OrgList:  org.NewListHandler(settingQueryRepo, p.OrgRepos.Query),

		// Team Settings
		TeamSet:   team.NewSetHandler(settingQueryRepo, p.TeamRepos.Command),
		TeamReset: team.NewResetHandler(p.TeamRepos.Command),
		TeamGet:   team.NewGetHandler(settingQueryRepo, p.TeamRepos.Query, p.OrgRepos.Query),
		TeamList:  team.NewListHandler(settingQueryRepo, p.TeamRepos.Query, p.OrgRepos.Query),
	}
}
