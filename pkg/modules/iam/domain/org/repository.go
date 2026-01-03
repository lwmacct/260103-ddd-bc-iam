package org

import "context"

// ============================================================================
// Organization Command Repository
// ============================================================================

// CommandRepository 组织命令仓储接口
type CommandRepository interface {
	// Create 创建组织
	Create(ctx context.Context, org *Org) error

	// Update 更新组织
	Update(ctx context.Context, org *Org) error

	// Delete 删除组织（软删除）
	Delete(ctx context.Context, id uint) error

	// UpdateStatus 更新组织状态
	UpdateStatus(ctx context.Context, id uint, status string) error
}

// ============================================================================
// Organization Query Repository
// ============================================================================

// QueryRepository 组织查询仓储接口
type QueryRepository interface {
	// GetByID 根据 ID 获取组织
	GetByID(ctx context.Context, id uint) (*Org, error)

	// GetByName 根据名称获取组织
	GetByName(ctx context.Context, name string) (*Org, error)

	// GetByIDWithTeams 根据 ID 获取组织（包含团队列表）
	GetByIDWithTeams(ctx context.Context, id uint) (*Org, error)

	// GetByIDWithMembers 根据 ID 获取组织（包含成员列表）
	GetByIDWithMembers(ctx context.Context, id uint) (*Org, error)

	// List 获取组织列表（分页）
	List(ctx context.Context, offset, limit int) ([]*Org, error)

	// Count 统计组织数量
	Count(ctx context.Context) (int64, error)

	// Exists 检查组织是否存在
	Exists(ctx context.Context, id uint) (bool, error)

	// ExistsByName 检查组织名称是否存在
	ExistsByName(ctx context.Context, name string) (bool, error)

	// Search 搜索组织（支持名称、显示名称模糊匹配）
	Search(ctx context.Context, keyword string, offset, limit int) ([]*Org, error)

	// CountBySearch 统计搜索结果数量
	CountBySearch(ctx context.Context, keyword string) (int64, error)
}
