package container

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/config"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/platform/cache"

	appcache "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application/cache"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application/setting"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/application/auth"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/application/user"
)

// CacheServicesResult 使用 fx.Out 批量返回所有缓存服务。
type CacheServicesResult struct {
	fx.Out

	UserSettingQuery setting.UserSettingQueryCacheService
	UserSetting      setting.UserSettingCacheService
	UserWithRoles    user.UserWithRolesCacheService
	Permission       auth.PermissionCacheService
	Settings         setting.SettingsCacheService
	Admin            appcache.AdminCacheService
}

// CacheModule 提供所有缓存服务。
var CacheModule = fx.Module("cache",
	fx.Provide(NewAllCacheServices),
)

// NewAllCacheServices 创建所有缓存服务。
func NewAllCacheServices(client *redis.Client, cfg *config.Config) CacheServicesResult {
	prefix := cfg.Data.RedisKeyPrefix
	return CacheServicesResult{
		UserSettingQuery: cache.NewUserSettingQueryCacheService(client, prefix),
		UserSetting:      cache.NewUserSettingCacheService(client, prefix),
		UserWithRoles:    cache.NewUserWithRolesCacheService(client, prefix),
		Permission:       cache.NewPermissionCacheService(client, prefix),
		Settings:         cache.NewSettingsCacheService(client, prefix),
		Admin:            cache.NewAdminCacheService(client, prefix),
	}
}
