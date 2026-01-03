package eventbus

import (
	"context"
	"log/slog"
	"strings"
	"sync"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/shared/event"
)

// InMemoryEventBus 内存事件总线实现
// 同步处理事件，适用于简单场景
type InMemoryEventBus struct {
	handlers map[string][]event.EventHandler
	mu       sync.RWMutex
}

// NewInMemoryEventBus 创建内存事件总线
func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]event.EventHandler),
	}
}

// Publish 发布事件
// 同步执行所有匹配的处理器
func (b *InMemoryEventBus) Publish(ctx context.Context, events ...event.Event) error {
	for _, e := range events {
		if err := b.publishOne(ctx, e); err != nil {
			return err
		}
	}
	return nil
}

// Subscribe 订阅事件
func (b *InMemoryEventBus) Subscribe(eventName string, handler event.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

// Unsubscribe 取消订阅
func (b *InMemoryEventBus) Unsubscribe(eventName string, handler event.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers := b.handlers[eventName]
	for i, h := range handlers {
		// 通过指针比较（简单实现）
		if &h == &handler {
			b.handlers[eventName] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

// Close 关闭事件总线
func (b *InMemoryEventBus) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers = make(map[string][]event.EventHandler)
	return nil
}

// HandlerCount 返回订阅的处理器数量（用于测试）
func (b *InMemoryEventBus) HandlerCount(eventName string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return len(b.handlers[eventName])
}

// publishOne 发布单个事件
func (b *InMemoryEventBus) publishOne(ctx context.Context, e event.Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	eventName := e.EventName()

	// 收集所有匹配的处理器
	var matchedHandlers []event.EventHandler

	// 精确匹配
	if handlers, ok := b.handlers[eventName]; ok {
		matchedHandlers = append(matchedHandlers, handlers...)
	}

	// 通配符匹配：{aggregate}.*
	if idx := strings.LastIndex(eventName, "."); idx > 0 {
		wildcardKey := eventName[:idx] + ".*"
		if handlers, ok := b.handlers[wildcardKey]; ok {
			matchedHandlers = append(matchedHandlers, handlers...)
		}
	}

	// 全局通配符匹配：*
	if handlers, ok := b.handlers["*"]; ok {
		matchedHandlers = append(matchedHandlers, handlers...)
	}

	// 执行所有处理器（错误不中断后续处理器）
	var lastErr error
	for _, handler := range matchedHandlers {
		if err := handler.Handle(ctx, e); err != nil {
			slog.Error("event handler failed", "event", e.EventName(), "error", err.Error())
			lastErr = err
			// 继续执行后续处理器
		}
	}

	return lastErr
}

// ============================================================================
// 空事件总线（用于禁用事件或测试）
// ============================================================================

// NoOpEventBus 空操作事件总线
// 所有操作都不执行，用于禁用事件或测试
type NoOpEventBus struct{}

func (b *NoOpEventBus) Publish(ctx context.Context, events ...event.Event) error { return nil }
func (b *NoOpEventBus) Subscribe(eventName string, handler event.EventHandler)   {}
func (b *NoOpEventBus) Unsubscribe(eventName string, handler event.EventHandler) {}
func (b *NoOpEventBus) Close() error                                             { return nil }

// Ensure interfaces are implemented
var (
	_ event.EventBus = (*InMemoryEventBus)(nil)
	_ event.EventBus = (*NoOpEventBus)(nil)
)
