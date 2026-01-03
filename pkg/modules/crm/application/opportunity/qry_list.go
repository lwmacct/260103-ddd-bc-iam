package opportunity

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/opportunity"
)

// ListHandler 商机列表处理器。
type ListHandler struct {
	queryRepo opportunity.QueryRepository
}

// NewListHandler 创建处理器实例。
func NewListHandler(queryRepo opportunity.QueryRepository) *ListHandler {
	return &ListHandler{queryRepo: queryRepo}
}

// Handle 执行商机列表查询。
func (h *ListHandler) Handle(ctx context.Context, query ListQuery) (*OpportunityListDTO, error) {
	var (
		items []*opportunity.Opportunity
		total int64
		err   error
	)

	switch {
	case query.Stage != nil:
		items, err = h.queryRepo.ListByStage(ctx, *query.Stage, query.Offset, query.Limit)
		if err != nil {
			return nil, err
		}
		total, err = h.queryRepo.CountByStage(ctx, *query.Stage)
	case query.OwnerID != nil:
		items, err = h.queryRepo.ListByOwner(ctx, *query.OwnerID, query.Offset, query.Limit)
		if err != nil {
			return nil, err
		}
		total, err = h.queryRepo.CountByOwner(ctx, *query.OwnerID)
	default:
		items, err = h.queryRepo.List(ctx, query.Offset, query.Limit)
		if err != nil {
			return nil, err
		}
		total, err = h.queryRepo.Count(ctx)
	}

	if err != nil {
		return nil, err
	}

	return &OpportunityListDTO{
		Items: ToOpportunityDTOs(items),
		Total: total,
	}, nil
}
