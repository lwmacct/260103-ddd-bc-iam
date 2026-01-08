package user

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/infra/persistence"
	settingpersistence "github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/infra/persistence"
)

// UserUseCases 用户配置用例处理器聚合
type UserUseCases struct {
	Set            *SetHandler
	BatchSet       *BatchSetHandler
	Reset          *ResetHandler
	ResetAll       *ResetAllHandler
	Get            *GetHandler
	List           *ListHandler
	ListCategories *ListCategoriesHandler
}

// Module 用户配置用例 Fx 模块
var Module = fx.Module("settings.user",
	fx.Provide(newUserUseCases),
)

type userUseCasesParams struct {
	fx.In

	Repos        persistence.Repositories               // User Settings BC 仓储
	SettingRepos settingpersistence.SettingRepositories // 跨 BC 依赖：Settings BC
}

func newUserUseCases(p userUseCasesParams) *UserUseCases {
	// 从 Settings BC 获取 QueryRepository
	settingQueryRepo := p.SettingRepos.Query
	categoryQueryRepo := p.SettingRepos.CategoryQuery

	return &UserUseCases{
		Set:            NewSetHandler(settingQueryRepo, categoryQueryRepo, p.Repos.Command),
		BatchSet:       NewBatchSetHandler(settingQueryRepo, categoryQueryRepo, p.Repos.Command),
		Reset:          NewResetHandler(p.Repos.Command),
		ResetAll:       NewResetAllHandler(p.Repos.Command),
		Get:            NewGetHandler(settingQueryRepo, categoryQueryRepo, p.Repos.Query),
		List:           NewListHandler(settingQueryRepo, categoryQueryRepo, p.Repos.Query),
		ListCategories: NewListCategoriesHandler(categoryQueryRepo),
	}
}
