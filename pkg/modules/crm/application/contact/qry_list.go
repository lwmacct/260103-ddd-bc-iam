package contact

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/contact"
)

// ListHandler 联系人列表处理器。
type ListHandler struct {
	queryRepo contact.QueryRepository
}

// NewListHandler 创建 ListHandler。
func NewListHandler(queryRepo contact.QueryRepository) *ListHandler {
	return &ListHandler{queryRepo: queryRepo}
}

// Handle 处理联系人列表查询。
func (h *ListHandler) Handle(ctx context.Context, query ListQuery) (*ListResultDTO, error) {
	var (
		items []*contact.Contact
		total int64
		err   error
	)

	switch {
	case query.CompanyID != nil:
		// 按公司筛选
		items, err = h.queryRepo.ListByCompany(ctx, *query.CompanyID, query.Offset, query.Limit)
		if err != nil {
			return nil, err
		}
		total, err = h.queryRepo.CountByCompany(ctx, *query.CompanyID)
		if err != nil {
			return nil, err
		}
	case query.OwnerID != nil:
		// 按负责人筛选
		items, err = h.queryRepo.ListByOwner(ctx, *query.OwnerID, query.Offset, query.Limit)
		if err != nil {
			return nil, err
		}
		total, err = h.queryRepo.CountByOwner(ctx, *query.OwnerID)
		if err != nil {
			return nil, err
		}
	default:
		// 查询所有
		items, err = h.queryRepo.List(ctx, query.Offset, query.Limit)
		if err != nil {
			return nil, err
		}
		total, err = h.queryRepo.Count(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &ListResultDTO{
		Items: ToContactDTOs(items),
		Total: total,
	}, nil
}
