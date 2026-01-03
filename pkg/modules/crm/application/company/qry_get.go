package company

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/company"
)

// GetHandler 获取公司处理器。
type GetHandler struct {
	queryRepo company.QueryRepository
}

// NewGetHandler 创建 GetHandler。
func NewGetHandler(queryRepo company.QueryRepository) *GetHandler {
	return &GetHandler{queryRepo: queryRepo}
}

// Handle 处理获取公司查询。
func (h *GetHandler) Handle(ctx context.Context, query GetQuery) (*CompanyDTO, error) {
	c, err := h.queryRepo.GetByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	return ToCompanyDTO(c), nil
}
