package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	appuser "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/application/user"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/user"
)

const userWithRolesCacheTTL = 5 * time.Minute

// userWithRolesCacheService 用户实体缓存服务的 Redis 实现。
//
// 提供用户实体的缓存操作：
//   - Key 格式：{prefix}user:entity:{userID}
//   - TTL：5 分钟
//   - RedisJSON 存储
//   - 直接序列化 Domain 实体（实体已有 json tags）
type userWithRolesCacheService struct {
	client    *redis.Client
	keyPrefix string
}

// NewUserWithRolesCacheService 创建用户实体缓存服务。
func NewUserWithRolesCacheService(client *redis.Client, keyPrefix string) appuser.UserWithRolesCacheService {
	return &userWithRolesCacheService{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// GetUserWithRoles 获取缓存的用户实体（使用 RedisJSON）。
// 缓存未命中返回 nil, nil。
func (s *userWithRolesCacheService) GetUserWithRoles(ctx context.Context, userID uint) (*user.User, error) {
	data, err := s.client.JSONGet(ctx, s.buildKey(userID), "$").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil //nolint:nilnil // cache miss
		}
		return nil, fmt.Errorf("redis json get error: %w", err)
	}

	// JSON.GET $ 返回数组包装：[actual_data]
	var wrapper []*user.User
	if err := json.Unmarshal([]byte(data), &wrapper); err != nil {
		// 缓存数据损坏，删除并返回未命中
		_ = s.client.Del(ctx, s.buildKey(userID))
		return nil, nil //nolint:nilnil,nilerr // corrupted cache treated as miss
	}

	if len(wrapper) == 0 {
		return nil, nil //nolint:nilnil // empty wrapper
	}

	return wrapper[0], nil
}

// SetUserWithRoles 设置用户实体缓存（使用 RedisJSON）。
func (s *userWithRolesCacheService) SetUserWithRoles(ctx context.Context, u *user.User) error {
	// 使用 Pipeline 执行 JSON.SET + EXPIRE
	key := s.buildKey(u.ID)
	pipe := s.client.Pipeline()
	pipe.JSONSet(ctx, key, "$", u)
	pipe.Expire(ctx, key, userWithRolesCacheTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set user entity cache: %w", err)
	}

	return nil
}

// InvalidateUser 失效单个用户缓存。
func (s *userWithRolesCacheService) InvalidateUser(ctx context.Context, userID uint) error {
	return s.client.Del(ctx, s.buildKey(userID)).Err()
}

func (s *userWithRolesCacheService) buildKey(userID uint) string {
	return fmt.Sprintf("%suser:entity:%d", s.keyPrefix, userID)
}

var _ appuser.UserWithRolesCacheService = (*userWithRolesCacheService)(nil)
