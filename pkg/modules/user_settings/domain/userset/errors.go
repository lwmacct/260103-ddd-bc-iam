package userset

import "errors"

// 用户配置相关错误
var (
	// ErrUserSettingNotFound 用户配置不存在
	ErrUserSettingNotFound = errors.New("user setting not found")

	// ErrInvalidSettingKey 无效的配置键名（在 Settings BC 中不存在）
	ErrInvalidSettingKey = errors.New("invalid setting key")

	// ErrInvalidSettingValue 无效的配置值（不符合 Schema 验证规则）
	ErrInvalidSettingValue = errors.New("invalid setting value")

	// ErrValidationFailed 验证失败
	ErrValidationFailed = errors.New("validation failed")
)
