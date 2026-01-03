package lead

import "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/lead"

// GetQuery 获取线索查询。
type GetQuery struct {
	ID uint
}

// ListQuery 线索列表查询。
type ListQuery struct {
	Status  *lead.Status
	OwnerID *uint
	Offset  int
	Limit   int
}
