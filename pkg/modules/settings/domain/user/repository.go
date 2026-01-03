package user

import "context"

// CommandRepository 用户配置命令仓储接口
// 负责所有修改状态的操作
type CommandRepository interface {
	// Upsert 创建或更新用户配置（基于 user_id + setting_key 唯一约束）
	Upsert(ctx context.Context, setting *UserSetting) error

	// Delete 删除指定用户的指定配置
	Delete(ctx context.Context, userID uint, key string) error

	// DeleteByUser 删除指定用户的所有配置
	DeleteByUser(ctx context.Context, userID uint) error

	// BatchUpsert 批量创建或更新用户配置
	BatchUpsert(ctx context.Context, settings []*UserSetting) error
}

// QueryRepository 用户配置查询仓储接口
// 负责所有只读查询操作
type QueryRepository interface {
	// FindByUserAndKey 根据用户 ID 和键名查找用户配置
	// 如果不存在返回 nil, nil
	FindByUserAndKey(ctx context.Context, userID uint, key string) (*UserSetting, error)

	// FindByUser 查找用户的所有自定义配置
	FindByUser(ctx context.Context, userID uint) ([]*UserSetting, error)

	// FindByUserAndKeys 批量查询用户的多个配置
	FindByUserAndKeys(ctx context.Context, userID uint, keys []string) ([]*UserSetting, error)
}
