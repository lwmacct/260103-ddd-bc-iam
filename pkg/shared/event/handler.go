package event

import "context"

// EventHandler 事件处理器接口
type EventHandler interface {
	// Handle 处理事件
	// 返回错误时：
	//   - 同步事件总线：可能影响主流程
	//   - 异步事件总线：记录错误并继续
	Handle(ctx context.Context, event Event) error
}

// EventHandlerFunc 函数类型的事件处理器
// 便于使用闭包作为处理器
type EventHandlerFunc func(ctx context.Context, event Event) error

func (f EventHandlerFunc) Handle(ctx context.Context, event Event) error {
	return f(ctx, event)
}

// EventHandlerMiddleware 事件处理器中间件
// 用于添加日志、监控、重试等横切关注点
type EventHandlerMiddleware func(next EventHandler) EventHandler

// ChainMiddlewares 链接多个中间件
func ChainMiddlewares(handler EventHandler, middlewares ...EventHandlerMiddleware) EventHandler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// LoggingMiddleware 日志中间件示例
// 记录事件处理的开始和结束
type LoggingMiddleware struct {
	Logger interface {
		Info(msg string, args ...any)
		Error(msg string, args ...any)
	}
}

func (m *LoggingMiddleware) Wrap(next EventHandler) EventHandler {
	return EventHandlerFunc(func(ctx context.Context, event Event) error {
		m.Logger.Info("handling event",
			"event", event.EventName(),
			"aggregate_id", event.AggregateID(),
		)

		err := next.Handle(ctx, event)

		if err != nil {
			m.Logger.Error("event handling failed",
				"event", event.EventName(),
				"error", err,
			)
		}

		return err
	})
}

// RetryMiddleware 重试中间件（占位）
// 可用于异步事件的重试逻辑
type RetryMiddleware struct {
	MaxRetries int
}

func (m *RetryMiddleware) Wrap(next EventHandler) EventHandler {
	return EventHandlerFunc(func(ctx context.Context, event Event) error {
		var lastErr error
		for i := 0; i <= m.MaxRetries; i++ {
			if err := next.Handle(ctx, event); err != nil {
				lastErr = err
				continue
			}
			return nil
		}
		return lastErr
	})
}
