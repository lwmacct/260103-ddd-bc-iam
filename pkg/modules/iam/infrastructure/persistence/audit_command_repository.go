package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/audit"
	"gorm.io/gorm"
)

// auditCommandRepository 审计日志命令仓储的 GORM 实现
// 嵌入 GenericCommandRepository 以复用 Create/Delete 操作
type auditCommandRepository struct {
	*GenericCommandRepository[audit.Audit, *AuditModel]
}

// NewAuditCommandRepository 创建审计日志命令仓储实例
func NewAuditCommandRepository(db *gorm.DB) audit.CommandRepository {
	return &auditCommandRepository{
		GenericCommandRepository: NewGenericCommandRepository(
			db, newAuditModelFromEntity,
		),
	}
}

// Create、Delete 方法由 GenericCommandRepository 提供

// DeleteOlderThan deletes audit logs older than the specified date
func (r *auditCommandRepository) DeleteOlderThan(ctx context.Context, days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	if err := r.DB().WithContext(ctx).
		Where("created_at < ?", cutoffDate).
		Delete(&AuditModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete old audit logs: %w", err)
	}
	return nil
}

// BatchCreate creates multiple audit log entries
func (r *auditCommandRepository) BatchCreate(ctx context.Context, logs []*audit.Audit) error {
	if len(logs) == 0 {
		return nil
	}
	models := make([]*AuditModel, 0, len(logs))
	for _, log := range logs {
		if model := newAuditModelFromEntity(log); model != nil {
			models = append(models, model)
		}
	}
	if err := r.DB().WithContext(ctx).Create(models).Error; err != nil {
		return fmt.Errorf("failed to batch create audit logs: %w", err)
	}

	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			*logs[i] = *entity
		}
	}

	return nil
}
