package cache

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

// NewClient 创建并初始化 Redis 客户端
// redisURL 格式: redis://[:password@]host:port[/db]
// 例如: redis://localhost:6379/0 或 redis://:password@localhost:6379/1
// enableTracing: 是否启用 OpenTelemetry 追踪
func NewClient(ctx context.Context, redisURL string, enableTracing bool) (*redis.Client, error) {
	if redisURL == "" {
		return nil, errors.New("redis URL cannot be empty")
	}

	// 解析 Redis URL
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	// 创建客户端
	client := redis.NewClient(opts)

	// 启用 OpenTelemetry 追踪
	if enableTracing {
		if tracingErr := redisotel.InstrumentTracing(client); tracingErr != nil {
			slog.Warn("Failed to enable Redis tracing", "error", tracingErr)
		} else {
			slog.Info("Redis OpenTelemetry tracing enabled")
		}
	}

	// 使用超时上下文进行健康检查
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 测试连接
	if err := client.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	slog.Info("Redis client initialized successfully",
		"addr", opts.Addr,
		"db", opts.DB,
	)

	return client, nil
}

// Close 关闭 Redis 客户端连接
func Close(client *redis.Client) error {
	if client == nil {
		return nil
	}

	if err := client.Close(); err != nil {
		slog.Error("Failed to close redis client", "error", err)
		return err
	}

	slog.Info("Redis client closed successfully")
	return nil
}

// HealthCheck 检查 Redis 连接健康状态
func HealthCheck(ctx context.Context, client *redis.Client) error {
	if client == nil {
		return errors.New("redis client is nil")
	}

	checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := client.Ping(checkCtx).Err(); err != nil {
		return fmt.Errorf("redis health check failed: %w", err)
	}

	return nil
}
