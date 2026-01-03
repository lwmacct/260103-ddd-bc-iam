// Package cache 提供 User Settings 模块的缓存服务。
//
// 本包包含：
//   - [CacheService]: 用户配置缓存服务接口
//   - Redis 缓存实现（后续实现）
//
// 缓存策略：
//   - Key 格式：{prefix}user_settings:{userID}
//   - TTL：5 分钟
//   - 写操作后失效缓存，读操作回源后写入
package cache
