// Package eventhandler 提供领域事件处理器实现。
//
// # Overview
//
// 本包实现领域事件的处理逻辑，订阅业务事件并执行副作用操作：
//   - [AuditEventHandler]: 审计日志事件处理器
//
// # 事件类型
//
// AuditEventHandler 处理以下事件：
//   - CommandExecutedEvent: 命令执行事件
//   - LoginSucceededEvent: 登录成功事件
//   - LoginFailedEvent: 登录失败事件
//   - UserCreatedEvent: 用户创建事件
//   - UserDeletedEvent: 用户删除事件
//   - UserRoleAssignedEvent: 角色分配事件
//   - RolePermissionsChangedEvent: 权限变更事件
//
// # 错误处理
//
// 事件处理器失败不阻塞业务流程：
//   - 错误仅记录日志，不返回给调用方
//   - 保证业务流程的完整性
//   - 丢失的审计日志不影响核心功能
//
// # Thread Safety
//
// 事件处理器是无状态的，仅依赖注入的 Repository（通过 Fx 管理）。
// Handler 本身是并发安全的，可以安全地在多个 goroutine 中调用。
//
// # 依赖关系
//
// 本包依赖：
//   - Domain 层的 audit.CommandRepository: 写入审计日志
//   - Shared 层的 event 包: 事件类型定义
package eventhandler
