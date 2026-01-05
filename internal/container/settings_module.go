package container

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app/setting"
	settingsconfig "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/config"
	settingsPersistence "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/infra/persistence"
)

// SettingsModule 提供 Settings Bounded Context 的 Fx 模块。
//
// 注意：此模块排除 cache.CacheModule，改用 IAM 项目提供的缓存服务。
func SettingsModule() fx.Option {
	return fx.Module("settings",
		fx.Supply(settingsconfig.DefaultConfig()),
		settingsPersistence.RepositoryModule,
		setting.UseCaseModule,
	)
}
