package usersetting

import "errors"

// 领域错误
var (
	// ErrUserSettingNotFound 用户设置不存在
	ErrUserSettingNotFound = errors.New("user setting not found")

	// ErrInvalidSettingKey 无效的设置键名（在 Schema 中不存在）
	ErrInvalidSettingKey = errors.New("invalid setting key")

	// ErrInvalidSettingValue 无效的设置值（不符合 Schema 验证规则）
	ErrInvalidSettingValue = errors.New("invalid setting value")
)
