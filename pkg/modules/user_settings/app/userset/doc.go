// Package userset 提供用户配置模块的用例处理器。
//
// 本包实现了用户配置的 CQRS 命令和查询：
//
// # 命令（写操作）
//   - [SetHandler]: 设置单个配置（覆盖系统默认值）
//   - [BatchSetHandler]: 批量设置配置
//   - [ResetHandler]: 重置单个配置（删除用户覆盖值）
//   - [ResetAllHandler]: 重置所有配置（删除所有用户覆盖值）
//
// # 查询（读操作）
//   - [GetHandler]: 获取单个配置（合并系统默认+用户覆盖）
//   - [ListHandler]: 获取配置列表（合并视图）
//   - [ListCategoriesHandler]: 获取配置分类列表
//
// # 跨 BC 依赖
//
// 本包依赖 Settings BC 进行配置校验：
//   - 校验 key 是否存在
//   - ValueType 类型校验
//   - InputType 格式校验（email/url/password 等）
package userset
