package cache

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	appauth "github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/auth"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/user"
)

// CacheModule 提供 IAM 模块的专属缓存服务。
var CacheModule = fx.Module("iam.cache",
	fx.Provide(
		newPermissionCacheService,
		newUserWithRolesCacheService,
	),
)

// newPermissionCacheService 创建权限缓存服务。
// Platform 配置（keyPrefix）由容器层直接注入。
func newPermissionCacheService(client *redis.Client, keyPrefix string) appauth.PermissionCacheService {
	return NewPermissionCacheService(client, keyPrefix)
}

// newUserWithRolesCacheService 创建用户实体缓存服务。
// Platform 配置（keyPrefix）由容器层直接注入。
func newUserWithRolesCacheService(client *redis.Client, keyPrefix string) user.UserWithRolesCacheService {
	return NewUserWithRolesCacheService(client, keyPrefix)
}
