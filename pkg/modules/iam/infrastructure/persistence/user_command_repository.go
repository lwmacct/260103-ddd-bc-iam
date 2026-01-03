package persistence

import (
	"context"
	"errors"
	"fmt"

	corepersistence "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/infrastructure/persistence"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/user"
	"gorm.io/gorm"
)

// userCommandRepository 用户命令仓储的 GORM 实现
// 嵌入 GenericCommandRepository 以复用基础 CRUD 操作
type userCommandRepository struct {
	*corepersistence.GenericCommandRepository[user.User, *UserModel]
}

// NewUserCommandRepository 创建用户命令仓储实例
func NewUserCommandRepository(db *gorm.DB) user.CommandRepository {
	return &userCommandRepository{
		GenericCommandRepository: corepersistence.NewGenericCommandRepository(
			db, newUserModelFromEntity,
		),
	}
}

// Create、Update、Delete 方法由 GenericCommandRepository 提供

// AssignRoles 为用户分配角色（替换现有角色）
func (r *userCommandRepository) AssignRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	var u UserModel
	if err := r.DB().WithContext(ctx).First(&u, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.ErrUserNotFound
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// 如果 roleIDs 为空，清空用户的所有角色
	if len(roleIDs) == 0 {
		if err := r.DB().WithContext(ctx).Model(&u).Association("Roles").Clear(); err != nil {
			return fmt.Errorf("failed to clear roles: %w", err)
		}
		return nil
	}

	var roles []RoleModel
	if err := r.DB().WithContext(ctx).Find(&roles, roleIDs).Error; err != nil {
		return fmt.Errorf("failed to find roles: %w", err)
	}

	if err := r.DB().WithContext(ctx).Model(&u).Association("Roles").Replace(roles); err != nil {
		return fmt.Errorf("failed to assign roles: %w", err)
	}

	return nil
}

// RemoveRoles 移除用户的角色
func (r *userCommandRepository) RemoveRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	var u UserModel
	if err := r.DB().WithContext(ctx).First(&u, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.ErrUserNotFound
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	var roles []RoleModel
	if err := r.DB().WithContext(ctx).Find(&roles, roleIDs).Error; err != nil {
		return fmt.Errorf("failed to find roles: %w", err)
	}

	if err := r.DB().WithContext(ctx).Model(&u).Association("Roles").Delete(roles); err != nil {
		return fmt.Errorf("failed to remove roles: %w", err)
	}

	return nil
}

// UpdatePassword 更新用户密码
func (r *userCommandRepository) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	if err := r.DB().WithContext(ctx).Model(&UserModel{}).
		Where("id = ?", userID).
		Update("password", hashedPassword).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

// UpdateStatus 更新用户状态
func (r *userCommandRepository) UpdateStatus(ctx context.Context, userID uint, status string) error {
	if err := r.DB().WithContext(ctx).Model(&UserModel{}).
		Where("id = ?", userID).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	return nil
}
