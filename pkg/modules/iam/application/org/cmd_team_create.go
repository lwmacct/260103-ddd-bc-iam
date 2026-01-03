package org

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/org"
)

// TeamCreateHandler 创建团队命令处理器
type TeamCreateHandler struct {
	teamCommand       org.TeamCommandRepository
	teamQuery         org.TeamQueryRepository
	orgQuery          org.QueryRepository
	teamMemberCommand org.TeamMemberCommandRepository
}

// NewTeamCreateHandler 创建团队命令处理器
func NewTeamCreateHandler(
	teamCommand org.TeamCommandRepository,
	teamQuery org.TeamQueryRepository,
	orgQuery org.QueryRepository,
	teamMemberCommand org.TeamMemberCommandRepository,
) *TeamCreateHandler {
	return &TeamCreateHandler{
		teamCommand:       teamCommand,
		teamQuery:         teamQuery,
		orgQuery:          orgQuery,
		teamMemberCommand: teamMemberCommand,
	}
}

// Handle 处理创建团队命令
func (h *TeamCreateHandler) Handle(ctx context.Context, cmd CreateTeamCommand) (*TeamDTO, error) {
	// 1. 检查组织是否存在
	exists, err := h.orgQuery.Exists(ctx, cmd.OrgID)
	if err != nil {
		return nil, fmt.Errorf("failed to check org existence: %w", err)
	}
	if !exists {
		return nil, org.ErrOrgNotFound
	}

	// 2. 检查团队名称是否已存在（组织内唯一）
	exists, err = h.teamQuery.ExistsByOrgAndName(ctx, cmd.OrgID, cmd.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check team name existence: %w", err)
	}
	if exists {
		return nil, org.ErrTeamNameAlreadyExists
	}

	// 3. 创建团队实体
	team := &org.Team{
		OrgID:       cmd.OrgID,
		Name:        cmd.Name,
		DisplayName: cmd.DisplayName,
		Description: cmd.Description,
		Avatar:      cmd.Avatar,
	}

	// 4. 保存团队
	if err := h.teamCommand.Create(ctx, team); err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	// 5. 添加团队负责人（如果提供）
	if cmd.LeadUserID > 0 {
		member := &org.TeamMember{
			TeamID:   team.ID,
			UserID:   cmd.LeadUserID,
			Role:     org.TeamMemberRoleLead,
			JoinedAt: time.Now(),
		}
		if err := h.teamMemberCommand.Add(ctx, member); err != nil {
			return nil, fmt.Errorf("failed to add team lead: %w", err)
		}
	}

	return ToTeamDTO(team), nil
}
