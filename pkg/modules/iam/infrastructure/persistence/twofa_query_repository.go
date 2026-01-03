package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/twofa"
	"gorm.io/gorm"
)

// twofaQueryRepository 2FA 查询仓储的 GORM 实现
type twofaQueryRepository struct {
	db *gorm.DB
}

// NewTwoFAQueryRepository 创建 2FA 查询仓储实例
func NewTwoFAQueryRepository(db *gorm.DB) twofa.QueryRepository {
	return &twofaQueryRepository{db: db}
}

// FindByUserID 根据用户 ID 查找 2FA 配置
func (r *twofaQueryRepository) FindByUserID(ctx context.Context, userID uint) (*twofa.TwoFA, error) {
	var model TwoFAModel
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // returns nil for not found, valid pattern
		}
		return nil, fmt.Errorf("failed to find 2FA by user ID: %w", result.Error)
	}

	return model.ToEntity(), nil
}

// IsEnabled 检查用户是否启用了 2FA
func (r *twofaQueryRepository) IsEnabled(ctx context.Context, userID uint) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&TwoFAModel{}).
		Where("user_id = ? AND enabled = ?", userID, true).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check 2FA status: %w", result.Error)
	}

	return count > 0, nil
}
