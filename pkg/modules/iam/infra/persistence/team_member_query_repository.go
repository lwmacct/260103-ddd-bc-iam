package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/org"
	"gorm.io/gorm"
)

// teamMemberQueryRepository 团队成员查询仓储的 GORM 实现
type teamMemberQueryRepository struct {
	db *gorm.DB
}

// NewTeamMemberQueryRepository 创建团队成员查询仓储实例
func NewTeamMemberQueryRepository(db *gorm.DB) org.TeamMemberQueryRepository {
	return &teamMemberQueryRepository{db: db}
}

// GetByTeamAndUser 获取指定团队的指定用户成员信息
func (r *teamMemberQueryRepository) GetByTeamAndUser(ctx context.Context, teamID, userID uint) (*org.TeamMember, error) {
	var model TeamMemberModel
	if err := r.db.WithContext(ctx).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, org.ErrNotTeamMember
		}
		return nil, fmt.Errorf("failed to get team member: %w", err)
	}
	return model.ToEntity(), nil
}

// ListByTeam 获取团队的所有成员
func (r *teamMemberQueryRepository) ListByTeam(ctx context.Context, teamID uint, offset, limit int) ([]*org.TeamMember, error) {
	var models []*TeamMemberModel
	if err := r.db.WithContext(ctx).
		Where("team_id = ?", teamID).
		Offset(offset).Limit(limit).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list team members: %w", err)
	}
	return mapTeamMemberModelPtrsToEntities(models), nil
}

// CountByTeam 统计团队成员数量
func (r *teamMemberQueryRepository) CountByTeam(ctx context.Context, teamID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&TeamMemberModel{}).
		Where("team_id = ?", teamID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count team members: %w", err)
	}
	return count, nil
}

// ListByUser 获取用户加入的所有团队的成员记录
func (r *teamMemberQueryRepository) ListByUser(ctx context.Context, userID uint) ([]*org.TeamMember, error) {
	var models []*TeamMemberModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list user team memberships: %w", err)
	}
	return mapTeamMemberModelPtrsToEntities(models), nil
}

// ListByUserInOrg 获取用户在指定组织内加入的所有团队的成员记录
func (r *teamMemberQueryRepository) ListByUserInOrg(ctx context.Context, userID, orgID uint) ([]*org.TeamMember, error) {
	var models []*TeamMemberModel
	if err := r.db.WithContext(ctx).
		Joins("JOIN teams ON teams.id = team_members.team_id").
		Where("team_members.user_id = ? AND teams.org_id = ?", userID, orgID).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list user team memberships in org: %w", err)
	}
	return mapTeamMemberModelPtrsToEntities(models), nil
}

// IsMember 检查用户是否为团队成员
func (r *teamMemberQueryRepository) IsMember(ctx context.Context, teamID, userID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&TeamMemberModel{}).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check team membership: %w", err)
	}
	return count > 0, nil
}
