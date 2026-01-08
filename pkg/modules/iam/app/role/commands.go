package role

import (
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/role"
)

// CreateCommand 创建角色命令
type CreateCommand struct {
	Name        string // 角色名称（唯一）
	DisplayName string // 显示名称
	Description string // 描述
}

// UpdateCommand 更新角色命令
type UpdateCommand struct {
	RoleID      uint
	DisplayName *string // 可选：显示名称
	Description *string // 可选：描述
}

// DeleteCommand 删除角色命令
type DeleteCommand struct {
	RoleID uint
}

// SetPermissionsCommand 设置角色权限命令
// 新 RBAC 模型：使用 Permission 模式（OperationPattern + ResourcePattern）
type SetPermissionsCommand struct {
	RoleID      uint
	Permissions []role.Permission
}
