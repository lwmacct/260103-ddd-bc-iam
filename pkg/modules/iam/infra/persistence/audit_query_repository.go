package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/audit"
	"gorm.io/gorm"
)

// auditQueryRepository 审计日志查询仓储的 GORM 实现
type auditQueryRepository struct {
	db *gorm.DB
}

// NewAuditQueryRepository 创建审计日志查询仓储实例
func NewAuditQueryRepository(db *gorm.DB) audit.QueryRepository {
	return &auditQueryRepository{db: db}
}

// FindByID finds an audit log by ID
func (r *auditQueryRepository) FindByID(ctx context.Context, id uint) (*audit.Audit, error) {
	var model AuditModel
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("audit log not found")
		}
		return nil, fmt.Errorf("failed to find audit log: %w", err)
	}
	return model.ToEntity(), nil
}

// List returns audit logs with filtering and pagination
func (r *auditQueryRepository) List(ctx context.Context, filter audit.FilterOptions) ([]audit.Audit, int64, error) {
	var models []AuditModel
	var total int64

	query := r.db.WithContext(ctx).Model(&AuditModel{})

	// Apply filters
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Resource != "" {
		query = query.Where("resource = ?", filter.Resource)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	// Apply pagination
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	// Order by created_at desc
	query = query.Order("created_at DESC")

	if err := query.Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list audit logs: %w", err)
	}

	logs := make([]audit.Audit, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			logs = append(logs, *entity)
		}
	}

	return logs, total, nil
}

// ListByUser returns audit logs for a specific user
func (r *auditQueryRepository) ListByUser(ctx context.Context, userID uint, page, limit int) ([]audit.Audit, int64, error) {
	return r.listByCondition(ctx, "user_id = ?", userID, page, limit,
		"failed to count audit logs by user", "failed to list audit logs by user")
}

// ListByResource returns audit logs for a specific resource
func (r *auditQueryRepository) ListByResource(ctx context.Context, resource string, page, limit int) ([]audit.Audit, int64, error) {
	return r.listByCondition(ctx, "resource = ?", resource, page, limit,
		"failed to count audit logs by resource", "failed to list audit logs by resource")
}

// ListByAction returns audit logs for a specific action
func (r *auditQueryRepository) ListByAction(ctx context.Context, action string, page, limit int) ([]audit.Audit, int64, error) {
	return r.listByCondition(ctx, "action = ?", action, page, limit,
		"failed to count audit logs by action", "failed to list audit logs by action")
}

// Count returns the total number of audit logs matching the filter
func (r *auditQueryRepository) Count(ctx context.Context, filter audit.FilterOptions) (int64, error) {
	var total int64

	query := r.db.WithContext(ctx).Model(&AuditModel{})

	// Apply filters
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Resource != "" {
		query = query.Where("resource = ?", filter.Resource)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	return total, nil
}

// Search searches audit logs by keyword
func (r *auditQueryRepository) Search(ctx context.Context, keyword string, page, limit int) ([]audit.Audit, int64, error) {
	var models []AuditModel
	var total int64

	query := r.db.WithContext(ctx).Model(&AuditModel{}).
		Where("resource LIKE ? OR action LIKE ? OR details LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search audit logs: %w", err)
	}

	logs := make([]audit.Audit, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			logs = append(logs, *entity)
		}
	}

	return logs, total, nil
}

// listByCondition 通用条件查询方法，减少重复代码
func (r *auditQueryRepository) listByCondition(ctx context.Context, condition string, value any, page, limit int, errMsgCount, errMsgList string) ([]audit.Audit, int64, error) {
	var models []AuditModel
	var total int64

	query := r.db.WithContext(ctx).Model(&AuditModel{}).
		Where(condition, value)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("%s: %w", errMsgCount, err)
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("%s: %w", errMsgList, err)
	}

	logs := make([]audit.Audit, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			logs = append(logs, *entity)
		}
	}

	return logs, total, nil
}
