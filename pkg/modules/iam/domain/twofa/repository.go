package twofa

import "context"

// ============================================================================
// Command Repository
// ============================================================================

// CommandRepository 定义 2FA 写操作接口
type CommandRepository interface {
	// CreateOrUpdate 创建或更新 2FA 配置
	CreateOrUpdate(ctx context.Context, twoFA *TwoFA) error

	// Delete 删除 2FA 配置
	Delete(ctx context.Context, userID uint) error
}

// ============================================================================
// Query Repository
// ============================================================================

// QueryRepository 定义 2FA 读操作接口
type QueryRepository interface {
	// FindByUserID 根据用户ID查找 2FA 配置
	FindByUserID(ctx context.Context, userID uint) (*TwoFA, error)

	// IsEnabled 检查用户是否启用了 2FA
	IsEnabled(ctx context.Context, userID uint) (bool, error)
}
