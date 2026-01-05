package container

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app/setting"
)

// SettingsModule 提供 Settings Bounded Context 的 Fx 模块。
//
// 注意：此模块排除 cache.CacheModule 和 handler.HandlerModule，
// 原因：
//  1. cache：改用 IAM 项目提供的缓存服务
//  2. handler：上游 HandlerModule 设计与 UseCaseModule 不兼容（聚合结构体 vs 单独实例）
//     因此 Handler 的创建移到 HTTPModule 中处理。
func SettingsModule() fx.Option {
	return fx.Module("settings",
		// 基础设施层（排除 cache，使用自定义提供）
		// 注意：persistence.RepositoryModule 需要单独导入，这里不包含它
		setting.UseCaseModule,
	)
}
