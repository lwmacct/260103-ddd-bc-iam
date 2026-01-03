package task

import "context"

// CommandRepository 任务写仓储接口。
type CommandRepository interface {
	// Create 创建任务，成功后回写 ID 到实体。
	Create(ctx context.Context, task *Task) error
	// Update 更新任务。
	Update(ctx context.Context, task *Task) error
	// Delete 删除任务。
	Delete(ctx context.Context, id uint) error
}

// QueryRepository 任务读仓储接口。
type QueryRepository interface {
	// GetByID 根据 ID 获取任务。
	GetByID(ctx context.Context, id uint) (*Task, error)
	// GetByIDAndTeam 根据 ID 获取任务（验证归属于指定组织和团队）。
	GetByIDAndTeam(ctx context.Context, id, orgID, teamID uint) (*Task, error)
	// ListByTeam 分页获取团队任务列表。
	ListByTeam(ctx context.Context, orgID, teamID uint, offset, limit int) ([]*Task, error)
	// CountByTeam 获取团队任务总数。
	CountByTeam(ctx context.Context, orgID, teamID uint) (int64, error)
}
