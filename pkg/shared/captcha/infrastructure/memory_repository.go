package infrastructure

import (
	"context"
	"sync"
	"time"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/shared/captcha"
)

// MemoryRepository 验证码内存存储实现。
// 线程安全，支持并发访问。
type MemoryRepository struct {
	mu     sync.RWMutex
	store  map[string]*captcha.CaptchaData
	ticker *time.Ticker
	done   chan struct{}
}

// NewMemoryRepository 创建新的内存验证码仓储。
func NewMemoryRepository() *MemoryRepository {
	repo := &MemoryRepository{
		store: make(map[string]*captcha.CaptchaData),
		done:  make(chan struct{}),
	}

	// 启动后台清理协程，每分钟清理过期验证码
	repo.ticker = time.NewTicker(1 * time.Minute)
	go repo.cleanupExpired()

	return repo
}

// ============================================================================
// CommandRepository 实现
// ============================================================================

// Create 创建验证码并存储
func (r *MemoryRepository) Create(ctx context.Context, captchaID string, code string, expiration time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store[captchaID] = &captcha.CaptchaData{
		Code:      code,
		ExpireAt:  time.Now().Add(expiration),
		CreatedAt: time.Now(),
	}

	return nil
}

// Verify 验证验证码（一次性使用，验证后需要从存储中移除）
func (r *MemoryRepository) Verify(ctx context.Context, captchaID string, code string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, exists := r.store[captchaID]
	if !exists {
		return false, captcha.ErrCaptchaNotFound
	}

	// 检查是否过期
	if data.IsExpired() {
		// 清理过期验证码
		delete(r.store, captchaID)
		return false, captcha.ErrCaptchaExpired
	}

	// 验证码是否正确
	if !data.Verify(code) {
		return false, captcha.ErrInvalidCaptcha
	}

	// 验证成功，删除验证码（一次性使用）
	delete(r.store, captchaID)
	return true, nil
}

// Delete 根据 ID 删除验证码
func (r *MemoryRepository) Delete(ctx context.Context, captchaID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.store, captchaID)
	return nil
}

// ============================================================================
// QueryRepository 实现
// ============================================================================

// GetStats 获取验证码存储统计信息
func (r *MemoryRepository) GetStats(ctx context.Context) map[string]any {
	r.mu.RLock()
	defer r.mu.RUnlock()

	total := len(r.store)
	expired := 0
	active := 0

	now := time.Now()
	for _, data := range r.store {
		if now.After(data.ExpireAt) {
			expired++
		} else {
			active++
		}
	}

	return map[string]any{
		"total":   total,
		"active":  active,
		"expired": expired,
	}
}

// ============================================================================
// 生命周期管理
// ============================================================================

// cleanupExpired 定期清理过期验证码
func (r *MemoryRepository) cleanupExpired() {
	for {
		select {
		case <-r.ticker.C:
			r.mu.Lock()
			now := time.Now()
			for id, data := range r.store {
				if now.After(data.ExpireAt) {
					delete(r.store, id)
				}
			}
			r.mu.Unlock()
		case <-r.done:
			return
		}
	}
}

// Close 关闭仓储，停止后台清理协程
func (r *MemoryRepository) Close() error {
	r.ticker.Stop()
	close(r.done)
	return nil
}
