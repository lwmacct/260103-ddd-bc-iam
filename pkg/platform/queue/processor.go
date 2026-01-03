package queue

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// JobHandler 任务处理器接口
type JobHandler interface {
	Handle(ctx context.Context, data []byte) error
}

// Processor 队列处理器
type Processor struct {
	queue       *RedisQueue
	handler     JobHandler
	concurrency int
	wg          sync.WaitGroup
	stopCh      chan struct{}
	stopped     bool
	mu          sync.Mutex
}

// NewProcessor 创建队列处理器
func NewProcessor(queue *RedisQueue, handler JobHandler, concurrency int) *Processor {
	if concurrency < 1 {
		concurrency = 1
	}
	return &Processor{
		queue:       queue,
		handler:     handler,
		concurrency: concurrency,
		stopCh:      make(chan struct{}),
	}
}

// Start 启动处理器
func (p *Processor) Start(ctx context.Context) {
	slog.Info("Starting queue processor", "concurrency", p.concurrency)

	for i := range p.concurrency {
		p.wg.Add(1)
		go p.worker(ctx, i)
	}

	p.wg.Wait()
	slog.Info("Queue processor stopped")
}

// Stop 停止处理器
func (p *Processor) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.stopped {
		return
	}

	slog.Info("Stopping queue processor...")
	close(p.stopCh)
	p.stopped = true
}

// worker 工作协程
func (p *Processor) worker(ctx context.Context, workerID int) {
	defer p.wg.Done()

	slog.Info("Worker started", "worker_id", workerID)

	for {
		select {
		case <-p.stopCh:
			slog.Info("Worker stopped", "worker_id", workerID)
			return
		case <-ctx.Done():
			slog.Info("Worker canceled", "worker_id", workerID)
			return
		default:
			p.processOne(ctx, workerID)
		}
	}
}

// processOne 处理一个任务
func (p *Processor) processOne(ctx context.Context, workerID int) {
	// 从队列中取出任务 (超时 5 秒)
	data, err := p.queue.Dequeue(ctx, 5*time.Second)
	if err != nil {
		slog.Error("Failed to dequeue job", "worker_id", workerID, "error", err)
		return
	}

	if data == nil {
		// 没有任务，继续等待
		return
	}

	// 处理任务
	slog.Info("Processing job", "worker_id", workerID, "data_size", len(data))

	if err := p.handler.Handle(ctx, data); err != nil {
		slog.Error("Failed to handle job", "worker_id", workerID, "error", err)
		// 这里可以实现重试逻辑或将失败的任务放入死信队列
		return
	}

	slog.Info("Job processed successfully", "worker_id", workerID)
}

// DefaultJobHandler 默认任务处理器 (仅记录日志)
type DefaultJobHandler struct{}

// Handle 处理任务
func (h *DefaultJobHandler) Handle(ctx context.Context, data []byte) error {
	slog.Info("Handling job", "data", string(data))
	// 模拟处理时间
	time.Sleep(100 * time.Millisecond)
	return nil
}
