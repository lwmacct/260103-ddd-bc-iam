package task

// GetTaskQuery 获取任务详情查询。
type GetTaskQuery struct {
	OrgID  uint
	TeamID uint
	ID     uint
}

// ListTasksQuery 任务列表查询。
type ListTasksQuery struct {
	OrgID  uint
	TeamID uint
	Offset int
	Limit  int
}
