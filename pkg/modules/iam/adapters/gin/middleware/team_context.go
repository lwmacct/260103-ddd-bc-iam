package middleware

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/org"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/role"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/ctxutil"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/response"
)

// TeamContext 团队上下文中间件。
// 必须在 OrgContext 之后使用。
// 从路由参数提取 team_id，验证团队属于组织。
// 组织管理员（owner/admin）可以访问所有团队，普通成员必须是团队成员。
// 验证通过后注入以下值到 Gin Context:
//   - team_id: uint - 团队 ID
//   - team_role: string - 用户在团队中的角色 (仅当是团队成员时)
//
// 权限注入：根据 team_role 动态添加团队级权限到 permissions 列表：
//   - lead: org:team:* 对 org.{oid}.team.{tid}:*:*（团队操作权限）
func TeamContext(
	teamQuery org.TeamQueryRepository,
	teamMemberQuery org.TeamMemberQueryRepository,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取组织上下文（必须先通过 OrgContext）
		orgID, exists := c.Get(ctxutil.OrgID)
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
		userID := c.MustGet(ctxutil.UserID).(uint)
		orgRole := c.MustGet(ctxutil.OrgRole).(string)

		// 组织管理员可以访问所有团队
		if orgRole == "owner" || orgRole == "admin" {
			c.Set(ctxutil.TeamID, uint(teamID))
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
		c.Set(ctxutil.TeamID, uint(teamID))
		c.Set(ctxutil.TeamRole, string(teamMember.Role))

		// 6. 动态注入基于 team_role 的权限
		// 团队负责人获得团队配置管理权限
		// 注意：URN 解析器只支持 3 段，org:team:settings:list 解析为 Identifier="settings:list"
		// 因此通配符必须用 org:team:* 而非 org:team:settings:*
		if teamMember.Role == "lead" {
			permissions, _ := ctxutil.Get[[]role.Permission](c, ctxutil.Permissions)
			permissions = append(permissions, role.Permission{
				OperationPattern: "org:team:*",
				ResourcePattern:  fmt.Sprintf("org.%d.team.%d:*:*", orgID.(uint), uint(teamID)),
			})
			c.Set(ctxutil.Permissions, permissions)
		}

		c.Next()
	}
}
