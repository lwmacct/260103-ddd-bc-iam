package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/audit"
	"gorm.io/gorm"
)

// auditCommandRepository 审计日志命令仓储的 GORM 实现
type auditCommandRepository struct {
	db *gorm.DB
}

// NewAuditCommandRepository 创建审计日志命令仓储实例
func NewAuditCommandRepository(db *gorm.DB) audit.CommandRepository {
	return &auditCommandRepository{db: db}
}

// Create 创建审计日志
func (r *auditCommandRepository) Create(ctx context.Context, log *audit.Audit) error {
	model := newAuditModelFromEntity(log)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	log.ID = model.ID
	return nil
}

// Update 更新审计日志
func (r *auditCommandRepository) Update(ctx context.Context, log *audit.Audit) error {
	model := newAuditModelFromEntity(log)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete 删除审计日志
func (r *auditCommandRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&AuditModel{}, id).Error
}

// DeleteOlderThan deletes audit logs older than the specified date
func (r *auditCommandRepository) DeleteOlderThan(ctx context.Context, days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	if err := r.db.WithContext(ctx).
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
	if err := r.db.WithContext(ctx).Create(models).Error; err != nil {
		return fmt.Errorf("failed to batch create audit logs: %w", err)
	}

	for i := range models {
		if models[i].ToEntity() != nil {
			*logs[i] = *models[i].ToEntity()
		}
	}

	return nil
}
