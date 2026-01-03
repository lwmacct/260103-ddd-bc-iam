package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	appsetting "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application/setting"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
)

const (
	userSettingQueryCacheTTL  = 30 * time.Minute
	userSettingQueryKeyPrefix = "usersetting:user:"
)

// userSettingQueryCacheService 用户设置查询缓存的 Redis 实现。
//
// 存储原始 UserSetting 记录，用于 Repository 层减少数据库查询。
// 采用用户维度全量缓存策略：一次查询缓存用户所有自定义配置。
// 直接序列化 Domain 实体（实体已有 json tags）。
type userSettingQueryCacheService struct {
	client    *redis.Client
	keyPrefix string
}

// NewUserSettingQueryCacheService 创建用户设置查询缓存服务。
func NewUserSettingQueryCacheService(client *redis.Client, keyPrefix string) appsetting.UserSettingQueryCacheService {
	return &userSettingQueryCacheService{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// GetByUser 获取用户的所有自定义配置缓存（使用 RedisJSON）。
func (s *userSettingQueryCacheService) GetByUser(ctx context.Context, userID uint) (map[string]*setting.UserSetting, error) {
	key := s.buildKey(userID)
	data, err := s.client.JSONGet(ctx, key, "$").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil //nolint:nilnil // cache miss
		}
		return nil, fmt.Errorf("redis json get error: %w", err)
	}

	// JSON.GET $ 返回数组包装：[actual_data]
	// actual_data 本身也是数组
	var wrapper [][]*setting.UserSetting
	if err := json.Unmarshal([]byte(data), &wrapper); err != nil {
		// 缓存数据损坏，删除并返回未命中
		_ = s.client.Del(ctx, key)
		return nil, nil //nolint:nilerr,nilnil // corrupted cache treated as miss
	}

	if len(wrapper) == 0 {
		return nil, nil //nolint:nilnil // empty wrapper
	}

	entities := wrapper[0]
	result := make(map[string]*setting.UserSetting, len(entities))
	for _, entity := range entities {
		result[entity.SettingKey] = entity
	}
	return result, nil
}

// SetByUser 设置用户的所有自定义配置缓存（使用 RedisJSON）。
func (s *userSettingQueryCacheService) SetByUser(ctx context.Context, userID uint, settings []*setting.UserSetting) error {
	// 使用 Pipeline 执行 JSON.SET + EXPIRE（直接序列化实体切片）
	key := s.buildKey(userID)
	pipe := s.client.Pipeline()
	pipe.JSONSet(ctx, key, "$", settings)
	pipe.Expire(ctx, key, userSettingQueryCacheTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set user setting query cache: %w", err)
	}

	return nil
}

// DeleteByUser 删除用户的所有配置缓存。
func (s *userSettingQueryCacheService) DeleteByUser(ctx context.Context, userID uint) error {
	key := s.buildKey(userID)
	return s.client.Del(ctx, key).Err()
}

// buildKey 构建缓存 key。
// 格式：{prefix}usersetting:user:{userID}
func (s *userSettingQueryCacheService) buildKey(userID uint) string {
	return fmt.Sprintf("%s%s%d", s.keyPrefix, userSettingQueryKeyPrefix, userID)
}

var _ appsetting.UserSettingQueryCacheService = (*userSettingQueryCacheService)(nil)
