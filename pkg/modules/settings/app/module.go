package app

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/app/user"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/infra/persistence"
	settingpersistence "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/infra/persistence"
)

// UseCases 聚合所有用例处理器
type UseCases struct {
	Set            *user.SetHandler
	BatchSet       *user.BatchSetHandler
	Reset          *user.ResetHandler
	ResetAll       *user.ResetAllHandler
	Get            *user.GetHandler
	List           *user.ListHandler
	ListCategories *user.ListCategoriesHandler
}

// UseCaseModule 用例 Fx 模块
var UseCaseModule = fx.Module("user_settings.usecase",
	fx.Provide(newUseCases),
)

// useCasesParams Fx 注入参数
type useCasesParams struct {
	fx.In

	Repos        persistence.Repositories               // User Settings BC 仓储
	SettingRepos settingpersistence.SettingRepositories // 跨 BC 依赖：Settings BC
}

// newUseCases 创建用例聚合实例
func newUseCases(p useCasesParams) *UseCases {
	// 从 Settings BC 获取 QueryRepository
	settingQueryRepo := p.SettingRepos.Query
	categoryQueryRepo := p.SettingRepos.CategoryQuery

	return &UseCases{
		Set:            user.NewSetHandler(settingQueryRepo, p.Repos.Command),
		BatchSet:       user.NewBatchSetHandler(settingQueryRepo, p.Repos.Command),
		Reset:          user.NewResetHandler(p.Repos.Command),
		ResetAll:       user.NewResetAllHandler(p.Repos.Command),
		Get:            user.NewGetHandler(settingQueryRepo, p.Repos.Query),
		List:           user.NewListHandler(settingQueryRepo, p.Repos.Query),
		ListCategories: user.NewListCategoriesHandler(categoryQueryRepo),
	}
}
