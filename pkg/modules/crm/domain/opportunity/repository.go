package opportunity

import (
	"context"
)

// CommandRepository 商机写仓储接口。
type CommandRepository interface {
	// Create 创建商机。
	Create(ctx context.Context, opp *Opportunity) error
	// Update 更新商机。
	Update(ctx context.Context, opp *Opportunity) error
	// Delete 删除商机。
	Delete(ctx context.Context, id uint) error
}

// QueryRepository 商机读仓储接口。
type QueryRepository interface {
	// GetByID 根据 ID 查询商机。
	GetByID(ctx context.Context, id uint) (*Opportunity, error)
	// ListByStage 根据阶段查询商机列表。
	ListByStage(ctx context.Context, stage Stage, offset, limit int) ([]*Opportunity, error)
	// CountByStage 统计指定阶段的商机数量。
	CountByStage(ctx context.Context, stage Stage) (int64, error)
	// ListByOwner 根据负责人查询商机列表。
	ListByOwner(ctx context.Context, ownerID uint, offset, limit int) ([]*Opportunity, error)
	// CountByOwner 统计指定负责人的商机数量。
	CountByOwner(ctx context.Context, ownerID uint) (int64, error)
	// List 分页查询商机列表。
	List(ctx context.Context, offset, limit int) ([]*Opportunity, error)
	// Count 统计商机总数。
	Count(ctx context.Context) (int64, error)
}
