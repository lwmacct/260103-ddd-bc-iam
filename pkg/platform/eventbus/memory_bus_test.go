package eventbus

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/shared/event"
)

// mockEvent 测试用事件
type mockEvent struct {
	event.BaseEvent

	Data string
}

// mockHandler 测试用事件处理器
type mockHandler struct {
	mu            sync.Mutex
	handledEvents []event.Event
	handleError   error
}

func newMockHandler() *mockHandler {
	return &mockHandler{
		handledEvents: make([]event.Event, 0),
	}
}

func (h *mockHandler) Handle(ctx context.Context, e event.Event) error {
	if h.handleError != nil {
		return h.handleError
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handledEvents = append(h.handledEvents, e)
	return nil
}

func (h *mockHandler) HandledCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.handledEvents)
}

func TestInMemoryEventBus_Publish(t *testing.T) {
	ctx := context.Background()

	t.Run("发布事件到订阅者", func(t *testing.T) {
		bus := NewInMemoryEventBus()
		handler := newMockHandler()

		bus.Subscribe("test.event", handler)

		evt := &mockEvent{
			BaseEvent: event.NewBaseEvent("test.event", "test", "1"),
			Data:      "test data",
		}

		err := bus.Publish(ctx, evt)

		require.NoError(t, err)
		assert.Equal(t, 1, handler.HandledCount())
	})

	t.Run("发布多个事件", func(t *testing.T) {
		bus := NewInMemoryEventBus()
		handler := newMockHandler()

		bus.Subscribe("test.event", handler)

		evt1 := &mockEvent{BaseEvent: event.NewBaseEvent("test.event", "test", "1")}
		evt2 := &mockEvent{BaseEvent: event.NewBaseEvent("test.event", "test", "2")}

		err := bus.Publish(ctx, evt1, evt2)

		require.NoError(t, err)
		assert.Equal(t, 2, handler.HandledCount())
	})

	t.Run("无订阅者时发布不报错", func(t *testing.T) {
		bus := NewInMemoryEventBus()

		evt := &mockEvent{BaseEvent: event.NewBaseEvent("test.event", "test", "1")}

		err := bus.Publish(ctx, evt)

		require.NoError(t, err)
	})

	t.Run("处理器错误传播", func(t *testing.T) {
		bus := NewInMemoryEventBus()
		handler := newMockHandler()
		handler.handleError = errors.New("handler error")

		bus.Subscribe("test.event", handler)

		evt := &mockEvent{BaseEvent: event.NewBaseEvent("test.event", "test", "1")}

		err := bus.Publish(ctx, evt)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "handler error")
	})
}

func TestInMemoryEventBus_WildcardSubscription(t *testing.T) {
	ctx := context.Background()

	t.Run("通配符匹配 aggregate.*", func(t *testing.T) {
		bus := NewInMemoryEventBus()
		handler := newMockHandler()

		bus.Subscribe("user.*", handler)

		evt1 := &mockEvent{BaseEvent: event.NewBaseEvent("user.created", "user", "1")}
		evt2 := &mockEvent{BaseEvent: event.NewBaseEvent("user.deleted", "user", "2")}

		_ = bus.Publish(ctx, evt1)
		_ = bus.Publish(ctx, evt2)

		assert.Equal(t, 2, handler.HandledCount())
	})

	t.Run("全局通配符 *", func(t *testing.T) {
		bus := NewInMemoryEventBus()
		handler := newMockHandler()

		bus.Subscribe("*", handler)

		evt1 := &mockEvent{BaseEvent: event.NewBaseEvent("user.created", "user", "1")}
		evt2 := &mockEvent{BaseEvent: event.NewBaseEvent("role.updated", "role", "1")}

		_ = bus.Publish(ctx, evt1)
		_ = bus.Publish(ctx, evt2)

		assert.Equal(t, 2, handler.HandledCount())
	})

	t.Run("精确匹配优先于通配符", func(t *testing.T) {
		bus := NewInMemoryEventBus()
		exactHandler := newMockHandler()
		wildcardHandler := newMockHandler()

		bus.Subscribe("user.created", exactHandler)
		bus.Subscribe("user.*", wildcardHandler)

		evt := &mockEvent{BaseEvent: event.NewBaseEvent("user.created", "user", "1")}

		_ = bus.Publish(ctx, evt)

		// 两个处理器都应该被调用
		assert.Equal(t, 1, exactHandler.HandledCount())
		assert.Equal(t, 1, wildcardHandler.HandledCount())
	})
}

func TestInMemoryEventBus_Subscribe(t *testing.T) {
	t.Run("订阅后处理器计数增加", func(t *testing.T) {
		bus := NewInMemoryEventBus()
		handler := newMockHandler()

		assert.Equal(t, 0, bus.HandlerCount("test.event"))

		bus.Subscribe("test.event", handler)

		assert.Equal(t, 1, bus.HandlerCount("test.event"))
	})

	t.Run("同一事件多个处理器", func(t *testing.T) {
		bus := NewInMemoryEventBus()
		handler1 := newMockHandler()
		handler2 := newMockHandler()

		bus.Subscribe("test.event", handler1)
		bus.Subscribe("test.event", handler2)

		assert.Equal(t, 2, bus.HandlerCount("test.event"))
	})
}

func TestInMemoryEventBus_Close(t *testing.T) {
	t.Run("关闭后清空处理器", func(t *testing.T) {
		bus := NewInMemoryEventBus()
		handler := newMockHandler()

		bus.Subscribe("test.event", handler)
		assert.Equal(t, 1, bus.HandlerCount("test.event"))

		err := bus.Close()

		require.NoError(t, err)
		assert.Equal(t, 0, bus.HandlerCount("test.event"))
	})
}

func TestNoOpEventBus(t *testing.T) {
	ctx := context.Background()

	t.Run("所有操作都不执行", func(t *testing.T) {
		bus := &NoOpEventBus{}
		handler := newMockHandler()

		bus.Subscribe("test.event", handler)

		evt := &mockEvent{BaseEvent: event.NewBaseEvent("test.event", "test", "1")}
		err := bus.Publish(ctx, evt)

		require.NoError(t, err)
		assert.Equal(t, 0, handler.HandledCount()) // 不会调用处理器
	})

	t.Run("Close 不报错", func(t *testing.T) {
		bus := &NoOpEventBus{}
		err := bus.Close()
		require.NoError(t, err)
	})
}

func TestInMemoryEventBus_Concurrency(t *testing.T) {
	t.Run("并发发布和订阅", func(t *testing.T) {
		bus := NewInMemoryEventBus()
		var counter int64

		handler := &countingHandler{counter: &counter}
		bus.Subscribe("test.event", handler)

		var wg sync.WaitGroup
		for range 100 {
			wg.Go(func() {
				evt := &mockEvent{BaseEvent: event.NewBaseEvent("test.event", "test", "1")}
				_ = bus.Publish(context.Background(), evt)
			})
		}

		wg.Wait()

		assert.Equal(t, int64(100), atomic.LoadInt64(&counter))
	})
}

// countingHandler 用于并发测试的处理器
type countingHandler struct {
	counter *int64
}

func (h *countingHandler) Handle(ctx context.Context, e event.Event) error {
	atomic.AddInt64(h.counter, 1)
	return nil
}
