package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/company"
)

// companyCommandRepository 公司写仓储实现。
type companyCommandRepository struct {
	db *gorm.DB
}

// NewCompanyCommandRepository 创建公司写仓储。
func NewCompanyCommandRepository(db *gorm.DB) company.CommandRepository {
	return &companyCommandRepository{db: db}
}

// Create 创建公司。
func (r *companyCommandRepository) Create(ctx context.Context, c *company.Company) error {
	model := toCompanyModel(c)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	c.ID = model.ID
	c.CreatedAt = model.CreatedAt
	c.UpdatedAt = model.UpdatedAt
	return nil
}

// Update 更新公司。
func (r *companyCommandRepository) Update(ctx context.Context, c *company.Company) error {
	model := toCompanyModel(c)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	c.UpdatedAt = model.UpdatedAt
	return nil
}

// Delete 删除公司。
func (r *companyCommandRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&CompanyModel{}, id).Error
}
