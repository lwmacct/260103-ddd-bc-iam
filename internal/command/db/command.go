// Package db 提供数据库管理命令。
package db

import (
	"github.com/urfave/cli/v3"
)

// Command 定义数据库管理命令组。
var Command = &cli.Command{
	Name:  "db",
	Usage: "数据库操作",
	Description: `数据库迁移和种子数据管理。

示例:
  db migrate              执行迁移
  db migrate -c dev.yaml  使用指定配置
  db reset                重置数据库
  db seed                 执行种子数据`,
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
	},
	Commands: []*cli.Command{
		{
			Name:    "migrate",
			Aliases: []string{"m"},
			Usage:   "执行数据库迁移",
			Action:  actionMigrate,
		},
		{
			Name:    "reset",
			Aliases: []string{"r"},
			Usage:   "重置数据库（删表+重建+种子数据）",
			Action:  actionReset,
		},
		{
			Name:    "seed",
			Aliases: []string{"s"},
			Usage:   "执行种子数据",
			Action:  actionSeed,
		},
	},
}
