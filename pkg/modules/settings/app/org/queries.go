package org

// GetQuery 获取单个配置查询
type GetQuery struct {
	OrgID uint
	Key   string
}

// ListQuery 获取配置列表查询
type ListQuery struct {
	OrgID      uint   `form:"-"`                      // 从上下文获取，不绑定
	CategoryID uint   `form:"category_id" binding:"omitempty"` // 可选：按分类 ID 过滤
}
