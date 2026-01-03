// Package eventhandler 提供领域事件处理器实现。
//
// 本包包含处理领域事件的各种处理器：
//   - [CacheInvalidationHandler]: 缓存失效处理器
//   - [AuditLogHandler]: 审计日志处理器
//
// # 缓存失效
//
// 当以下事件发生时，相关缓存会被自动失效：
//   - user.role_assigned: 用户权限缓存失效
//   - user.deleted: 删除用户的权限缓存
//   - role.permissions_changed: 拥有该角色的所有用户权限缓存失效
//
// # 审计日志
//
// 以下事件会自动记录审计日志：
//   - audit.command_executed: 通用命令执行审计
//   - auth.login_succeeded: 登录成功
//   - auth.login_failed: 登录失败
//   - user.created: 用户创建
//   - user.deleted: 用户删除
//   - user.role_assigned: 用户角色分配
//   - role.permissions_changed: 角色权限变更
//
// # 使用示例
//
//	// 缓存失效处理器
//	eventBus.Subscribe("user.role_assigned", cacheHandler)
//	eventBus.Subscribe("role.permissions_changed", cacheHandler)
//
//	// 审计日志处理器（使用通配符订阅所有事件）
//	eventBus.Subscribe("*", auditHandler)
package eventhandler
