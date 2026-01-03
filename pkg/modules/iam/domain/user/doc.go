// Package user 定义用户领域模型和仓储接口。
//
// 本包是用户管理的领域层核心，定义了：
//   - [User]: 用户实体（富领域模型，包含 RBAC 角色关联）
//   - [CommandRepository]: 写仓储接口（创建、更新、删除、角色分配）
//   - [QueryRepository]: 读仓储接口（查询、搜索、统计）
//   - 用户领域错误（见 errors.go）
//
// 用户状态：
// 用户通过 Status 字段管理生命周期状态：
//   - active: 正常状态，可登录
//   - inactive: 未激活状态
//   - banned: 禁用状态
//
// RBAC 集成：
// [User] 实体通过 Roles 字段关联 [role.Role]，提供：
//   - [User.HasRole]: 检查用户是否拥有指定角色
//   - [User.HasPermission]: 检查用户是否拥有指定权限
//   - [User.GetPermissionCodes]: 获取用户所有权限代码
//
// 依赖倒置：
// 本包仅定义接口，实现位于 infrastructure/persistence 包。
package user
