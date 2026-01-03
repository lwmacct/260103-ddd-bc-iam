package container

import (
	"go.uber.org/fx"

	settingsHandler "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/adapters/gin/handler"
	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app/setting"
	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/infra/persistence"
)

// SettingsModule 提供 Settings Bounded Context 的 Fx 模块。
//
// 注意：此模块排除 cache.CacheModule，改用 IAM 项目提供的缓存服务
// 以避免依赖 Settings 包的 internal/config.Config。
func SettingsModule() fx.Option {
	return fx.Module("settings",
		// 基础设施层（排除 cache，使用自定义提供）
		persistence.RepositoryModule,

		// 应用层
		setting.UseCaseModule,

		// 适配器层
		settingsHandler.HandlerModule,
	)
}
