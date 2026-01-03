// Package telemetry 提供 OpenTelemetry 分布式追踪支持。
//
// 本包负责初始化和配置 OpenTelemetry SDK，是应用的横切关注点（Cross-Cutting Concern）。
//
// # 功能特性
//
//   - OTLP gRPC 导出：支持 Jaeger、Tempo、SigNoz 等后端
//   - Stdout 导出：开发调试用
//   - 自动传播 trace context（W3C TraceContext + Baggage）
//   - 可配置采样率
//
// # 支持的 Exporter
//
//   - "otlp": 生产环境，导出到 OTLP 兼容后端
//   - "stdout": 开发调试，输出到控制台
//   - "none": 禁用导出（保留 API 兼容性）
//
// # 使用示例
//
//	cfg := telemetry.Config{
//	    ServiceName:  "my-service",
//	    Enabled:      true,
//	    ExporterType: "otlp",
//	    OTLPEndpoint: "localhost:4317",
//	    SampleRate:   1.0,
//	}
//
//	shutdown, err := telemetry.InitTracer(ctx, cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer shutdown(ctx) // 5 秒超时，确保 span 导出完成
//
// # 初始化顺序
//
// Telemetry 应在所有其他基础设施之前初始化，确保追踪完整：
//
//	Telemetry → Database → Redis → EventBus → Repositories → ...
//
// # 与其他模块集成
//
// Database 和 Redis 通过配置启用追踪：
//
//	dbConfig.EnableTracing = cfg.Telemetry.Enabled
//	redis.NewClient(ctx, url, cfg.Telemetry.Enabled)
//
// HTTP 层通过 OTEL Gin 中间件自动追踪请求。
package telemetry
