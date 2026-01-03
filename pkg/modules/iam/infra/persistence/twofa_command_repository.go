package persistence

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/twofa"
	"gorm.io/gorm"
)

// twofaCommandRepository 2FA 命令仓储的 GORM 实现
type twofaCommandRepository struct {
	db *gorm.DB
}

// NewTwoFACommandRepository 创建 2FA 命令仓储实例
func NewTwoFACommandRepository(db *gorm.DB) twofa.CommandRepository {
	return &twofaCommandRepository{db: db}
}

// CreateOrUpdate 创建或更新 2FA 配置
func (r *twofaCommandRepository) CreateOrUpdate(ctx context.Context, tfa *twofa.TwoFA) error {
	model := newTwoFAModelFromEntity(tfa)
	var persisted TwoFAModel
	result := r.db.WithContext(ctx).
		Where("user_id = ?", model.UserID).
		Assign(model).
		FirstOrCreate(&persisted)

	if result.Error != nil {
		return fmt.Errorf("failed to create or update 2FA: %w", result.Error)
	}

	if entity := persisted.ToEntity(); entity != nil {
		*tfa = *entity
	}

	return nil
}

// Delete 删除 2FA 配置
func (r *twofaCommandRepository) Delete(ctx context.Context, userID uint) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&TwoFAModel{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete 2FA: %w", result.Error)
	}

	return nil
}
