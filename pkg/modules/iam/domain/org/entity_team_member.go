package org

import "time"

// TeamMember 团队成员实体（用户-团队关联）
type TeamMember struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`

	// TeamID 所属团队 ID
	TeamID uint `json:"team_id"`

	// UserID 用户 ID
	UserID uint `json:"user_id"`

	// Role 团队成员角色: lead, member
	Role TeamMemberRole `json:"role"`

	// JoinedAt 加入时间
	JoinedAt time.Time `json:"joined_at"`
}

// TeamMemberRole 团队成员角色
type TeamMemberRole string

const (
	// TeamMemberRoleLead 团队负责人，可管理团队成员。
	TeamMemberRoleLead TeamMemberRole = "lead"

	// TeamMemberRoleMember 团队普通成员。
	TeamMemberRoleMember TeamMemberRole = "member"
)

// IsLead 报告成员是否为团队负责人。
func (tm *TeamMember) IsLead() bool {
	return tm.Role == TeamMemberRoleLead
}

// CanManageTeamMembers 报告成员是否可以管理团队成员。
// 只有 lead 可以添加/移除团队成员。
func (tm *TeamMember) CanManageTeamMembers() bool {
	return tm.IsLead()
}

// ValidTeamMemberRoles 返回所有有效的团队成员角色。
func ValidTeamMemberRoles() []TeamMemberRole {
	return []TeamMemberRole{TeamMemberRoleLead, TeamMemberRoleMember}
}

// IsValidTeamMemberRole 检查团队成员角色是否有效。
func IsValidTeamMemberRole(role TeamMemberRole) bool {
	switch role {
	case TeamMemberRoleLead, TeamMemberRoleMember:
		return true
	default:
		return false
	}
}
