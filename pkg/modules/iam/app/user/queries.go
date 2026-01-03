package user

// GetQuery 获取用户查询
type GetQuery struct {
	UserID    uint
	WithRoles bool // 是否包含角色信息
}

// ListQuery 获取用户列表查询
type ListQuery struct {
	Page   int
	Limit  int
	Search string // 搜索关键词（用户名或邮箱）
}

// GetOffset 计算数据库查询偏移量
func (q ListQuery) GetOffset() int {
	page := max(q.Page, 1)
	return (page - 1) * q.Limit
}
