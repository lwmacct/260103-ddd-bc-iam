package persistence

import (
	"context"
	"log/slog"

	appuser "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/user"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/user"
)

// cachedUserCommandRepository 带缓存失效的用户命令仓储（装饰器模式）。
//
// 在写操作完成后同步失效用户缓存，确保缓存与数据库一致性。
//
// 需要失效缓存的操作：
//   - Update: 用户信息变更
//   - Delete: 用户删除
//   - AssignRoles: 角色分配（影响权限）
//   - RemoveRoles: 角色移除（影响权限）
//   - UpdateStatus: 状态变更
type cachedUserCommandRepository struct {
	delegate user.CommandRepository
	cache    appuser.UserWithRolesCacheService
}

// NewCachedUserCommandRepository 创建带缓存失效的用户命令仓储。
func NewCachedUserCommandRepository(
	delegate user.CommandRepository,
	cache appuser.UserWithRolesCacheService,
) user.CommandRepository {
	return &cachedUserCommandRepository{
		delegate: delegate,
		cache:    cache,
	}
}

// Create 创建用户（不需要失效缓存，新用户没有缓存）。
func (r *cachedUserCommandRepository) Create(ctx context.Context, u *user.User) error {
	return r.delegate.Create(ctx, u)
}

// Update 更新用户。
func (r *cachedUserCommandRepository) Update(ctx context.Context, u *user.User) error {
	if err := r.delegate.Update(ctx, u); err != nil {
		return err
	}
	r.invalidateUserCache(ctx, u.ID)
	return nil
}

// Delete 删除用户。
func (r *cachedUserCommandRepository) Delete(ctx context.Context, id uint) error {
	if err := r.delegate.Delete(ctx, id); err != nil {
		return err
	}
	r.invalidateUserCache(ctx, id)
	return nil
}

// AssignRoles 为用户分配角色。
func (r *cachedUserCommandRepository) AssignRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	if err := r.delegate.AssignRoles(ctx, userID, roleIDs); err != nil {
		return err
	}
	r.invalidateUserCache(ctx, userID)
	return nil
}

// RemoveRoles 移除用户的角色。
func (r *cachedUserCommandRepository) RemoveRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	if err := r.delegate.RemoveRoles(ctx, userID, roleIDs); err != nil {
		return err
	}
	r.invalidateUserCache(ctx, userID)
	return nil
}

// UpdatePassword 更新用户密码（不需要失效缓存，密码不在缓存中）。
func (r *cachedUserCommandRepository) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	return r.delegate.UpdatePassword(ctx, userID, hashedPassword)
}

// UpdateStatus 更新用户状态。
func (r *cachedUserCommandRepository) UpdateStatus(ctx context.Context, userID uint, status string) error {
	if err := r.delegate.UpdateStatus(ctx, userID, status); err != nil {
		return err
	}
	r.invalidateUserCache(ctx, userID)
	return nil
}

// invalidateUserCache 同步失效用户缓存。
func (r *cachedUserCommandRepository) invalidateUserCache(ctx context.Context, userID uint) {
	if err := r.cache.InvalidateUser(ctx, userID); err != nil {
		slog.Warn("failed to invalidate user cache", "user_id", userID, "error", err.Error())
	}
}

var _ user.CommandRepository = (*cachedUserCommandRepository)(nil)
