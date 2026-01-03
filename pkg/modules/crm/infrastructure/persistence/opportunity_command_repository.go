package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/opportunity"
)

// opportunityCommandRepository 商机写仓储实现。
type opportunityCommandRepository struct {
	db *gorm.DB
}

// NewOpportunityCommandRepository 创建商机写仓储。
func NewOpportunityCommandRepository(db *gorm.DB) opportunity.CommandRepository {
	return &opportunityCommandRepository{db: db}
}

// Create 创建商机。
func (r *opportunityCommandRepository) Create(ctx context.Context, o *opportunity.Opportunity) error {
	model := toOpportunityModel(o)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	o.ID = model.ID
	o.CreatedAt = model.CreatedAt
	o.UpdatedAt = model.UpdatedAt
	return nil
}

// Update 更新商机。
func (r *opportunityCommandRepository) Update(ctx context.Context, o *opportunity.Opportunity) error {
	model := toOpportunityModel(o)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	o.UpdatedAt = model.UpdatedAt
	return nil
}

// Delete 删除商机。
func (r *opportunityCommandRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&OpportunityModel{}, id).Error
}
