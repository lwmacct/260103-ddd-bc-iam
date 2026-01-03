package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/contact"
)

// contactQueryRepository 联系人读仓储实现。
type contactQueryRepository struct {
	db *gorm.DB
}

// NewContactQueryRepository 创建联系人读仓储。
func NewContactQueryRepository(db *gorm.DB) contact.QueryRepository {
	return &contactQueryRepository{db: db}
}

// GetByID 根据 ID 获取联系人。
func (r *contactQueryRepository) GetByID(ctx context.Context, id uint) (*contact.Contact, error) {
	var model ContactModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, contact.ErrContactNotFound
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

// GetByEmail 根据邮箱获取联系人。
func (r *contactQueryRepository) GetByEmail(ctx context.Context, email string) (*contact.Contact, error) {
	var model ContactModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, contact.ErrContactNotFound
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

// ListByCompany 获取公司下的联系人列表。
func (r *contactQueryRepository) ListByCompany(ctx context.Context, companyID uint, offset, limit int) ([]*contact.Contact, error) {
	var models []ContactModel
	if err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}
	return mapContactModelsToEntities(models), nil
}

// CountByCompany 统计公司下的联系人数量。
func (r *contactQueryRepository) CountByCompany(ctx context.Context, companyID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&ContactModel{}).
		Where("company_id = ?", companyID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ListByOwner 获取负责人的联系人列表。
func (r *contactQueryRepository) ListByOwner(ctx context.Context, ownerID uint, offset, limit int) ([]*contact.Contact, error) {
	var models []ContactModel
	if err := r.db.WithContext(ctx).
		Where("owner_id = ?", ownerID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}
	return mapContactModelsToEntities(models), nil
}

// CountByOwner 统计负责人的联系人数量。
func (r *contactQueryRepository) CountByOwner(ctx context.Context, ownerID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&ContactModel{}).
		Where("owner_id = ?", ownerID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// List 获取联系人列表。
func (r *contactQueryRepository) List(ctx context.Context, offset, limit int) ([]*contact.Contact, error) {
	var models []ContactModel
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}
	return mapContactModelsToEntities(models), nil
}

// Count 统计联系人总数。
func (r *contactQueryRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&ContactModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
