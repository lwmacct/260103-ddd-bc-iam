package contact

// GetQuery 获取联系人详情查询。
type GetQuery struct {
	ID uint
}

// ListQuery 联系人列表查询。
type ListQuery struct {
	CompanyID *uint // 按公司筛选
	OwnerID   *uint // 按负责人筛选
	Offset    int
	Limit     int
}
