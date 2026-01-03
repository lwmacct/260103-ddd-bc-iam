package pat

import "context"

// ============================================================================
// Command Repository
// ============================================================================

// CommandRepository 定义 PAT 写操作接口
type CommandRepository interface {
	// Create 创建新的个人访问令牌
	Create(ctx context.Context, pat *PersonalAccessToken) error

	// Update 更新令牌（主要用于 LastUsedAt）
	Update(ctx context.Context, pat *PersonalAccessToken) error

	// Delete 硬删除令牌
	Delete(ctx context.Context, id uint) error

	// Disable 禁用令牌（设置状态为 disabled）
	Disable(ctx context.Context, id uint) error

	// Enable 启用令牌（设置状态为 active）
	Enable(ctx context.Context, id uint) error

	// DeleteByUserID 删除指定用户的所有令牌
	DeleteByUserID(ctx context.Context, userID uint) error

	// CleanupExpired 清理过期令牌
	CleanupExpired(ctx context.Context) error
}

// ============================================================================
// Query Repository
// ============================================================================

// QueryRepository 定义 PAT 读操作接口
type QueryRepository interface {
	// FindByToken 通过令牌哈希查找（用于认证）
	FindByToken(ctx context.Context, tokenHash string) (*PersonalAccessToken, error)

	// FindByID 通过 ID 查找令牌
	FindByID(ctx context.Context, id uint) (*PersonalAccessToken, error)

	// FindByPrefix 通过前缀查找令牌
	FindByPrefix(ctx context.Context, prefix string) (*PersonalAccessToken, error)

	// ListByUser 获取指定用户的所有令牌
	ListByUser(ctx context.Context, userID uint) ([]*PersonalAccessToken, error)
}
