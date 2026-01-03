package user

import (
	"context"
)

// ============================================================================
// Command Repository
// ============================================================================

// CommandRepository 用户命令仓储接口
// 负责所有修改状态的操作（Create, Update, Delete）
type CommandRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *User) error

	// Update 更新用户
	Update(ctx context.Context, user *User) error

	// Delete 删除用户 (软删除)
	Delete(ctx context.Context, id uint) error

	// AssignRoles 为用户分配角色
	AssignRoles(ctx context.Context, userID uint, roleIDs []uint) error

	// RemoveRoles 移除用户的角色
	RemoveRoles(ctx context.Context, userID uint, roleIDs []uint) error

	// UpdatePassword 更新用户密码
	UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error

	// UpdateStatus 更新用户状态
	UpdateStatus(ctx context.Context, userID uint, status string) error
}

// ============================================================================
// Query Repository - 细粒度接口定义（遵循接口隔离原则 ISP）
// ============================================================================

// BaseQueryRepository 基础查询接口
// 提供最基本的 ID 查询能力
type BaseQueryRepository interface {
	// GetByID 根据 ID 获取用户
	GetByID(ctx context.Context, id uint) (*User, error)
}

// AuthQueryRepository 认证查询接口
// 用于登录、认证等场景的用户查找
type AuthQueryRepository interface {
	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*User, error)

	// GetByEmail 根据邮箱获取用户
	GetByEmail(ctx context.Context, email string) (*User, error)

	// GetByUsernameWithRoles 根据用户名获取用户（包含角色和权限信息）
	GetByUsernameWithRoles(ctx context.Context, username string) (*User, error)

	// GetByEmailWithRoles 根据邮箱获取用户（包含角色和权限信息）
	GetByEmailWithRoles(ctx context.Context, email string) (*User, error)
}

// DetailQueryRepository 详情查询接口
// 用于需要关联数据的场景
type DetailQueryRepository interface {
	// GetByIDWithRoles 根据 ID 获取用户（包含角色和权限信息）
	GetByIDWithRoles(ctx context.Context, id uint) (*User, error)
}

// ValidationQueryRepository 验证查询接口
// 用于存在性检查、唯一性验证等场景
type ValidationQueryRepository interface {
	// Exists 检查用户是否存在
	Exists(ctx context.Context, id uint) (bool, error)

	// ExistsByUsername 检查用户名是否存在
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// ListQueryRepository 列表查询接口
// 用于分页列表、搜索等场景
type ListQueryRepository interface {
	// List 获取用户列表 (分页)
	List(ctx context.Context, offset, limit int) ([]*User, error)

	// Count 统计用户数量
	Count(ctx context.Context) (int64, error)

	// Search 搜索用户（支持用户名、邮箱、全名模糊匹配）
	Search(ctx context.Context, keyword string, offset, limit int) ([]*User, error)

	// CountBySearch 统计搜索结果数量
	CountBySearch(ctx context.Context, keyword string) (int64, error)
}

// RoleQueryRepository 角色关联查询接口
// 用于角色相关的用户查询
type RoleQueryRepository interface {
	// GetRoles 获取用户的所有角色 ID
	GetRoles(ctx context.Context, userID uint) ([]uint, error)

	// GetUserIDsByRole 获取拥有指定角色的所有用户 ID
	// 用于权限缓存失效场景（角色权限变更时需要清除所有相关用户的缓存）
	GetUserIDsByRole(ctx context.Context, roleID uint) ([]uint, error)
}

// ============================================================================
// Query Repository - 聚合接口
// ============================================================================

// QueryRepository 用户查询仓储接口（聚合接口）。
// 组合所有细粒度接口，实现类只需实现此接口即可满足所有子接口的类型约束。
type QueryRepository interface {
	BaseQueryRepository
	AuthQueryRepository
	DetailQueryRepository
	ValidationQueryRepository
	ListQueryRepository
	RoleQueryRepository
}
