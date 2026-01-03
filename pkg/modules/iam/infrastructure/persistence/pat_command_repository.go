package persistence

import (
	"context"
	"fmt"
	"time"

	corepersistence "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/infrastructure/persistence"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/pat"
	"gorm.io/gorm"
)

// patCommandRepository PAT 命令仓储的 GORM 实现
// 嵌入 GenericCommandRepository 以复用 Create/Update 操作
type patCommandRepository struct {
	*corepersistence.GenericCommandRepository[pat.PersonalAccessToken, *PersonalAccessTokenModel]
}

// NewPATCommandRepository 创建 PAT 命令仓储实例
func NewPATCommandRepository(db *gorm.DB) pat.CommandRepository {
	return &patCommandRepository{
		GenericCommandRepository: corepersistence.NewGenericCommandRepository(
			db, newPATModelFromEntity,
		),
	}
}

// Create、Update 方法由 GenericCommandRepository 提供

// Delete 硬删除令牌（覆盖泛型的软删除行为）
func (r *patCommandRepository) Delete(ctx context.Context, id uint) error {
	if err := r.DB().WithContext(ctx).
		Unscoped(). // Hard delete
		Delete(&PersonalAccessTokenModel{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete PAT: %w", err)
	}
	return nil
}

// Disable 禁用令牌（设置状态为 disabled）
func (r *patCommandRepository) Disable(ctx context.Context, id uint) error {
	if err := r.DB().WithContext(ctx).
		Model(&PersonalAccessTokenModel{}).
		Where("id = ?", id).
		Update("status", "disabled").Error; err != nil {
		return fmt.Errorf("failed to disable PAT: %w", err)
	}
	return nil
}

// Enable 重新启用令牌
func (r *patCommandRepository) Enable(ctx context.Context, id uint) error {
	if err := r.DB().WithContext(ctx).
		Model(&PersonalAccessTokenModel{}).
		Where("id = ?", id).
		Update("status", "active").Error; err != nil {
		return fmt.Errorf("failed to enable PAT: %w", err)
	}
	return nil
}

// DeleteByUserID 删除指定用户的所有令牌
func (r *patCommandRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	if err := r.DB().WithContext(ctx).
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

	if err := r.DB().WithContext(ctx).
		Model(&PersonalAccessTokenModel{}).
		Where("expires_at IS NOT NULL AND expires_at < ? AND status != ?", now, "expired").
		Update("status", "expired").Error; err != nil {
		return fmt.Errorf("failed to cleanup expired PATs: %w", err)
	}
	return nil
}
