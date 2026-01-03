// Package cache 提供 IAM 模块的专属缓存服务。
//
// 缓存服务实现 Application 层定义的缓存接口：
//   - PermissionCacheService: 用户角色和权限缓存
//   - UserWithRolesCacheService: 用户实体缓存
//
// 使用 RedisJSON 进行序列化，TTL 为 5 分钟。
package cache
