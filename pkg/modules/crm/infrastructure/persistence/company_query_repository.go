package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/company"
)

// companyQueryRepository 公司读仓储实现。
type companyQueryRepository struct {
	db *gorm.DB
}

// NewCompanyQueryRepository 创建公司读仓储。
func NewCompanyQueryRepository(db *gorm.DB) company.QueryRepository {
	return &companyQueryRepository{db: db}
}

// GetByID 根据 ID 查询公司。
func (r *companyQueryRepository) GetByID(ctx context.Context, id uint) (*company.Company, error) {
	var model CompanyModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, company.ErrCompanyNotFound
		}
		return nil, err
	}
	return toCompanyEntity(&model), nil
}

// GetByName 根据名称查询公司。
func (r *companyQueryRepository) GetByName(ctx context.Context, name string) (*company.Company, error) {
	var model CompanyModel
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, company.ErrCompanyNotFound
		}
		return nil, err
	}
	return toCompanyEntity(&model), nil
}

// ListByIndustry 根据行业查询公司列表。
func (r *companyQueryRepository) ListByIndustry(ctx context.Context, industry string, offset, limit int) ([]*company.Company, error) {
	var models []*CompanyModel
	if err := r.db.WithContext(ctx).
		Where("industry = ?", industry).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}
	return toCompanyEntities(models), nil
}

// CountByIndustry 统计指定行业的公司数量。
func (r *companyQueryRepository) CountByIndustry(ctx context.Context, industry string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&CompanyModel{}).
		Where("industry = ?", industry).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ListByOwner 根据负责人查询公司列表。
func (r *companyQueryRepository) ListByOwner(ctx context.Context, ownerID uint, offset, limit int) ([]*company.Company, error) {
	var models []*CompanyModel
	if err := r.db.WithContext(ctx).
		Where("owner_id = ?", ownerID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}
	return toCompanyEntities(models), nil
}

// CountByOwner 统计指定负责人的公司数量。
func (r *companyQueryRepository) CountByOwner(ctx context.Context, ownerID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&CompanyModel{}).
		Where("owner_id = ?", ownerID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// List 分页查询公司列表。
func (r *companyQueryRepository) List(ctx context.Context, offset, limit int) ([]*company.Company, error) {
	var models []*CompanyModel
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}
	return toCompanyEntities(models), nil
}

// Count 统计公司总数。
func (r *companyQueryRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&CompanyModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
