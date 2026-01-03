package captcha

import (
	"context"
	"time"
)

// ============================================================================
// Command Repository
// ============================================================================

// CommandRepository 定义所有会修改验证码状态的操作
type CommandRepository interface {
	// Create 创建验证码并存储
	Create(ctx context.Context, captchaID string, code string, expiration time.Duration) error

	// Verify 验证验证码（一次性使用，验证后需要从存储中移除）
	Verify(ctx context.Context, captchaID string, code string) (bool, error)

	// Delete 根据 ID 删除验证码
	Delete(ctx context.Context, captchaID string) error
}

// ============================================================================
// Query Repository
// ============================================================================

// QueryRepository 定义只读操作
type QueryRepository interface {
	// GetStats 获取验证码存储统计信息
	GetStats(ctx context.Context) map[string]any
}
