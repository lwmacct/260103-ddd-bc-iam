package org

import "context"

// TeamMemberCommandRepository 团队成员命令仓储接口
type TeamMemberCommandRepository interface {
	// Add 添加成员到团队
	Add(ctx context.Context, member *TeamMember) error

	// Remove 从团队移除成员
	Remove(ctx context.Context, teamID, userID uint) error

	// UpdateRole 更新团队成员角色
	UpdateRole(ctx context.Context, teamID, userID uint, role TeamMemberRole) error
}
