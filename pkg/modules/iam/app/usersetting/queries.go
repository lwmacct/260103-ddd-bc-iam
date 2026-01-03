package usersetting

// GetQuery 获取单个用户设置查询
type GetQuery struct {
	UserID uint
	Key    string
}

// ListQuery 获取用户设置列表查询
type ListQuery struct {
	UserID   uint
	Category string // 可选：按分类过滤
}
