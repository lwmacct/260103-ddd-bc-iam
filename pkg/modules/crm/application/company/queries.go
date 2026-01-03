package company

// GetQuery 获取公司查询。
type GetQuery struct {
	ID uint
}

// ListQuery 公司列表查询。
type ListQuery struct {
	Industry *string
	OwnerID  *uint
	Offset   int
	Limit    int
}
