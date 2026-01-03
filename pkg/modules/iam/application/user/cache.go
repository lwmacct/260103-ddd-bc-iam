package user

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/user"
)

// UserWithRolesCacheService 用户实体缓存服务接口。
//
// 缓存完整用户实体（包含角色和权限信息），用于优化高频读取场景。
//
// Key 命名规范：{prefix}user:entity:{userID}
// 默认 TTL：5 分钟（与权限缓存保持一致）
//
// 缓存失效策略：
//   - TTL 过期自动失效
//   - 用户信息变更时主动失效
//   - 用户角色变更时主动失效
//
// 实现位于 [infrastructure/cache.userWithRolesCacheService]。
type UserWithRolesCacheService interface {
	// GetUserWithRoles 获取缓存的用户实体（含角色和权限）。
	// 缓存未命中返回 nil, nil（不返回错误）。
	// 缓存数据损坏时自动清除并返回 nil, nil。
	GetUserWithRoles(ctx context.Context, userID uint) (*user.User, error)

	// SetUserWithRoles 设置用户实体缓存。
	// 使用默认 TTL（5 分钟）。
	SetUserWithRoles(ctx context.Context, u *user.User) error

	// InvalidateUser 失效单个用户缓存。
	// 用于用户信息变更、角色变更等场景。
	InvalidateUser(ctx context.Context, userID uint) error
}
