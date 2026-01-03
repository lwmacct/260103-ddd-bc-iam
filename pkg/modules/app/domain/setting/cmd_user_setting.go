package setting

import "context"

// UserSettingCommandRepository 用户配置写操作接口。
type UserSettingCommandRepository interface {
	// Upsert 插入或更新用户配置
	// 如果 (UserID, SettingKey) 已存在则更新，否则创建
	Upsert(ctx context.Context, setting *UserSetting) error

	// Delete 删除用户配置（恢复为默认值）
	Delete(ctx context.Context, userID uint, key string) error

	// DeleteByUser 删除用户的所有配置
	DeleteByUser(ctx context.Context, userID uint) error

	// BatchUpsert 批量插入或更新用户配置
	BatchUpsert(ctx context.Context, settings []*UserSetting) error
}
