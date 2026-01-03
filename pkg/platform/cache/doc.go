// Package cache 提供缓存服务的 Redis 实现。
//
// # 概述
//
// 本包实现 Domain 层定义的缓存服务接口，使用 RedisJSON 原生 JSON 类型存储。
//
// # 客户端管理
//
// [NewClient] 创建 Redis 客户端连接（定义在 redis_client.go）：
//   - 支持单节点和集群模式
//   - 连接池管理
//   - OpenTelemetry 追踪集成
//
// # 缓存服务
//
// 本包提供的缓存服务实现：
//   - [NewSettingCacheService]: 系统设置缓存
//   - [NewSettingCategoryCacheService]: 设置分类缓存
//   - [NewUserSettingCacheService]: 用户设置缓存
//   - [NewPermissionCacheService]: 权限缓存
//   - [NewUserWithRolesCacheService]: 用户实体缓存（含角色和权限）
//   - [NewSettingsCacheService]: Settings API 响应缓存
//
// # 存储格式
//
// 所有缓存使用 RedisJSON 模块：
//   - 写入：JSON.SET key $ value + EXPIRE（Pipeline 原子执行）
//   - 读取：JSON.GET key $（返回数组包装 [actual_data]）
//
// # 键命名规范
//
// 格式：{prefix}{模块}:{scope}:{id}
//   - 用户权限：{prefix}user:perms:{user_id}
//   - 用户实体：{prefix}user:entity:{user_id}
//   - 系统设置：{prefix}setting:{key}
package cache
