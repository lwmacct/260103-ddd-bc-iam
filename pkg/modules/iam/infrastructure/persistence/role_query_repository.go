package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/role"
	"gorm.io/gorm"
)

// roleQueryRepository 角色查询仓储的 GORM 实现
type roleQueryRepository struct {
	db *gorm.DB
}

// NewRoleQueryRepository 创建角色查询仓储实例
func NewRoleQueryRepository(db *gorm.DB) role.QueryRepository {
	return &roleQueryRepository{db: db}
}

// FindByID 根据 ID 查找角色（含权限，JSONB 字段自动加载）
func (r *roleQueryRepository) FindByID(ctx context.Context, id uint) (*role.Role, error) {
	var model RoleModel
	err := r.db.WithContext(ctx).
		First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // returns nil for not found, valid pattern
		}
		return nil, fmt.Errorf("failed to find role by id: %w", err)
	}
	return model.ToEntity(), nil
}

// FindByName 根据名称查找角色（含权限，JSONB 字段自动加载）
func (r *roleQueryRepository) FindByName(ctx context.Context, name string) (*role.Role, error) {
	var model RoleModel
	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // returns nil for not found, valid pattern
		}
		return nil, fmt.Errorf("failed to find role by name: %w", err)
	}
	return model.ToEntity(), nil
}

// FindByIDWithPermissions 根据 ID 查找角色（含权限）
//
// 由于权限已存储在 JSONB 字段中，此方法与 FindByID 等效。
// 保留此方法是为了保持接口兼容性。
func (r *roleQueryRepository) FindByIDWithPermissions(ctx context.Context, id uint) (*role.Role, error) {
	return r.FindByID(ctx, id)
}

// List 获取角色列表（分页，含权限）
//
// 权限存储在 JSONB 字段中，单次查询即可获取完整数据。
func (r *roleQueryRepository) List(ctx context.Context, page, limit int) ([]role.Role, int64, error) {
	var models []RoleModel
	var total int64

	query := r.db.WithContext(ctx).Model(&RoleModel{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count roles: %w", err)
	}

	offset := (page - 1) * limit
	err := query.
		Offset(offset).
		Limit(limit).
		Find(&models).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list roles: %w", err)
	}

	// 直接转换，权限已包含在 JSONB 字段中
	return mapRoleModelsToEntities(models), total, nil
}

// GetPermissions 获取角色的所有权限
func (r *roleQueryRepository) GetPermissions(ctx context.Context, roleID uint) ([]role.Permission, error) {
	var model RoleModel
	if err := r.db.WithContext(ctx).
		Select("id", "permissions").
		First(&model, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("role not found with id: %d", roleID)
		}
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	return unmarshalPermissions(model.Permissions), nil
}

// Exists 检查角色是否存在
func (r *roleQueryRepository) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&RoleModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check role existence: %w", err)
	}
	return count > 0, nil
}

// ExistsByName 检查角色名称是否存在
func (r *roleQueryRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&RoleModel{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check role name existence: %w", err)
	}
	return count > 0, nil
}
