package org

// GetQuery 获取单个配置查询
type GetQuery struct {
	OrgID uint
	Key   string
}

// ListQuery 获取配置列表查询
type ListQuery struct {
	OrgID    uint   `form:"-" swaggerignore:"true"`       // 从上下文获取，不绑定
	Category string `form:"category" binding:"omitempty"` // 可选：按分类 Key 过滤（如 "general"）
}
