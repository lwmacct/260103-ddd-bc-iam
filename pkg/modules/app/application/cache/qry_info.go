package cache

import (
	"context"
	"fmt"
)

// InfoHandler 缓存信息查询 Handler（类似 redis-cli INFO）。
type InfoHandler struct {
	svc AdminCacheService
}

// NewInfoHandler 创建缓存信息查询 Handler。
func NewInfoHandler(svc AdminCacheService) *InfoHandler {
	return &InfoHandler{svc: svc}
}

// Handle 查询缓存信息。
func (h *InfoHandler) Handle(ctx context.Context) (*CacheInfoDTO, error) {
	keyPrefix := h.svc.KeyPrefix()
	result := &CacheInfoDTO{
		KeyPrefix: keyPrefix,
	}

	// 统计应用前缀下的 key 数量
	pattern := keyPrefix + "*"
	count, err := h.countKeys(ctx, pattern)
	if err != nil {
		return nil, fmt.Errorf("count keys failed: %w", err)
	}
	result.DBSize = count

	// 获取 Redis 服务器信息
	version, memoryUsage, err := h.svc.GetInfo(ctx)
	if err == nil {
		result.RedisVersion = version
		result.MemoryUsage = memoryUsage
	}

	return result, nil
}

// countKeys 使用 SCAN 统计匹配的 key 数量
func (h *InfoHandler) countKeys(ctx context.Context, pattern string) (int64, error) {
	var count int64
	var cursor uint64

	for {
		keys, nextCursor, err := h.svc.ScanKeys(ctx, pattern, cursor, 1000)
		if err != nil {
			return count, err
		}

		count += int64(len(keys))

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return count, nil
}
