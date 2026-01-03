package org

import "time"

// ============================================================================
// Organization DTOs
// ============================================================================

// CreateOrgDTO 创建组织 DTO
type CreateOrgDTO struct {
	Name        string `json:"name" binding:"required,min=2,max=50,loweralphanumhyphen"`
	DisplayName string `json:"display_name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"max=500"`
	Avatar      string `json:"avatar" binding:"omitempty,max=255,url"`
}

// UpdateOrgDTO 更新组织 DTO
type UpdateOrgDTO struct {
	DisplayName *string `json:"display_name" binding:"omitempty,min=2,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	Avatar      *string `json:"avatar" binding:"omitempty,max=255"`
	Status      *string `json:"status" binding:"omitempty,oneof=active suspended"`
}

// OrgDTO 组织响应 DTO
type OrgDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	Avatar      string    `json:"avatar"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OrgWithStatsDTO 组织响应 DTO（包含统计信息）
type OrgWithStatsDTO struct {
	OrgDTO

	MemberCount int64 `json:"member_count"`
	TeamCount   int64 `json:"team_count"`
}

// ============================================================================
// Team DTOs
// ============================================================================

// CreateTeamDTO 创建团队 DTO
type CreateTeamDTO struct {
	Name        string `json:"name" binding:"required,min=2,max=50,loweralphanumhyphen"`
	DisplayName string `json:"display_name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"max=500"`
	Avatar      string `json:"avatar" binding:"omitempty,max=255,url"`
}

// UpdateTeamDTO 更新团队 DTO
type UpdateTeamDTO struct {
	DisplayName *string `json:"display_name" binding:"omitempty,min=2,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	Avatar      *string `json:"avatar" binding:"omitempty,max=255"`
}

// TeamDTO 团队响应 DTO
type TeamDTO struct {
	ID          uint      `json:"id"`
	OrgID       uint      `json:"org_id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	Avatar      string    `json:"avatar"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TeamWithStatsDTO 团队响应 DTO（包含统计信息）
type TeamWithStatsDTO struct {
	TeamDTO

	MemberCount int64 `json:"member_count"`
}

// ============================================================================
// Member DTOs
// ============================================================================

// AddMemberDTO 添加成员 DTO
type AddMemberDTO struct {
	UserID uint   `json:"user_id" binding:"required,gt=0"`
	Role   string `json:"role" binding:"required,oneof=owner admin member"`
}

// UpdateMemberRoleDTO 更新成员角色 DTO
type UpdateMemberRoleDTO struct {
	Role string `json:"role" binding:"required,oneof=owner admin member"`
}

// MemberDTO 组织成员响应 DTO
type MemberDTO struct {
	ID       uint      `json:"id"`
	OrgID    uint      `json:"org_id"`
	UserID   uint      `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
	// 关联的用户信息（可选加载）
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

// ============================================================================
// Team Member DTOs
// ============================================================================

// AddTeamMemberDTO 添加团队成员 DTO
type AddTeamMemberDTO struct {
	UserID uint   `json:"user_id" binding:"required,gt=0"`
	Role   string `json:"role" binding:"required,oneof=lead member"`
}

// UpdateTeamMemberRoleDTO 更新团队成员角色 DTO
type UpdateTeamMemberRoleDTO struct {
	Role string `json:"role" binding:"required,oneof=lead member"`
}

// TeamMemberDTO 团队成员响应 DTO
type TeamMemberDTO struct {
	ID       uint      `json:"id"`
	TeamID   uint      `json:"team_id"`
	UserID   uint      `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
	// 关联的用户信息（可选加载）
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

// ============================================================================
// User View DTOs (用户视角)
// ============================================================================

// UserOrgDTO 用户加入的组织 DTO（用户视角）
type UserOrgDTO struct {
	OrgDTO

	Role     string    `json:"role"`      // 用户在该组织中的角色
	JoinedAt time.Time `json:"joined_at"` // 加入时间
}

// UserTeamDTO 用户加入的团队 DTO（用户视角）
type UserTeamDTO struct {
	TeamDTO

	OrgName  string    `json:"org_name"`  // 所属组织名称
	Role     string    `json:"role"`      // 用户在该团队中的角色
	JoinedAt time.Time `json:"joined_at"` // 加入时间
}
