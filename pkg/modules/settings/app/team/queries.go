package team

// GetQuery 获取单个配置查询
type GetQuery struct {
	TeamID uint
	OrgID  uint // 所属组织 ID（用于继承查询）
	Key    string
}

// ListQuery 获取配置列表查询
type ListQuery struct {
	TeamID     uint
	OrgID      uint // 所属组织 ID（用于继承查询）
	CategoryID uint // 可选：按分类 ID 过滤
}
