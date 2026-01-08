package user

import "errors"

// 用户相关错误
var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists 用户已存在
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrUsernameAlreadyExists 用户名已存在
	ErrUsernameAlreadyExists = errors.New("username already exists")

	// ErrEmailAlreadyExists 邮箱已存在
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrRoleAlreadyAssigned 角色已分配
	ErrRoleAlreadyAssigned = errors.New("role already assigned to user")

	// ErrRoleNotFound 角色不存在
	ErrRoleNotFound = errors.New("role not found")

	// ErrInvalidUserStatus 无效的用户状态
	ErrInvalidUserStatus = errors.New("invalid user status")

	// ErrCannotDeleteSelf 不能删除自己
	ErrCannotDeleteSelf = errors.New("cannot delete yourself")

	// ErrCannotModifyAdmin 不能修改管理员
	ErrCannotModifyAdmin = errors.New("cannot modify admin user")

	// ErrInvalidPassword 密码错误
	ErrInvalidPassword = errors.New("invalid password")
)

// 系统用户保护相关错误
var (
	// ErrCannotDeleteSystemUser 不能删除系统用户
	ErrCannotDeleteSystemUser = errors.New("cannot delete system user")

	// ErrCannotModifySystemUsername 不能修改系统用户的用户名
	ErrCannotModifySystemUsername = errors.New("cannot modify system username")

	// ErrCannotModifyRootStatus 不能修改 root 用户状态
	ErrCannotModifyRootStatus = errors.New("cannot modify root user status")

	// ErrCannotModifyRootRoles 不能修改 root 用户的角色
	ErrCannotModifyRootRoles = errors.New("cannot modify root user roles")
)

// 用户类型相关错误
var (
	// ErrInvalidUserType 无效的用户类型
	ErrInvalidUserType = errors.New("invalid user type")

	// ErrCannotModifyUserType 用户类型不可修改
	ErrCannotModifyUserType = errors.New("user type cannot be modified")

	// ErrServiceAccountPasswordLogin 服务账户不能密码登录
	ErrServiceAccountPasswordLogin = errors.New("service account cannot login with password")
)

// 通用输入验证错误
var (
	// ErrInvalidUserID 无效的用户 ID
	ErrInvalidUserID = errors.New("invalid user ID")
)
