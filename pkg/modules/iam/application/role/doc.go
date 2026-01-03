// Package role 实现角色管理的应用层用例。
//
// 本包提供 CQRS 模式的 Command 和 Query Handler：
//
// # Command（写操作）
//
//   - [CreateHandler]: 创建角色
//   - [UpdateHandler]: 更新角色信息
//   - [DeleteHandler]: 删除角色
//   - [SetPermissionsHandler]: 设置角色权限
//
// # Query（读操作）
//
//   - [GetHandler]: 获取角色详情（含权限）
//   - [ListHandler]: 角色列表查询
//   - [ListPermissionsHandler]: 权限列表查询
//
// # DTO 与映射
//
// 请求 DTO：
//   - [CreateDTO]: 创建角色请求
//   - [UpdateDTO]: 更新角色请求
//   - [SetPermissionsDTO]: 设置权限请求
//
// 响应 DTO：
//   - [RoleResponse]: 角色信息响应
//   - [RoleWithPermissionsResponse]: 角色详情响应（含权限）
//   - [PermissionResponse]: 权限信息响应
//
// 映射函数：
//   - [ToRoleResponse]: Role -> RoleResponse
//   - [ToRoleWithPermissionsResponse]: Role -> RoleWithPermissionsResponse
//
// 系统角色保护：
// 系统内置角色（IsSystem=true）不可删除或修改，由领域层强制约束。
//
// 依赖注入：所有 Handler 通过 [bootstrap.Container] 注册。
package role
