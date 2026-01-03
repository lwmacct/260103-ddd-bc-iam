// Package user 实现用户管理的应用层用例。
//
// 本包提供 CQRS 模式的 Command 和 Query Handler：
//
// # Command（写操作）
//
//   - [command.CreateHandler]: 创建用户
//   - [command.UpdateHandler]: 更新用户信息
//   - [command.DeleteHandler]: 删除用户
//   - [command.AssignRolesHandler]: 分配角色
//   - [command.ChangePasswordHandler]: 修改密码
//   - [command.BatchCreateHandler]: 批量创建用户
//
// # Query（读操作）
//
//   - [query.GetHandler]: 获取用户详情（支持角色关联）
//   - [query.ListHandler]: 用户列表查询（支持分页、搜索）
//
// # DTO 与映射
//
// 请求 DTO：
//   - [CreateUserDTO]: 创建用户请求
//   - [UpdateUserDTO]: 更新用户请求
//   - [ChangePasswordDTO]: 修改密码请求
//   - [AssignRolesDTO]: 分配角色请求
//   - [BatchCreateUserDTO]: 批量创建用户请求
//
// 响应 DTO：
//   - [UserResponse]: 用户基本信息响应
//   - [UserWithRolesResponse]: 用户详情响应（含角色）
//   - [BatchCreateUserResponse]: 批量创建结果响应
//
// 映射函数：
//   - [ToUserResponse]: User -> UserResponse
//   - [ToUserWithRolesResponse]: User -> UserWithRolesResponse
//
// 依赖注入：所有 Handler 通过 [bootstrap.Container] 注册。
package user
