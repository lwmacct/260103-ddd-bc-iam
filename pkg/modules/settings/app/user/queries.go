package user

// GetQuery 获取单个配置查询
type GetQuery struct {
	UserID uint
	Key    string
}

// ListQuery 获取配置列表查询
type ListQuery struct {
	UserID   uint   `form:"-" swaggerignore:"true"`       // 从上下文获取，不绑定
	Category string `form:"category" binding:"omitempty"` // 可选：按分类 Key 过滤（如 "general"）
}

// ListCategoriesQuery 获取分类列表查询
type ListCategoriesQuery struct {
	UserID uint // 预留：后续可用于过滤用户可见分类
}
