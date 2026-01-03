package company

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/company"
)

// ListHandler 公司列表处理器。
type ListHandler struct {
	queryRepo company.QueryRepository
}

// NewListHandler 创建 ListHandler。
func NewListHandler(queryRepo company.QueryRepository) *ListHandler {
	return &ListHandler{queryRepo: queryRepo}
}

// Handle 处理公司列表查询。
func (h *ListHandler) Handle(ctx context.Context, query ListQuery) (*ListResultDTO, error) {
	var companies []*company.Company
	var total int64
	var err error

	switch {
	case query.Industry != nil:
		companies, err = h.queryRepo.ListByIndustry(ctx, *query.Industry, query.Offset, query.Limit)
		if err != nil {
			return nil, err
		}
		total, err = h.queryRepo.CountByIndustry(ctx, *query.Industry)
	case query.OwnerID != nil:
		companies, err = h.queryRepo.ListByOwner(ctx, *query.OwnerID, query.Offset, query.Limit)
		if err != nil {
			return nil, err
		}
		total, err = h.queryRepo.CountByOwner(ctx, *query.OwnerID)
	default:
		companies, err = h.queryRepo.List(ctx, query.Offset, query.Limit)
		if err != nil {
			return nil, err
		}
		total, err = h.queryRepo.Count(ctx)
	}

	if err != nil {
		return nil, err
	}

	return &ListResultDTO{
		Items: ToCompanyDTOs(companies),
		Total: total,
	}, nil
}
