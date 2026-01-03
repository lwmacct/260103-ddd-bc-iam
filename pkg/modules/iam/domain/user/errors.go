package user

import "errors"

// 用户相关错误
var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("用户不存在")

	// ErrUserAlreadyExists 用户已存在
	ErrUserAlreadyExists = errors.New("用户已存在")

	// ErrUsernameAlreadyExists 用户名已存在
	ErrUsernameAlreadyExists = errors.New("用户名已存在")

	// ErrEmailAlreadyExists 邮箱已存在
	ErrEmailAlreadyExists = errors.New("邮箱已存在")

	// ErrRoleAlreadyAssigned 角色已分配
	ErrRoleAlreadyAssigned = errors.New("角色已分配给用户")

	// ErrRoleNotFound 角色不存在
	ErrRoleNotFound = errors.New("角色不存在")

	// ErrInvalidUserStatus 无效的用户状态
	ErrInvalidUserStatus = errors.New("无效的用户状态")

	// ErrCannotDeleteSelf 不能删除自己
	ErrCannotDeleteSelf = errors.New("不能删除自己")

	// ErrCannotModifyAdmin 不能修改管理员
	ErrCannotModifyAdmin = errors.New("不能修改管理员用户")

	// ErrInvalidPassword 密码错误
	ErrInvalidPassword = errors.New("密码错误")
)

// 系统用户保护相关错误
var (
	// ErrCannotDeleteSystemUser 不能删除系统用户
	ErrCannotDeleteSystemUser = errors.New("不能删除系统用户")

	// ErrCannotModifySystemUsername 不能修改系统用户的用户名
	ErrCannotModifySystemUsername = errors.New("不能修改系统用户的用户名")

	// ErrCannotModifyRootStatus 不能修改 root 用户状态
	ErrCannotModifyRootStatus = errors.New("不能修改 root 用户状态")

	// ErrCannotModifyRootRoles 不能修改 root 用户的角色
	ErrCannotModifyRootRoles = errors.New("不能修改 root 用户角色")
)

// 用户类型相关错误
var (
	// ErrInvalidUserType 无效的用户类型
	ErrInvalidUserType = errors.New("无效的用户类型")

	// ErrCannotModifyUserType 用户类型不可修改
	ErrCannotModifyUserType = errors.New("用户类型不可修改")

	// ErrServiceAccountPasswordLogin 服务账户不能密码登录
	ErrServiceAccountPasswordLogin = errors.New("服务账户不能使用密码登录")
)

// 通用输入验证错误
var (
	// ErrInvalidUserID 无效的用户 ID
	ErrInvalidUserID = errors.New("无效的用户 ID")
)
