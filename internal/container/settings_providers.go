package container

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/config"
	settingApp "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app/setting"
	settingDomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
	settingsCache "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/infra/cache"
)

// SettingsCacheResult 使用 fx.Out 提供缓存服务接口。
type SettingsCacheResult struct {
	fx.Out

	SettingsCacheService settingApp.SettingsCacheService
	CacheInvalidator     settingDomain.CacheInvalidator
}

// newSettingsCacheService 创建 Settings 缓存服务。
//
// 绕过 Settings 包的 internal/config.Config 依赖，直接从 IAM 配置中提取 Redis Key Prefix。
func newSettingsCacheService(client *redis.Client, iamCfg *config.Config) SettingsCacheResult {
	service := settingsCache.NewSettingsCacheService(client, iamCfg.Redis.KeyPrefix)

	return SettingsCacheResult{
		SettingsCacheService: service,
		CacheInvalidator:     service,
	}
}
