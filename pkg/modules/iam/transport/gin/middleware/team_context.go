package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/org"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
)

// TeamContext 团队上下文中间件。
// 必须在 OrgContext 之后使用。
// 从路由参数提取 team_id，验证团队属于组织。
// 组织管理员（owner/admin）可以访问所有团队，普通成员必须是团队成员。
// 验证通过后注入以下值到 Gin Context:
//   - team_id: uint - 团队 ID
//   - team_role: string - 用户在团队中的角色 (仅当是团队成员时)
func TeamContext(
	teamQuery org.TeamQueryRepository,
	teamMemberQuery org.TeamMemberQueryRepository,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取组织上下文（必须先通过 OrgContext）
		orgID, exists := c.Get("org_id")
		if !exists {
			response.InternalError(c, "org context not found, OrgContext middleware required")
			c.Abort()
			return
		}

		// 2. 解析路由参数中的 team_id
		teamIDStr := c.Param("team_id")
		if teamIDStr == "" {
			response.BadRequest(c, "team_id is required")
			c.Abort()
			return
		}
		teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "invalid team_id")
			c.Abort()
			return
		}

		// 3. 验证团队属于该组织
		team, err := teamQuery.GetByID(c.Request.Context(), uint(teamID))
		if err != nil {
			response.NotFoundMessage(c, "team not found")
			c.Abort()
			return
		}
		if team.OrgID != orgID.(uint) {
			response.Forbidden(c, "team does not belong to this organization")
			c.Abort()
			return
		}

		// 4. 检查用户权限
		userID := c.MustGet("user_id").(uint)
		orgRole := c.MustGet("org_role").(string)

		// 组织管理员可以访问所有团队
		if orgRole == "owner" || orgRole == "admin" {
			c.Set("team_id", uint(teamID))
			c.Next()
			return
		}

		// 普通成员必须是团队成员
		teamMember, err := teamMemberQuery.GetByTeamAndUser(c.Request.Context(), uint(teamID), userID)
		if err != nil {
			response.Forbidden(c, "not a member of this team")
			c.Abort()
			return
		}

		// 5. 注入团队上下文
		c.Set("team_id", uint(teamID))
		c.Set("team_role", string(teamMember.Role))

		c.Next()
	}
}
