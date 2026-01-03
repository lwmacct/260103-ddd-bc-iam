package opportunity

import "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/opportunity"

// GetQuery 获取商机查询。
type GetQuery struct {
	ID uint
}

// ListQuery 商机列表查询。
type ListQuery struct {
	Stage   *opportunity.Stage
	OwnerID *uint
	Offset  int
	Limit   int
}
