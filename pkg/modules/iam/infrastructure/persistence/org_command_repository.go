package persistence

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/org"
	"gorm.io/gorm"
)

// organizationCommandRepository 组织命令仓储的 GORM 实现
type organizationCommandRepository struct {
	db *gorm.DB
}

// NewOrganizationCommandRepository 创建组织命令仓储实例
func NewOrganizationCommandRepository(db *gorm.DB) org.CommandRepository {
	return &organizationCommandRepository{db: db}
}

// Create 创建组织
func (r *organizationCommandRepository) Create(ctx context.Context, org *org.Org) error {
	model := newOrgModelFromEntity(org)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	org.ID = model.ID
	return nil
}

// Update 更新组织
func (r *organizationCommandRepository) Update(ctx context.Context, org *org.Org) error {
	model := newOrgModelFromEntity(org)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete 删除组织
func (r *organizationCommandRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&OrgModel{}, id).Error
}

// UpdateStatus 更新组织状态
func (r *organizationCommandRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	if err := r.db.WithContext(ctx).Model(&OrgModel{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update organization status: %w", err)
	}
	return nil
}
