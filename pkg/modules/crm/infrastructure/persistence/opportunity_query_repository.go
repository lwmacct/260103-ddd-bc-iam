package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/opportunity"
)

// opportunityQueryRepository 商机读仓储实现。
type opportunityQueryRepository struct {
	db *gorm.DB
}

// NewOpportunityQueryRepository 创建商机读仓储。
func NewOpportunityQueryRepository(db *gorm.DB) opportunity.QueryRepository {
	return &opportunityQueryRepository{db: db}
}

// GetByID 根据 ID 查询商机。
func (r *opportunityQueryRepository) GetByID(ctx context.Context, id uint) (*opportunity.Opportunity, error) {
	var model OpportunityModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, opportunity.ErrOpportunityNotFound
		}
		return nil, err
	}
	return toOpportunityEntity(&model), nil
}

// ListByStage 根据阶段查询商机列表。
func (r *opportunityQueryRepository) ListByStage(ctx context.Context, stage opportunity.Stage, offset, limit int) ([]*opportunity.Opportunity, error) {
	var models []*OpportunityModel
	if err := r.db.WithContext(ctx).
		Where("stage = ?", string(stage)).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}
	return toOpportunityEntities(models), nil
}

// CountByStage 统计指定阶段的商机数量。
func (r *opportunityQueryRepository) CountByStage(ctx context.Context, stage opportunity.Stage) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&OpportunityModel{}).
		Where("stage = ?", string(stage)).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ListByOwner 根据负责人查询商机列表。
func (r *opportunityQueryRepository) ListByOwner(ctx context.Context, ownerID uint, offset, limit int) ([]*opportunity.Opportunity, error) {
	var models []*OpportunityModel
	if err := r.db.WithContext(ctx).
		Where("owner_id = ?", ownerID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}
	return toOpportunityEntities(models), nil
}

// CountByOwner 统计指定负责人的商机数量。
func (r *opportunityQueryRepository) CountByOwner(ctx context.Context, ownerID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&OpportunityModel{}).
		Where("owner_id = ?", ownerID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// List 分页查询商机列表。
func (r *opportunityQueryRepository) List(ctx context.Context, offset, limit int) ([]*opportunity.Opportunity, error) {
	var models []*OpportunityModel
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}
	return toOpportunityEntities(models), nil
}

// Count 统计商机总数。
func (r *opportunityQueryRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&OpportunityModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
