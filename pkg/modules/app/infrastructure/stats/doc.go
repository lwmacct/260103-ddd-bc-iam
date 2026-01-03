// Package stats 提供系统统计信息的聚合查询实现。
//
// 本包实现 domain/stats.QueryRepository 接口，提供跨域聚合查询能力。
// 特点：
//   - 聚合查询多个表（users, roles, permissions, menus, audit_logs）
//   - 只读操作，无写入
//   - 直接使用 GORM Table() 查询，无独立 Model
//
// # 组件职责
//
//   - [NewQueryRepository]: 创建统计查询仓储实例
//
// # 使用示例
//
//	repo := stats.NewQueryRepository(db)
//	systemStats, err := repo.GetSystemStats(10) // 获取最近 10 条审计日志
package stats
