// Package middleware 提供 IAM 模块的 HTTP 认证和授权中间件。
//
// 本包实现了基于 JWT/PAT 的认证和基于 URN 的 RBAC 授权：
//
// 认证中间件：
//   - [Auth]: 统一认证（支持 JWT 和 PAT 双模式）
//
// 授权中间件：
//   - [RequireOperation]: 权限检查（URN 风格 RBAC）
//   - [RequireOperationWithResource]: 权限 + 资源检查
//
// # 权限缓存机制
//
// JWT/PAT 仅存储 user_id，权限信息从 PermissionCacheService
// 实时查询，支持权限变更后立即生效。
//
// # PAT Scope 过滤
//
// PAT 认证时，根据 PAT 的 Scopes 字段过滤用户权限。
// 例如 Scope 为 ["self"] 时，只保留 self:* 前缀的权限。
//
// # 依赖注入
//
// 中间件需要通过应用层依赖注入获取服务实例：
//
//	// 应用层（internal/container）
//	authMW := middleware.Auth(jwtManager, patService, permCacheService)
//	rbacMW := middleware.RequireOperation(permission.Operation("admin:users:create"))
package middleware
