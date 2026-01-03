package db

import (
	"fmt"
	"log/slog"
	"strings"

	"gorm.io/gorm"
)

// MigrationManager 迁移管理器
type MigrationManager struct {
	db     *gorm.DB
	models []any
}

// NewMigrationManager 创建迁移管理器
func NewMigrationManager(db *gorm.DB, models []any) *MigrationManager {
	return &MigrationManager{
		db:     db,
		models: models,
	}
}

// IndexMigration 索引迁移配置
type IndexMigration struct {
	Model   any
	Indexes []string
}

// Reset 删除所有表并重新迁移
func (m *MigrationManager) Reset() error {
	return m.ResetWithIndexes(nil)
}

// ResetWithIndexes 删除所有表并重新迁移（包含索引）
func (m *MigrationManager) ResetWithIndexes(indexMigrations []IndexMigration) error {
	// 1. 删除所有表
	if err := m.dropAllTables(); err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	// 2. 重新创建所有表
	if err := m.db.AutoMigrate(m.models...); err != nil {
		return fmt.Errorf("failed to migrate models: %w", err)
	}

	// 3. 创建索引
	for _, im := range indexMigrations {
		if err := m.createIndexes(im.Model, im.Indexes); err != nil {
			return fmt.Errorf("failed to create indexes: %w", err)
		}
	}

	return nil
}

// createIndexes 为模型创建索引
func (m *MigrationManager) createIndexes(model any, indexes []string) error {
	migrator := m.db.Migrator()

	for _, idx := range indexes {
		if !migrator.HasIndex(model, idx) {
			if err := migrator.CreateIndex(model, idx); err != nil {
				return fmt.Errorf("failed to create index %s: %w", idx, err)
			}
			slog.Info("Index created", "index", idx)
		}
	}

	return nil
}

// dropAllTables 使用原生 SQL 删除所有用户表
func (m *MigrationManager) dropAllTables() error {
	const listTablesSQL = `
SELECT table_schema, table_name
FROM information_schema.tables
WHERE table_type = 'BASE TABLE'
  AND table_schema NOT IN ('pg_catalog', 'information_schema')
`

	type tableInfo struct {
		Schema string `gorm:"column:table_schema"`
		Name   string `gorm:"column:table_name"`
	}

	var tables []tableInfo
	if err := m.db.Raw(listTablesSQL).Scan(&tables).Error; err != nil {
		return fmt.Errorf("failed to list tables: %w", err)
	}

	if len(tables) == 0 {
		return nil
	}

	return m.db.Transaction(func(tx *gorm.DB) error {
		for _, tbl := range tables {
			stmt := fmt.Sprintf("DROP TABLE IF EXISTS %s.%s CASCADE",
				quoteIdentifier(tbl.Schema), quoteIdentifier(tbl.Name))
			if err := tx.Exec(stmt).Error; err != nil {
				return fmt.Errorf("failed to drop table %s.%s: %w", tbl.Schema, tbl.Name, err)
			}
		}
		return nil
	})
}

// quoteIdentifier 安全地引用标识符
func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// CreateIndexes 为已存在的表创建索引（供 API 启动时使用）
// GORM AutoMigrate 只在建表时创建索引，此函数用于为已存在的表添加新索引
func CreateIndexes(db *gorm.DB, model any, indexes []string) error {
	migrator := db.Migrator()

	for _, idx := range indexes {
		if !migrator.HasIndex(model, idx) {
			if err := migrator.CreateIndex(model, idx); err != nil {
				return fmt.Errorf("failed to create index %s: %w", idx, err)
			}
			slog.Info("Index created", "index", idx)
		}
	}

	return nil
}

// JoinTableIndex 定义关联表索引
type JoinTableIndex struct {
	Table   string // 表名
	Name    string // 索引名
	Columns string // 索引列（可多列，逗号分隔）
}

// CreateJoinTableIndexes 为 GORM many2many 关联表创建索引
//
// GORM AutoMigrate 创建关联表时只创建复合主键，不会为各外键列创建单独索引。
// 这导致 WHERE user_id = ? 类型的查询需要全表扫描，严重影响性能。
//
// 此函数使用 CREATE INDEX IF NOT EXISTS 确保幂等性。
func CreateJoinTableIndexes(db *gorm.DB, indexes []JoinTableIndex) error {
	for _, idx := range indexes {
		sql := fmt.Sprintf(
			"CREATE INDEX IF NOT EXISTS %s ON %s (%s)",
			quoteIdentifier(idx.Name),
			quoteIdentifier(idx.Table),
			idx.Columns, // 列名不引用，因为可能是多列
		)
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("failed to create join table index %s: %w", idx.Name, err)
		}
		slog.Info("Join table index ensured", "table", idx.Table, "index", idx.Name)
	}
	return nil
}
