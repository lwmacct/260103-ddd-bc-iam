package contact

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/contact"
)

// GetHandler 获取联系人详情处理器。
type GetHandler struct {
	queryRepo contact.QueryRepository
}

// NewGetHandler 创建 GetHandler。
func NewGetHandler(queryRepo contact.QueryRepository) *GetHandler {
	return &GetHandler{queryRepo: queryRepo}
}

// Handle 处理获取联系人详情查询。
func (h *GetHandler) Handle(ctx context.Context, query GetQuery) (*ContactDTO, error) {
	entity, err := h.queryRepo.GetByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	return ToContactDTO(entity), nil
}
