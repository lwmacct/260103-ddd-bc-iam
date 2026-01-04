package team

import "errors"

// 团队配置相关错误
var (
	// ErrTeamSettingNotFound 团队配置不存在
	ErrTeamSettingNotFound = errors.New("team setting not found")

	// ErrInvalidSettingKey 无效的配置键名（在 Settings BC 中不存在）
	ErrInvalidSettingKey = errors.New("invalid setting key")

	// ErrInvalidSettingValue 无效的配置值（不符合 Schema 验证规则）
	ErrInvalidSettingValue = errors.New("invalid setting value")

	// ErrValidationFailed 验证失败
	ErrValidationFailed = errors.New("validation failed")
)

// 配置定义可见性和可配置性相关错误
var (
	// ErrSettingNotVisibleAtTeam 该设置在团队级别不可见
	ErrSettingNotVisibleAtTeam = errors.New("setting is not visible at team level")

	// ErrSettingNotConfigurableAtTeam 该设置不允许在团队级别配置
	ErrSettingNotConfigurableAtTeam = errors.New("setting is not configurable at team level")
)
