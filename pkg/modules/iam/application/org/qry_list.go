package org

import (
	"context"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/org"
)

// ListHandler 组织列表查询处理器
type ListHandler struct {
	orgQueryRepo org.QueryRepository
}

// NewListHandler 创建组织列表查询处理器
func NewListHandler(orgQueryRepo org.QueryRepository) *ListHandler {
	return &ListHandler{orgQueryRepo: orgQueryRepo}
}

// ListResult 列表查询结果
type ListResult struct {
	Items []*OrgDTO
	Total int64
}

// Handle 处理组织列表查询
func (h *ListHandler) Handle(ctx context.Context, query ListOrgsQuery) (*ListResult, error) {
	if query.Keyword != "" {
		return h.handleSearch(ctx, query)
	}
	return h.handleList(ctx, query)
}

func (h *ListHandler) handleSearch(ctx context.Context, query ListOrgsQuery) (*ListResult, error) {
	orgs, err := h.orgQueryRepo.Search(ctx, query.Keyword, query.Offset, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search organizations: %w", err)
	}
	total, err := h.orgQueryRepo.CountBySearch(ctx, query.Keyword)
	if err != nil {
		return nil, fmt.Errorf("failed to count search results: %w", err)
	}
	return &ListResult{Items: ToOrgDTOs(orgs), Total: total}, nil
}

func (h *ListHandler) handleList(ctx context.Context, query ListOrgsQuery) (*ListResult, error) {
	orgs, err := h.orgQueryRepo.List(ctx, query.Offset, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}
	total, err := h.orgQueryRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count organizations: %w", err)
	}
	return &ListResult{Items: ToOrgDTOs(orgs), Total: total}, nil
}
