// Package container 提供 Fx 模块定义，用于依赖注入。
//
// 每个模块聚合相关的 Provider 和生命周期钩子：
//   - [InfraModule]: 数据库、Redis、EventBus、链路追踪
//   - [CacheModule]: 所有缓存服务
//   - [RepositoryModule]: 带缓存装饰器的 CQRS 仓储
//   - [ServiceModule]: 领域服务和基础设施服务
//   - [UseCaseModule]: 应用层用例处理器
//   - [HTTPModule]: HTTP 处理器和路由
package container

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/gorm"

	internalConfig "github.com/lwmacct/260103-ddd-bc-iam/internal/config"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/config"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/cache"
	dbpkg "github.com/lwmacct/260103-ddd-shared/pkg/platform/db"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/eventbus"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/telemetry"
	"github.com/lwmacct/260103-ddd-shared/pkg/shared/event"
)

// InfraModule 提供基础设施组件。
//
// 组件：
//   - OpenTelemetry 链路追踪
//   - PostgreSQL 数据库连接
//   - Redis 客户端
//   - 内存事件总线
//
// 生命周期：
//   - OnStart: 初始化链路追踪、连接数据库、连接 Redis
//   - OnStop: 按相反顺序关闭连接
var InfraModule = fx.Module("infra",
	fx.Provide(
		newTelemetry,
		newDatabase,
		newRedisClient,
		newEventBus,
		newIAMConfig,
		newSettingsCacheService,
	),
)

// TelemetryResult 包装链路追踪的关闭函数，用于 Fx 生命周期管理。
type TelemetryResult struct {
	fx.Out

	Shutdown telemetry.ShutdownFunc
}

func newTelemetry(lc fx.Lifecycle, cfg *internalConfig.Config) (TelemetryResult, error) {
	ctx := context.Background()
	shutdown, err := telemetry.InitTracer(ctx, telemetry.Config{
		ServiceName:    "go-ddd-pkg-lib",
		ServiceVersion: "1.0.0",
		Environment:    cfg.Server.Env,
		Enabled:        cfg.Telemetry.Enabled,
		ExporterType:   cfg.Telemetry.ExporterType,
		OTLPEndpoint:   cfg.Telemetry.OTLPEndpoint,
		SampleRate:     cfg.Telemetry.SampleRate,
	})
	if err != nil {
		return TelemetryResult{}, err
	}

	if cfg.Telemetry.Enabled {
		slog.Info("OpenTelemetry tracing initialized",
			"exporter", cfg.Telemetry.ExporterType,
			"sample_rate", cfg.Telemetry.SampleRate)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if shutdown != nil {
				if err := shutdown(ctx); err != nil {
					slog.Error("failed to shutdown telemetry", "error", err)
					return err
				}
				slog.Info("OpenTelemetry shutdown completed")
			}
			return nil
		},
	})

	return TelemetryResult{Shutdown: shutdown}, nil
}

func newDatabase(lc fx.Lifecycle, cfg *internalConfig.Config) (*gorm.DB, error) {
	ctx := context.Background()
	dbConfig := dbpkg.DefaultConfig(cfg.Data.PgsqlURL)
	dbConfig.EnableTracing = cfg.Telemetry.Enabled

	db, err := dbpkg.NewConnection(ctx, dbConfig)
	if err != nil {
		slog.Error("Failed to connect to database",
			"error", err,
			"hint", "Ensure PostgreSQL is running and APP_PGSQL_URL is correct",
		)
		return nil, err
	}

	// 如果启用了自动迁移则执行
	if cfg.Data.AutoMigrate {
		if err := runAutoMigrate(db); err != nil {
			return nil, err
		}
	} else {
		slog.Info("Auto-migration disabled, skipping database migration",
			"hint", "Set APP_AUTO_MIGRATE=true to enable",
		)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return dbpkg.Close(db)
		},
	})

	return db, nil
}

// runAutoMigrate 执行数据库自动迁移和索引创建。
func runAutoMigrate(db *gorm.DB) error {
	slog.Info("Auto-migration enabled, migrating database...")

	if err := db.AutoMigrate(GetAllModels()...); err != nil {
		return err
	}

	// 为多对多关联表创建索引
	// role_permissions 使用复合主键，PostgreSQL 自动利用前缀索引
	if err := dbpkg.CreateJoinTableIndexes(db, []dbpkg.JoinTableIndex{
		{Table: "user_roles", Name: "idx_user_roles_user_id", Columns: "user_id"},
		{Table: "user_roles", Name: "idx_user_roles_role_id", Columns: "role_id"},
		{Table: "role_permissions", Name: "idx_role_permissions_role_model_id", Columns: "role_model_id"},
		{Table: "role_permissions", Name: "idx_role_permissions_permission_model_id", Columns: "permission_model_id"},
	}); err != nil {
		return err
	}

	slog.Info("Database migration completed")
	return nil
}

func newRedisClient(lc fx.Lifecycle, cfg *internalConfig.Config) (*redis.Client, error) {
	ctx := context.Background()
	client, err := cache.NewClient(ctx, cfg.Data.RedisURL, cfg.Telemetry.Enabled)
	if err != nil {
		slog.Error("Failed to connect to Redis",
			"error", err,
			"hint", "Ensure Redis is running and APP_REDIS_URL is correct",
		)
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return cache.Close(client)
		},
	})

	return client, nil
}

func newEventBus(lc fx.Lifecycle) event.EventBus {
	bus := eventbus.NewInMemoryEventBus()
	slog.Info("Event bus initialized", "type", "InMemoryEventBus")

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return bus.Close()
		},
	})

	return bus
}

// newIAMConfig 从通用配置转换到 IAM 专属配置。
func newIAMConfig(cfg *internalConfig.Config) *config.Config {
	iamCfg := ToIAMConfig(cfg)
	return &iamCfg
}
