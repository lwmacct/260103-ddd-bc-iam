package role

import "errors"

// 角色相关错误
var (
	// ErrRoleNotFound 角色不存在
	ErrRoleNotFound = errors.New("角色不存在")

	// ErrRoleNameExists 角色名称已存在
	ErrRoleNameExists = errors.New("角色名称已存在")

	// ErrCannotDeleteSystemRole 不能删除系统角色
	ErrCannotDeleteSystemRole = errors.New("不能删除系统角色")

	// ErrCannotModifySystemRole 不能修改系统角色
	ErrCannotModifySystemRole = errors.New("不能修改系统角色")

	// ErrInvalidRoleName 无效的角色名称
	ErrInvalidRoleName = errors.New("无效的角色名称")

	// ErrInvalidRoleID 无效的角色 ID
	ErrInvalidRoleID = errors.New("无效的角色 ID")

	// ErrRoleHasUsers 角色下有关联用户
	ErrRoleHasUsers = errors.New("角色下有关联用户")
)

// 权限相关错误
var (
	// ErrPermissionNotFound 权限不存在
	ErrPermissionNotFound = errors.New("权限不存在")

	// ErrPermissionCodeExists 权限代码已存在
	ErrPermissionCodeExists = errors.New("权限代码已存在")

	// ErrInvalidPermissionCode 无效的权限代码格式
	ErrInvalidPermissionCode = errors.New("无效的权限代码格式")

	// ErrPermissionInUse 权限正在被使用
	ErrPermissionInUse = errors.New("权限正在被角色使用")
)
