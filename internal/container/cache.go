package container

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/config"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/application/auth"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/application/user"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/platform/cache"
)

// CacheServicesResult 使用 fx.Out 批量返回所有缓存服务。
type CacheServicesResult struct {
	fx.Out

	UserWithRoles user.UserWithRolesCacheService
	Permission    auth.PermissionCacheService
}

// CacheModule 提供所有缓存服务。
var CacheModule = fx.Module("cache",
	fx.Provide(NewAllCacheServices),
)

// NewAllCacheServices 创建所有缓存服务。
func NewAllCacheServices(client *redis.Client, cfg *config.Config) CacheServicesResult {
	prefix := cfg.Data.RedisKeyPrefix
	return CacheServicesResult{
		UserWithRoles: cache.NewUserWithRolesCacheService(client, prefix),
		Permission:    cache.NewPermissionCacheService(client, prefix),
	}
}
