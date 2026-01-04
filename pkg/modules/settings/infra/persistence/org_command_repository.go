package persistence

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/org"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// orgCommandRepository 组织配置命令仓储的 GORM 实现
type orgCommandRepository struct {
	db *gorm.DB
}

// NewOrgCommandRepository 创建组织配置命令仓储实例
func NewOrgCommandRepository(db *gorm.DB) org.CommandRepository {
	return &orgCommandRepository{db: db}
}

// Upsert 创建或更新组织配置（基于 org_id + setting_key 唯一约束）
func (r *orgCommandRepository) Upsert(ctx context.Context, setting *org.OrgSetting) error {
	model := newOrgModelFromEntity(setting)
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "org_id"}, {Name: "setting_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(model).Error

	if err != nil {
		return fmt.Errorf("failed to upsert org setting: %w", err)
	}

	// 回写生成的 ID
	setting.ID = model.ID
	return nil
}

// Delete 删除指定组织的指定配置
func (r *orgCommandRepository) Delete(ctx context.Context, orgID uint, key string) error {
	result := r.db.WithContext(ctx).
		Where("org_id = ? AND setting_key = ?", orgID, key).
		Delete(&OrgSettingModel{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete org setting: %w", result.Error)
	}
	return nil
}

// DeleteByOrg 删除指定组织的所有配置
func (r *orgCommandRepository) DeleteByOrg(ctx context.Context, orgID uint) error {
	result := r.db.WithContext(ctx).
		Where("org_id = ?", orgID).
		Delete(&OrgSettingModel{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete org settings: %w", result.Error)
	}
	return nil
}
