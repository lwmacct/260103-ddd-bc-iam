package lead

import "context"

// CommandRepository 线索写仓储接口。
type CommandRepository interface {
	// Create 创建线索，创建后回写 ID。
	Create(ctx context.Context, lead *Lead) error

	// Update 更新线索信息。
	Update(ctx context.Context, lead *Lead) error

	// Delete 删除线索。
	Delete(ctx context.Context, id uint) error
}

// QueryRepository 线索读仓储接口。
type QueryRepository interface {
	// GetByID 根据 ID 查询线索。
	GetByID(ctx context.Context, id uint) (*Lead, error)

	// ListByStatus 根据状态查询线索列表。
	ListByStatus(ctx context.Context, status Status, offset, limit int) ([]*Lead, error)

	// CountByStatus 统计指定状态的线索数量。
	CountByStatus(ctx context.Context, status Status) (int64, error)

	// ListByOwner 根据负责人查询线索列表。
	ListByOwner(ctx context.Context, ownerID uint, offset, limit int) ([]*Lead, error)

	// CountByOwner 统计指定负责人的线索数量。
	CountByOwner(ctx context.Context, ownerID uint) (int64, error)

	// List 分页查询线索列表。
	List(ctx context.Context, offset, limit int) ([]*Lead, error)

	// Count 统计线索总数。
	Count(ctx context.Context) (int64, error)
}
