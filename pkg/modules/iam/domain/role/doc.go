// Package role 定义角色和权限领域模型。
//
// # Overview
//
// 本包是 RBAC（基于角色的访问控制）系统的核心领域层，定义了：
//   - [Role]: 角色实体，支持多对多关联权限
//   - [Permission]: 权限实体，采用 domain:resource:action 三段式格式
//   - [CommandRepository]: 写仓储接口（创建、更新、删除、权限分配）
//   - [QueryRepository]: 读仓储接口（查询、搜索）
//   - 角色领域错误（见 errors.go）
//
// 系统角色：
// [Role.IsSystem] 字段标识系统内置角色，系统角色：
//   - 不可删除（[Role.CanBeDeleted] 返回 false）
//   - 不可修改（[Role.CanBeModified] 返回 false）
//
// 权限管理：
// [Role] 实体通过 Permissions 字段关联权限，提供：
//   - [Role.HasPermission]: 检查是否拥有指定权限
//   - [Role.AddPermission]: 添加权限
//   - [Role.RemovePermission]: 移除权限
//
// # Usage
//
//	// 创建角色实体
//	role := &role.Role{
//	    Name:    "editor",
//	    IsSystem: false,
//	}
//
//	// 添加权限
//	perm := &role.Permission{
//	    Code: "sys:posts:write",
//	}
//	role.AddPermission(perm)
//
//	// 检查权限
//	if role.HasPermission("sys:posts:write") {
//	    // 允许操作
//	}
//
// # Thread Safety
//
// 角色和权限实体都是值类型，是并发安全的。
// Repository 接口的实现需要保证并发安全性（由基础设施层负责）。
//
// # 依赖关系
//
// 本包仅定义接口，实现位于 infrastructure/persistence 包。
package role
