package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/usersetting"
	"gorm.io/gorm"
)

// userSettingQueryRepository 用户设置查询仓储的 GORM 实现
type userSettingQueryRepository struct {
	db *gorm.DB
}

// NewUserSettingQueryRepository 创建用户设置查询仓储实例
func NewUserSettingQueryRepository(db *gorm.DB) usersetting.QueryRepository {
	return &userSettingQueryRepository{db: db}
}

// FindByUserAndKey 根据用户 ID 和键名查找用户设置
// 返回 nil 表示未找到（这是预期情况，不是错误）
func (r *userSettingQueryRepository) FindByUserAndKey(ctx context.Context, userID uint, key string) (*usersetting.UserSetting, error) {
	var model UserSettingModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND key = ?", userID, key).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//nolint:nilnil // 返回 nil, nil 表示未找到，这是预期情况
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user setting by user_id and key: %w", err)
	}

	return newUserSettingEntityFromModel(&model), nil
}

// FindByUser 查找用户的所有自定义设置
func (r *userSettingQueryRepository) FindByUser(ctx context.Context, userID uint) ([]*usersetting.UserSetting, error) {
	var models []*UserSettingModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find user settings by user_id: %w", err)
	}

	return newUserSettingEntitiesFromModels(models), nil
}

// FindByUserAndCategory 根据用户 ID 和分类查找用户设置
func (r *userSettingQueryRepository) FindByUserAndCategory(ctx context.Context, userID uint, category string) ([]*usersetting.UserSetting, error) {
	var models []*UserSettingModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND category = ?", userID, category).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find user settings by user_id and category: %w", err)
	}

	return newUserSettingEntitiesFromModels(models), nil
}

// ExistsByUserAndKey 检查用户是否对指定键有自定义值
func (r *userSettingQueryRepository) ExistsByUserAndKey(ctx context.Context, userID uint, key string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&UserSettingModel{}).
		Where("user_id = ? AND key = ?", userID, key).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user setting existence: %w", err)
	}

	return count > 0, nil
}
