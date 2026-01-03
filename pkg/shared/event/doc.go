// Package event 定义领域事件机制。
//
// 本包提供事件驱动架构的核心接口，用于模块间解耦通信。
//
// # 核心接口
//
//   - [Event]: 领域事件基础接口
//   - [EventBus]: 事件总线接口（发布/订阅）
//   - [EventHandler]: 事件处理器接口
//
// # 事件定义
//
// 具体事件定义在 events 子包中：
//   - [events.UserCreatedEvent]: 用户创建事件
//   - [events.UserRoleAssignedEvent]: 用户角色分配事件
//   - [events.RolePermissionsChangedEvent]: 角色权限变更事件
//   - [events.LoginSucceededEvent]: 登录成功事件
//   - [events.LoginFailedEvent]: 登录失败事件
//
// # 使用示例
//
// 发布事件：
//
//	event := events.NewUserCreatedEvent(user.ID, user.Username)
//	if err := eventBus.Publish(ctx, event); err != nil {
//	    log.Error("failed to publish event", err)
//	}
//
// 订阅事件：
//
//	eventBus.Subscribe("user.created", &AuditLogHandler{})
//	eventBus.Subscribe("role.permissions_changed", &CacheInvalidationHandler{})
//
// # 同步 vs 异步
//
//   - 同步事件：立即处理，失败会影响主流程
//   - 异步事件：队列处理，主流程不受影响
//
// 实现位于 infrastructure/event 包：
//   - [InMemoryEventBus]: 内存事件总线（同步）
//   - [AsyncEventBus]: 异步事件总线（可选）
package event
