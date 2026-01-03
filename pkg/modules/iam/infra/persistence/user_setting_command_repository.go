package persistence

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/usersetting"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// userSettingCommandRepository 用户设置命令仓储的 GORM 实现
type userSettingCommandRepository struct {
	db *gorm.DB
}

// NewUserSettingCommandRepository 创建用户设置命令仓储实例
func NewUserSettingCommandRepository(db *gorm.DB) usersetting.CommandRepository {
	return &userSettingCommandRepository{db: db}
}

// Create 创建用户设置
func (r *userSettingCommandRepository) Create(ctx context.Context, setting *usersetting.UserSetting) error {
	model := newUserSettingModelFromEntity(setting)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create user setting: %w", err)
	}
	// 回写生成的 ID
	setting.ID = model.ID
	return nil
}

// Update 更新用户设置
func (r *userSettingCommandRepository) Update(ctx context.Context, setting *usersetting.UserSetting) error {
	model := newUserSettingModelFromEntity(setting)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to update user setting: %w", err)
	}
	return nil
}

// Delete 删除用户设置 (软删除)
func (r *userSettingCommandRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&UserSettingModel{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user setting: %w", err)
	}
	return nil
}

// Upsert 创建或更新用户设置（基于 user_id + key 唯一约束）
func (r *userSettingCommandRepository) Upsert(ctx context.Context, userID uint, key string, value string) error {
	category := extractCategoryFromKey(key)

	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "category", "updated_at"}),
	}).Create(&UserSettingModel{
		UserID:   userID,
		Key:      key,
		Value:    value,
		Category: category,
	}).Error; err != nil {
		return fmt.Errorf("failed to upsert user setting: %w", err)
	}
	return nil
}
