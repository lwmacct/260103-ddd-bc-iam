// Package app 提供 User Settings 模块的应用层用例处理器。
//
// # Overview
//
// 本包实现了用户配置管理的核心业务逻辑：
//   - [SetHandler]: 设置单个用户配置（含校验）
//   - [BatchSetHandler]: 批量设置用户配置
//   - [ResetHandler]: 重置单个配置（恢复默认值）
//   - [ResetAllHandler]: 重置所有配置
//   - [GetHandler]: 获取单个配置（合并视图）
//   - [ListHandler]: 获取配置列表（合并视图）
//   - [ListCategoriesHandler]: 获取分类列表
//
// # 跨 BC 依赖
//
// 本包依赖 Settings BC 获取配置定义：
//   - [setting.QueryRepository]: 查询配置 Schema
//   - [setting.SettingCategoryQueryRepository]: 查询配置分类
//
// # 合并视图
//
// 返回给前端的配置合并了系统默认值和用户自定义值：
//   - Value: 实际生效值（用户值 > 默认值）
//   - DefaultValue: 系统默认值
//   - IsCustomized: 是否用户自定义
//
// # Thread Safety
//
// 所有 Handler 是无状态的，通过依赖注入获取仓储实例。
// 并发安全性由仓储实现保证。
package app
