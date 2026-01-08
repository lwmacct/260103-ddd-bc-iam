// Package server 提供 HTTP 服务器启动命令。
//
// Swagger 总体配置 - 使用者自定义
//
//	@title           Go DDD Package Library API
//	@version         1.0
//	@description     基于 DDD + CQRS 架构的可复用模块库
//	@host            localhost:8080
//	@BasePath        /
//
//	@contact.name    API Support
//	@contact.url     https://github.com/lwmacct/260103-ddd-iam-bc
//
//	@license.name    MIT
//	@license.url     https://opensource.org/licenses/MIT
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Bearer token authentication
package server

import (
	// Swagger docs - 空白导入触发 docs.go 的 init() 函数
	_ "github.com/lwmacct/260103-ddd-iam-bc/docs/swagger"

	"github.com/urfave/cli/v3"
)

// Command 定义 HTTP 服务器启动命令。
var Command = &cli.Command{
	Name:    "server",
	Aliases: []string{"serve", "api"},
	Usage:   "启动 HTTP 服务器",
	Description: `启动 HTTP API 服务器，提供 RESTful API 接口。

示例:
  server                  启动 HTTP 服务器
  server -c config.yaml   指定配置文件
  server --fx-log         启用 Fx 日志`,
	Action: action,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "配置文件路径",
		},
		&cli.StringFlag{
			Name:  "env-prefix",
			Usage: "环境变量前缀",
			Value: "APP_",
		},
		&cli.BoolFlag{
			Name:  "fx-log",
			Usage: "启用 Fx 依赖注入日志",
		},
	},
}
