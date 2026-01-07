package auth

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	appauth "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/auth"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/role"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/user"
)

// mergePermissions 合并权限列表（去重）。
// 使用 OperationPattern + ResourcePattern 作为去重键。
func mergePermissions(base, additional []role.Permission) []role.Permission {
	seen := make(map[string]bool, len(base))
	for _, p := range base {
		seen[p.OperationPattern+"|"+p.ResourcePattern] = true
	}

	result := append([]role.Permission{}, base...)
	for _, p := range additional {
		key := p.OperationPattern + "|" + p.ResourcePattern
		if !seen[key] {
			seen[key] = true
			result = append(result, p)
		}
	}

	return result
}

// PermissionCacheService 权限缓存服务（Cache-Aside 模式）
//
// 本服务封装 Cache-Aside 逻辑，供中间件使用：
//   - GetUserPermissions: 先查缓存，未命中则查数据库并回写缓存
//   - InvalidateUsersWithRole: 按角色批量失效（需查询数据库获取用户列表）
//
// 新 RBAC 模型：返回 []role.Permission 支持 Operation + Resource Pattern 匹配。
//
// 隐性角色：
// 所有已认证用户自动拥有 [role.DefaultUserRoleName] 角色及其权限，
// 无需在数据库中显式分配。
//
// 性能优化：
// 查询数据库时，会同时写入用户实体缓存（UserWithRolesCacheService），
// 供后续 Handler 通过 Repository 缓存装饰器直接命中，避免重复查询。
//
// 底层缓存操作委托给 [appauth.PermissionCacheService] 接口实现。
type PermissionCacheService struct {
	cache         appauth.PermissionCacheService
	userCache     user.UserWithRolesCacheService
	userQueryRepo user.QueryRepository
	roleQueryRepo role.QueryRepository // 用于查询默认角色权限
}

// NewPermissionCacheService 创建权限缓存服务
func NewPermissionCacheService(
	cacheService appauth.PermissionCacheService,
	userCacheService user.UserWithRolesCacheService,
	userQueryRepo user.QueryRepository,
	roleQueryRepo role.QueryRepository,
) *PermissionCacheService {
	return &PermissionCacheService{
		cache:         cacheService,
		userCache:     userCacheService,
		userQueryRepo: userQueryRepo,
		roleQueryRepo: roleQueryRepo,
	}
}

// GetUserPermissions 获取用户权限（Cache-Aside 模式）
//
// 执行流程：
//  1. 尝试从缓存读取
//  2. 缓存未命中，查询数据库
//  3. root 用户：直接返回最高权限（硬编码）
//  4. 普通用户：追加默认角色及其权限（隐性角色）
//  5. 同步写入权限缓存
//  6. 同步写入用户实体缓存（供后续 Handler 通过 Repository 装饰器命中）
//
// 返回 []role.Permission 支持 Operation + Resource Pattern 匹配。
func (s *PermissionCacheService) GetUserPermissions(ctx context.Context, userID uint) ([]string, []role.Permission, error) {
	// 1. 尝试从缓存读取（缓存已包含默认角色权限）
	roles, permissions, err := s.cache.GetUserPermissions(ctx, userID)
	// 注意：空切片不等于 nil，需要检查长度
	if err == nil && len(roles) > 0 {
		return roles, permissions, nil
	}

	// 2. 缓存未命中，查询数据库
	u, err := s.userQueryRepo.GetByIDWithRoles(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 3. root 用户：直接返回最高权限（硬编码，无需数据库角色）
	if u.Username == user.RootUsername {
		roles = []string{"root"}
		permissions = []role.Permission{
			{OperationPattern: "*:*:*", ResourcePattern: "*"},
		}
	} else {
		// 4. 普通用户：从角色获取权限 + 追加默认角色
		roles = u.GetRoleNames()
		permissions = u.GetPermissions()
		roles, permissions = s.appendDefaultRole(ctx, roles, permissions)
	}

	// 5. 同步写入权限缓存（Redis 写入 < 1ms，延迟可忽略）
	if err := s.cache.SetUserPermissions(ctx, userID, roles, permissions); err != nil {
		slog.WarnContext(ctx, "Failed to cache user permissions",
			"user_id", userID,
			"error", err.Error(),
		)
	}

	// 6. 同步写入用户实体缓存（关键！供后续 Handler 通过 Repository 装饰器命中）
	if err := s.userCache.SetUserWithRoles(ctx, u); err != nil {
		slog.WarnContext(ctx, "Failed to cache user entity",
			"user_id", userID,
			"error", err.Error(),
		)
	}

	return roles, permissions, nil
}

// InvalidateUser 清除指定用户的权限缓存
// 用于用户角色变更、用户状态变更等场景
func (s *PermissionCacheService) InvalidateUser(ctx context.Context, userID uint) error {
	if err := s.cache.InvalidateUser(ctx, userID); err != nil {
		slog.Warn("Failed to invalidate user cache",
			"user_id", userID,
			"error", err,
		)
		return err
	}

	slog.Debug("User permissions cache invalidated",
		"user_id", userID,
	)

	return nil
}

// InvalidateUsersWithRole 清除拥有指定角色的所有用户缓存
// 用于角色权限变更场景（需要使所有拥有该角色的用户缓存失效）
func (s *PermissionCacheService) InvalidateUsersWithRole(ctx context.Context, roleID uint) error {
	// 查询所有拥有该角色的用户
	userIDs, err := s.userQueryRepo.GetUserIDsByRole(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to get users with role %d: %w", roleID, err)
	}

	// 批量清除缓存
	if len(userIDs) > 0 {
		if err := s.cache.InvalidateUsers(ctx, userIDs); err != nil {
			slog.Warn("Failed to invalidate users cache",
				"role_id", roleID,
				"user_count", len(userIDs),
				"error", err,
			)
			return err
		}

		slog.Info("Users permissions cache invalidated",
			"role_id", roleID,
			"user_count", len(userIDs),
		)
	}

	return nil
}

// InvalidateAllUsers 清除所有用户权限缓存（慎用，仅用于全局权限系统变更）
func (s *PermissionCacheService) InvalidateAllUsers(ctx context.Context) error {
	if err := s.cache.InvalidateAll(ctx); err != nil {
		return fmt.Errorf("failed to invalidate all users cache: %w", err)
	}

	slog.Info("All users permissions cache invalidated")
	return nil
}

// appendDefaultRole 追加默认用户角色及其权限。
//
// 所有已认证用户隐性拥有 [role.DefaultUserRoleName] 角色。
// 如果用户已有该角色，不重复追加。
func (s *PermissionCacheService) appendDefaultRole(
	ctx context.Context,
	roles []string,
	permissions []role.Permission,
) ([]string, []role.Permission) {
	// 检查用户是否已有默认角色
	if slices.Contains(roles, role.DefaultUserRoleName) {
		return roles, permissions // 已有，无需追加
	}

	// 追加角色名
	roles = append(roles, role.DefaultUserRoleName)

	// 查询默认角色的权限
	defaultRole, err := s.roleQueryRepo.FindByName(ctx, role.DefaultUserRoleName)
	if err != nil {
		slog.Warn("Failed to get default user role, skipping",
			"role_name", role.DefaultUserRoleName,
			"error", err,
		)
		return roles, permissions
	}

	// 合并权限（去重）
	permissions = mergePermissions(permissions, defaultRole.Permissions)

	return roles, permissions
}
