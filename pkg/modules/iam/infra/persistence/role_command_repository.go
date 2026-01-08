package persistence

import (
	"context"
	"encoding/json"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/role"
	"gorm.io/gorm"
)

// SetPermissions 设置角色权限 (替换现有权限)
func (r *roleCommandRepository) SetPermissions(ctx context.Context, roleID uint, permissions []role.Permission) error {
	var roleModel RoleModel
	if err := r.db.WithContext(ctx).First(&roleModel, roleID).Error; err != nil {
		return err
	}

	// 直接更新 JSONB 字段
	permissionsJSON, err := json.Marshal(permissions)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).
		Model(&RoleModel{}).
		Where("id = ?", roleID).
		Update("permissions", permissionsJSON).
		Error
}

// roleCommandRepository 角色命令仓储的 GORM 实现
type roleCommandRepository struct {
	db *gorm.DB
}

// NewRoleCommandRepository 创建角色命令仓储实例
func NewRoleCommandRepository(db *gorm.DB) role.CommandRepository {
	return &roleCommandRepository{db: db}
}

// Create 创建角色
func (r *roleCommandRepository) Create(ctx context.Context, role *role.Role) error {
	model := newRoleModelFromEntity(role)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	role.ID = model.ID
	return nil
}

// Update 更新角色
func (r *roleCommandRepository) Update(ctx context.Context, role *role.Role) error {
	model := newRoleModelFromEntity(role)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete 删除角色
func (r *roleCommandRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&RoleModel{}, id).Error
}
