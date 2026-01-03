package org

import "context"

// TeamMemberQueryRepository 团队成员查询仓储接口
type TeamMemberQueryRepository interface {
	// GetByTeamAndUser 获取指定团队的指定用户成员信息
	GetByTeamAndUser(ctx context.Context, teamID, userID uint) (*TeamMember, error)

	// ListByTeam 获取团队的所有成员
	ListByTeam(ctx context.Context, teamID uint, offset, limit int) ([]*TeamMember, error)

	// CountByTeam 统计团队成员数量
	CountByTeam(ctx context.Context, teamID uint) (int64, error)

	// ListByUser 获取用户加入的所有团队的成员记录
	ListByUser(ctx context.Context, userID uint) ([]*TeamMember, error)

	// ListByUserInOrg 获取用户在指定组织内加入的所有团队的成员记录
	ListByUserInOrg(ctx context.Context, userID, orgID uint) ([]*TeamMember, error)

	// IsMember 检查用户是否为团队成员
	IsMember(ctx context.Context, teamID, userID uint) (bool, error)
}
