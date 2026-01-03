package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/org"
	"gorm.io/gorm"
)

// teamMemberCommandRepository 团队成员命令仓储的 GORM 实现
type teamMemberCommandRepository struct {
	db *gorm.DB
}

// NewTeamMemberCommandRepository 创建团队成员命令仓储实例
func NewTeamMemberCommandRepository(db *gorm.DB) org.TeamMemberCommandRepository {
	return &teamMemberCommandRepository{db: db}
}

// Add 添加成员到团队
func (r *teamMemberCommandRepository) Add(ctx context.Context, member *org.TeamMember) error {
	model := newTeamMemberModelFromEntity(member)
	if model.JoinedAt.IsZero() {
		model.JoinedAt = time.Now()
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to add team member: %w", err)
	}

	// 回写生成的 ID
	member.ID = model.ID
	member.CreatedAt = model.CreatedAt

	return nil
}

// Remove 从团队移除成员
func (r *teamMemberCommandRepository) Remove(ctx context.Context, teamID, userID uint) error {
	result := r.db.WithContext(ctx).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Delete(&TeamMemberModel{})

	if result.Error != nil {
		return fmt.Errorf("failed to remove team member: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return org.ErrNotTeamMember
	}

	return nil
}

// UpdateRole 更新团队成员角色
func (r *teamMemberCommandRepository) UpdateRole(ctx context.Context, teamID, userID uint, role org.TeamMemberRole) error {
	result := r.db.WithContext(ctx).Model(&TeamMemberModel{}).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Update("role", string(role))

	if result.Error != nil {
		return fmt.Errorf("failed to update team member role: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return org.ErrNotTeamMember
	}

	return nil
}
