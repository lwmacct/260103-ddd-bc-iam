package cache

import "encoding/json"

// CacheInfoDTO 缓存信息响应（类似 redis-cli INFO）
type CacheInfoDTO struct {
	// DBSize 当前 key 数量（仅计算应用前缀下的）
	DBSize int64 `json:"db_size"`

	// KeyPrefix 应用使用的 key 前缀
	KeyPrefix string `json:"key_prefix"`

	// MemoryUsage 内存使用量（如果可用）
	MemoryUsage string `json:"memory_usage,omitempty"`

	// RedisVersion Redis 版本
	RedisVersion string `json:"redis_version,omitempty"`
}

// CacheKeyDTO 单个缓存 Key 信息
type CacheKeyDTO struct {
	// Key 完整 key 名称
	Key string `json:"key"`

	// Type key 类型（string, ReJSON-RL, hash 等）
	Type string `json:"type"`

	// TTL 剩余过期时间（秒），-1 表示永不过期，-2 表示不存在
	TTL int64 `json:"ttl"`
}

// ScanKeysResultDTO SCAN 结果
type ScanKeysResultDTO struct {
	// Keys key 列表
	Keys []CacheKeyDTO `json:"keys"`

	// Cursor 下一个游标（用于分页），为 "0" 表示结束
	Cursor string `json:"cursor"`

	// TotalScanned 本次扫描的 key 数量
	TotalScanned int64 `json:"total_scanned"`
}

// CacheValueDTO 单个 Key 的完整信息（含值）
type CacheValueDTO struct {
	CacheKeyDTO

	// Value key 的值（JSON 格式）
	Value json.RawMessage `json:"value"`
}

// DeleteResultDTO 删除结果
type DeleteResultDTO struct {
	// DeletedCount 删除的 key 数量
	DeletedCount int64 `json:"deleted_count"`
}
