// Package stats 定义系统统计的领域模型和接口。
//
// 本包提供系统级统计数据查询能力，定义了：
//   - [SystemStats]: 系统统计信息值对象（用户数、角色数等）
//   - [AuditLogSummary]: 审计日志摘要
//   - [QueryRepository]: 统计查询仓储接口
//
// 统计维度：
// [SystemStats] 汇总以下数据：
//   - 用户统计：总数、活跃、未激活、禁用
//   - 角色统计：角色总数、权限总数
//   - 菜单统计：菜单总数
//   - 近期审计日志
//
// 设计说明：
// 本包仅提供 [QueryRepository]（只读），不涉及数据修改操作。
// 统计数据通常用于管理后台仪表盘展示。
//
// 依赖倒置：
// 本包仅定义接口，实现位于 infrastructure/persistence 包。
package stats
