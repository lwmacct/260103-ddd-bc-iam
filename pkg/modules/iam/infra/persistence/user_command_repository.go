package persistence

import (
	"context"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// AssignRoles 为用户分配角色（使用批量插入，避免 N+1 查询）
func (r *userCommandRepository) AssignRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	// 如果 roleIDs 为空，清空用户的所有角色
	if len(roleIDs) == 0 {
		// 批量删除：直接操作关联表
		if err := r.db.WithContext(ctx).
			Table("user_roles").
			Where("user_id = ?", userID).
			Delete(&struct{}{}).Error; err != nil {
			return err
		}
		return nil
	}

	// 先删除现有角色关联（使用批量删除）
	if err := r.db.WithContext(ctx).
		Table("user_roles").
		Where("user_id = ?", userID).
		Delete(&struct{}{}).Error; err != nil {
		return err
	}

	// 构建批量插入数据
	records := make([]map[string]any, len(roleIDs))
	for i, roleID := range roleIDs {
		records[i] = map[string]any{
			"user_id": userID,
			"role_id": roleID,
		}
	}

	// 批量插入关联表（使用 OnConflict 避免重复）
	if err := r.db.WithContext(ctx).
		Table("user_roles").
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&records).Error; err != nil {
		return err
	}

	return nil
}

// RemoveRoles 移除用户的角色（使用批量删除，避免 N+1 查询）
func (r *userCommandRepository) RemoveRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	if len(roleIDs) == 0 {
		return nil
	}

	// 批量删除：直接操作关联表
	if err := r.db.WithContext(ctx).
		Table("user_roles").
		Where("user_id = ? AND role_id IN ?", userID, roleIDs).
		Delete(&struct{}{}).Error; err != nil {
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
