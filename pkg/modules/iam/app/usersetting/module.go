package usersetting

import (
	"go.uber.org/fx"

	usersettingpersistence "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infra/persistence"
	settingpersistence "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/infra/persistence"
)

// UserSettingUseCases 用户设置用例处理器聚合
type UserSettingUseCases struct {
	Get    *GetHandler
	List   *ListHandler
	Update *UpdateHandler
	Delete *DeleteHandler
}

// UseCaseModule Fx 模块注册
var UseCaseModule = fx.Module("iam.usersetting.usecase",
	fx.Provide(newUserSettingUseCases),
)

// userSettingUseCasesParams 依赖参数
type userSettingUseCasesParams struct {
	fx.In

	UserSettingRepos usersettingpersistence.UserSettingRepositories
	SettingRepos     settingpersistence.SettingRepositories // 跨 BC 依赖：Settings QueryRepository
}

// newUserSettingUseCases 创建用户设置用例处理器
func newUserSettingUseCases(p userSettingUseCasesParams) *UserSettingUseCases {
	return &UserSettingUseCases{
		Get:    NewGetHandler(p.UserSettingRepos.Query, p.SettingRepos.Query),
		List:   NewListHandler(p.UserSettingRepos.Query, p.SettingRepos.Query),
		Update: NewUpdateHandler(p.UserSettingRepos.Command, p.SettingRepos.Query),
		Delete: NewDeleteHandler(p.UserSettingRepos.Query, p.UserSettingRepos.Command),
	}
}
