package cache

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/config"
	appauth "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/application/auth"
	appuser "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/application/user"
)

// CacheModule 提供 IAM 模块的专属缓存服务。
var CacheModule = fx.Module("iam.cache",
	fx.Provide(
		newPermissionCacheService,
		newUserWithRolesCacheService,
	),
)

// newPermissionCacheService 创建权限缓存服务。
func newPermissionCacheService(client *redis.Client, cfg *config.Config) appauth.PermissionCacheService {
	return NewPermissionCacheService(client, cfg.Data.RedisKeyPrefix)
}

// newUserWithRolesCacheService 创建用户实体缓存服务。
func newUserWithRolesCacheService(client *redis.Client, cfg *config.Config) appuser.UserWithRolesCacheService {
	return NewUserWithRolesCacheService(client, cfg.Data.RedisKeyPrefix)
}
