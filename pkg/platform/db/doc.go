// Package database 提供数据库连接和迁移管理。
//
// # 连接管理
//
// [NewConnection] 创建 GORM 数据库连接：
//   - 支持 PostgreSQL（主要）和 SQLite（测试）
//   - 连接池配置（最大连接数、空闲连接、连接生命周期）
//   - 自动日志记录（可配置日志级别）
//
// 使用示例：
//
//	db, err := database.NewConnection(cfg.Database)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # 数据库迁移
//
// [MigrationManager] 管理数据库 Schema 迁移：
//   - 自动创建表结构
//   - 支持增量迁移
//   - 记录迁移历史
//
// [Migrator] 执行具体的迁移操作：
//   - AutoMigrate: 自动迁移所有模型
//   - 支持自定义迁移脚本
//
// # 数据种子
//
// [Seeder] 填充初始数据：
//   - 创建默认管理员用户
//   - 创建系统角色和权限
//   - 创建默认菜单结构
//   - 创建系统设置项
//
// 使用示例：
//
//	seeder := database.NewSeeder(db, authService)
//	if err := seeder.Seed(); err != nil {
//	    log.Fatal(err)
//	}
//
// # 配置项
//
// 通过 config.Database 配置：
//   - Host: 数据库主机
//   - Port: 数据库端口
//   - User: 用户名
//   - Password: 密码
//   - DBName: 数据库名
//   - SSLMode: SSL 模式
//   - MaxOpenConns: 最大打开连接数
//   - MaxIdleConns: 最大空闲连接数
//   - ConnMaxLifetime: 连接最大生命周期
package db
