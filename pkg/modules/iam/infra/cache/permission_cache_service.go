package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	appauth "github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/auth"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/role"
)

const permissionCacheTTL = 5 * time.Minute

// userPermissionsCache 用户权限缓存数据结构。
// 直接使用 Domain 类型（role.Permission 已有 json tags）。
type userPermissionsCache struct {
	Roles       []string          `json:"roles"`
	Permissions []role.Permission `json:"permissions"`
}

// permissionCacheService 权限缓存服务的 Redis 实现。
//
// 提供用户权限的缓存操作：
//   - Key 格式：{prefix}user:perms:{userID}
//   - TTL：5 分钟
//   - JSON 序列化存储
//   - 直接序列化 Domain 实体（实体已有 json tags）
type permissionCacheService struct {
	client    *redis.Client
	keyPrefix string
}

// NewPermissionCacheService 创建权限缓存服务。
func NewPermissionCacheService(client *redis.Client, keyPrefix string) appauth.PermissionCacheService {
	return &permissionCacheService{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// GetUserPermissions 获取用户权限（使用 RedisJSON）。
// 缓存未命中返回三个 nil。
func (s *permissionCacheService) GetUserPermissions(ctx context.Context, userID uint) ([]string, []role.Permission, error) {
	data, err := s.client.JSONGet(ctx, s.buildKey(userID), "$").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil, nil // cache miss
		}
		return nil, nil, fmt.Errorf("redis json get error: %w", err)
	}

	// JSON.GET $ 返回数组包装：[actual_data]
	var wrapper []userPermissionsCache
	if err := json.Unmarshal([]byte(data), &wrapper); err != nil {
		// 缓存数据损坏，删除并返回未命中
		_ = s.client.Del(ctx, s.buildKey(userID))
		return nil, nil, nil //nolint:nilerr // corrupted cache treated as miss
	}

	if len(wrapper) == 0 {
		return nil, nil, nil // empty wrapper
	}

	return wrapper[0].Roles, wrapper[0].Permissions, nil
}

// SetUserPermissions 设置用户权限缓存（使用 RedisJSON）。
func (s *permissionCacheService) SetUserPermissions(ctx context.Context, userID uint, roles []string, permissions []role.Permission) error {
	perms := userPermissionsCache{Roles: roles, Permissions: permissions}

	// 使用 Pipeline 执行 JSON.SET + EXPIRE
	key := s.buildKey(userID)
	pipe := s.client.Pipeline()
	pipe.JSONSet(ctx, key, "$", perms)
	pipe.Expire(ctx, key, permissionCacheTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set permission cache: %w", err)
	}

	return nil
}

// InvalidateUser 失效单个用户缓存。
func (s *permissionCacheService) InvalidateUser(ctx context.Context, userID uint) error {
	return s.client.Del(ctx, s.buildKey(userID)).Err()
}

// InvalidateUsers 批量失效用户缓存。
func (s *permissionCacheService) InvalidateUsers(ctx context.Context, userIDs []uint) error {
	if len(userIDs) == 0 {
		return nil
	}

	keys := make([]string, 0, len(userIDs))
	for _, id := range userIDs {
		keys = append(keys, s.buildKey(id))
	}

	return s.client.Del(ctx, keys...).Err()
}

// InvalidateAll 失效所有用户权限缓存。
func (s *permissionCacheService) InvalidateAll(ctx context.Context) error {
	pattern := s.keyPrefix + "user:perms:*"

	iter := s.client.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan keys: %w", err)
	}

	if len(keys) > 0 {
		return s.client.Del(ctx, keys...).Err()
	}

	return nil
}

func (s *permissionCacheService) buildKey(userID uint) string {
	return fmt.Sprintf("%suser:perms:%d", s.keyPrefix, userID)
}

var _ appauth.PermissionCacheService = (*permissionCacheService)(nil)
