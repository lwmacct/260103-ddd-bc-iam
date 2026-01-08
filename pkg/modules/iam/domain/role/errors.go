package role

import "errors"

// 角色相关错误
var (
	// ErrRoleNotFound 角色不存在
	ErrRoleNotFound = errors.New("role not found")

	// ErrRoleNameExists 角色名称已存在
	ErrRoleNameExists = errors.New("role name already exists")

	// ErrCannotDeleteSystemRole 不能删除系统角色
	ErrCannotDeleteSystemRole = errors.New("cannot delete system role")

	// ErrCannotModifySystemRole 不能修改系统角色
	ErrCannotModifySystemRole = errors.New("cannot modify system role")

	// ErrInvalidRoleName 无效的角色名称
	ErrInvalidRoleName = errors.New("invalid role name")

	// ErrInvalidRoleID 无效的角色 ID
	ErrInvalidRoleID = errors.New("invalid role ID")

	// ErrRoleHasUsers 角色下有关联用户
	ErrRoleHasUsers = errors.New("role has associated users")
)

// 权限相关错误
var (
	// ErrPermissionNotFound 权限不存在
	ErrPermissionNotFound = errors.New("permission not found")

	// ErrPermissionCodeExists 权限代码已存在
	ErrPermissionCodeExists = errors.New("permission code already exists")

	// ErrInvalidPermissionCode 无效的权限代码格式
	ErrInvalidPermissionCode = errors.New("invalid permission code format")

	// ErrPermissionInUse 权限正在被使用
	ErrPermissionInUse = errors.New("permission is in use by roles")
)
