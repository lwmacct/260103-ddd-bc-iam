package middleware

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/org"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/role"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/ctxutil"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/response"
)

// OrgContext 组织上下文中间件。
// 从路由参数提取 org_id，验证当前用户是否为组织成员。
// 验证通过后注入以下值到 Gin Context:
//   - org_id: uint - 组织 ID
//   - org_role: string - 用户在组织中的角色 (owner/admin/member)
//
// 权限注入：根据 org_role 动态添加组织级权限到 permissions 列表：
//   - owner: org:*:* 对 org.{id}:*:*（组织内所有操作）
//   - admin: org:settings:* + org:team:* 对 org.{id}:*:*（配置和团队管理）
func OrgContext(memberQuery org.MemberQueryRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取当前用户 ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.Unauthorized(c, "user not authenticated")
			c.Abort()
			return
		}
		uid, ok := userID.(uint)
		if !ok {
			response.InternalError(c, "invalid user ID format")
			c.Abort()
			return
		}

		// 2. 解析路由参数中的 org_id
		orgIDStr := c.Param("org_id")
		if orgIDStr == "" {
			response.BadRequest(c, "org_id is required")
			c.Abort()
			return
		}
		orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "invalid org_id")
			c.Abort()
			return
		}

		// 3. 验证用户是否是组织成员
		member, err := memberQuery.GetByOrgAndUser(c.Request.Context(), uint(orgID), uid)
		if err != nil {
			response.Forbidden(c, "not a member of this organization")
			c.Abort()
			return
		}

		// 4. 注入组织上下文
		c.Set("org_id", uint(orgID))
		c.Set("org_role", string(member.Role))

		// 5. 动态注入基于 org_role 的权限
		// 组织管理员获得组织级操作权限，无需在全局 RBAC 中预配置
		permissions, _ := ctxutil.Get[[]role.Permission](c, ctxutil.Permissions)
		orgIDUint := uint(orgID)
		switch member.Role {
		case "owner":
			// owner 获得组织内所有操作权限
			permissions = append(permissions, role.Permission{
				OperationPattern: "org:*:*",
				ResourcePattern:  fmt.Sprintf("org.%d:*:*", orgIDUint),
			})
		case "admin":
			// admin 获得组织配置和团队管理权限
			// 注意：URN 解析器只支持 3 段，需用 org:team:* 而非 org:team:*:*
			permissions = append(permissions,
				role.Permission{
					OperationPattern: "org:settings:*",
					ResourcePattern:  fmt.Sprintf("org.%d:*:*", orgIDUint),
				},
				role.Permission{
					OperationPattern: "org:team:*",
					ResourcePattern:  fmt.Sprintf("org.%d.team.*:*:*", orgIDUint),
				},
			)
		default:
			// member 角色不注入额外权限
		}
		c.Set(ctxutil.Permissions, permissions)

		c.Next()
	}
}
