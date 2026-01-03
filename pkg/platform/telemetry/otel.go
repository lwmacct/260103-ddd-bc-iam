package telemetry

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// Config OpenTelemetry 配置
type Config struct {
	// ServiceName 服务名称（必填）
	ServiceName string

	// ServiceVersion 服务版本
	ServiceVersion string

	// Environment 环境标识 (development, staging, production)
	Environment string

	// Enabled 是否启用追踪
	Enabled bool

	// ExporterType 导出器类型: "otlp", "stdout", "none"
	ExporterType string

	// OTLPEndpoint OTLP gRPC 端点 (如 "localhost:4317")
	OTLPEndpoint string

	// SampleRate 采样率 (0.0-1.0)，1.0 表示全部采样
	SampleRate float64
}

// ShutdownFunc 关闭函数类型
type ShutdownFunc func(context.Context) error

// InitTracer 初始化 OpenTelemetry Tracer
// 返回 shutdown 函数，应在程序退出时调用
func InitTracer(ctx context.Context, cfg Config) (ShutdownFunc, error) {
	if !cfg.Enabled {
		// 返回空操作的 shutdown 函数
		return func(context.Context) error { return nil }, nil
	}

	// 创建 resource
	res, err := newResource(cfg)
	if err != nil {
		return nil, err
	}

	// 创建 exporter
	exporter, err := newExporter(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// 创建 sampler
	sampler := newSampler(cfg.SampleRate)

	// 创建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	// 设置全局 TracerProvider
	otel.SetTracerProvider(tp)

	// 设置全局 Propagator（支持 W3C TraceContext 和 Baggage）
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// 返回 shutdown 函数
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return tp.Shutdown(ctx)
	}, nil
}

// newResource 创建 OpenTelemetry Resource
func newResource(cfg Config) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentName(cfg.Environment),
		),
	)
}

// newExporter 根据配置创建 Exporter
func newExporter(ctx context.Context, cfg Config) (sdktrace.SpanExporter, error) {
	switch cfg.ExporterType {
	case "otlp":
		return newOTLPExporter(ctx, cfg.OTLPEndpoint)
	case "stdout":
		return stdouttrace.New(stdouttrace.WithPrettyPrint())
	case "none", "":
		// 返回一个空的 exporter
		return &noopExporter{}, nil
	default:
		return nil, errors.New("unsupported exporter type: " + cfg.ExporterType)
	}
}

// newOTLPExporter 创建 OTLP gRPC Exporter
func newOTLPExporter(ctx context.Context, endpoint string) (sdktrace.SpanExporter, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithInsecure(), // 开发环境使用不安全连接
	}
	if endpoint != "" {
		opts = append(opts, otlptracegrpc.WithEndpoint(endpoint))
	}
	return otlptracegrpc.New(ctx, opts...)
}

// newSampler 根据采样率创建 Sampler
func newSampler(rate float64) sdktrace.Sampler {
	if rate <= 0 {
		return sdktrace.NeverSample()
	}
	if rate >= 1 {
		return sdktrace.AlwaysSample()
	}
	return sdktrace.TraceIDRatioBased(rate)
}

// noopExporter 空操作 Exporter（用于禁用导出）
type noopExporter struct{}

func (e *noopExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	return nil
}

func (e *noopExporter) Shutdown(ctx context.Context) error {
	return nil
}
