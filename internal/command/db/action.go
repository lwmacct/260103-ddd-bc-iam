package db

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/260103-ddd-bc-iam/internal/config"
	"github.com/lwmacct/260103-ddd-bc-iam/internal/container"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// loadConfig 加载配置。
func loadConfig(cmd *cli.Command) *config.Config {
	prefix := cmd.String("env-prefix")
	if prefix == "" {
		prefix = "APP_"
	}

	opts := []cfgm.Option{
		cfgm.WithEnvPrefix(prefix),
	}

	if path := cmd.String("config"); path != "" {
		opts = append(opts, cfgm.WithConfigPaths(path))
	}

	return cfgm.MustLoadCmd(cmd, config.DefaultConfig(), "", opts...)
}

// actionMigrate 执行数据库迁移。
func actionMigrate(ctx context.Context, cmd *cli.Command) error {
	cfg := loadConfig(cmd)

	appCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fxApp := fx.New(
		fx.Supply(cfg),
		fx.StartTimeout(5*time.Minute),
		fx.StopTimeout(10*time.Second),
		container.InfraModule,
		fx.Invoke(container.RunMigration),
		fx.WithLogger(func() fxevent.Logger { return nopLogger{} }),
	)

	if err := fxApp.Err(); err != nil {
		return fmt.Errorf("create fx app: %w", err)
	}

	if err := fxApp.Start(appCtx); err != nil {
		return err
	}
	return fxApp.Stop(appCtx)
}

// actionReset 重置数据库。
func actionReset(ctx context.Context, cmd *cli.Command) error {
	cfg := loadConfig(cmd)

	appCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fxApp := fx.New(
		fx.Supply(cfg),
		fx.StartTimeout(5*time.Minute),
		fx.StopTimeout(10*time.Second),
		container.InfraModule,
		fx.Invoke(container.RunReset),
		fx.WithLogger(func() fxevent.Logger { return nopLogger{} }),
	)

	if err := fxApp.Err(); err != nil {
		return fmt.Errorf("create fx app: %w", err)
	}

	if err := fxApp.Start(appCtx); err != nil {
		return err
	}
	return fxApp.Stop(appCtx)
}

// actionSeed 执行种子数据。
func actionSeed(ctx context.Context, cmd *cli.Command) error {
	cfg := loadConfig(cmd)

	appCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fxApp := fx.New(
		fx.Supply(cfg),
		fx.StartTimeout(5*time.Minute),
		fx.StopTimeout(10*time.Second),
		container.InfraModule,
		fx.Invoke(container.RunSeed),
		fx.WithLogger(func() fxevent.Logger { return nopLogger{} }),
	)

	if err := fxApp.Err(); err != nil {
		return fmt.Errorf("create fx app: %w", err)
	}

	if err := fxApp.Start(appCtx); err != nil {
		return err
	}
	return fxApp.Stop(appCtx)
}

// nopLogger 空日志记录器。
type nopLogger struct{}

func (nopLogger) LogEvent(fxevent.Event) {}
