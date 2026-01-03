package auth

import "errors"

// 认证相关错误
var (
	// ErrInvalidCredentials 无效的凭证
	ErrInvalidCredentials = errors.New("无效的凭证")

	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("用户不存在")

	// ErrUserBanned 用户已被禁用
	ErrUserBanned = errors.New("用户已被禁用")

	// ErrUserInactive 用户未激活
	ErrUserInactive = errors.New("用户未激活")

	// ErrWeakPassword 密码强度不足
	ErrWeakPassword = errors.New("密码强度不足")

	// ErrPasswordMismatch 密码不匹配
	ErrPasswordMismatch = errors.New("密码不匹配")

	// ErrInvalidToken Token 无效
	ErrInvalidToken = errors.New("token 无效")

	// ErrTokenExpired Token 已过期
	ErrTokenExpired = errors.New("token 已过期")

	// ErrInvalidCaptcha 验证码无效
	ErrInvalidCaptcha = errors.New("验证码无效")

	// Err2FARequired 需要双因素认证
	Err2FARequired = errors.New("需要双因素认证")

	// ErrInvalid2FACode 无效的 2FA 验证码
	ErrInvalid2FACode = errors.New("无效的双因素认证验证码")

	// ErrSessionNotFound Session 不存在
	ErrSessionNotFound = errors.New("会话不存在")

	// ErrSessionExpired Session 已过期
	ErrSessionExpired = errors.New("会话已过期")

	// ErrUserNotAuthenticated 用户未认证
	ErrUserNotAuthenticated = errors.New("用户未认证")

	// ErrUserIDNotFound 用户 ID 未找到
	ErrUserIDNotFound = errors.New("用户 ID 未找到")

	// ErrInvalidUserIDType 无效的用户 ID 类型
	ErrInvalidUserIDType = errors.New("无效的用户 ID 类型")

	// ErrInvalidUserContext 无效的用户上下文
	ErrInvalidUserContext = errors.New("无效的用户上下文")
)
