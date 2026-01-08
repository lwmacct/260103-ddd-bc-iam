package middleware

import (
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/role"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/ctxutil"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/permission"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/response"
)

// RequireOperation 检查用户是否有执行指定 Operation 的权限。
// URN 风格 RBAC：权限为 Operation Pattern + Resource Pattern 组合。
// 根据 Operation 的 Scope 自动检查对应资源：
//   - self: 域 → 检查 self:user:@me
//   - org: 域 → 检查 org.{org_id}:*:* 和 org.{org_id}.team.{team_id}:*:*
func RequireOperation(op permission.Operation) gin.HandlerFunc {
	return func(c *gin.Context) {
		permList, ok := ctxutil.Get[[]role.Permission](c, ctxutil.Permissions)
		if !ok {
			response.Unauthorized(c, "No permissions found")
			c.Abort()
			return
		}

		// 检查是否有匹配的权限
		hasPermission := false
		operationStr := string(op)

		// 确定要检查的资源
		resourcesToCheck := []string{"*:*:*"}
		scope := op.Scope()
		switch scope {
		case "self":
			resourcesToCheck = appendSelfResources(c, resourcesToCheck)
		case "org":
			resourcesToCheck = appendOrgResources(c, resourcesToCheck)
		}

		for _, p := range permList {
			// 匹配 Operation Pattern
			if permission.MatchOperation(p.OperationPattern, operationStr) {
				// 检查资源是否匹配任一候选资源
				for _, res := range resourcesToCheck {
					if permission.MatchResource(p.ResourcePattern, res) {
						hasPermission = true
						break
					}
				}
				if hasPermission {
					break
				}
			}
		}

		if !hasPermission {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireOperationWithResource 检查用户是否有对指定资源执行指定 Operation 的权限。
// 支持细粒度资源控制，如 sys:user:123、sys:role:*。
func RequireOperationWithResource(op permission.Operation, resource permission.Resource) gin.HandlerFunc {
	return func(c *gin.Context) {
		permList, ok := ctxutil.Get[[]role.Permission](c, ctxutil.Permissions)
		if !ok {
			response.Unauthorized(c, "No permissions found")
			c.Abort()
			return
		}

		hasPermission := false
		operationStr := string(op)
		resourceStr := string(resource)
		for _, p := range permList {
			if permission.MatchOperation(p.OperationPattern, operationStr) &&
				permission.MatchResource(p.ResourcePattern, resourceStr) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole creates a middleware that checks if the user has a specific role
func RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesList, ok := ctxutil.Get[[]string](c, ctxutil.Roles)
		if !ok {
			response.Unauthorized(c, "No roles found")
			c.Abort()
			return
		}

		if !slices.Contains(rolesList, roleName) {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole creates a middleware that checks if the user has any of the specified roles
func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesList, ok := ctxutil.Get[[]string](c, ctxutil.Roles)
		if !ok {
			response.Unauthorized(c, "No roles found")
			c.Abort()
			return
		}

		hasRole := slices.ContainsFunc(roles, func(requiredRole string) bool {
			return slices.Contains(rolesList, requiredRole)
		})

		if !hasRole {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireOwnership creates a middleware that checks if the user is accessing their own resource
// The resource ID should be in the URL parameter specified by paramName (default: "id")
func RequireOwnership(paramName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
		if !ok {
			response.Unauthorized(c, "No user ID found")
			c.Abort()
			return
		}

		param := "id"
		if len(paramName) > 0 {
			param = paramName[0]
		}

		resourceIDStr := c.Param(param)
		resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid resource ID")
			c.Abort()
			return
		}

		if uint(resourceID) != uid {
			response.Forbidden(c, "Can only access own resources")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdminOrOwnership combines admin role check with ownership check
// Allows access if user is admin OR owns the resource
func RequireAdminOrOwnership(paramName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is admin
		if isAdmin(c) {
			c.Next()
			return
		}

		// If not admin, check ownership
		uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
		if !ok {
			response.Unauthorized(c, "No user ID found")
			c.Abort()
			return
		}

		param := "id"
		if len(paramName) > 0 {
			param = paramName[0]
		}

		resourceIDStr := c.Param(param)
		resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid resource ID")
			c.Abort()
			return
		}

		if uint(resourceID) != uid {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// isAdmin 检查当前用户是否具有 admin 角色
func isAdmin(c *gin.Context) bool {
	roles, ok := ctxutil.Get[[]string](c, ctxutil.Roles)
	if !ok {
		return false
	}
	return slices.Contains(roles, "admin")
}

// buildContextVars 从 Gin Context 构建变量映射。
// 收集 @me, @org, @team 等运行时变量。
func buildContextVars(c *gin.Context) map[string]string {
	vars := make(map[string]string)

	// @me - 当前用户 ID
	if userID, ok := ctxutil.Get[uint](c, ctxutil.UserID); ok {
		vars["@me"] = strconv.FormatUint(uint64(userID), 10)
	}

	// @org - 当前组织 ID（由 OrgContext 中间件注入）
	if orgID, ok := ctxutil.Get[uint](c, ctxutil.OrgID); ok {
		vars["@org"] = strconv.FormatUint(uint64(orgID), 10)
	}

	// @team - 当前团队 ID（由 TeamContext 中间件注入）
	if teamID, ok := ctxutil.Get[uint](c, ctxutil.TeamID); ok {
		vars["@team"] = strconv.FormatUint(uint64(teamID), 10)
	}

	return vars
}

// appendSelfResources 为 self: 域操作添加 self:user:@me 资源（解析后）
func appendSelfResources(c *gin.Context, resources []string) []string {
	vars := buildContextVars(c)
	if vars["@me"] == "" {
		return resources
	}

	resolver := permission.NewResolver(vars)
	selfUserResource := resolver.ResolveString("self:user:@me")

	return append(resources,
		"self:user:@me",  // 原始模式（用于匹配 self:user:@me 模式的权限）
		selfUserResource, // 解析后的具体资源（如 self:user:123）
	)
}

// appendOrgResources 为 org: 域操作添加组织级资源。
// 返回 org.{org_id}:*:* 形式的资源标识符。
func appendOrgResources(c *gin.Context, resources []string) []string {
	vars := buildContextVars(c)
	if vars["@org"] == "" {
		return resources
	}

	resolver := permission.NewResolver(vars)

	// 组织级资源：org.{org_id}:*:*
	orgResource := resolver.ResolveString("org.@org:*:*")
	resources = append(resources,
		"org.@org:*:*", // 原始模式
		orgResource,    // 解析后（如 org.123:*:*）
	)

	// 如果有团队上下文，添加团队级资源
	if vars["@team"] != "" {
		teamResource := resolver.ResolveString("org.@org.team.@team:*:*")
		resources = append(resources,
			"org.@org.team.@team:*:*", // 原始模式
			teamResource,              // 解析后（如 org.123.team.456:*:*）
		)
	}

	return resources
}
