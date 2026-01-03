package role

import (
	"time"

	"github.com/lwmacct/260101-go-pkg-gin/pkg/permission"
)

// ============================================================================
// Permission 值对象（Operation + Resource 组合）
// ============================================================================

// Permission 权限条目，定义 Operation 对 Resource 的访问权限。
//
// AWS IAM 风格设计：
//   - OperationPattern: 操作模式，如 "sys:users.create" 或 "sys:*"
//   - ResourcePattern:  资源模式，如 "user/123" 或 "*"
//
// 通配符支持：
//   - Operation: "sys:*", "sys:users.*", "*:*.read"
//   - Resource:  "user/*", "*"
type Permission struct {
	OperationPattern string `json:"operation_pattern"` // 操作模式
	ResourcePattern  string `json:"resource_pattern"`  // 资源模式，默认 "*"
}

// NewPermission 创建权限条目。
func NewPermission(opPattern, resPattern string) Permission {
	if resPattern == "" {
		resPattern = "*"
	}
	return Permission{
		OperationPattern: opPattern,
		ResourcePattern:  resPattern,
	}
}

// Matches 检查权限是否匹配指定操作和资源。
func (p Permission) Matches(op permission.Operation, res permission.Resource) bool {
	return permission.MatchOperation(p.OperationPattern, string(op)) &&
		permission.MatchResource(p.ResourcePattern, string(res))
}

// ============================================================================
// Role 实体
// ============================================================================

// Role 角色实体，RBAC 系统的核心组件。
//
// Operation-Centric RBAC：
//   - 角色直接关联操作模式，无需独立的权限实体
//   - 支持通配符实现细粒度和粗粒度权限控制
type Role struct {
	ID          uint         `json:"id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   *time.Time   `json:"deleted_at,omitempty"`
	Name        string       `json:"name"`
	DisplayName string       `json:"display_name"`
	Description string       `json:"description"`
	IsSystem    bool         `json:"is_system"`
	Permissions []Permission `json:"permissions"`
}

// IsSystemRole 检查是否为系统角色。
func (r *Role) IsSystemRole() bool {
	return r.IsSystem
}

// CanBeDeleted 检查角色是否可以被删除（系统角色不可删除）。
func (r *Role) CanBeDeleted() bool {
	return !r.IsSystem
}

// CanBeModified 检查角色是否可以被修改（系统角色不可修改）。
func (r *Role) CanBeModified() bool {
	return !r.IsSystem
}

// HasPermission 检查角色是否有指定操作对指定资源的权限。
func (r *Role) HasPermission(op permission.Operation, res permission.Resource) bool {
	for _, p := range r.Permissions {
		if p.Matches(op, res) {
			return true
		}
	}
	return false
}

// HasOperationPermission 检查角色是否有指定操作的权限（资源为 *）。
// 这是 HasPermission 的简化版本，用于不需要资源级别控制的场景。
func (r *Role) HasOperationPermission(op permission.Operation) bool {
	return r.HasPermission(op, permission.ResourceAll)
}

// GetPermissionCount 获取权限数量。
func (r *Role) GetPermissionCount() int {
	return len(r.Permissions)
}

// IsEmpty 检查角色是否没有权限。
func (r *Role) IsEmpty() bool {
	return len(r.Permissions) == 0
}

// AddPermission 添加权限到角色。
func (r *Role) AddPermission(p Permission) {
	// 检查是否已存在相同的权限
	for _, existing := range r.Permissions {
		if existing.OperationPattern == p.OperationPattern &&
			existing.ResourcePattern == p.ResourcePattern {
			return
		}
	}
	r.Permissions = append(r.Permissions, p)
}

// RemovePermission 从角色移除权限。
func (r *Role) RemovePermission(opPattern, resPattern string) bool {
	for i, p := range r.Permissions {
		if p.OperationPattern == opPattern && p.ResourcePattern == resPattern {
			r.Permissions = append(r.Permissions[:i], r.Permissions[i+1:]...)
			return true
		}
	}
	return false
}

// ClearPermissions 清空所有权限。
func (r *Role) ClearPermissions() {
	r.Permissions = nil
}

// SetPermissions 设置角色的权限列表（替换现有权限）。
func (r *Role) SetPermissions(permissions []Permission) {
	r.Permissions = permissions
}
