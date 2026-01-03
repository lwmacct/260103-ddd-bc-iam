package company

import "context"

// CommandRepository 公司写仓储接口。
type CommandRepository interface {
	// Create 创建公司，创建后回写 ID。
	Create(ctx context.Context, company *Company) error

	// Update 更新公司信息。
	Update(ctx context.Context, company *Company) error

	// Delete 删除公司。
	Delete(ctx context.Context, id uint) error
}

// QueryRepository 公司读仓储接口。
type QueryRepository interface {
	// GetByID 根据 ID 查询公司。
	GetByID(ctx context.Context, id uint) (*Company, error)

	// GetByName 根据名称查询公司。
	GetByName(ctx context.Context, name string) (*Company, error)

	// ListByIndustry 根据行业查询公司列表。
	ListByIndustry(ctx context.Context, industry string, offset, limit int) ([]*Company, error)

	// CountByIndustry 统计指定行业的公司数量。
	CountByIndustry(ctx context.Context, industry string) (int64, error)

	// ListByOwner 根据负责人查询公司列表。
	ListByOwner(ctx context.Context, ownerID uint, offset, limit int) ([]*Company, error)

	// CountByOwner 统计指定负责人的公司数量。
	CountByOwner(ctx context.Context, ownerID uint) (int64, error)

	// List 分页查询公司列表。
	List(ctx context.Context, offset, limit int) ([]*Company, error)

	// Count 统计公司总数。
	Count(ctx context.Context) (int64, error)
}
