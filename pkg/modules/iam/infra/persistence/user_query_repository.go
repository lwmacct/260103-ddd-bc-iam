package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/role"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/user"
	"gorm.io/gorm"
)

// userQueryRepository 用户查询仓储的 GORM 实现
type userQueryRepository struct {
	db *gorm.DB
}

// NewUserQueryRepository 创建用户查询仓储实例
func NewUserQueryRepository(db *gorm.DB) user.QueryRepository {
	return &userQueryRepository{db: db}
}

// GetByID 根据 ID 获取用户
func (r *userQueryRepository) GetByID(ctx context.Context, id uint) (*user.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return model.ToEntity(), nil
}

// GetByUsername 根据用户名获取用户
func (r *userQueryRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return model.ToEntity(), nil
}

// GetByEmail 根据邮箱获取用户
func (r *userQueryRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return model.ToEntity(), nil
}

// GetByIDWithRoles 根据 ID 获取用户（包含角色和权限信息）
//
// 优化策略：单次 JOIN 查询
//
// 针对高延迟远程数据库优化，将多次查询合并为单次 JOIN，
// 最大限度减少网络往返次数（5 次 → 1 次）。
func (r *userQueryRepository) GetByIDWithRoles(ctx context.Context, id uint) (*user.User, error) {
	return r.getUserWithRolesByCondition(ctx, "u.id = ?", id)
}

// GetByUsernameWithRoles 根据用户名获取用户（包含角色和权限信息）
//
// 优化策略同 [GetByIDWithRoles]。
func (r *userQueryRepository) GetByUsernameWithRoles(ctx context.Context, username string) (*user.User, error) {
	return r.getUserWithRolesByCondition(ctx, "u.username = ?", username)
}

// GetByEmailWithRoles 根据邮箱获取用户（包含角色和权限信息）
//
// 优化策略同 [GetByIDWithRoles]。
func (r *userQueryRepository) GetByEmailWithRoles(ctx context.Context, email string) (*user.User, error) {
	return r.getUserWithRolesByCondition(ctx, "u.email = ?", email)
}

// List 获取用户列表 (分页)
func (r *userQueryRepository) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	var models []UserModel
	query := r.db.WithContext(ctx).Offset(offset).Limit(limit)

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return mapUserModelsToEntities(models), nil
}

// Count 统计用户数量
func (r *userQueryRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// GetRoles 获取用户的所有角色 ID
func (r *userQueryRepository) GetRoles(ctx context.Context, userID uint) ([]uint, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Preload("Roles").First(&model, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleIDs := make([]uint, 0, len(model.Roles))
	for _, r := range model.Roles {
		roleIDs = append(roleIDs, r.ID)
	}

	return roleIDs, nil
}

// Search 搜索用户（支持用户名、邮箱、真实姓名、昵称、手机号模糊匹配）
func (r *userQueryRepository) Search(ctx context.Context, keyword string, offset, limit int) ([]*user.User, error) {
	var models []UserModel
	query := r.db.WithContext(ctx).
		Where("username LIKE ? OR email LIKE ? OR real_name LIKE ? OR nickname LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Offset(offset).
		Limit(limit)

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return mapUserModelsToEntities(models), nil
}

// CountBySearch 统计搜索结果数量
func (r *userQueryRepository) CountBySearch(ctx context.Context, keyword string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("username LIKE ? OR email LIKE ? OR real_name LIKE ? OR nickname LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count search results: %w", err)
	}
	return count, nil
}

// Exists 检查用户是否存在
func (r *userQueryRepository) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}

// ExistsByUsername 检查用户名是否存在
func (r *userQueryRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}
	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *userQueryRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}

// GetUserIDsByRole 获取拥有指定角色的所有用户 ID
func (r *userQueryRepository) GetUserIDsByRole(ctx context.Context, roleID uint) ([]uint, error) {
	var userIDs []uint

	// 通过 user_roles 关联表查询
	if err := r.db.WithContext(ctx).
		Table("user_roles").
		Where("role_id = ?", roleID).
		Pluck("user_id", &userIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to get user IDs by role: %w", err)
	}

	return userIDs, nil
}

// =========================================================================
// 内部辅助方法
// =========================================================================

// userWithRolesRow 单次 JOIN 查询的结果行结构
//
// 由于是 LEFT JOIN，每个角色对应一行，用户字段会重复。
// 权限存储在 roles.permissions JSONB 字段中，无需额外 JOIN。
type userWithRolesRow struct {
	// User 字段
	UserID        uint       `gorm:"column:user_id"`
	UserCreatedAt time.Time  `gorm:"column:user_created_at"`
	UserUpdatedAt time.Time  `gorm:"column:user_updated_at"`
	UserDeletedAt *time.Time `gorm:"column:user_deleted_at"`
	Username      string     `gorm:"column:username"`
	Email         *string    `gorm:"column:email"` // nullable，对应数据库列
	Password      string     `gorm:"column:password"`
	RealName      string     `gorm:"column:real_name"`
	Nickname      string     `gorm:"column:nickname"`
	Phone         *string    `gorm:"column:phone"` // nullable，对应数据库列
	Signature     string     `gorm:"column:signature"`
	Avatar        string     `gorm:"column:avatar"`
	Bio           string     `gorm:"column:bio"`
	Status        string     `gorm:"column:status"`
	Type          string     `gorm:"column:type"`

	// Role 字段（可能为 NULL）
	RoleID          *uint      `gorm:"column:role_id"`
	RoleCreatedAt   *time.Time `gorm:"column:role_created_at"`
	RoleUpdatedAt   *time.Time `gorm:"column:role_updated_at"`
	RoleName        *string    `gorm:"column:role_name"`
	RoleDisplayName *string    `gorm:"column:role_display_name"`
	RoleDescription *string    `gorm:"column:role_description"`
	RoleIsSystem    *bool      `gorm:"column:role_is_system"`

	// Permission JSONB 字段（可能为 NULL）
	RolePermissions []byte `gorm:"column:role_permissions"`
}

// getUserWithRolesByCondition 通过单次 JOIN 查询获取用户及其角色权限
//
// 将 User、Role 的查询合并为 1 次 JOIN 查询（3 表），
// 权限从 roles.permissions JSONB 字段获取，无需 JOIN role_permissions 表。
func (r *userQueryRepository) getUserWithRolesByCondition(ctx context.Context, condition string, args ...any) (*user.User, error) {
	var rows []userWithRolesRow

	// 单次 JOIN 查询获取所有数据（3 表：users + user_roles + roles）
	query := r.db.WithContext(ctx).
		Table("users u").
		Select(`
			u.id as user_id, u.created_at as user_created_at, u.updated_at as user_updated_at,
			u.deleted_at as user_deleted_at, u.username, u.email, u.password,
			u.real_name, u.nickname, u.phone, u.signature, u.avatar, u.bio, u.status, u.type,
			r.id as role_id, r.created_at as role_created_at, r.updated_at as role_updated_at,
			r.name as role_name, r.display_name as role_display_name,
			r.description as role_description, r.is_system as role_is_system,
			r.permissions as role_permissions
		`).
		Joins("LEFT JOIN user_roles ur ON u.id = ur.user_id").
		Joins("LEFT JOIN roles r ON ur.role_id = r.id AND r.deleted_at IS NULL").
		Where("u.deleted_at IS NULL").
		Where(condition, args...)

	if err := query.Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("failed to get user with roles: %w", err)
	}

	if len(rows) == 0 {
		return nil, user.ErrUserNotFound
	}

	// 从扁平结果重建嵌套结构
	return r.buildUserFromRows(rows), nil
}

// buildUserFromRows 从 JOIN 查询结果构建 User 实体
//
// 将扁平化的行数据重建为 User → []Role → []Permission 的嵌套结构。
// 权限从 JSONB 字段解析。
func (r *userQueryRepository) buildUserFromRows(rows []userWithRolesRow) *user.User {
	if len(rows) == 0 {
		return nil
	}

	// 第一行包含用户信息
	first := rows[0]
	result := &user.User{
		ID:        first.UserID,
		CreatedAt: first.UserCreatedAt,
		UpdatedAt: first.UserUpdatedAt,
		DeletedAt: first.UserDeletedAt,
		Username:  first.Username,
		Email:     first.Email,
		Password:  first.Password,
		RealName:  first.RealName,
		Nickname:  first.Nickname,
		Phone:     first.Phone,
		Signature: first.Signature,
		Avatar:    first.Avatar,
		Bio:       first.Bio,
		Status:    first.Status,
		Type:      user.UserType(first.Type),
	}

	// 收集角色（去重）
	roleMap := make(map[uint]*role.Role)

	for _, row := range rows {
		// 跳过无角色的行
		if row.RoleID == nil {
			continue
		}

		roleID := *row.RoleID

		// 确保角色存在（每个角色只处理一次）
		if _, exists := roleMap[roleID]; !exists {
			roleMap[roleID] = &role.Role{
				ID:          roleID,
				CreatedAt:   derefTime(row.RoleCreatedAt),
				UpdatedAt:   derefTime(row.RoleUpdatedAt),
				Name:        derefString(row.RoleName),
				DisplayName: derefString(row.RoleDisplayName),
				Description: derefString(row.RoleDescription),
				IsSystem:    derefBool(row.RoleIsSystem),
				Permissions: unmarshalPermissions(row.RolePermissions),
			}
		}
	}

	// 转换为切片
	result.Roles = make([]role.Role, 0, len(roleMap))
	for _, r := range roleMap {
		result.Roles = append(result.Roles, *r)
	}

	return result
}

// 辅助函数：解引用指针
func derefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func derefBool(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

func derefTime(p *time.Time) time.Time {
	if p == nil {
		return time.Time{}
	}
	return *p
}
