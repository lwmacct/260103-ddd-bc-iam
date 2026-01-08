package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/pat"
	"gorm.io/gorm"
)

// patCommandRepository PAT 命令仓储的 GORM 实现
type patCommandRepository struct {
	db *gorm.DB
}

// NewPATCommandRepository 创建 PAT 命令仓储实例
func NewPATCommandRepository(db *gorm.DB) pat.CommandRepository {
	return &patCommandRepository{db: db}
}

// Create 创建 PAT 令牌
func (r *patCommandRepository) Create(ctx context.Context, patToken *pat.PersonalAccessToken) error {
	model := newPATModelFromEntity(patToken)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create PAT: %w", err)
	}
	// 回写生成的 ID
	patToken.ID = model.ID
	return nil
}

// Update 更新 PAT 令牌
func (r *patCommandRepository) Update(ctx context.Context, patToken *pat.PersonalAccessToken) error {
	model := newPATModelFromEntity(patToken)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete 硬删除令牌（覆盖泛型的软删除行为）
func (r *patCommandRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).
		Unscoped(). // Hard delete
		Delete(&PersonalAccessTokenModel{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete PAT: %w", err)
	}
	return nil
}

// Disable 禁用令牌（设置状态为 disabled）
func (r *patCommandRepository) Disable(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).
		Model(&PersonalAccessTokenModel{}).
		Where("id = ?", id).
		Update("status", "disabled").Error; err != nil {
		return fmt.Errorf("failed to disable PAT: %w", err)
	}
	return nil
}

// Enable 重新启用令牌
func (r *patCommandRepository) Enable(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).
		Model(&PersonalAccessTokenModel{}).
		Where("id = ?", id).
		Update("status", "active").Error; err != nil {
		return fmt.Errorf("failed to enable PAT: %w", err)
	}
	return nil
}

// DeleteByUserID 删除指定用户的所有令牌
func (r *patCommandRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("user_id = ?", userID).
		Delete(&PersonalAccessTokenModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete PATs by user ID: %w", err)
	}
	return nil
}

// CleanupExpired 清理过期令牌
func (r *patCommandRepository) CleanupExpired(ctx context.Context) error {
	now := time.Now()

	if err := r.db.WithContext(ctx).
		Model(&PersonalAccessTokenModel{}).
		Where("expires_at IS NOT NULL AND expires_at < ? AND status != ?", now, "expired").
		Update("status", "expired").Error; err != nil {
		return fmt.Errorf("failed to cleanup expired PATs: %w", err)
	}
	return nil
}
