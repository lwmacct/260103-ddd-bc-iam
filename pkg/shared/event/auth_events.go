package event

import (
	"strconv"
)

// ============================================================================
// 认证事件
// ============================================================================

// LoginSucceededEvent 登录成功事件
type LoginSucceededEvent struct {
	BaseEvent

	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
}

// NewLoginSucceededEvent 创建登录成功事件
func NewLoginSucceededEvent(userID uint, username, ipAddress, userAgent string) *LoginSucceededEvent {
	return &LoginSucceededEvent{
		BaseEvent: NewBaseEvent("auth.login_succeeded", "auth", strconv.FormatUint(uint64(userID), 10)),
		UserID:    userID,
		Username:  username,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
}

// LoginFailedEvent 登录失败事件
type LoginFailedEvent struct {
	BaseEvent

	Username  string `json:"username"`
	IPAddress string `json:"ip_address"`
	Reason    string `json:"reason"`
}

// NewLoginFailedEvent 创建登录失败事件
func NewLoginFailedEvent(username, ipAddress, reason string) *LoginFailedEvent {
	return &LoginFailedEvent{
		BaseEvent: NewBaseEvent("auth.login_failed", "auth", username),
		Username:  username,
		IPAddress: ipAddress,
		Reason:    reason,
	}
}

// LogoutEvent 登出事件
type LogoutEvent struct {
	BaseEvent

	UserID uint `json:"user_id"`
}

// NewLogoutEvent 创建登出事件
func NewLogoutEvent(userID uint) *LogoutEvent {
	return &LogoutEvent{
		BaseEvent: NewBaseEvent("auth.logout", "auth", strconv.FormatUint(uint64(userID), 10)),
		UserID:    userID,
	}
}

// TokenRefreshedEvent 令牌刷新事件
type TokenRefreshedEvent struct {
	BaseEvent

	UserID uint `json:"user_id"`
}

// NewTokenRefreshedEvent 创建令牌刷新事件
func NewTokenRefreshedEvent(userID uint) *TokenRefreshedEvent {
	return &TokenRefreshedEvent{
		BaseEvent: NewBaseEvent("auth.token_refreshed", "auth", strconv.FormatUint(uint64(userID), 10)),
		UserID:    userID,
	}
}

// TwoFAEnabledEvent 2FA 启用事件
type TwoFAEnabledEvent struct {
	BaseEvent

	UserID uint `json:"user_id"`
}

// NewTwoFAEnabledEvent 创建 2FA 启用事件
func NewTwoFAEnabledEvent(userID uint) *TwoFAEnabledEvent {
	return &TwoFAEnabledEvent{
		BaseEvent: NewBaseEvent("auth.twofa_enabled", "auth", strconv.FormatUint(uint64(userID), 10)),
		UserID:    userID,
	}
}

// TwoFADisabledEvent 2FA 禁用事件
type TwoFADisabledEvent struct {
	BaseEvent

	UserID uint `json:"user_id"`
}

// NewTwoFADisabledEvent 创建 2FA 禁用事件
func NewTwoFADisabledEvent(userID uint) *TwoFADisabledEvent {
	return &TwoFADisabledEvent{
		BaseEvent: NewBaseEvent("auth.twofa_disabled", "auth", strconv.FormatUint(uint64(userID), 10)),
		UserID:    userID,
	}
}
