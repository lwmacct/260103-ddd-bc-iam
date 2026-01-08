// Package seeds 提供团队配置值的种子数据
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

// TeamSettingSeeder 团队配置值种子数据
type TeamSettingSeeder struct{}

// Seed 执行团队配置值种子数据填充
//
// 为示例团队创建示例配置值，展示三级继承功能：
// - acme 的 engineering 团队使用浅色主题（覆盖组织配置）
// - acme 的 product 团队设置特定时区
// - globex 的 engineering 团队设置特定语言
func (s *TeamSettingSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)

	// 种子数据：团队配置值
	// 格式：(组织名称, 团队名称, 配置键, 值)
	seedData := []struct {
		OrgName    string
		TeamName   string
		SettingKey string
		Value      any
	}{
		// acme 的 engineering 团队 - 使用浅色主题（覆盖组织的深色配置）
		{"acme", "engineering", "general.theme", "light"},
		{"acme", "engineering", "notification.enable_slack", true},

		// acme 的 product 团队 - 设置特定时区
		{"acme", "product", "general.timezone", "Asia/Shanghai"},

		// globex 的 engineering 团队
		{"globex", "engineering", "general.language", "en-US"},
	}

	insertedCount := 0

	for _, seed := range seedData {
		// 1. 查找团队 ID（通过组织名称和团队名称）
		var teamResult struct {
			ID uint
		}
		err := db.Table("teams").
			Select("teams.id").
			Joins("JOIN orgs ON orgs.id = teams.org_id").
			Where("orgs.name = ? AND teams.name = ?", seed.OrgName, seed.TeamName).
			First(&teamResult).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				slog.Info("Team not found, skipping setting",
					"org_name", seed.OrgName,
					"team_name", seed.TeamName,
					"setting_key", seed.SettingKey)
				continue
			}
			return err
		}
		teamID := teamResult.ID

		// 2. 检查是否已存在
		var existing persistence.TeamSettingModel
		err = db.Where("team_id = ? AND setting_key = ?", teamID, seed.SettingKey).First(&existing).Error
		if err == nil {
			slog.Info("Team setting already exists, skipping",
				"team_id", teamID,
				"setting_key", seed.SettingKey)
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
		teamSetting := persistence.TeamSettingModel{
			TeamID:     teamID,
			SettingKey: seed.SettingKey,
			Value:      valueJSON,
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "team_id"}, {Name: "setting_key"}},
			DoNothing: true,
		}).Create(&teamSetting).Error; err != nil {
			return err
		}

		insertedCount++
		slog.Info("Seeded team setting",
			"org_name", seed.OrgName,
			"team_name", seed.TeamName,
			"setting_key", seed.SettingKey,
			"value", seed.Value)
	}

	slog.Info("Seeded team settings", "attempted", len(seedData), "inserted", insertedCount)

	return nil
}

// Name 返回种子器名称。
func (s *TeamSettingSeeder) Name() string {
	return "TeamSettingSeeder"
}
