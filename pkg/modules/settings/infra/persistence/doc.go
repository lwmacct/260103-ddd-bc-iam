// Package persistence 提供 Settings 模块的数据库持久化实现。
//
// # Overview
//
// 本包实现 Domain 层定义的所有 Repository 接口，提供配置管理功能：
//   - Setting Repositories: 系统配置仓储（全局配置、配置分类）
//   - UserSetting Repositories: 用户配置仓储
//   - OrgSetting Repositories: 组织配置仓储
//   - TeamSetting Repositories: 团队配置仓储
//
// # CQRS 模式
//
// 每个领域模块都有独立的 Command 和 Query 仓储：
//   - Command Repository: 写操作（Create、Update、Delete、Upsert）
//   - Query Repository: 读操作（Get、List、Exists）
//   - 聚合结构体: 同时提供 Command 和 Query，便于依赖注入
//
// # 缓存策略
//
// Query Repository 使用缓存装饰器提升性能：
//   - Cache-Aside 模式：先查缓存，未命中查库，同步回写
//   - Invalidate 策略：写操作后失效相关缓存
//   - TTL: 5 分钟过期
//   - 缓存服务见 cache 包
//
// # 数据模型
//
// 每个领域都有对应的 GORM Model：
//   - Model 与 Entity 分离：Model 包含 GORM tags，Entity 是纯领域模型
//   - 映射函数: toModel() / toEntity() 实现 Model ↔ Entity 转换
//   - 禁止物理外键：使用逻辑关联，由应用层保证数据一致性
//
// # Thread Safety
//
// 所有 Repository 实现都是并发安全的，依赖 GORM 的连接池管理。
// Repository 本身是无状态的，可以安全地在多个 goroutine 中共享。
//
// # 依赖关系
//
// 本包实现 Domain 层的 Repository 接口，依赖：
//   - GORM DB: 数据库连接（通过 Fx 注入）
//   - Redis Client: 缓存存储（通过 cache 服务注入）
package persistence
