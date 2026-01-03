package audit

import "context"

// ============================================================================
// Command Repository
// ============================================================================

// CommandRepository 审计日志命令仓储接口（写操作）
type CommandRepository interface {
	// Create creates a new audit log entry
	Create(ctx context.Context, log *Audit) error

	// Delete deletes an audit log (soft delete, for data retention policy)
	Delete(ctx context.Context, id uint) error

	// DeleteOlderThan deletes audit logs older than the specified date
	DeleteOlderThan(ctx context.Context, days int) error

	// BatchCreate creates multiple audit log entries
	BatchCreate(ctx context.Context, logs []*Audit) error
}

// ============================================================================
// Query Repository
// ============================================================================

// QueryRepository 审计日志查询仓储接口（读操作）
type QueryRepository interface {
	// FindByID finds an audit log by ID
	FindByID(ctx context.Context, id uint) (*Audit, error)

	// List returns audit logs with filtering and pagination
	List(ctx context.Context, filter FilterOptions) ([]Audit, int64, error)

	// ListByUser returns audit logs for a specific user
	ListByUser(ctx context.Context, userID uint, page, limit int) ([]Audit, int64, error)

	// ListByResource returns audit logs for a specific resource
	ListByResource(ctx context.Context, resource string, page, limit int) ([]Audit, int64, error)

	// ListByAction returns audit logs for a specific action
	ListByAction(ctx context.Context, action string, page, limit int) ([]Audit, int64, error)

	// Count returns the total number of audit logs matching the filter
	Count(ctx context.Context, filter FilterOptions) (int64, error)

	// Search searches audit logs by keyword
	Search(ctx context.Context, keyword string, page, limit int) ([]Audit, int64, error)
}
