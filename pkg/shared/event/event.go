package event

import "time"

// Event 领域事件接口
// 所有领域事件都应实现此接口
type Event interface {
	// EventName 返回事件名称，用于订阅匹配
	// 格式：{aggregate}.{action}，如 "user.created"
	EventName() string

	// OccurredAt 返回事件发生时间
	OccurredAt() time.Time

	// AggregateID 返回聚合根 ID
	AggregateID() string

	// AggregateType 返回聚合类型
	// 如 "user"、"role"、"menu"
	AggregateType() string
}

// BaseEvent 事件基础结构
// 提供 Event 接口的通用实现
type BaseEvent struct {
	Name        string
	Timestamp   time.Time
	AggregateId string
	Aggregate   string
}

// NewBaseEvent 创建基础事件
func NewBaseEvent(name, aggregateType, aggregateID string) BaseEvent {
	return BaseEvent{
		Name:        name,
		Timestamp:   time.Now(),
		AggregateId: aggregateID,
		Aggregate:   aggregateType,
	}
}

func (e BaseEvent) EventName() string     { return e.Name }
func (e BaseEvent) OccurredAt() time.Time { return e.Timestamp }
func (e BaseEvent) AggregateID() string   { return e.AggregateId }
func (e BaseEvent) AggregateType() string { return e.Aggregate }

// EventWithPayload 带有效载荷的事件
type EventWithPayload interface {
	Event
	Payload() any
}
