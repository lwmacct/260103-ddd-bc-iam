package cache

import (
	"context"
	"encoding/json"
)

// AdminCacheService 缓存管理服务接口。
//
// 为管理员提供 Redis 风格的缓存运维能力，支持键扫描、值获取和删除操作。
// Infrastructure 层实现此接口，Application 层处理器依赖此抽象。
type AdminCacheService interface {
	// ScanKeys 扫描匹配 pattern 的 keys。
	// cursor 用于分页，首次传 0。
	// 返回匹配的 keys、下一个 cursor 和错误。
	ScanKeys(ctx context.Context, pattern string, cursor uint64, count int64) (keys []string, nextCursor uint64, err error)

	// GetKeyInfo 获取单个 key 的元信息（类型、TTL）。
	GetKeyInfo(ctx context.Context, key string) (keyType string, ttl int64, err error)

	// GetKeyValue 获取单个 key 的值。
	// 根据 key 类型返回不同格式的值（string、JSON、hash 等）。
	GetKeyValue(ctx context.Context, key string) (value json.RawMessage, keyType string, ttl int64, err error)

	// KeyExists 检查 key 是否存在。
	KeyExists(ctx context.Context, key string) (bool, error)

	// DeleteKey 删除单个 key。
	DeleteKey(ctx context.Context, key string) (deleted int64, err error)

	// DeleteKeys 批量删除 keys。
	DeleteKeys(ctx context.Context, keys ...string) (deleted int64, err error)

	// GetInfo 获取缓存服务器信息。
	// 返回服务器版本、内存使用等信息。
	GetInfo(ctx context.Context) (version string, memoryUsage string, err error)

	// KeyPrefix 返回配置的 key 前缀。
	KeyPrefix() string
}
