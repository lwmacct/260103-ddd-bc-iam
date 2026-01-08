// Package cache 提供 IAM 模块的专属缓存服务。
//
// # Overview
//
// 本包实现 Application 层定义的缓存接口，使用 Redis 作为缓存存储：
//   - [PermissionCacheService]: 用户角色和权限缓存
//   - [UserWithRolesCacheService]: 用户实体缓存
//
// 缓存策略：
//   - Cache-Aside 模式：先查缓存，未命中查库，同步回写
//   - Invalidate 策略：写操作后失效缓存
//   - TTL: 5 分钟过期
//   - 序列化：RedisJSON (JSON.GET/JSON.SET)
//
// # Usage
//
// 缓存服务通过 Fx 依赖注入自动注册：
//
//	// 查询权限（先查缓存）
//	perms, err := s.permCache.GetByUser(ctx, userID)
//
//	// 更新权限后失效缓存
//	_ = s.permCache.InvalidateByUser(ctx, userID)
//
// # Thread Safety
//
// 所有缓存服务方法都是并发安全的，可以安全地在多个 goroutine 中调用。
// Redis 客户端本身是连接池，支持并发访问。
//
// # 依赖关系
//
// 本包实现 Application 层定义的缓存接口，依赖 Redis 客户端（通过 Fx 注入）。
package cache
