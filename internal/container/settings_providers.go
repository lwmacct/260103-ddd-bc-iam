package container

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	settingApp "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app/setting"
	settingDomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
	settingsCache "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/infra/cache"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/config"
)

// SettingsCacheResult 使用 fx.Out 提供缓存服务接口。
type SettingsCacheResult struct {
	fx.Out

	// Application 层使用的完整缓存服务
	SettingsCacheService settingApp.SettingsCacheService

	// Infrastructure 层使用的配置变更通知接口（业务语义）
	SettingChangeNotifier settingDomain.SettingChangeNotifier
}

// newSettingsCacheService 创建 Settings 缓存服务。
//
// 绕过 Settings 包的 internal/config.Config 依赖，直接从 IAM 配置中提取 Redis Key Prefix。
//
// 导出两个接口：
//   - SettingsCacheService：Application 层使用
//   - SettingChangeNotifier：Infrastructure 层的缓存装饰器使用
func newSettingsCacheService(client *redis.Client, iamCfg *config.Config) SettingsCacheResult {
	// 创建缓存服务（内部类型）
	service := settingsCache.NewSettingsCacheService(client, iamCfg.Redis.KeyPrefix)

	// 类型断言：确保实现了 SettingChangeNotifier 接口
	// （外部模块的 settingsCacheService 实现了该接口）
	notifier, ok := service.(settingDomain.SettingChangeNotifier)
	if !ok {
		// 如果类型断言失败，说明外部模块实现有变化
		// 这种情况下应该直接使用 external module 的 cache.CacheModule
		panic("settingsCacheService does not implement SettingChangeNotifier interface, " +
			"consider using external module's cache.CacheModule instead")
	}

	return SettingsCacheResult{
		SettingsCacheService:  service,
		SettingChangeNotifier: notifier,
	}
}
