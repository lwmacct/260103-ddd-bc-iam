package opportunity

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/opportunity"
)

// GetHandler 获取商机处理器。
type GetHandler struct {
	queryRepo opportunity.QueryRepository
}

// NewGetHandler 创建处理器实例。
func NewGetHandler(queryRepo opportunity.QueryRepository) *GetHandler {
	return &GetHandler{queryRepo: queryRepo}
}

// Handle 执行获取商机查询。
func (h *GetHandler) Handle(ctx context.Context, query GetQuery) (*OpportunityDTO, error) {
	opp, err := h.queryRepo.GetByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	return ToOpportunityDTO(opp), nil
}
