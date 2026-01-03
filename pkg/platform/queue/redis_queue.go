package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisQueue 基于 Redis 的队列实现
type RedisQueue struct {
	client    *redis.Client
	queueName string
}

// NewRedisQueue 创建 Redis 队列
func NewRedisQueue(client *redis.Client, queueName string) *RedisQueue {
	return &RedisQueue{
		client:    client,
		queueName: queueName,
	}
}

// Enqueue 将任务入队
func (q *RedisQueue) Enqueue(ctx context.Context, job any) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	if err := q.client.LPush(ctx, q.queueName, data).Err(); err != nil {
		return fmt.Errorf("failed to enqueue job: %w", err)
	}

	return nil
}

// Dequeue 从队列中取出一个任务 (阻塞式)
func (q *RedisQueue) Dequeue(ctx context.Context, timeout time.Duration) ([]byte, error) {
	result, err := q.client.BRPop(ctx, timeout, q.queueName).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // 超时，没有任务
		}
		return nil, fmt.Errorf("failed to dequeue job: %w", err)
	}

	if len(result) < 2 {
		return nil, nil
	}

	return []byte(result[1]), nil
}

// Length 获取队列长度
func (q *RedisQueue) Length(ctx context.Context) (int64, error) {
	length, err := q.client.LLen(ctx, q.queueName).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get queue length: %w", err)
	}
	return length, nil
}

// Clear 清空队列
func (q *RedisQueue) Clear(ctx context.Context) error {
	if err := q.client.Del(ctx, q.queueName).Err(); err != nil {
		return fmt.Errorf("failed to clear queue: %w", err)
	}
	return nil
}
