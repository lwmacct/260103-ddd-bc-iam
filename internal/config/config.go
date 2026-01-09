// Package config 提供应用配置管理
package config

import (
	"strings"
	"time"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/config"
)

// Server 服务器配置
type Server struct {
	Addr         string `koanf:"addr" desc:"监听地址，格式: host:port，例如 '0.0.0.0:8080' 或 ':8080'"`
	Env          string `koanf:"env" desc:"运行环境: development | production"`
	WebDist      string `koanf:"web-dist" desc:"静态资源目录路径，用于提供前端文件服务 (如 SPA 应用)"`
	DocsDist     string `koanf:"docs-dist" desc:"文档目录路径，用于提供 VitePress 构建的文档服务，通过 /docs 路由访问"`
	FxLogEnabled bool   `koanf:"fx-log-enabled" desc:"是否显示 Fx 框架的依赖注入日志 (默认 false，减少启动输出噪音)"`
}

// Data 数据源配置
type Data struct {
	PgsqlURL       string `koanf:"pgsql-url" desc:"PostgreSQL 连接 URL，格式: postgresql://user:password@host:port/dbname?sslmode=disable"`
	RedisURL       string `koanf:"redis-url" desc:"Redis 连接 URL，格式: redis://:password@host:port/db"`
	RedisKeyPrefix string `koanf:"redis-key-prefix" desc:"Redis key 前缀，所有 key 读写都会自动拼接此前缀，例如 'app:'"`
	AutoMigrate    bool   `koanf:"auto-migrate" desc:"是否在应用启动时自动执行数据库迁移 (仅推荐在开发环境使用，生产环境应使用 migrate 命令)"`
}

// Telemetry OpenTelemetry 追踪配置
type Telemetry struct {
	Enabled      bool    `koanf:"enabled" desc:"是否启用分布式追踪"`
	ExporterType string  `koanf:"exporter-type" desc:"导出器类型: otlp, stdout, none"`
	OTLPEndpoint string  `koanf:"otlp-endpoint" desc:"OTLP gRPC 端点 (如 localhost:4317)"`
	SampleRate   float64 `koanf:"sample-rate" desc:"采样率 (0.0-1.0)，1.0 表示全部采样"`
}

// Config 应用配置
type Config struct {
	Server    Server       `koanf:"server" desc:"服务器配置"`
	Data      Data         `koanf:"data" desc:"数据源配置"`
	JWT       config.JWT   `koanf:"jwt" desc:"JWT 认证配置"`
	Auth      config.Auth  `koanf:"auth" desc:"认证配置"`
	Telemetry Telemetry    `koanf:"telemetry" desc:"OpenTelemetry 追踪配置"`
	Redis     config.Redis `koanf:"redis" desc:"Redis 配置"`
}

// GetBaseUrl 返回服务的基础URL
// 注意：当 Addr 为 0.0.0.0 时，自动替换为 localhost（客户端无法连接到0.0.0.0）
func (c *Config) GetBaseUrl(https bool) string {
	addr := c.Server.Addr
	// 0.0.0.0 是服务器监听地址，客户端应该连接到 localhost
	if strings.HasPrefix(addr, "0.0.0.0:") {
		addr = "localhost" + addr[7:] // 保留端口部分
	}

	if https {
		return "https://" + addr
	}
	return "http://" + addr
}

// DefaultConfig 返回默认配置
// 注意：internal/command/command.go 中的 Defaults 变量引用此函数以实现单一配置来源。
func DefaultConfig() Config {
	return Config{
		Server: Server{
			Addr:         "0.0.0.0:8080",
			Env:          "development",
			WebDist:      "dist",
			DocsDist:     "docs/.vitepress/dist",
			FxLogEnabled: false, // 默认关闭 Fx 日志，减少启动输出噪音
		},
		Data: Data{
			PgsqlURL:       "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable",
			RedisURL:       "redis://localhost:6379/0",
			RedisKeyPrefix: "app:",
			AutoMigrate:    false, // 默认关闭自动迁移，生产环境使用 migrate 命令
		},
		JWT: config.JWT{
			Secret:             "change-me-in-production",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		},
		Auth: config.Auth{
			DevSecret:       "dev-secret-change-me",
			TwoFAIssuer:     "Go-DDD-Package-Lib",
			CaptchaRequired: true, // 默认开启验证码
		},
		Telemetry: Telemetry{
			Enabled:      false,  // 默认关闭，按需开启
			ExporterType: "none", // 默认不导出
			OTLPEndpoint: "localhost:4317",
			SampleRate:   1.0, // 默认全部采样
		},
		Redis: config.Redis{
			KeyPrefix: "app:",
		},
	}
}
