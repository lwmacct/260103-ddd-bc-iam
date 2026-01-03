// Package queue 提供基于 Redis 的异步任务队列实现。
//
// 本包提供生产者-消费者模式的队列处理能力，适用于：
//   - 异步任务处理（如发送邮件、生成报告）
//   - 解耦高延迟操作
//   - 削峰填谷
//
// # 组件
//
//   - [RedisQueue]: 基于 Redis List 的 FIFO 队列
//   - [Processor]: 并发工作池，支持多 worker 消费
//   - [JobHandler]: 任务处理器接口
//
// # 使用示例
//
//	// 创建队列
//	queue := queue.NewRedisQueue(redisClient, "my-queue")
//
//	// 入队任务
//	queue.Enqueue(ctx, map[string]any{"type": "email", "to": "user@example.com"})
//
//	// 创建处理器（4 个并发 worker）
//	processor := queue.NewProcessor(queue, myHandler, 4)
//
//	// 启动处理（阻塞）
//	go processor.Start(ctx)
//
//	// 优雅关闭
//	processor.Stop()
//
// # 队列特性
//
//   - FIFO 顺序：使用 Redis LPUSH/BRPOP 保证先进先出
//   - 阻塞消费：worker 空闲时阻塞等待，避免轮询
//   - 并发处理：支持配置多 worker 并行消费
//   - 优雅关闭：Stop() 等待所有 worker 完成当前任务
//
// # 限制
//
// 当前实现不支持：
//   - 失败重试（需自行在 JobHandler 中实现）
//   - 死信队列（失败任务直接丢弃）
//   - 延迟队列（需使用 Redis Sorted Set 扩展）
//   - 任务去重
//
// 对于复杂场景，请考虑使用专业消息队列（如 RabbitMQ、Kafka）。
package queue
