package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/domain/org"
	"gorm.io/gorm"
)

// orgQueryRepository 组织配置查询仓储的 GORM 实现
type orgQueryRepository struct {
	db *gorm.DB
}

// NewOrgQueryRepository 创建组织配置查询仓储实例
func NewOrgQueryRepository(db *gorm.DB) org.QueryRepository {
	return &orgQueryRepository{db: db}
}

// FindByOrgAndKey 根据组织 ID 和键名查找组织配置
// 如果不存在返回 nil, nil
func (r *orgQueryRepository) FindByOrgAndKey(ctx context.Context, orgID uint, key string) (*org.OrgSetting, error) {
	var model OrgSettingModel
	err := r.db.WithContext(ctx).
		Where("org_id = ? AND setting_key = ?", orgID, key).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//nolint:nilnil // 返回 nil, nil 表示未找到，这是预期情况
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find org setting: %w", err)
	}

	return model.ToEntity(), nil
}

// FindByOrg 查找组织的所有自定义配置
func (r *orgQueryRepository) FindByOrg(ctx context.Context, orgID uint) ([]*org.OrgSetting, error) {
	var models []*OrgSettingModel
	err := r.db.WithContext(ctx).
		Where("org_id = ?", orgID).
		Order("setting_key ASC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find org settings: %w", err)
	}

	return toOrgEntities(models), nil
}

// FindByOrgAndKeys 批量查询组织的多个配置
func (r *orgQueryRepository) FindByOrgAndKeys(ctx context.Context, orgID uint, keys []string) ([]*org.OrgSetting, error) {
	if len(keys) == 0 {
		return []*org.OrgSetting{}, nil
	}

	var models []*OrgSettingModel
	err := r.db.WithContext(ctx).
		Where("org_id = ? AND setting_key IN ?", orgID, keys).
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find org settings by keys: %w", err)
	}

	return toOrgEntities(models), nil
}

// CountByOrg 统计组织的自定义配置数量
func (r *orgQueryRepository) CountByOrg(ctx context.Context, orgID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&OrgSettingModel{}).
		Where("org_id = ?", orgID).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count org settings: %w", err)
	}

	return count, nil
}
