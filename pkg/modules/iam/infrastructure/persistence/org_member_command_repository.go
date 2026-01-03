package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/org"
	"gorm.io/gorm"
)

// orgMemberCommandRepository 组织成员命令仓储的 GORM 实现
type orgMemberCommandRepository struct {
	db *gorm.DB
}

// NewOrgMemberCommandRepository 创建组织成员命令仓储实例
func NewOrgMemberCommandRepository(db *gorm.DB) org.MemberCommandRepository {
	return &orgMemberCommandRepository{db: db}
}

// Add 添加成员到组织
func (r *orgMemberCommandRepository) Add(ctx context.Context, member *org.Member) error {
	model := newOrgMemberModelFromEntity(member)
	if model.JoinedAt.IsZero() {
		model.JoinedAt = time.Now()
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to add organization member: %w", err)
	}

	// 回写生成的 ID
	member.ID = model.ID
	member.CreatedAt = model.CreatedAt
	member.UpdatedAt = model.UpdatedAt

	return nil
}

// Remove 从组织移除成员
func (r *orgMemberCommandRepository) Remove(ctx context.Context, orgID, userID uint) error {
	result := r.db.WithContext(ctx).
		Where("org_id = ? AND user_id = ?", orgID, userID).
		Delete(&OrgMemberModel{})

	if result.Error != nil {
		return fmt.Errorf("failed to remove organization member: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return org.ErrMemberNotFound
	}

	return nil
}

// UpdateRole 更新成员角色
func (r *orgMemberCommandRepository) UpdateRole(ctx context.Context, orgID, userID uint, role org.MemberRole) error {
	result := r.db.WithContext(ctx).Model(&OrgMemberModel{}).
		Where("org_id = ? AND user_id = ?", orgID, userID).
		Update("role", string(role))

	if result.Error != nil {
		return fmt.Errorf("failed to update member role: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return org.ErrMemberNotFound
	}

	return nil
}
