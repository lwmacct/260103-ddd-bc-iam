// Package cache 提供缓存管理的应用层服务。
//
// 本包为管理员提供 Redis 风格的缓存运维 API：
//   - [InfoHandler]: 查看缓存信息（类似 redis-cli INFO）
//   - [ScanKeysHandler]: 扫描 Keys（类似 redis-cli SCAN）
//   - [GetKeyHandler]: 获取 Key 值（类似 redis-cli GET/JSON.GET）
//   - [DeleteHandler]: 删除 Keys（类似 redis-cli DEL）
//
// # API 端点
//
//   - GET /api/admin/cache/info - 缓存状态信息
//   - GET /api/admin/cache/keys - 按 pattern 扫描 keys
//   - GET /api/admin/cache/keys/{key} - 获取单个 key 的值
//   - DELETE /api/admin/cache/keys/{key} - 删除单个 key
//   - DELETE /api/admin/cache/keys?pattern=xxx - 按 pattern 批量删除
//
// # 使用场景
//
//   - 调试：检查缓存状态是否符合预期
//   - 运维：数据不一致时手动清缓存
//   - 测试：集成测试前清空缓存确保干净环境
//
// 本包是 DDD 架构的 Application 层。
package cache
