package org

// GetQuery 获取单个配置查询
type GetQuery struct {
	OrgID uint
	Key   string
}

// ListQuery 获取配置列表查询
type ListQuery struct {
	OrgID      uint
	CategoryID uint // 可选：按分类 ID 过滤
}
