package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
	"gorm.io/gorm"
)

// userSettingQueryRepository 用户配置查询仓储的 GORM 实现
type userSettingQueryRepository struct {
	db *gorm.DB
}

// NewUserSettingQueryRepository 创建用户配置查询仓储实例
func NewUserSettingQueryRepository(db *gorm.DB) setting.UserSettingQueryRepository {
	return &userSettingQueryRepository{db: db}
}

// FindByUserAndKey 根据用户 ID 和 Key 查找用户配置
func (r *userSettingQueryRepository) FindByUserAndKey(ctx context.Context, userID uint, key string) (*setting.UserSetting, error) {
	var model UserSettingModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND setting_key = ?", userID, key).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // returns nil for not found, valid pattern
		}
		return nil, fmt.Errorf("failed to find user setting: %w", err)
	}
	return model.ToEntity(), nil
}

// FindByUser 查找用户的所有自定义配置
func (r *userSettingQueryRepository) FindByUser(ctx context.Context, userID uint) ([]*setting.UserSetting, error) {
	var models []UserSettingModel
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("setting_key ASC").
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find user settings: %w", err)
	}
	return mapUserSettingModelsToEntities(models), nil
}

// FindByUserAndKeys 根据用户 ID 和多个 Key 批量查找用户配置
func (r *userSettingQueryRepository) FindByUserAndKeys(ctx context.Context, userID uint, keys []string) ([]*setting.UserSetting, error) {
	if len(keys) == 0 {
		return []*setting.UserSetting{}, nil
	}
	var models []UserSettingModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND setting_key IN ?", userID, keys).
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find user settings by keys: %w", err)
	}
	return mapUserSettingModelsToEntities(models), nil
}
