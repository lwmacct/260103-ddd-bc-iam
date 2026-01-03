package lead

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/lead"
)

// GetHandler 获取线索处理器。
type GetHandler struct {
	queryRepo lead.QueryRepository
}

// NewGetHandler 创建 GetHandler。
func NewGetHandler(queryRepo lead.QueryRepository) *GetHandler {
	return &GetHandler{queryRepo: queryRepo}
}

// Handle 处理获取线索查询。
func (h *GetHandler) Handle(ctx context.Context, query GetQuery) (*LeadDTO, error) {
	l, err := h.queryRepo.GetByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	return ToLeadDTO(l), nil
}
