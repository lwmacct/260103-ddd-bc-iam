package cache

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	appauth "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/auth"
	appuser "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/user"
	config "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/config"
)

// CacheModule 提供 IAM 模块的专属缓存服务。
var CacheModule = fx.Module("iam.cache",
	fx.Provide(
		newPermissionCacheService,
		newUserWithRolesCacheService,
	),
)

// newPermissionCacheService 创建权限缓存服务。
func newPermissionCacheService(client *redis.Client, iamCfg *config.Config) appauth.PermissionCacheService {
	return NewPermissionCacheService(client, iamCfg.RedisCache.KeyPrefix)
}

// newUserWithRolesCacheService 创建用户实体缓存服务。
func newUserWithRolesCacheService(client *redis.Client, iamCfg *config.Config) appuser.UserWithRolesCacheService {
	return NewUserWithRolesCacheService(client, iamCfg.RedisCache.KeyPrefix)
}
