package task

// CreateTaskCommand 创建任务命令。
type CreateTaskCommand struct {
	OrgID       uint
	TeamID      uint
	Title       string
	Description string
	AssigneeID  *uint
	CreatedBy   uint
}

// UpdateTaskCommand 更新任务命令。
type UpdateTaskCommand struct {
	OrgID       uint
	TeamID      uint
	ID          uint
	Title       *string
	Description *string
	Status      *string
	AssigneeID  *uint
}

// DeleteTaskCommand 删除任务命令。
type DeleteTaskCommand struct {
	OrgID  uint
	TeamID uint
	ID     uint
}
