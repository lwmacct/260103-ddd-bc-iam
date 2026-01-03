// Package userset 定义用户配置领域模型和仓储接口。
//
// 本包定义了：
//   - [UserSetting]: 用户配置实体（用户对系统配置的覆盖值）
//   - [CommandRepository]: 写仓储接口（Upsert、Delete、BatchUpsert）
//   - [QueryRepository]: 读仓储接口（FindByUserAndKey、FindByUser）
//
// 依赖倒置：本包仅定义接口，实现位于 infra/persistence。
//
// # 覆盖模式
//
// 用户配置采用覆盖（Override）模式：
//   - 只存储用户修改过的配置值
//   - 删除配置即恢复系统默认值
//   - 依赖 Settings BC 获取配置 Schema 和默认值
package userset
