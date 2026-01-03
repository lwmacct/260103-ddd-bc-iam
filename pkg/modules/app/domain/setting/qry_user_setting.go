package setting

import "context"

// UserSettingQueryRepository 用户配置读操作接口。
type UserSettingQueryRepository interface {
	// FindByUserAndKey 根据用户 ID 和 Key 查找用户配置
	FindByUserAndKey(ctx context.Context, userID uint, key string) (*UserSetting, error)

	// FindByUser 查找用户的所有自定义配置
	FindByUser(ctx context.Context, userID uint) ([]*UserSetting, error)

	// FindByUserAndKeys 根据用户 ID 和多个 Key 批量查找用户配置
	FindByUserAndKeys(ctx context.Context, userID uint, keys []string) ([]*UserSetting, error)
}
