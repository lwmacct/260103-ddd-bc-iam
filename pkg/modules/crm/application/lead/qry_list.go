package lead

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/lead"
)

// ListHandler 线索列表处理器。
type ListHandler struct {
	queryRepo lead.QueryRepository
}

// NewListHandler 创建 ListHandler。
func NewListHandler(queryRepo lead.QueryRepository) *ListHandler {
	return &ListHandler{queryRepo: queryRepo}
}

// Handle 处理线索列表查询。
func (h *ListHandler) Handle(ctx context.Context, query ListQuery) (*ListResultDTO, error) {
	var leads []*lead.Lead
	var total int64
	var err error

	switch {
	case query.Status != nil:
		leads, err = h.queryRepo.ListByStatus(ctx, *query.Status, query.Offset, query.Limit)
		if err != nil {
			return nil, err
		}
		total, err = h.queryRepo.CountByStatus(ctx, *query.Status)
	case query.OwnerID != nil:
		leads, err = h.queryRepo.ListByOwner(ctx, *query.OwnerID, query.Offset, query.Limit)
		if err != nil {
			return nil, err
		}
		total, err = h.queryRepo.CountByOwner(ctx, *query.OwnerID)
	default:
		leads, err = h.queryRepo.List(ctx, query.Offset, query.Limit)
		if err != nil {
			return nil, err
		}
		total, err = h.queryRepo.Count(ctx)
	}

	if err != nil {
		return nil, err
	}

	return &ListResultDTO{
		Items: ToLeadDTOs(leads),
		Total: total,
	}, nil
}
