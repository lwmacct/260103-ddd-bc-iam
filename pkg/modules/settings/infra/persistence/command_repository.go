package persistence

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// commandRepository 用户配置命令仓储的 GORM 实现
type commandRepository struct {
	db *gorm.DB
}

// NewCommandRepository 创建命令仓储实例
func NewCommandRepository(db *gorm.DB) user.CommandRepository {
	return &commandRepository{db: db}
}

// Upsert 创建或更新用户配置（基于 user_id + setting_key 唯一约束）
func (r *commandRepository) Upsert(ctx context.Context, setting *user.UserSetting) error {
	model := newModelFromEntity(setting)
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "setting_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(model).Error

	if err != nil {
		return fmt.Errorf("failed to upsert user setting: %w", err)
	}

	// 回写生成的 ID
	setting.ID = model.ID
	return nil
}

// Delete 删除指定用户的指定配置
func (r *commandRepository) Delete(ctx context.Context, userID uint, key string) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND setting_key = ?", userID, key).
		Delete(&UserSettingModel{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete user setting: %w", result.Error)
	}
	return nil
}

// DeleteByUser 删除指定用户的所有配置
func (r *commandRepository) DeleteByUser(ctx context.Context, userID uint) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&UserSettingModel{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete user settings: %w", result.Error)
	}
	return nil
}

// BatchUpsert 批量创建或更新用户配置
func (r *commandRepository) BatchUpsert(ctx context.Context, settings []*user.UserSetting) error {
	if len(settings) == 0 {
		return nil
	}

	models := make([]*UserSettingModel, len(settings))
	for i, s := range settings {
		models[i] = newModelFromEntity(s)
	}

	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "setting_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&models).Error

	if err != nil {
		return fmt.Errorf("failed to batch upsert user settings: %w", err)
	}

	// 回写生成的 IDs
	for i, m := range models {
		settings[i].ID = m.ID
	}
	return nil
}
