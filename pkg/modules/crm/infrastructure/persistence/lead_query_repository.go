package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/lead"
)

// leadQueryRepository 线索读仓储实现。
type leadQueryRepository struct {
	db *gorm.DB
}

// NewLeadQueryRepository 创建线索读仓储。
func NewLeadQueryRepository(db *gorm.DB) lead.QueryRepository {
	return &leadQueryRepository{db: db}
}

// GetByID 根据 ID 查询线索。
func (r *leadQueryRepository) GetByID(ctx context.Context, id uint) (*lead.Lead, error) {
	var model LeadModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, lead.ErrLeadNotFound
		}
		return nil, err
	}
	return toLeadEntity(&model), nil
}

// ListByStatus 根据状态查询线索列表。
func (r *leadQueryRepository) ListByStatus(ctx context.Context, status lead.Status, offset, limit int) ([]*lead.Lead, error) {
	var models []*LeadModel
	if err := r.db.WithContext(ctx).
		Where("status = ?", string(status)).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}
	return toLeadEntities(models), nil
}

// CountByStatus 统计指定状态的线索数量。
func (r *leadQueryRepository) CountByStatus(ctx context.Context, status lead.Status) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&LeadModel{}).
		Where("status = ?", string(status)).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ListByOwner 根据负责人查询线索列表。
func (r *leadQueryRepository) ListByOwner(ctx context.Context, ownerID uint, offset, limit int) ([]*lead.Lead, error) {
	var models []*LeadModel
	if err := r.db.WithContext(ctx).
		Where("owner_id = ?", ownerID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}
	return toLeadEntities(models), nil
}

// CountByOwner 统计指定负责人的线索数量。
func (r *leadQueryRepository) CountByOwner(ctx context.Context, ownerID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&LeadModel{}).
		Where("owner_id = ?", ownerID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// List 分页查询线索列表。
func (r *leadQueryRepository) List(ctx context.Context, offset, limit int) ([]*lead.Lead, error) {
	var models []*LeadModel
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}
	return toLeadEntities(models), nil
}

// Count 统计线索总数。
func (r *leadQueryRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&LeadModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
