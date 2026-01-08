package persistence

import (
	"context"
	"log/slog"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/user"
)

// cachedUserQueryRepository 带缓存的用户查询仓储（装饰器模式）。
//
// 对 GetByIDWithRoles 等高频查询方法添加缓存层，其他方法直接委托。
// 采用 Cache-Aside 模式：
//  1. 先查缓存
//  2. 缓存未命中则查数据库
//  3. 同步回写缓存
type cachedUserQueryRepository struct {
	delegate user.QueryRepository
	cache    user.UserWithRolesCacheService
}

// NewCachedUserQueryRepository 创建带缓存的用户查询仓储。
func NewCachedUserQueryRepository(
	delegate user.QueryRepository,
	cache user.UserWithRolesCacheService,
) user.QueryRepository {
	return &cachedUserQueryRepository{
		delegate: delegate,
		cache:    cache,
	}
}

// ============================================================================
// 带缓存的方法
// ============================================================================

// GetByIDWithRoles 根据 ID 获取用户（带缓存）。
func (r *cachedUserQueryRepository) GetByIDWithRoles(ctx context.Context, id uint) (*user.User, error) {
	// 1. 查缓存
	cached, err := r.cache.GetUserWithRoles(ctx, id)
	if err != nil {
		slog.WarnContext(ctx, "user cache get error, falling back to db", "user_id", id, "error", err.Error())
	}
	if cached != nil {
		return cached, nil // 缓存命中
	}

	// 2. 查数据库
	result, err := r.delegate.GetByIDWithRoles(ctx, id)
	if err != nil {
		return nil, err
	}

	// 3. 同步回写缓存
	if err := r.cache.SetUserWithRoles(ctx, result); err != nil {
		slog.WarnContext(ctx, "user cache set failed", "user_id", id, "error", err.Error())
	}

	return result, nil
}

// ============================================================================
// 直接委托的方法（不加缓存）
// ============================================================================

// GetByID 根据 ID 获取用户。
func (r *cachedUserQueryRepository) GetByID(ctx context.Context, id uint) (*user.User, error) {
	return r.delegate.GetByID(ctx, id)
}

// GetByUsername 根据用户名获取用户。
func (r *cachedUserQueryRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	return r.delegate.GetByUsername(ctx, username)
}

// GetByEmail 根据邮箱获取用户。
func (r *cachedUserQueryRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	return r.delegate.GetByEmail(ctx, email)
}

// GetByUsernameWithRoles 根据用户名获取用户（包含角色和权限信息）。
func (r *cachedUserQueryRepository) GetByUsernameWithRoles(ctx context.Context, username string) (*user.User, error) {
	return r.delegate.GetByUsernameWithRoles(ctx, username)
}

// GetByEmailWithRoles 根据邮箱获取用户（包含角色和权限信息）。
func (r *cachedUserQueryRepository) GetByEmailWithRoles(ctx context.Context, email string) (*user.User, error) {
	return r.delegate.GetByEmailWithRoles(ctx, email)
}

// Exists 检查用户是否存在。
func (r *cachedUserQueryRepository) Exists(ctx context.Context, id uint) (bool, error) {
	return r.delegate.Exists(ctx, id)
}

// ExistsByUsername 检查用户名是否存在。
func (r *cachedUserQueryRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return r.delegate.ExistsByUsername(ctx, username)
}

// ExistsByEmail 检查邮箱是否存在。
func (r *cachedUserQueryRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.delegate.ExistsByEmail(ctx, email)
}

// List 获取用户列表。
func (r *cachedUserQueryRepository) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	return r.delegate.List(ctx, offset, limit)
}

// Count 统计用户数量。
func (r *cachedUserQueryRepository) Count(ctx context.Context) (int64, error) {
	return r.delegate.Count(ctx)
}

// Search 搜索用户。
func (r *cachedUserQueryRepository) Search(ctx context.Context, keyword string, offset, limit int) ([]*user.User, error) {
	return r.delegate.Search(ctx, keyword, offset, limit)
}

// CountBySearch 统计搜索结果数量。
func (r *cachedUserQueryRepository) CountBySearch(ctx context.Context, keyword string) (int64, error) {
	return r.delegate.CountBySearch(ctx, keyword)
}

// GetRoles 获取用户的所有角色 ID。
func (r *cachedUserQueryRepository) GetRoles(ctx context.Context, userID uint) ([]uint, error) {
	return r.delegate.GetRoles(ctx, userID)
}

// GetUserIDsByRole 获取拥有指定角色的所有用户 ID。
func (r *cachedUserQueryRepository) GetUserIDsByRole(ctx context.Context, roleID uint) ([]uint, error) {
	return r.delegate.GetUserIDsByRole(ctx, roleID)
}

var _ user.QueryRepository = (*cachedUserQueryRepository)(nil)
