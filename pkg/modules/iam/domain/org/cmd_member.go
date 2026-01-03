package org

import "context"

// MemberCommandRepository 组织成员命令仓储接口
type MemberCommandRepository interface {
	// Add 添加成员到组织
	Add(ctx context.Context, member *Member) error

	// Remove 从组织移除成员
	Remove(ctx context.Context, orgID, userID uint) error

	// UpdateRole 更新成员角色
	UpdateRole(ctx context.Context, orgID, userID uint, role MemberRole) error
}
