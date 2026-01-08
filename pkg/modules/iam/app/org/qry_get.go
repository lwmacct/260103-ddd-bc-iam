package org

import (
	"context"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/org"
)

// GetHandler 获取组织查询处理器
type GetHandler struct {
	orgQueryRepo org.QueryRepository
}

// NewGetHandler 创建获取组织查询处理器
func NewGetHandler(orgQueryRepo org.QueryRepository) *GetHandler {
	return &GetHandler{orgQueryRepo: orgQueryRepo}
}

// Handle 处理获取组织查询
func (h *GetHandler) Handle(ctx context.Context, query GetOrgQuery) (*OrgDTO, error) {
	org, err := h.orgQueryRepo.GetByID(ctx, query.OrgID)
	if err != nil {
		return nil, err
	}
	return ToOrgDTO(org), nil
}
