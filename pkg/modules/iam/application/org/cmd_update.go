package org

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/org"
)

// UpdateHandler 更新组织命令处理器
type UpdateHandler struct {
	orgCommandRepo org.CommandRepository
	orgQueryRepo   org.QueryRepository
}

// NewUpdateHandler 创建更新组织命令处理器
func NewUpdateHandler(
	orgCommandRepo org.CommandRepository,
	orgQueryRepo org.QueryRepository,
) *UpdateHandler {
	return &UpdateHandler{
		orgCommandRepo: orgCommandRepo,
		orgQueryRepo:   orgQueryRepo,
	}
}

// Handle 处理更新组织命令
func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateOrgCommand) (*OrgDTO, error) {
	// 1. 获取现有组织
	org, err := h.orgQueryRepo.GetByID(ctx, cmd.OrgID)
	if err != nil {
		return nil, err
	}

	// 2. 更新字段
	if cmd.DisplayName != nil {
		org.DisplayName = *cmd.DisplayName
	}
	if cmd.Description != nil {
		org.Description = *cmd.Description
	}
	if cmd.Avatar != nil {
		org.Avatar = *cmd.Avatar
	}
	if cmd.Status != nil {
		org.Status = *cmd.Status
	}

	// 3. 保存更新
	if err := h.orgCommandRepo.Update(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	return ToOrgDTO(org), nil
}
