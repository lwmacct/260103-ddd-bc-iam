package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	appsetting "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application/setting"
)

const userSettingCacheTTL = 30 * time.Minute

// userSettingCacheService 用户设置缓存服务的 Redis 实现。
//
// Key 格式：{prefix}user:{userID}:setting:{key}
// TTL：30 分钟
type userSettingCacheService struct {
	client    *redis.Client
	keyPrefix string
}

// NewUserSettingCacheService 创建用户设置缓存服务。
func NewUserSettingCacheService(client *redis.Client, keyPrefix string) appsetting.UserSettingCacheService {
	return &userSettingCacheService{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// =========================================================================
// 单条操作
// =========================================================================

// Get 获取用户的有效设置值（使用 RedisJSON）。
func (s *userSettingCacheService) Get(ctx context.Context, userID uint, key string) (*appsetting.EffectiveUserSetting, error) {
	data, err := s.client.JSONGet(ctx, s.buildKey(userID, key), "$").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil //nolint:nilnil // cache miss
		}
		return nil, fmt.Errorf("redis json get error: %w", err)
	}

	// JSON.GET $ 返回数组包装：[actual_data]
	var wrapper []appsetting.EffectiveUserSetting
	if err := json.Unmarshal([]byte(data), &wrapper); err != nil {
		// 缓存数据损坏，删除并返回未命中
		_ = s.client.Del(ctx, s.buildKey(userID, key))
		return nil, nil //nolint:nilerr,nilnil // corrupted cache
	}

	if len(wrapper) == 0 {
		return nil, nil //nolint:nilnil // empty wrapper
	}

	return &wrapper[0], nil
}

// Set 缓存用户的有效设置值（使用 RedisJSON）。
func (s *userSettingCacheService) Set(ctx context.Context, userID uint, value *appsetting.EffectiveUserSetting) error {
	if value == nil {
		return nil
	}

	// 使用 Pipeline 执行 JSON.SET + EXPIRE
	key := s.buildKey(userID, value.Key)
	pipe := s.client.Pipeline()
	pipe.JSONSet(ctx, key, "$", value)
	pipe.Expire(ctx, key, userSettingCacheTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set user setting cache: %w", err)
	}

	return nil
}

// =========================================================================
// 批量操作
// =========================================================================

// GetByKeys 批量获取用户的有效设置值（使用 RedisJSON）。
func (s *userSettingCacheService) GetByKeys(ctx context.Context, userID uint, keys []string) (map[string]*appsetting.EffectiveUserSetting, error) {
	if len(keys) == 0 {
		return map[string]*appsetting.EffectiveUserSetting{}, nil
	}

	redisKeys := make([]string, len(keys))
	for i, k := range keys {
		redisKeys[i] = s.buildKey(userID, k)
	}

	// JSON.MGET 返回每个 key 的 JSON 字符串
	values, err := s.client.JSONMGet(ctx, "$", redisKeys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to json mget: %w", err)
	}

	result := make(map[string]*appsetting.EffectiveUserSetting)
	for i, v := range values {
		if v == nil {
			continue
		}

		data, ok := v.(string)
		if !ok {
			continue
		}

		// JSON.MGET 返回数组包装
		var wrapper []appsetting.EffectiveUserSetting
		if json.Unmarshal([]byte(data), &wrapper) == nil && len(wrapper) > 0 {
			result[keys[i]] = &wrapper[0]
		}
	}

	return result, nil
}

// SetBatch 批量设置用户的有效设置值（使用 RedisJSON）。
func (s *userSettingCacheService) SetBatch(ctx context.Context, userID uint, values []*appsetting.EffectiveUserSetting) error {
	if len(values) == 0 {
		return nil
	}

	pipe := s.client.Pipeline()
	for _, v := range values {
		if v == nil {
			continue
		}

		key := s.buildKey(userID, v.Key)
		pipe.JSONSet(ctx, key, "$", v)
		pipe.Expire(ctx, key, userSettingCacheTTL)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute pipeline: %w", err)
	}

	return nil
}

// =========================================================================
// 删除操作
// =========================================================================

// Delete 删除用户的指定设置缓存。
func (s *userSettingCacheService) Delete(ctx context.Context, userID uint, key string) error {
	return s.client.Del(ctx, s.buildKey(userID, key)).Err()
}

// DeleteByKeys 批量删除用户的指定设置缓存。
func (s *userSettingCacheService) DeleteByKeys(ctx context.Context, userID uint, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	redisKeys := make([]string, len(keys))
	for i, k := range keys {
		redisKeys[i] = s.buildKey(userID, k)
	}

	return s.client.Del(ctx, redisKeys...).Err()
}

// DeleteByUser 删除用户的所有设置缓存。
func (s *userSettingCacheService) DeleteByUser(ctx context.Context, userID uint) error {
	pattern := s.buildUserPattern(userID)
	return s.deleteByPattern(ctx, pattern)
}

// DeleteBySettingKey 删除所有用户的某个设置缓存。
func (s *userSettingCacheService) DeleteBySettingKey(ctx context.Context, key string) error {
	pattern := s.buildSettingKeyPattern(key)
	return s.deleteByPattern(ctx, pattern)
}

// DeleteBySettingKeys 批量删除所有用户的多个设置缓存。
func (s *userSettingCacheService) DeleteBySettingKeys(ctx context.Context, keys []string) error {
	for _, key := range keys {
		if err := s.DeleteBySettingKey(ctx, key); err != nil {
			slog.Warn("failed to delete user setting cache by key", "key", key, "error", err.Error())
		}
	}
	return nil
}

// deleteByPattern 使用 SCAN 删除匹配模式的所有 key。
func (s *userSettingCacheService) deleteByPattern(ctx context.Context, pattern string) error {
	iter := s.client.Scan(ctx, 0, pattern, 100).Iterator()
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

// =========================================================================
// 内部辅助方法
// =========================================================================

func (s *userSettingCacheService) buildKey(userID uint, settingKey string) string {
	return fmt.Sprintf("%suser:%d:setting:%s", s.keyPrefix, userID, settingKey)
}

func (s *userSettingCacheService) buildUserPattern(userID uint) string {
	return fmt.Sprintf("%suser:%d:setting:*", s.keyPrefix, userID)
}

func (s *userSettingCacheService) buildSettingKeyPattern(settingKey string) string {
	return fmt.Sprintf("%suser:*:setting:%s", s.keyPrefix, settingKey)
}

var _ appsetting.UserSettingCacheService = (*userSettingCacheService)(nil)
