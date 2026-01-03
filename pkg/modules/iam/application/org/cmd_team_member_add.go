package org

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/org"
)

// TeamMemberAddHandler 添加团队成员命令处理器
type TeamMemberAddHandler struct {
	teamMemberCommand org.TeamMemberCommandRepository
	teamMemberQuery   org.TeamMemberQueryRepository
	teamQuery         org.TeamQueryRepository
	memberQuery       org.MemberQueryRepository
}

// NewTeamMemberAddHandler 创建添加团队成员命令处理器
func NewTeamMemberAddHandler(
	teamMemberCommand org.TeamMemberCommandRepository,
	teamMemberQuery org.TeamMemberQueryRepository,
	teamQuery org.TeamQueryRepository,
	memberQuery org.MemberQueryRepository,
) *TeamMemberAddHandler {
	return &TeamMemberAddHandler{
		teamMemberCommand: teamMemberCommand,
		teamMemberQuery:   teamMemberQuery,
		teamQuery:         teamQuery,
		memberQuery:       memberQuery,
	}
}

// Handle 处理添加团队成员命令
func (h *TeamMemberAddHandler) Handle(ctx context.Context, cmd AddTeamMemberCommand) (*TeamMemberDTO, error) {
	// 1. 验证团队属于组织
	belongsTo, err := h.teamQuery.BelongsToOrg(ctx, cmd.TeamID, cmd.OrgID)
	if err != nil {
		return nil, fmt.Errorf("failed to check team org: %w", err)
	}
	if !belongsTo {
		return nil, org.ErrTeamNotInOrg
	}

	// 2. 检查用户是否是组织成员（必须先是组织成员才能加入团队）
	isOrgMember, err := h.memberQuery.IsMember(ctx, cmd.OrgID, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check org membership: %w", err)
	}
	if !isOrgMember {
		return nil, org.ErrMustBeOrgMemberFirst
	}

	// 3. 检查是否已是团队成员
	isTeamMember, err := h.teamMemberQuery.IsMember(ctx, cmd.TeamID, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check team membership: %w", err)
	}
	if isTeamMember {
		return nil, org.ErrMemberAlreadyExists
	}

	// 4. 验证角色
	role := org.TeamMemberRole(cmd.Role)
	if !org.IsValidTeamMemberRole(role) {
		return nil, org.ErrInvalidTeamMemberRole
	}

	// 5. 添加团队成员
	member := &org.TeamMember{
		TeamID:   cmd.TeamID,
		UserID:   cmd.UserID,
		Role:     role,
		JoinedAt: time.Now(),
	}

	if err := h.teamMemberCommand.Add(ctx, member); err != nil {
		return nil, fmt.Errorf("failed to add team member: %w", err)
	}

	return ToTeamMemberDTO(member), nil
}
