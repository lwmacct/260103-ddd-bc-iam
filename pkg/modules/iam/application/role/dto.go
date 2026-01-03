package role

import "time"

// CreateDTO 创建角色请求 DTO
type CreateDTO struct {
	Name        string `json:"name" binding:"required,min=2,max=50" example:"developer"`
	DisplayName string `json:"display_name" binding:"required,max=100" example:"开发者"`
	Description string `json:"description" binding:"max=255" example:"系统开发人员角色"`
}

// UpdateDTO 更新角色请求 DTO
type UpdateDTO struct {
	DisplayName *string `json:"display_name,omitempty" binding:"omitempty,max=100"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=255"`
}

// PermissionInputDTO 权限输入 DTO（用于设置权限）
type PermissionInputDTO struct {
	OperationPattern string `json:"operation_pattern" binding:"required" example:"sys:users.*"`
	ResourcePattern  string `json:"resource_pattern" binding:"omitempty" example:"user/*"`
}

// SetPermissionsDTO 设置角色权限请求 DTO
// 新 RBAC 模型：使用 Operation + Resource Pattern
type SetPermissionsDTO struct {
	Permissions []PermissionInputDTO `json:"permissions" binding:"required"`
}

// CreateResultDTO 创建角色响应 DTO
type CreateResultDTO struct {
	RoleID      uint   `json:"role_id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

// RoleDTO 角色响应 DTO
type RoleDTO struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	DisplayName string           `json:"display_name"`
	Description string           `json:"description"`
	IsSystem    bool             `json:"is_system"`
	Permissions []*PermissionDTO `json:"permissions,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// PermissionDTO 权限响应 DTO
// 新 RBAC 模型：Operation + Resource Pattern
type PermissionDTO struct {
	OperationPattern string `json:"operation_pattern" example:"sys:users.*"`
	ResourcePattern  string `json:"resource_pattern" example:"user/*"`
}

// ListRolesDTO 角色列表响应 DTO
type ListRolesDTO struct {
	Roles []*RoleDTO `json:"roles"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
}
