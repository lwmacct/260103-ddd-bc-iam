package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/org"
	"gorm.io/gorm"
)

// orgMemberQueryRepository 组织成员查询仓储的 GORM 实现
type orgMemberQueryRepository struct {
	db *gorm.DB
}

// NewOrgMemberQueryRepository 创建组织成员查询仓储实例
func NewOrgMemberQueryRepository(db *gorm.DB) org.MemberQueryRepository {
	return &orgMemberQueryRepository{db: db}
}

// GetByOrgAndUser 获取指定组织的指定用户成员信息
func (r *orgMemberQueryRepository) GetByOrgAndUser(ctx context.Context, orgID, userID uint) (*org.Member, error) {
	var model OrgMemberModel
	if err := r.db.WithContext(ctx).
		Where("org_id = ? AND user_id = ?", orgID, userID).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, org.ErrMemberNotFound
		}
		return nil, fmt.Errorf("failed to get org member: %w", err)
	}
	return model.ToEntity(), nil
}

// ListByOrg 获取组织的所有成员
func (r *orgMemberQueryRepository) ListByOrg(ctx context.Context, orgID uint, offset, limit int) ([]*org.Member, error) {
	var models []*OrgMemberModel
	if err := r.db.WithContext(ctx).
		Where("org_id = ?", orgID).
		Offset(offset).Limit(limit).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list org members: %w", err)
	}
	return mapOrgMemberModelPtrsToEntities(models), nil
}

// CountByOrg 统计组织成员数量
func (r *orgMemberQueryRepository) CountByOrg(ctx context.Context, orgID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&OrgMemberModel{}).
		Where("org_id = ?", orgID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count org members: %w", err)
	}
	return count, nil
}

// ListByUser 获取用户加入的所有组织的成员记录
func (r *orgMemberQueryRepository) ListByUser(ctx context.Context, userID uint) ([]*org.Member, error) {
	var models []*OrgMemberModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list user org memberships: %w", err)
	}
	return mapOrgMemberModelPtrsToEntities(models), nil
}

// IsMember 检查用户是否为组织成员
func (r *orgMemberQueryRepository) IsMember(ctx context.Context, orgID, userID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&OrgMemberModel{}).
		Where("org_id = ? AND user_id = ?", orgID, userID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check membership: %w", err)
	}
	return count > 0, nil
}

// GetOwner 获取组织所有者
func (r *orgMemberQueryRepository) GetOwner(ctx context.Context, orgID uint) (*org.Member, error) {
	var model OrgMemberModel
	if err := r.db.WithContext(ctx).
		Where("org_id = ? AND role = ?", orgID, org.MemberRoleOwner).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, org.ErrMemberNotFound
		}
		return nil, fmt.Errorf("failed to get org owner: %w", err)
	}
	return model.ToEntity(), nil
}

// CountOwners 统计组织所有者数量
func (r *orgMemberQueryRepository) CountOwners(ctx context.Context, orgID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&OrgMemberModel{}).
		Where("org_id = ? AND role = ?", orgID, org.MemberRoleOwner).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count org owners: %w", err)
	}
	return count, nil
}

// memberWithUserRow 包含成员和用户信息的查询结果行
type memberWithUserRow struct {
	// OrgMember 字段
	ID       uint
	OrgID    uint
	UserID   uint
	Role     string
	JoinedAt time.Time
	// User 字段
	Username string
	Email    string
	FullName string
	Avatar   string
}

// ListByOrgWithUsers 获取组织成员列表（包含用户信息）
func (r *orgMemberQueryRepository) ListByOrgWithUsers(ctx context.Context, orgID uint, offset, limit int) ([]*org.MemberWithUser, error) {
	var rows []memberWithUserRow

	err := r.db.WithContext(ctx).
		Table("org_members").
		Select(`
			org_members.id,
			org_members.org_id,
			org_members.user_id,
			org_members.role,
			org_members.joined_at,
			u.username,
			u.email,
			u.real_name,
			u.avatar
		`).
		Joins("LEFT JOIN users u ON org_members.user_id = u.id AND u.deleted_at IS NULL").
		Where("org_members.org_id = ?", orgID).
		Order("org_members.created_at ASC").
		Offset(offset).Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list org members with users: %w", err)
	}

	return mapRowsToMemberWithUser(rows), nil
}

// mapRowsToMemberWithUser 将查询结果行映射为 MemberWithUser 值对象
func mapRowsToMemberWithUser(rows []memberWithUserRow) []*org.MemberWithUser {
	if len(rows) == 0 {
		return []*org.MemberWithUser{}
	}

	result := make([]*org.MemberWithUser, 0, len(rows))
	for _, row := range rows {
		result = append(result, &org.MemberWithUser{
			Member: org.Member{
				ID:        row.ID,
				OrgID:     row.OrgID,
				UserID:    row.UserID,
				Role:      org.MemberRole(row.Role),
				JoinedAt:  row.JoinedAt,
				CreatedAt: time.Time{}, // 查询结果不包含时间戳
				UpdatedAt: time.Time{},
			},
			Username: row.Username,
			Email:    row.Email,
			FullName: row.FullName,
			Avatar:   row.Avatar,
		})
	}
	return result
}
