package usersetting

import "context"

// ============================================================================
// Command Repository
// ============================================================================

// CommandRepository 用户设置命令仓储接口
// 负责所有修改状态的操作（Create, Update, Delete, Upsert）
type CommandRepository interface {
	// Create 创建用户设置
	Create(ctx context.Context, setting *UserSetting) error

	// Update 更新用户设置
	Update(ctx context.Context, setting *UserSetting) error

	// Delete 删除用户设置 (软删除)
	Delete(ctx context.Context, id uint) error

	// Upsert 创建或更新用户设置（基于 user_id + key 唯一约束）
	Upsert(ctx context.Context, userID uint, key string, value string) error
}

// ============================================================================
// Query Repository
// ============================================================================

// QueryRepository 用户设置查询仓储接口
// 负责所有只读查询操作
type QueryRepository interface {
	// FindByUserAndKey 根据用户 ID 和键名查找用户设置
	FindByUserAndKey(ctx context.Context, userID uint, key string) (*UserSetting, error)

	// FindByUser 查找用户的所有自定义设置
	FindByUser(ctx context.Context, userID uint) ([]*UserSetting, error)

	// FindByUserAndCategory 根据用户 ID 和分类查找用户设置
	FindByUserAndCategory(ctx context.Context, userID uint, category string) ([]*UserSetting, error)

	// ExistsByUserAndKey 检查用户是否对指定键有自定义值
	ExistsByUserAndKey(ctx context.Context, userID uint, key string) (bool, error)
}
