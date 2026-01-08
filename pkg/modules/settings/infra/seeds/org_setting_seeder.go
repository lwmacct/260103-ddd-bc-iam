// Package seeds 提供组织配置值的种子数据
package seeds

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	persistence "github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/infra/persistence"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// OrgSettingSeeder 组织配置值种子数据
type OrgSettingSeeder struct{}

// Seed 执行组织配置值种子数据填充
//
// 为示例组织创建示例配置值，展示系统功能：
// - acme 组织使用深色主题
// - globex 组织使用英语语言
func (s *OrgSettingSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)

	// 种子数据：组织配置值
	// 格式：(组织名称, 配置键, 值)
	seedData := []struct {
		OrgName    string
		SettingKey string
		Value      any
	}{
		// acme 组织 - 深色主题偏好
		{"acme", "general.theme", "dark"},
		{"acme", "general.language", "en-US"},

		// globex 组织 - 浅色主题
		{"globex", "general.theme", "light"},
		{"globex", "general.language", "zh-CN"},
	}

	insertedCount := 0

	for _, seed := range seedData {
		// 1. 查找组织 ID
		var orgResult struct {
			ID uint
		}
		err := db.Table("orgs").Select("id").Where("name = ?", seed.OrgName).First(&orgResult).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				slog.Info("Org not found, skipping setting", "org_name", seed.OrgName, "setting_key", seed.SettingKey)
				continue
			}
			return err
		}
		orgID := orgResult.ID

		// 2. 检查是否已存在
		var existing persistence.OrgSettingModel
		err = db.Where("org_id = ? AND setting_key = ?", orgID, seed.SettingKey).First(&existing).Error
		if err == nil {
			slog.Info("Org setting already exists, skipping", "org_id", orgID, "setting_key", seed.SettingKey)
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 3. 序列化值
		valueJSON, err := json.Marshal(seed.Value)
		if err != nil {
			return err
		}

		// 4. 插入记录
		orgSetting := persistence.OrgSettingModel{
			OrgID:      orgID,
			SettingKey: seed.SettingKey,
			Value:      valueJSON,
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "org_id"}, {Name: "setting_key"}},
			DoNothing: true,
		}).Create(&orgSetting).Error; err != nil {
			return err
		}

		insertedCount++
		slog.Info("Seeded org setting", "org_name", seed.OrgName, "setting_key", seed.SettingKey, "value", seed.Value)
	}

	slog.Info("Seeded org settings", "attempted", len(seedData), "inserted", insertedCount)

	return nil
}

// Name 返回种子器名称。
func (s *OrgSettingSeeder) Name() string {
	return "OrgSettingSeeder"
}
