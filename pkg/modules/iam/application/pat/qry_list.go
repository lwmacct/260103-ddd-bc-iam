package pat

import (
	"context"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/pat"
)

// ListHandler 获取 Token 列表查询处理器
type ListHandler struct {
	patQueryRepo pat.QueryRepository
}

// NewListHandler 创建 ListHandler 实例
func NewListHandler(patQueryRepo pat.QueryRepository) *ListHandler {
	return &ListHandler{
		patQueryRepo: patQueryRepo,
	}
}

// Handle 处理获取 Token 列表查询
func (h *ListHandler) Handle(ctx context.Context, query ListQuery) ([]*TokenInfoDTO, error) {
	tokens, err := h.patQueryRepo.ListByUser(ctx, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tokens: %w", err)
	}

	// 转换为 DTO
	tokenResponses := make([]*TokenInfoDTO, 0, len(tokens))
	for _, token := range tokens {
		tokenResponses = append(tokenResponses, ToTokenInfoDTO(token))
	}

	return tokenResponses, nil
}
