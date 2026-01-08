package org

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/org"
)

// CreateHandler 创建组织命令处理器
type CreateHandler struct {
	orgCommandRepo    org.CommandRepository
	orgQueryRepo      org.QueryRepository
	memberCommandRepo org.MemberCommandRepository
}

// NewCreateHandler 创建组织命令处理器
func NewCreateHandler(
	orgCommandRepo org.CommandRepository,
	orgQueryRepo org.QueryRepository,
	memberCommandRepo org.MemberCommandRepository,
) *CreateHandler {
	return &CreateHandler{
		orgCommandRepo:    orgCommandRepo,
		orgQueryRepo:      orgQueryRepo,
		memberCommandRepo: memberCommandRepo,
	}
}

// Handle 处理创建组织命令
func (h *CreateHandler) Handle(ctx context.Context, cmd CreateOrgCommand) (*OrgDTO, error) {
	// 1. 检查组织名称是否已存在
	exists, err := h.orgQueryRepo.ExistsByName(ctx, cmd.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check org name existence: %w", err)
	}
	if exists {
		return nil, org.ErrOrgNameAlreadyExists
	}

	// 2. 创建组织实体
	organization := &org.Org{
		Name:        cmd.Name,
		DisplayName: cmd.DisplayName,
		Description: cmd.Description,
		Avatar:      cmd.Avatar,
		Status:      org.StatusActive,
	}

	// 3. 保存组织
	if err := h.orgCommandRepo.Create(ctx, organization); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// 4. 添加创建者为组织 owner
	if cmd.OwnerUserID > 0 {
		member := &org.Member{
			OrgID:    organization.ID,
			UserID:   cmd.OwnerUserID,
			Role:     org.MemberRoleOwner,
			JoinedAt: time.Now(),
		}
		if err := h.memberCommandRepo.Add(ctx, member); err != nil {
			return nil, fmt.Errorf("failed to add owner: %w", err)
		}
	}

	return ToOrgDTO(organization), nil
}
