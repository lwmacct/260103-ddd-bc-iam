package persistence

import (
	"context"
	"errors"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/user"
	"gorm.io/gorm"
)

// userCommandRepository 用户命令仓储的 GORM 实现
type userCommandRepository struct {
	db *gorm.DB
}

// NewUserCommandRepository 创建用户命令仓储实例
func NewUserCommandRepository(db *gorm.DB) user.CommandRepository {
	return &userCommandRepository{db: db}
}

// Create 创建用户
func (r *userCommandRepository) Create(ctx context.Context, user *user.User) error {
	model := newUserModelFromEntity(user)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	user.ID = model.ID
	return nil
}

// Update 更新用户
func (r *userCommandRepository) Update(ctx context.Context, user *user.User) error {
	model := newUserModelFromEntity(user)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete 删除用户
func (r *userCommandRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&UserModel{}, id).Error
}

// AssignRoles 为用户分配角色（替换现有角色）
func (r *userCommandRepository) AssignRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	var u UserModel
	if err := r.db.WithContext(ctx).First(&u, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.ErrUserNotFound
		}
		return err
	}

	// 如果 roleIDs 为空，清空用户的所有角色
	if len(roleIDs) == 0 {
		if err := r.db.WithContext(ctx).Model(&u).Association("Roles").Clear(); err != nil {
			return err
		}
		return nil
	}

	var roles []RoleModel
	if err := r.db.WithContext(ctx).Find(&roles, roleIDs).Error; err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Model(&u).Association("Roles").Replace(roles); err != nil {
		return err
	}

	return nil
}

// RemoveRoles 移除用户的角色
func (r *userCommandRepository) RemoveRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	var u UserModel
	if err := r.db.WithContext(ctx).First(&u, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.ErrUserNotFound
		}
		return err
	}

	var roles []RoleModel
	if err := r.db.WithContext(ctx).Find(&roles, roleIDs).Error; err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Model(&u).Association("Roles").Delete(roles); err != nil {
		return err
	}

	return nil
}

// UpdatePassword 更新用户密码
func (r *userCommandRepository) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	if err := r.db.WithContext(ctx).Model(&UserModel{}).
		Where("id = ?", userID).
		Update("password", hashedPassword).Error; err != nil {
		return err
	}
	return nil
}

// UpdateStatus 更新用户状态
func (r *userCommandRepository) UpdateStatus(ctx context.Context, userID uint, status string) error {
	if err := r.db.WithContext(ctx).Model(&UserModel{}).
		Where("id = ?", userID).
		Update("status", status).Error; err != nil {
		return err
	}
	return nil
}
