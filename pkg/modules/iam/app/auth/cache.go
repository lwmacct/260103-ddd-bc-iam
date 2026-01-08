package auth

import (
	"context"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/role"
)

// PermissionCacheService 权限缓存服务接口。
//
// 提供用户权限的缓存操作，用于高并发场景下的权限检查优化。
// 新 RBAC 模型：权限为 Operation + Resource Pattern 组合。
//
// Key 命名规范：{prefix}user:perms:{userID}
// 默认 TTL：5 分钟
//
// 缓存失效策略：
//   - TTL 过期自动失效
//   - 用户角色变更时主动失效（通过 EventBus）
//   - 角色权限变更时批量失效关联用户
//
// 实现位于 infrastructure/cache 包。
type PermissionCacheService interface {
	// GetUserPermissions 获取用户权限（角色名 + 权限对象）。
	// permissions 为 role.Permission 切片，支持 OperationPattern/ResourcePattern 匹配。
	// 缓存未命中返回三个 nil（不返回错误）。
	// 缓存数据损坏时自动清除并返回三个 nil。
	GetUserPermissions(ctx context.Context, userID uint) (roles []string, permissions []role.Permission, err error)

	// SetUserPermissions 设置用户权限缓存。
	// 使用默认 TTL（5 分钟）。
	SetUserPermissions(ctx context.Context, userID uint, roles []string, permissions []role.Permission) error

	// InvalidateUser 失效单个用户缓存。
	// 用于用户角色变更、用户状态变更等场景。
	InvalidateUser(ctx context.Context, userID uint) error

	// InvalidateUsers 批量失效用户缓存。
	// userIDs 为空时直接返回 nil。
	InvalidateUsers(ctx context.Context, userIDs []uint) error

	// InvalidateAll 失效所有用户权限缓存。
	// 使用 SCAN 命令遍历，慎用（仅用于全局权限系统变更）。
	InvalidateAll(ctx context.Context) error
}
