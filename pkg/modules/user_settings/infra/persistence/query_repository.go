package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/user_settings/domain/userset"
	"gorm.io/gorm"
)

// queryRepository 用户配置查询仓储的 GORM 实现
type queryRepository struct {
	db *gorm.DB
}

// NewQueryRepository 创建查询仓储实例
func NewQueryRepository(db *gorm.DB) userset.QueryRepository {
	return &queryRepository{db: db}
}

// FindByUserAndKey 根据用户 ID 和键名查找用户配置
// 如果不存在返回 nil, nil
func (r *queryRepository) FindByUserAndKey(ctx context.Context, userID uint, key string) (*userset.UserSetting, error) {
	var model UserSettingModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND setting_key = ?", userID, key).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//nolint:nilnil // 返回 nil, nil 表示未找到，这是预期情况
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user setting: %w", err)
	}

	return model.ToEntity(), nil
}

// FindByUser 查找用户的所有自定义配置
func (r *queryRepository) FindByUser(ctx context.Context, userID uint) ([]*userset.UserSetting, error) {
	var models []*UserSettingModel
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find user settings: %w", err)
	}

	return toEntities(models), nil
}

// FindByUserAndKeys 批量查询用户的多个配置
func (r *queryRepository) FindByUserAndKeys(ctx context.Context, userID uint, keys []string) ([]*userset.UserSetting, error) {
	if len(keys) == 0 {
		return []*userset.UserSetting{}, nil
	}

	var models []*UserSettingModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND setting_key IN ?", userID, keys).
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find user settings by keys: %w", err)
	}

	return toEntities(models), nil
}
