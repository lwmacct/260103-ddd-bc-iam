package persistence

import (
	"context"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// userSettingCommandRepository 用户配置命令仓储的 GORM 实现
type userSettingCommandRepository struct {
	db *gorm.DB
}

// NewUserSettingCommandRepository 创建用户配置命令仓储实例
func NewUserSettingCommandRepository(db *gorm.DB) setting.UserSettingCommandRepository {
	return &userSettingCommandRepository{db: db}
}

// Upsert 插入或更新用户配置
func (r *userSettingCommandRepository) Upsert(ctx context.Context, s *setting.UserSetting) error {
	model := newUserSettingModelFromEntity(s)

	// 使用 ON CONFLICT 实现 Upsert
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "setting_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(model).Error; err != nil {
		return fmt.Errorf("failed to upsert user setting: %w", err)
	}

	// 回写生成的 ID
	if entity := model.ToEntity(); entity != nil {
		*s = *entity
	}

	return nil
}

// Delete 删除用户配置
func (r *userSettingCommandRepository) Delete(ctx context.Context, userID uint, key string) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND setting_key = ?", userID, key).
		Delete(&UserSettingModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete user setting: %w", err)
	}
	return nil
}

// DeleteByUser 删除用户的所有配置
func (r *userSettingCommandRepository) DeleteByUser(ctx context.Context, userID uint) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&UserSettingModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete user settings: %w", err)
	}
	return nil
}

// BatchUpsert 批量插入或更新用户配置
func (r *userSettingCommandRepository) BatchUpsert(ctx context.Context, settings []*setting.UserSetting) error {
	if len(settings) == 0 {
		return nil
	}

	models := make([]*UserSettingModel, 0, len(settings))
	for _, s := range settings {
		if model := newUserSettingModelFromEntity(s); model != nil {
			models = append(models, model)
		}
	}

	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "setting_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(models).Error; err != nil {
		return fmt.Errorf("failed to batch upsert user settings: %w", err)
	}

	// 回写生成的 ID
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			*settings[i] = *entity
		}
	}

	return nil
}
