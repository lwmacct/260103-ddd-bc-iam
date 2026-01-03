package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"

	appcache "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application/cache"
)

// adminCacheService AdminCacheService 的 Redis 实现。
type adminCacheService struct {
	client    *redis.Client
	keyPrefix string
}

// NewAdminCacheService 创建 AdminCacheService 实现。
func NewAdminCacheService(client *redis.Client, keyPrefix string) appcache.AdminCacheService {
	return &adminCacheService{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// ScanKeys 扫描匹配 pattern 的 keys。
func (s *adminCacheService) ScanKeys(ctx context.Context, pattern string, cursor uint64, count int64) ([]string, uint64, error) {
	return s.client.Scan(ctx, cursor, pattern, count).Result()
}

// GetKeyInfo 获取单个 key 的元信息。
func (s *adminCacheService) GetKeyInfo(ctx context.Context, key string) (string, int64, error) {
	keyType, err := s.client.Type(ctx, key).Result()
	if err != nil {
		return "", 0, fmt.Errorf("get key type failed: %w", err)
	}

	ttl, err := s.client.TTL(ctx, key).Result()
	if err != nil {
		return keyType, 0, fmt.Errorf("get key ttl failed: %w", err)
	}

	return keyType, int64(ttl.Seconds()), nil
}

// GetKeyValue 获取单个 key 的值。
func (s *adminCacheService) GetKeyValue(ctx context.Context, key string) (json.RawMessage, string, int64, error) {
	keyType, ttl, err := s.GetKeyInfo(ctx, key)
	if err != nil {
		return nil, "", 0, err
	}

	var value any
	switch keyType {
	case "string":
		value, err = s.client.Get(ctx, key).Result()
	case "ReJSON-RL":
		// RedisJSON 类型
		jsonStr, jerr := s.client.JSONGet(ctx, key, "$").Result()
		if jerr == nil {
			value = json.RawMessage(jsonStr)
		}
		err = jerr
	case "hash":
		value, err = s.client.HGetAll(ctx, key).Result()
	case "list":
		value, err = s.client.LRange(ctx, key, 0, -1).Result()
	case "set":
		value, err = s.client.SMembers(ctx, key).Result()
	case "zset":
		value, err = s.client.ZRangeWithScores(ctx, key, 0, -1).Result()
	default:
		value = "unsupported type: " + keyType
	}

	if err != nil {
		return nil, keyType, ttl, fmt.Errorf("get key value failed: %w", err)
	}

	// 序列化值为 JSON
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return nil, keyType, ttl, fmt.Errorf("marshal value failed: %w", err)
	}

	return valueBytes, keyType, ttl, nil
}

// KeyExists 检查 key 是否存在。
func (s *adminCacheService) KeyExists(ctx context.Context, key string) (bool, error) {
	n, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// DeleteKey 删除单个 key。
func (s *adminCacheService) DeleteKey(ctx context.Context, key string) (int64, error) {
	return s.client.Del(ctx, key).Result()
}

// DeleteKeys 批量删除 keys。
func (s *adminCacheService) DeleteKeys(ctx context.Context, keys ...string) (int64, error) {
	if len(keys) == 0 {
		return 0, nil
	}
	return s.client.Del(ctx, keys...).Result()
}

// GetInfo 获取缓存服务器信息。
func (s *adminCacheService) GetInfo(ctx context.Context) (string, string, error) {
	// 获取 Redis 服务器信息
	serverInfo, err := s.client.Info(ctx, "server").Result()
	if err != nil {
		return "", "", fmt.Errorf("get server info failed: %w", err)
	}
	version := extractValue(serverInfo, "redis_version:")

	// 获取内存使用量
	memInfo, err := s.client.Info(ctx, "memory").Result()
	if err != nil {
		return version, "", fmt.Errorf("get memory info failed: %w", err)
	}
	memoryUsage := extractValue(memInfo, "used_memory_human:")

	return version, memoryUsage, nil
}

// KeyPrefix 返回配置的 key 前缀。
func (s *adminCacheService) KeyPrefix() string {
	return s.keyPrefix
}

// extractValue 从 INFO 输出中提取值。
func extractValue(info, prefix string) string {
	for line := range strings.SplitSeq(info, "\n") {
		line = strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, prefix); ok {
			return after
		}
	}
	return ""
}
