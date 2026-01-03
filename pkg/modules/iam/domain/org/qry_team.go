package org

import "context"

// TeamQueryRepository 团队查询仓储接口
type TeamQueryRepository interface {
	// GetByID 根据 ID 获取团队
	GetByID(ctx context.Context, id uint) (*Team, error)

	// GetByOrgAndName 根据组织 ID 和团队名称获取团队
	GetByOrgAndName(ctx context.Context, orgID uint, name string) (*Team, error)

	// GetByIDWithMembers 根据 ID 获取团队（包含成员列表）
	GetByIDWithMembers(ctx context.Context, id uint) (*Team, error)

	// ListByOrg 获取组织的所有团队
	ListByOrg(ctx context.Context, orgID uint, offset, limit int) ([]*Team, error)

	// CountByOrg 统计组织的团队数量
	CountByOrg(ctx context.Context, orgID uint) (int64, error)

	// Exists 检查团队是否存在
	Exists(ctx context.Context, id uint) (bool, error)

	// ExistsByOrgAndName 检查组织内团队名称是否存在
	ExistsByOrgAndName(ctx context.Context, orgID uint, name string) (bool, error)

	// BelongsToOrg 检查团队是否属于指定组织
	BelongsToOrg(ctx context.Context, teamID, orgID uint) (bool, error)
}
