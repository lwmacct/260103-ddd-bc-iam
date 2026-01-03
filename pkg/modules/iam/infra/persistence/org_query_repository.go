package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/org"
	"gorm.io/gorm"
)

// organizationQueryRepository 组织查询仓储的 GORM 实现
type organizationQueryRepository struct {
	db *gorm.DB
}

// NewOrganizationQueryRepository 创建组织查询仓储实例
func NewOrganizationQueryRepository(db *gorm.DB) org.QueryRepository {
	return &organizationQueryRepository{db: db}
}

// GetByID 根据 ID 获取组织
func (r *organizationQueryRepository) GetByID(ctx context.Context, id uint) (*org.Org, error) {
	var model OrgModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, org.ErrOrgNotFound
		}
		return nil, fmt.Errorf("failed to get organization by id: %w", err)
	}
	return model.ToEntity(), nil
}

// GetByName 根据名称获取组织
func (r *organizationQueryRepository) GetByName(ctx context.Context, name string) (*org.Org, error) {
	var model OrgModel
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, org.ErrOrgNotFound
		}
		return nil, fmt.Errorf("failed to get organization by name: %w", err)
	}
	return model.ToEntity(), nil
}

// GetByIDWithTeams 根据 ID 获取组织（包含团队列表）
func (r *organizationQueryRepository) GetByIDWithTeams(ctx context.Context, id uint) (*org.Org, error) {
	var model OrgModel
	if err := r.db.WithContext(ctx).Preload("Teams").First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, org.ErrOrgNotFound
		}
		return nil, fmt.Errorf("failed to get organization with teams: %w", err)
	}
	return model.ToEntity(), nil
}

// GetByIDWithMembers 根据 ID 获取组织（包含成员列表）
func (r *organizationQueryRepository) GetByIDWithMembers(ctx context.Context, id uint) (*org.Org, error) {
	var model OrgModel
	if err := r.db.WithContext(ctx).Preload("Members").First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, org.ErrOrgNotFound
		}
		return nil, fmt.Errorf("failed to get organization with members: %w", err)
	}
	return model.ToEntity(), nil
}

// List 获取组织列表（分页）
func (r *organizationQueryRepository) List(ctx context.Context, offset, limit int) ([]*org.Org, error) {
	var models []OrgModel
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}
	return mapOrgModelsToEntities(models), nil
}

// Count 统计组织数量
func (r *organizationQueryRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&OrgModel{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count organizations: %w", err)
	}
	return count, nil
}

// Exists 检查组织是否存在
func (r *organizationQueryRepository) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&OrgModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check organization existence: %w", err)
	}
	return count > 0, nil
}

// ExistsByName 检查组织名称是否存在
func (r *organizationQueryRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&OrgModel{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check organization name existence: %w", err)
	}
	return count > 0, nil
}

// Search 搜索组织（支持名称、显示名称模糊匹配）
func (r *organizationQueryRepository) Search(ctx context.Context, keyword string, offset, limit int) ([]*org.Org, error) {
	var models []OrgModel
	query := r.db.WithContext(ctx).
		Where("name ILIKE ? OR display_name ILIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Offset(offset).Limit(limit)

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to search organizations: %w", err)
	}
	return mapOrgModelsToEntities(models), nil
}

// CountBySearch 统计搜索结果数量
func (r *organizationQueryRepository) CountBySearch(ctx context.Context, keyword string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&OrgModel{}).
		Where("name ILIKE ? OR display_name ILIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count search results: %w", err)
	}
	return count, nil
}
