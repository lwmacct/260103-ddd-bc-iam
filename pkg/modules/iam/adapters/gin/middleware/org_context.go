package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/org"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/response"
)

// OrgContext 组织上下文中间件。
// 从路由参数提取 org_id，验证当前用户是否为组织成员。
// 验证通过后注入以下值到 Gin Context:
//   - org_id: uint - 组织 ID
//   - org_role: string - 用户在组织中的角色 (owner/admin/member)
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

		c.Next()
	}
}
