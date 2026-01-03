package event

import (
	"strconv"
)

// ============================================================================
// 角色事件
// ============================================================================

// RoleCreatedEvent 角色创建事件
type RoleCreatedEvent struct {
	BaseEvent

	RoleID   uint   `json:"role_id"`
	RoleName string `json:"role_name"`
}

// NewRoleCreatedEvent 创建角色创建事件
func NewRoleCreatedEvent(roleID uint, roleName string) *RoleCreatedEvent {
	return &RoleCreatedEvent{
		BaseEvent: NewBaseEvent("role.created", "role", strconv.FormatUint(uint64(roleID), 10)),
		RoleID:    roleID,
		RoleName:  roleName,
	}
}

// RoleUpdatedEvent 角色更新事件
type RoleUpdatedEvent struct {
	BaseEvent

	RoleID uint `json:"role_id"`
}

// NewRoleUpdatedEvent 创建角色更新事件
func NewRoleUpdatedEvent(roleID uint) *RoleUpdatedEvent {
	return &RoleUpdatedEvent{
		BaseEvent: NewBaseEvent("role.updated", "role", strconv.FormatUint(uint64(roleID), 10)),
		RoleID:    roleID,
	}
}

// RoleDeletedEvent 角色删除事件
type RoleDeletedEvent struct {
	BaseEvent

	RoleID   uint   `json:"role_id"`
	RoleName string `json:"role_name"`
}

// NewRoleDeletedEvent 创建角色删除事件
func NewRoleDeletedEvent(roleID uint, roleName string) *RoleDeletedEvent {
	return &RoleDeletedEvent{
		BaseEvent: NewBaseEvent("role.deleted", "role", strconv.FormatUint(uint64(roleID), 10)),
		RoleID:    roleID,
		RoleName:  roleName,
	}
}

// RolePermissionsChangedEvent 角色权限变更事件
// 这是重要的缓存失效触发事件
type RolePermissionsChangedEvent struct {
	BaseEvent

	RoleID        uint   `json:"role_id"`
	PermissionIDs []uint `json:"permission_ids"`
}

// NewRolePermissionsChangedEvent 创建角色权限变更事件
func NewRolePermissionsChangedEvent(roleID uint, permissionIDs []uint) *RolePermissionsChangedEvent {
	return &RolePermissionsChangedEvent{
		BaseEvent:     NewBaseEvent("role.permissions_changed", "role", strconv.FormatUint(uint64(roleID), 10)),
		RoleID:        roleID,
		PermissionIDs: permissionIDs,
	}
}
