package team

// GetQuery 获取单个配置查询
type GetQuery struct {
	TeamID uint
	OrgID  uint // 所属组织 ID（用于继承查询）
	Key    string
}

// ListQuery 获取配置列表查询
type ListQuery struct {
	TeamID   uint   `form:"-" swaggerignore:"true"`       // 从上下文获取，不绑定
	OrgID    uint   `form:"-" swaggerignore:"true"`       // 从上下文获取，不绑定
	Category string `form:"category" binding:"omitempty"` // 可选：按分类 Key 过滤（如 "general"）
}
