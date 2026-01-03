package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/org"
	"gorm.io/gorm"
)

// teamQueryRepository 团队查询仓储的 GORM 实现
type teamQueryRepository struct {
	db *gorm.DB
}

// NewTeamQueryRepository 创建团队查询仓储实例
func NewTeamQueryRepository(db *gorm.DB) org.TeamQueryRepository {
	return &teamQueryRepository{db: db}
}

// GetByID 根据 ID 获取团队
func (r *teamQueryRepository) GetByID(ctx context.Context, id uint) (*org.Team, error) {
	var model TeamModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, org.ErrTeamNotFound
		}
		return nil, fmt.Errorf("failed to get team by id: %w", err)
	}
	return model.ToEntity(), nil
}

// GetByOrgAndName 根据组织 ID 和团队名称获取团队
func (r *teamQueryRepository) GetByOrgAndName(ctx context.Context, orgID uint, name string) (*org.Team, error) {
	var model TeamModel
	if err := r.db.WithContext(ctx).
		Where("org_id = ? AND name = ?", orgID, name).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, org.ErrTeamNotFound
		}
		return nil, fmt.Errorf("failed to get team by org and name: %w", err)
	}
	return model.ToEntity(), nil
}

// GetByIDWithMembers 根据 ID 获取团队（包含成员列表）
func (r *teamQueryRepository) GetByIDWithMembers(ctx context.Context, id uint) (*org.Team, error) {
	var model TeamModel
	if err := r.db.WithContext(ctx).Preload("Members").First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, org.ErrTeamNotFound
		}
		return nil, fmt.Errorf("failed to get team with members: %w", err)
	}
	return model.ToEntity(), nil
}

// ListByOrg 获取组织的所有团队
func (r *teamQueryRepository) ListByOrg(ctx context.Context, orgID uint, offset, limit int) ([]*org.Team, error) {
	var models []*TeamModel
	if err := r.db.WithContext(ctx).
		Where("org_id = ?", orgID).
		Offset(offset).Limit(limit).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list teams by org: %w", err)
	}
	return mapTeamModelPtrsToEntities(models), nil
}

// CountByOrg 统计组织的团队数量
func (r *teamQueryRepository) CountByOrg(ctx context.Context, orgID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&TeamModel{}).
		Where("org_id = ?", orgID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count teams by org: %w", err)
	}
	return count, nil
}

// Exists 检查团队是否存在
func (r *teamQueryRepository) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&TeamModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check team existence: %w", err)
	}
	return count > 0, nil
}

// ExistsByOrgAndName 检查组织内团队名称是否存在
func (r *teamQueryRepository) ExistsByOrgAndName(ctx context.Context, orgID uint, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&TeamModel{}).
		Where("org_id = ? AND name = ?", orgID, name).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check team name existence: %w", err)
	}
	return count > 0, nil
}

// BelongsToOrg 检查团队是否属于指定组织
func (r *teamQueryRepository) BelongsToOrg(ctx context.Context, teamID, orgID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&TeamModel{}).
		Where("id = ? AND org_id = ?", teamID, orgID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check team belongs to org: %w", err)
	}
	return count > 0, nil
}
