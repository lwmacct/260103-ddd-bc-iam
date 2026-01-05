package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/lwmacct/260103-ddd-bc-iam/internal/command/db"
	"github.com/lwmacct/260103-ddd-bc-iam/internal/command/server"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:    "iam-server",
		Version: "1.0.0",
		Usage:   "Go DDD Package Library Server - 完整的生产级 HTTP 服务器",
		Description: `基于 DDD + CQRS 架构的可复用模块库。

示例:
  iam-server              启动 HTTP 服务器
  iam-server db migrate   执行数据库迁移
  iam-server db reset     重置数据库`,
		Commands: []*cli.Command{
			server.Command,
			db.Command,
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		slog.Error("Application failed to run", "error", err)
		os.Exit(1)
	}
}
