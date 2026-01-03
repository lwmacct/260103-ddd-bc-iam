// Package eventbus 提供内存事件总线实现。
//
// 本包实现 domain/event.EventBus 接口，提供同步事件发布/订阅功能。
//
// # 特性
//
//   - 同步事件处理（适用于简单场景）
//   - 支持通配符订阅：
//   - "*" 订阅所有事件
//   - "{aggregate}.*" 订阅指定聚合根的所有事件（如 "user.*"）
//   - 线程安全（使用 sync.RWMutex）
//
// # 使用示例
//
//	bus := eventbus.NewInMemoryEventBus()
//
//	// 订阅特定事件
//	bus.Subscribe("user.created", handler)
//
//	// 订阅所有用户事件
//	bus.Subscribe("user.*", userEventHandler)
//
//	// 订阅所有事件
//	bus.Subscribe("*", auditLogHandler)
//
//	// 发布事件
//	bus.Publish(ctx, event.NewUserCreatedEvent(userID))
//
// # 限制
//
// 本实现为同步模式，适用于简单场景。对于高吞吐量场景，
// 请考虑使用基于消息队列（如 Redis/RabbitMQ）的异步实现。
package eventbus
