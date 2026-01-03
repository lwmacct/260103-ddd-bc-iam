// Package usersetting 定义用户设置领域模型和仓储接口。
//
// 本包定义了：
//   - [UserSetting]: 用户设置实体（用户对系统配置的自定义值）
//   - [CommandRepository]: 写仓储接口（创建、更新、删除、Upsert）
//   - [QueryRepository]: 读仓储接口（查询用户自定义值）
//
// 设计说明：
// UserSetting 是用户对系统配置的自定义覆盖值，与 Settings BC 的 Setting 实体配合使用：
//   - Settings BC (setting.Setting): 定义配置 Schema（系统默认值、验证规则）
//   - IAM BC (usersetting.UserSetting): 存储用户自定义值
//
// 数据流：
//  1. GET /api/user/settings → 合并 Setting.Schema + UserSetting.Value
//  2. PUT /api/user/settings/{key} → Upsert 用户值（验证 Schema 存在）
//
// 依赖倒置：本包仅定义接口，实现位于 infrastructure/persistence/usersetting。
package usersetting
