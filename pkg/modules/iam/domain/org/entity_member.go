package org

import "time"

// Member 组织成员实体（用户-组织关联）
type Member struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// OrgID 所属组织 ID
	OrgID uint `json:"org_id"`

	// UserID 用户 ID
	UserID uint `json:"user_id"`

	// Role 成员角色: owner, admin, member
	Role MemberRole `json:"role"`

	// JoinedAt 加入时间
	JoinedAt time.Time `json:"joined_at"`
}

// MemberRole 组织成员角色
type MemberRole string

const (
	// MemberRoleOwner 组织所有者，拥有最高权限，可管理所有成员和设置。
	MemberRoleOwner MemberRole = "owner"

	// MemberRoleAdmin 组织管理员，可管理成员和团队。
	MemberRoleAdmin MemberRole = "admin"

	// MemberRoleMember 普通成员，基础访问权限。
	MemberRoleMember MemberRole = "member"
)

// IsOwner 报告成员是否为组织所有者。
func (m *Member) IsOwner() bool {
	return m.Role == MemberRoleOwner
}

// IsAdmin 报告成员是否为管理员（包括所有者）。
func (m *Member) IsAdmin() bool {
	return m.Role == MemberRoleOwner || m.Role == MemberRoleAdmin
}

// CanManageMembers 报告成员是否可以管理其他成员。
// 只有 owner 和 admin 可以管理成员。
func (m *Member) CanManageMembers() bool {
	return m.IsAdmin()
}

// CanManageTeams 报告成员是否可以管理团队。
// 只有 owner 和 admin 可以创建/删除团队。
func (m *Member) CanManageTeams() bool {
	return m.IsAdmin()
}

// CanTransferOwnership 报告成员是否可以转让所有权。
// 只有 owner 可以转让所有权。
func (m *Member) CanTransferOwnership() bool {
	return m.IsOwner()
}

// ValidMemberRoles 返回所有有效的成员角色。
func ValidMemberRoles() []MemberRole {
	return []MemberRole{MemberRoleOwner, MemberRoleAdmin, MemberRoleMember}
}

// IsValidMemberRole 检查角色是否有效。
func IsValidMemberRole(role MemberRole) bool {
	switch role {
	case MemberRoleOwner, MemberRoleAdmin, MemberRoleMember:
		return true
	default:
		return false
	}
}
