package org

// ============================================================================
// Organization Commands
// ============================================================================

// CreateOrgCommand 创建组织命令
type CreateOrgCommand struct {
	Name        string
	DisplayName string
	Description string
	Avatar      string
	OwnerUserID uint // 创建者自动成为 owner
}

// UpdateOrgCommand 更新组织命令
type UpdateOrgCommand struct {
	OrgID       uint
	DisplayName *string
	Description *string
	Avatar      *string
	Status      *string
}

// DeleteOrgCommand 删除组织命令
type DeleteOrgCommand struct {
	OrgID uint
}

// ============================================================================
// Team Commands
// ============================================================================

// CreateTeamCommand 创建团队命令
type CreateTeamCommand struct {
	OrgID       uint
	Name        string
	DisplayName string
	Description string
	Avatar      string
	LeadUserID  uint // 可选：指定团队负责人
}

// UpdateTeamCommand 更新团队命令
type UpdateTeamCommand struct {
	OrgID       uint // 用于验证团队归属
	TeamID      uint
	DisplayName *string
	Description *string
	Avatar      *string
}

// DeleteTeamCommand 删除团队命令
type DeleteTeamCommand struct {
	OrgID  uint // 用于验证团队归属
	TeamID uint
}

// ============================================================================
// Member Commands
// ============================================================================

// AddMemberCommand 添加组织成员命令
type AddMemberCommand struct {
	OrgID  uint
	UserID uint
	Role   string // owner, admin, member
}

// RemoveMemberCommand 移除组织成员命令
type RemoveMemberCommand struct {
	OrgID  uint
	UserID uint
}

// UpdateMemberRoleCommand 更新成员角色命令
type UpdateMemberRoleCommand struct {
	OrgID  uint
	UserID uint
	Role   string
}

// ============================================================================
// Team Member Commands
// ============================================================================

// AddTeamMemberCommand 添加团队成员命令
type AddTeamMemberCommand struct {
	OrgID  uint // 用于验证
	TeamID uint
	UserID uint
	Role   string // lead, member
}

// RemoveTeamMemberCommand 移除团队成员命令
type RemoveTeamMemberCommand struct {
	OrgID  uint // 用于验证
	TeamID uint
	UserID uint
}

// UpdateTeamMemberRoleCommand 更新团队成员角色命令
type UpdateTeamMemberRoleCommand struct {
	OrgID  uint
	TeamID uint
	UserID uint
	Role   string
}
