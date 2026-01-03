package container

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/gorm"

	eventhandler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infra/eventhandler"
	persistence "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infra/persistence"
	iamSeeds "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infra/seeds"
	localSeeds "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/infra/seeds"
	settingsSeeds "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/infra/seeds"
	dbpkg "github.com/lwmacct/260103-ddd-shared/pkg/platform/db"
	"github.com/lwmacct/260103-ddd-shared/pkg/shared/event"
)

// HooksModule 提供生命周期钩子和事件处理器注册。
var HooksModule = fx.Module("hooks",
	fx.Invoke(RegisterEventHandlers),
)

// eventHandlersParams 聚合事件处理器所需的依赖。
type eventHandlersParams struct {
	fx.In

	EventBus   event.EventBus
	AuditRepos persistence.AuditRepositories
}

// RegisterEventHandlers 设置审计日志的事件订阅。
//
// 订阅事件：
//   - *（所有事件）→ 审计日志
func RegisterEventHandlers(p eventHandlersParams) {
	// 审计日志处理器
	auditHandler := eventhandler.NewAuditEventHandler(p.AuditRepos.Command)

	// 订阅所有事件用于审计日志
	p.EventBus.Subscribe("*", auditHandler)

	slog.Info("Event handlers initialized",
		"handlers", []string{"AuditEventHandler"},
		"audit_subscriptions", []string{"*"},
	)
}

// --- CLI 命令函数 ---

// RunMigration 执行数据库迁移。
func RunMigration(lc fx.Lifecycle, db *gorm.DB) error {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			slog.Info("Running database migration...")

			models := GetAllModels()
			if err := db.AutoMigrate(models...); err != nil {
				return err
			}

			// 创建索引
			if err := createAllIndexes(db); err != nil {
				return err
			}

			slog.Info("Database migration completed successfully")
			return nil
		},
	})
	return nil
}

// RunReset 重置数据库（删表+迁移+种子数据）。
func RunReset(lc fx.Lifecycle, db *gorm.DB, redis *redis.Client) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.Info("Resetting db...")

			// 1. 清空 Redis 缓存
			if err := redis.FlushAll(ctx).Err(); err != nil {
				slog.Warn("Failed to flush Redis", "error", err)
			} else {
				slog.Info("Redis cache flushed")
			}

			// 2. 重置数据库
			migrator := dbpkg.NewMigrationManager(db, GetAllModels())
			if err := migrator.ResetWithIndexes(getIndexMigrations()); err != nil {
				return err
			}

			// 3. 合并 IAM、外部 Settings 和本地 Settings 的种子数据
			iamSeeders := iamSeeds.DefaultSeeders()
			settingsSeederList := settingsSeeds.DefaultSeeders()
			localSeeders := localSeeds.DefaultSeeders()

			// 转换 Settings Seeder 到 shared db.Seeder
			allSeeders := make([]dbpkg.Seeder, 0, len(iamSeeders)+len(settingsSeederList)+len(localSeeders))
			allSeeders = append(allSeeders, iamSeeders...)
			for _, s := range settingsSeederList {
				allSeeders = append(allSeeders, settingsSeederAdapter{s})
			}
			allSeeders = append(allSeeders, localSeeders...)

			// 4. 执行种子数据
			seeder := dbpkg.NewSeederManager(db, allSeeders)
			if err := seeder.Run(ctx); err != nil {
				return err
			}

			slog.Info("Database reset completed successfully")
			return nil
		},
	})
	return nil
}

// RunSeed 执行种子数据。
func RunSeed(lc fx.Lifecycle, db *gorm.DB) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.Info("Running database seeders...")

			// 合并 IAM、外部 Settings 和本地 Settings 的种子数据
			iamSeeders := iamSeeds.DefaultSeeders()
			settingsSeederList := settingsSeeds.DefaultSeeders()
			localSeeders := localSeeds.DefaultSeeders()

			// 转换 Settings Seeder 到 shared db.Seeder
			allSeeders := make([]dbpkg.Seeder, 0, len(iamSeeders)+len(settingsSeederList)+len(localSeeders))
			allSeeders = append(allSeeders, iamSeeders...)
			for _, s := range settingsSeederList {
				allSeeders = append(allSeeders, settingsSeederAdapter{s})
			}
			allSeeders = append(allSeeders, localSeeders...)

			seeder := dbpkg.NewSeederManager(db, allSeeders)
			if err := seeder.Run(ctx); err != nil {
				return err
			}

			slog.Info("Database seeding completed successfully")
			return nil
		},
	})
	return nil
}

// settingsSeederAdapter 适配 Settings Seeder 到 shared db.Seeder。
type settingsSeederAdapter struct {
	seeder settingsSeeds.Seeder
}

// Seed 执行种子数据填充。
func (a settingsSeederAdapter) Seed(ctx context.Context, db *gorm.DB) error {
	return a.seeder.Seed(ctx, db)
}

// Name 返回 Seeder 名称。
func (a settingsSeederAdapter) Name() string {
	return a.seeder.Name()
}

// createAllIndexes 创建所有索引。
func createAllIndexes(db *gorm.DB) error {
	// Model 索引
	for _, im := range getIndexMigrations() {
		if err := dbpkg.CreateIndexes(db, im.Model, im.Indexes); err != nil {
			return err
		}
	}

	// 关联表索引
	if err := dbpkg.CreateJoinTableIndexes(db, getJoinTableIndexes()); err != nil {
		return err
	}

	return nil
}

// getIndexMigrations 返回所有 Model 索引配置。
func getIndexMigrations() []dbpkg.IndexMigration {
	return []dbpkg.IndexMigration{
		// IAM module has no custom indexes
	}
}

// getJoinTableIndexes 返回所有关联表索引配置。
func getJoinTableIndexes() []dbpkg.JoinTableIndex {
	return []dbpkg.JoinTableIndex{
		{Table: "user_roles", Name: "idx_user_roles_user_id", Columns: "user_id"},
		{Table: "user_roles", Name: "idx_user_roles_role_id", Columns: "role_id"},
		{Table: "role_permissions", Name: "idx_role_permissions_role_model_id", Columns: "role_model_id"},
		{Table: "role_permissions", Name: "idx_role_permissions_permission_model_id", Columns: "permission_model_id"},
	}
}
