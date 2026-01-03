package org

import "context"

// MemberQueryRepository 组织成员查询仓储接口
type MemberQueryRepository interface {
	// GetByOrgAndUser 获取指定组织的指定用户成员信息
	GetByOrgAndUser(ctx context.Context, orgID, userID uint) (*Member, error)

	// ListByOrg 获取组织的所有成员
	ListByOrg(ctx context.Context, orgID uint, offset, limit int) ([]*Member, error)

	// CountByOrg 统计组织成员数量
	CountByOrg(ctx context.Context, orgID uint) (int64, error)

	// ListByUser 获取用户加入的所有组织的成员记录
	ListByUser(ctx context.Context, userID uint) ([]*Member, error)

	// IsMember 检查用户是否为组织成员
	IsMember(ctx context.Context, orgID, userID uint) (bool, error)

	// GetOwner 获取组织所有者
	GetOwner(ctx context.Context, orgID uint) (*Member, error)

	// CountOwners 统计组织所有者数量
	CountOwners(ctx context.Context, orgID uint) (int64, error)

	// ListByOrgWithUsers 获取组织成员列表（包含用户信息）
	// 用于成员列表展示，一次 JOIN 查询获取完整数据
	ListByOrgWithUsers(ctx context.Context, orgID uint, offset, limit int) ([]*MemberWithUser, error)
}
