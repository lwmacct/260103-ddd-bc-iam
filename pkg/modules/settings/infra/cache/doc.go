// Package cache 提供 User Settings 模块的缓存服务。
//
// # Overview
//
// 本包提供用户配置的缓存服务实现：
//   - [CacheService]: 用户配置缓存服务接口
//   - 使用 Redis 作为缓存存储
//
// 缓存策略：
//   - Key 格式：{prefix}user_settings:{userID}
//   - TTL：5 分钟
//   - Cache-Aside 模式：先查缓存，未命中查库，同步回写
//   - Invalidate 策略：写操作后失效缓存
//
// # Usage
//
//	// 查询用户配置（先查缓存）
//	settings, err := s.cache.GetByUser(ctx, userID)
//	if err != nil {
//	    // 缓存未命中，查询数据库
//	    settings, err = s.queryRepo.FindByUser(ctx, userID)
//	    // 回写缓存
//	    _ = s.cache.Set(ctx, settings)
//	}
//
//	// 更新配置后失效缓存
//	_ = s.cache.InvalidateByUser(ctx, userID)
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
