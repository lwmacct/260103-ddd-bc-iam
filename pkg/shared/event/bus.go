package event

import "context"

// EventBus 事件总线接口
// 负责事件的发布和订阅
type EventBus interface {
	// Publish 发布一个或多个事件
	// 同步事件总线会立即执行所有处理器
	// 异步事件总线会将事件放入队列
	Publish(ctx context.Context, events ...Event) error

	// Subscribe 订阅事件
	// eventName 支持通配符：
	//   - "user.created": 精确匹配
	//   - "user.*": 匹配 user 下所有事件
	//   - "*": 匹配所有事件
	Subscribe(eventName string, handler EventHandler)

	// Unsubscribe 取消订阅
	Unsubscribe(eventName string, handler EventHandler)

	// Close 关闭事件总线，释放资源
	Close() error
}

// EventPublisher 事件发布者接口
// 仅提供发布能力，用于 Command Handler
type EventPublisher interface {
	Publish(ctx context.Context, events ...Event) error
}

// EventSubscriber 事件订阅者接口
// 仅提供订阅能力
type EventSubscriber interface {
	Subscribe(eventName string, handler EventHandler)
	Unsubscribe(eventName string, handler EventHandler)
}
