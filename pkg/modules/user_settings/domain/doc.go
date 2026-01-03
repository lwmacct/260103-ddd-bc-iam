// Package domain 定义用户配置领域模型和仓储接口。
//
// # Overview
//
// 本包是 User Settings Bounded Context 的领域层核心，定义了：
//   - [UserSetting]: 用户配置实体（存储用户自定义覆盖值）
//   - [CommandRepository]: 写仓储接口（创建、更新、删除、批量操作）
//   - [QueryRepository]: 读仓储接口（查询用户自定义值）
//   - 领域错误（见 errors.go）
//
// # 设计说明
//
// UserSetting 是用户对系统配置的自定义覆盖值：
//   - 只存储用户修改过的配置（覆盖模式）
//   - 删除配置即恢复系统默认值
//   - 依赖 Settings BC 获取配置定义进行校验
//
// # 跨 BC 依赖
//
//	User Settings BC ──依赖──> Settings BC (QueryRepository)
//	                          ├─ 校验 key 是否存在
//	                          ├─ 类型校验（ValueType）
//	                          └─ 格式校验（InputType）
//
// # Thread Safety
//
// 实体是值对象风格，不包含内部状态同步机制。
// 并发安全性由仓储实现和调用方保证。
//
// # 依赖倒置
//
// 本包仅定义接口，实现位于 infra/persistence 包。
package domain
