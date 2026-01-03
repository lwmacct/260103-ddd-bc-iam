// Package seeds 提供用户配置值的种子数据
package seeds

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	persistence "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/infra/persistence"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserSettingSeeder 用户配置值种子数据
type UserSettingSeeder struct{}

// Seed 执行用户配置值种子数据填充
//
// 为演示用户创建示例配置值，展示系统功能：
// - admin 用户使用深色主题
// - acme_user 使用英语语言
// - demo 用户关闭邮件通知
func (s *UserSettingSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)

	// 种子数据：用户配置值
	// 格式：(用户名, 配置键, 值)
	seedData := []struct {
		Username   string
		SettingKey string
		Value      any
	}{
		// admin 用户 - 深色主题偏好
		{"admin", "general.theme", "dark"},
		{"admin", "general.language", "en-US"},

		// acme_user - 浅色主题 + 简体中文
		{"acme_user", "general.theme", "light"},
		{"acme_user", "general.language", "zh-CN"},

		// globex_user - 跟随系统主题
		{"globex_user", "general.theme", "system"},

		// demo 用户 - 关闭邮件通知
		{"demo", "notification.enable_email", false},
		{"demo", "general.theme", "dark"},
	}

	insertedCount := 0

	for _, seed := range seedData {
		// 1. 查找用户 ID
		var userResult struct {
			ID uint
		}
		err := db.Table("users").Select("id").Where("username = ?", seed.Username).First(&userResult).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				slog.Info("User not found, skipping setting", "username", seed.Username, "setting_key", seed.SettingKey)
				continue
			}
			return err
		}
		userID := userResult.ID

		// 2. 检查是否已存在
		var existing persistence.UserSettingModel
		err = db.Where("user_id = ? AND setting_key = ?", userID, seed.SettingKey).First(&existing).Error
		if err == nil {
			slog.Info("User setting already exists, skipping", "user_id", userID, "setting_key", seed.SettingKey)
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
		userSetting := persistence.UserSettingModel{
			UserID:     userID,
			SettingKey: seed.SettingKey,
			Value:      valueJSON,
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "setting_key"}},
			DoNothing: true,
		}).Create(&userSetting).Error; err != nil {
			return err
		}

		insertedCount++
		slog.Info("Seeded user setting", "username", seed.Username, "setting_key", seed.SettingKey, "value", seed.Value)
	}

	slog.Info("Seeded user settings", "attempted", len(seedData), "inserted", insertedCount)

	return nil
}

// Name 返回种子器名称。
func (s *UserSettingSeeder) Name() string {
	return "UserSettingSeeder"
}
