package db

import (
	"context"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

// Seeder 种子数据接口
// 所有种子实现都需要实现此接口
type Seeder interface {
	// Seed 执行种子数据填充
	Seed(ctx context.Context, db *gorm.DB) error
}

// SeederManager 种子数据管理器
type SeederManager struct {
	db      *gorm.DB
	seeders []Seeder
}

// NewSeederManager 创建种子管理器
func NewSeederManager(db *gorm.DB, seeders []Seeder) *SeederManager {
	return &SeederManager{
		db:      db,
		seeders: seeders,
	}
}

// Run 执行所有种子
func (sm *SeederManager) Run(ctx context.Context) error {
	for _, seeder := range sm.seeders {
		seederName := fmt.Sprintf("%T", seeder)
		slog.Info("Running seeder", "name", seederName)

		if err := seeder.Seed(ctx, sm.db); err != nil {
			return fmt.Errorf("seeder %s failed: %w", seederName, err)
		}

		slog.Info("Seeder completed", "name", seederName)
	}
	return nil
}

// RunOne 执行指定的种子
func (sm *SeederManager) RunOne(ctx context.Context, index int) error {
	if index < 0 || index >= len(sm.seeders) {
		return fmt.Errorf("invalid seeder index: %d (available: 0-%d)", index, len(sm.seeders)-1)
	}

	seeder := sm.seeders[index]
	seederName := fmt.Sprintf("%T", seeder)
	slog.Info("Running seeder", "name", seederName)

	if err := seeder.Seed(ctx, sm.db); err != nil {
		return fmt.Errorf("seeder %s failed: %w", seederName, err)
	}

	slog.Info("Seeder completed", "name", seederName)
	return nil
}

// List 列出所有可用的种子
func (sm *SeederManager) List() []string {
	names := make([]string, len(sm.seeders))
	for i, seeder := range sm.seeders {
		names[i] = fmt.Sprintf("%T", seeder)
	}
	return names
}
