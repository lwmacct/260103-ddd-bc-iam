package userset

// GetQuery 获取单个配置查询
type GetQuery struct {
	UserID uint
	Key    string
}

// ListQuery 获取配置列表查询
type ListQuery struct {
	UserID     uint
	CategoryID uint // 可选：按分类 ID 过滤
}

// ListCategoriesQuery 获取分类列表查询
type ListCategoriesQuery struct {
	UserID uint // 预留：后续可用于过滤用户可见分类
}
