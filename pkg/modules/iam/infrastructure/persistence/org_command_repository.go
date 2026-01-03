package persistence

import (
	"context"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/org"
	"gorm.io/gorm"
)

// organizationCommandRepository 组织命令仓储的 GORM 实现
type organizationCommandRepository struct {
	*GenericCommandRepository[org.Org, *OrgModel]
}

// NewOrganizationCommandRepository 创建组织命令仓储实例
func NewOrganizationCommandRepository(db *gorm.DB) org.CommandRepository {
	return &organizationCommandRepository{
		GenericCommandRepository: NewGenericCommandRepository(
			db, newOrgModelFromEntity,
		),
	}
}

// Create、Update、Delete 方法由 GenericCommandRepository 提供

// UpdateStatus 更新组织状态
func (r *organizationCommandRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	if err := r.DB().WithContext(ctx).Model(&OrgModel{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update organization status: %w", err)
	}
	return nil
}
