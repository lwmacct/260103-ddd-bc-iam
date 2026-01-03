package contact

import "context"

// CommandRepository 联系人写仓储接口。
type CommandRepository interface {
	// Create 创建联系人，成功后回写 ID 和时间戳。
	Create(ctx context.Context, contact *Contact) error
	// Update 更新联系人。
	Update(ctx context.Context, contact *Contact) error
	// Delete 删除联系人。
	Delete(ctx context.Context, id uint) error
}

// QueryRepository 联系人读仓储接口。
type QueryRepository interface {
	// GetByID 根据 ID 获取联系人。
	GetByID(ctx context.Context, id uint) (*Contact, error)
	// GetByEmail 根据邮箱获取联系人。
	GetByEmail(ctx context.Context, email string) (*Contact, error)
	// ListByCompany 获取公司下的联系人列表。
	ListByCompany(ctx context.Context, companyID uint, offset, limit int) ([]*Contact, error)
	// CountByCompany 统计公司下的联系人数量。
	CountByCompany(ctx context.Context, companyID uint) (int64, error)
	// ListByOwner 获取负责人的联系人列表。
	ListByOwner(ctx context.Context, ownerID uint, offset, limit int) ([]*Contact, error)
	// CountByOwner 统计负责人的联系人数量。
	CountByOwner(ctx context.Context, ownerID uint) (int64, error)
	// List 获取联系人列表（管理员）。
	List(ctx context.Context, offset, limit int) ([]*Contact, error)
	// Count 统计联系人总数。
	Count(ctx context.Context) (int64, error)
}
