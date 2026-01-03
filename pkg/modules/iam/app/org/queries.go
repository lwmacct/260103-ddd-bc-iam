package org

// ============================================================================
// Organization Queries
// ============================================================================

// GetOrgQuery 获取组织查询
type GetOrgQuery struct {
	OrgID uint
}

// ListOrgsQuery 组织列表查询
type ListOrgsQuery struct {
	Offset  int
	Limit   int
	Keyword string // 可选：搜索关键词
}

// ============================================================================
// Team Queries
// ============================================================================

// GetTeamQuery 获取团队查询
type GetTeamQuery struct {
	OrgID  uint
	TeamID uint
}

// ListTeamsQuery 团队列表查询
type ListTeamsQuery struct {
	OrgID  uint
	Offset int
	Limit  int
}

// ============================================================================
// Member Queries
// ============================================================================

// ListMembersQuery 组织成员列表查询
type ListMembersQuery struct {
	OrgID  uint
	Offset int
	Limit  int
}

// ============================================================================
// Team Member Queries
// ============================================================================

// ListTeamMembersQuery 团队成员列表查询
type ListTeamMembersQuery struct {
	OrgID  uint
	TeamID uint
	Offset int
	Limit  int
}

// ============================================================================
// User View Queries (用户视角)
// ============================================================================

// ListUserOrgsQuery 查询用户加入的组织
type ListUserOrgsQuery struct {
	UserID uint
}

// ListUserTeamsQuery 查询用户加入的团队
type ListUserTeamsQuery struct {
	UserID uint
	OrgID  uint // 可选：限定在某个组织内
}
