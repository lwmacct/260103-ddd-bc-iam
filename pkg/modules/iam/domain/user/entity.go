package user

import (
	"slices"
	"time"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/role"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/permission"
)

// User 用户实体
type User struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// Extra 扩展数据（JSONB 存储），留空备用
	Extra map[string]any `json:"extra,omitempty"`

	Username  string  `json:"username"`
	Email     *string `json:"email,omitempty"` // nullable，用于登录/注册
	Password  string  `json:"-"`               // 敏感字段，不序列化
	RealName  string  `json:"real_name"`       // 真实姓名
	Nickname  string  `json:"nickname"`        // 昵称
	Phone     *string `json:"phone,omitempty"` // nullable，用于登录/注册
	Signature string  `json:"signature"`       // 个性签名
	Avatar    string  `json:"avatar"`
	Bio       string  `json:"bio"` // 个人简介
	Status    string  `json:"status"`

	// Type 用户类型：human（人类用户）、service（服务账户）、system（系统用户）。
	// 系统用户（Type=system）不可删除，部分字段不可修改。
	Type UserType `json:"type"`

	// RBAC: Many-to-Many relationship with roles
	Roles []role.Role `json:"roles"`
}

// HasRole 检查用户是否拥有指定角色
func (u *User) HasRole(roleName string) bool {
	for _, r := range u.Roles {
		if r.Name == roleName {
			return true
		}
	}
	return false
}

// HasAnyRole 检查用户是否拥有任一指定角色
func (u *User) HasAnyRole(roleNames ...string) bool {
	return slices.ContainsFunc(roleNames, u.HasRole)
}

// HasPermission 检查用户是否有指定操作对指定资源的权限。
// 遍历用户所有角色的权限进行模式匹配。
func (u *User) HasPermission(op permission.Operation, res permission.Resource) bool {
	for _, r := range u.Roles {
		if r.HasPermission(op, res) {
			return true
		}
	}
	return false
}

// HasOperationPermission 检查用户是否有指定操作的权限（资源为 *）。
func (u *User) HasOperationPermission(op permission.Operation) bool {
	return u.HasPermission(op, permission.ResourceAll)
}

// GetRoleNames 获取用户所有角色名称
func (u *User) GetRoleNames() []string {
	names := make([]string, 0, len(u.Roles))
	for _, r := range u.Roles {
		names = append(names, r.Name)
	}
	return names
}

// GetPermissions 获取用户所有去重后的权限。
// 基于 OperationPattern + ResourcePattern 去重。
func (u *User) GetPermissions() []role.Permission {
	seen := make(map[string]bool)
	var permissions []role.Permission

	for _, r := range u.Roles {
		for _, p := range r.Permissions {
			key := p.OperationPattern + ":" + p.ResourcePattern
			if !seen[key] {
				seen[key] = true
				permissions = append(permissions, p)
			}
		}
	}

	return permissions
}

// IsAdmin 检查用户是否拥有管理员角色
func (u *User) IsAdmin() bool {
	return u.HasRole("admin")
}

// CanLogin 检查用户是否可以登录
func (u *User) CanLogin() bool {
	return u.Status == "active"
}

// IsBanned 检查用户是否被禁用
func (u *User) IsBanned() bool {
	return u.Status == "banned"
}

// IsInactive 检查用户是否未激活
func (u *User) IsInactive() bool {
	return u.Status == "inactive"
}

// Activate 激活用户
func (u *User) Activate() {
	u.Status = "active"
}

// Deactivate 停用用户
func (u *User) Deactivate() {
	u.Status = "inactive"
}

// Ban 禁用用户
func (u *User) Ban() {
	u.Status = "banned"
}

// AssignRole 分配角色（领域行为）
func (u *User) AssignRole(r role.Role) error {
	if u.HasRole(r.Name) {
		return ErrRoleAlreadyAssigned
	}
	u.Roles = append(u.Roles, r)
	return nil
}

// RemoveRole 移除角色（领域行为）
func (u *User) RemoveRole(roleName string) error {
	for i, r := range u.Roles {
		if r.Name == roleName {
			u.Roles = append(u.Roles[:i], u.Roles[i+1:]...)
			return nil
		}
	}
	return ErrRoleNotFound
}

// ClearRoles 清空所有角色
func (u *User) ClearRoles() {
	u.Roles = []role.Role{}
}

// UpdateProfile 更新用户资料（领域行为）
func (u *User) UpdateProfile(realName, nickname, phone, signature, avatar, bio string) {
	if realName != "" {
		u.RealName = realName
	}
	if nickname != "" {
		u.Nickname = nickname
	}
	if phone != "" {
		u.Phone = &phone
	}
	u.Signature = signature // Signature 可以为空
	if avatar != "" {
		u.Avatar = avatar
	}
	u.Bio = bio // Bio 可以为空
}

// ============================================================================
// 类型判断方法
// ============================================================================

// IsHuman 报告用户是否为人类用户。
func (u *User) IsHuman() bool {
	return u.Type == UserTypeHuman
}

// IsServiceAccount 报告用户是否为服务账户。
func (u *User) IsServiceAccount() bool {
	return u.Type == UserTypeService
}

// IsSystemUser 报告用户是否为系统预置用户。
// 系统用户不可删除，部分字段不可修改。
func (u *User) IsSystemUser() bool {
	return u.Type == UserTypeSystem
}

// IsSystem 报告用户是否为系统用户（IsSystemUser 的别名）
func (u *User) IsSystem() bool {
	return u.Type == UserTypeSystem
}

// IsRoot 报告用户是否为 root 超级管理员。
func (u *User) IsRoot() bool {
	return u.Username == RootUsername
}

// ============================================================================
// 保护策略方法
// ============================================================================

// CanBeDeleted 报告用户是否可以被删除。
// 系统用户不可删除。
func (u *User) CanBeDeleted() bool {
	return !u.IsSystem()
}

// CanModifyUsername 报告用户名是否可以被修改。
// 系统用户的用户名不可修改。
func (u *User) CanModifyUsername() bool {
	return !u.IsSystem()
}

// CanModifyStatus 报告用户状态是否可以被修改。
// 仅 root 用户状态不可修改。
func (u *User) CanModifyStatus() bool {
	return u.Username != RootUsername
}

// CanModifyRoles 报告用户角色是否可以被修改。
// root 用户角色不可修改（始终拥有 *:*:* 权限）。
func (u *User) CanModifyRoles() bool {
	return u.Username != RootUsername
}

// ============================================================================
// 认证相关方法
// ============================================================================

// CanPasswordLogin 报告用户是否可以使用密码登录。
// 服务账户不支持密码登录。
func (u *User) CanPasswordLogin() bool {
	return u.IsHuman() && u.CanLogin()
}

// RequiresPAT 报告用户是否必须使用 PAT 认证。
// 服务账户仅支持 PAT 认证。
func (u *User) RequiresPAT() bool {
	return u.IsServiceAccount()
}
