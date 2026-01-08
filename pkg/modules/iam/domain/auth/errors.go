package auth

import "errors"

// 认证相关错误
var (
	// ErrInvalidCredentials 无效的凭证
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")

	// ErrUserBanned 用户已被禁用
	ErrUserBanned = errors.New("user has been banned")

	// ErrUserInactive 用户未激活
	ErrUserInactive = errors.New("user is inactive")

	// ErrWeakPassword 密码强度不足
	ErrWeakPassword = errors.New("password is too weak")

	// ErrPasswordMismatch 密码不匹配
	ErrPasswordMismatch = errors.New("password mismatch")

	// ErrInvalidToken Token 无效
	ErrInvalidToken = errors.New("invalid token")

	// ErrTokenExpired Token 已过期
	ErrTokenExpired = errors.New("token has expired")

	// ErrInvalidCaptcha 验证码无效
	ErrInvalidCaptcha = errors.New("invalid captcha")

	// Err2FARequired 需要双因素认证
	Err2FARequired = errors.New("two-factor authentication required")

	// ErrInvalid2FACode 无效的 2FA 验证码
	ErrInvalid2FACode = errors.New("invalid two-factor authentication code")

	// ErrSessionNotFound Session 不存在
	ErrSessionNotFound = errors.New("session not found")

	// ErrSessionExpired Session 已过期
	ErrSessionExpired = errors.New("session has expired")

	// ErrUserNotAuthenticated 用户未认证
	ErrUserNotAuthenticated = errors.New("user not authenticated")

	// ErrUserIDNotFound 用户 ID 未找到
	ErrUserIDNotFound = errors.New("user ID not found")

	// ErrInvalidUserIDType 无效的用户 ID 类型
	ErrInvalidUserIDType = errors.New("invalid user ID type")

	// ErrInvalidUserContext 无效的用户上下文
	ErrInvalidUserContext = errors.New("invalid user context")
)
