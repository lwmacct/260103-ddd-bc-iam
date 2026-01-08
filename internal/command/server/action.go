package server

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/251219-go-pkg-logm/pkg/logm"
	"github.com/lwmacct/251219-go-pkg-logm/pkg/logm/formatter"
	"github.com/lwmacct/251219-go-pkg-logm/pkg/logm/writer"
	"github.com/lwmacct/260103-ddd-iam-bc/internal/config"
	"github.com/lwmacct/260103-ddd-iam-bc/internal/container"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings"
	_ "github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/response"
)

// action 执行 HTTP 服务器启动。
func action(ctx context.Context, cmd *cli.Command) error {
	initLogger()
	cfg := loadConfig(cmd)

	fxOptions := buildFxOptions(cfg, cmd.Bool("fx-log"))

	fxApp := fx.New(fxOptions...)
	if err := fxApp.Err(); err != nil {
		return fmt.Errorf("create fx app: %w", err)
	}

	fxApp.Run()
	return nil
}

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

// buildFxOptions 构建 Fx 选项。
func buildFxOptions(cfg *config.Config, fxLogEnabled bool) []fx.Option {
	fxOptions := []fx.Option{
		fx.Supply(cfg),
		fx.StartTimeout(30 * time.Second),
		fx.StopTimeout(10 * time.Second),
		// Platform 层 (基础设施)
		container.InfraModule,
		container.ServiceModule,
		// 业务模块 (Bounded Contexts) - 完全自治
		iam.Module(),
		settings.Module(),
		container.SettingsModule(),
		// HTTP 层 (跨模块handler + 路由)
		container.HTTPModule,
		container.HooksModule,
		// Swagger 端点注册 - 使用者决定
		fx.Invoke(func(r *gin.Engine) {
			r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}),
	}

	// CLI --fx-log 优先级高于配置文件
	if !cfg.Server.FxLogEnabled && !fxLogEnabled {
		fxOptions = append(fxOptions, fx.WithLogger(func() fxevent.Logger {
			return nopLogger{}
		}))
	}

	return fxOptions
}

// nopLogger 空日志记录器，不输出任何 Fx 框架日志。
type nopLogger struct{}

func (nopLogger) LogEvent(fxevent.Event) {}

// initLogger 初始化日志系统：终端彩色 + 文件纯文本。
func initLogger() {
	stdoutHandler := logm.NewHandler(&logm.HandlerConfig{
		LevelVar:   logm.GetLevelVar(),
		Formatter:  formatter.ColorText(),
		Writers:    []logm.Writer{writer.Stdout()},
		AddSource:  true,
		TimeFormat: "15:04:05.000",
	})

	fileHandler := logm.NewHandler(&logm.HandlerConfig{
		LevelVar:   logm.GetLevelVar(),
		Formatter:  formatter.Text(),
		Writers:    []logm.Writer{writer.File("/tmp/app.log")},
		AddSource:  true,
		TimeFormat: "2006-01-02 15:04:05.000",
	})

	logger := slog.New(&multiHandler{
		handlers: []slog.Handler{stdoutHandler, fileHandler},
	})
	slog.SetDefault(logger)
}

// multiHandler 将日志记录到多个 Handler。
type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		_ = h.Handle(ctx, r)
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}
