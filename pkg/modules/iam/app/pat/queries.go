package pat

// GetQuery 获取 Token 查询
type GetQuery struct {
	UserID  uint
	TokenID uint
}

// ListQuery 获取 Token 列表查询
type ListQuery struct {
	UserID uint
}
