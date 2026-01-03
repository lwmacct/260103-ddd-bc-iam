package cache

import (
	"context"
	"errors"
	"fmt"
)

// ErrKeyNotFound key 不存在错误
var ErrKeyNotFound = errors.New("key not found")

// GetKeyHandler 获取单个 Key 值的 Handler（类似 redis-cli GET/JSON.GET）。
type GetKeyHandler struct {
	svc AdminCacheService
}

// NewGetKeyHandler 创建获取 Key 值的 Handler。
func NewGetKeyHandler(svc AdminCacheService) *GetKeyHandler {
	return &GetKeyHandler{svc: svc}
}

// Handle 获取单个 Key 的值。
//
// key 参数是完整的 key 名称（含应用前缀）。
func (h *GetKeyHandler) Handle(ctx context.Context, key string) (*CacheValueDTO, error) {
	// 检查 key 是否存在
	exists, err := h.svc.KeyExists(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("check key exists failed: %w", err)
	}
	if !exists {
		return nil, ErrKeyNotFound
	}

	// 获取值及元信息
	value, keyType, ttl, err := h.svc.GetKeyValue(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("get key value failed: %w", err)
	}

	return &CacheValueDTO{
		CacheKeyDTO: CacheKeyDTO{
			Key:  key,
			Type: keyType,
			TTL:  ttl,
		},
		Value: value,
	}, nil
}
