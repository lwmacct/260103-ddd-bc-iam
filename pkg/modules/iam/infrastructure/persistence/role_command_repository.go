package persistence

import (
	"context"
	"errors"
	"fmt"

	corepersistence "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/infrastructure/persistence"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/role"
	"gorm.io/gorm"
)

// roleCommandRepository 角色命令仓储的 GORM 实现
// 嵌入 GenericCommandRepository 以复用基础 CRUD 操作
type roleCommandRepository struct {
	*corepersistence.GenericCommandRepository[role.Role, *RoleModel]
}

// NewRoleCommandRepository 创建角色命令仓储实例
func NewRoleCommandRepository(db *gorm.DB) role.CommandRepository {
	return &roleCommandRepository{
		GenericCommandRepository: corepersistence.NewGenericCommandRepository(
			db, newRoleModelFromEntity,
		),
	}
}

// Create、Update、Delete 方法由 GenericCommandRepository 提供

// SetPermissions 设置角色权限 (替换现有权限)
//
// 使用 JSONB 字段存储，直接更新 permissions 列。
func (r *roleCommandRepository) SetPermissions(ctx context.Context, roleID uint, permissions []role.Permission) error {
	// 验证角色存在
	var roleModel RoleModel
	if err := r.DB().WithContext(ctx).First(&roleModel, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("role not found with id: %d", roleID)
		}
		return fmt.Errorf("failed to find role: %w", err)
	}

	// 直接更新 JSONB 字段
	return r.DB().WithContext(ctx).
		Model(&RoleModel{}).
		Where("id = ?", roleID).
		Update("permissions", marshalPermissions(permissions)).
		Error
}
