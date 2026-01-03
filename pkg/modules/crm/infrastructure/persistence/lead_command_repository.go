package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/lead"
)

// leadCommandRepository 线索写仓储实现。
type leadCommandRepository struct {
	db *gorm.DB
}

// NewLeadCommandRepository 创建线索写仓储。
func NewLeadCommandRepository(db *gorm.DB) lead.CommandRepository {
	return &leadCommandRepository{db: db}
}

// Create 创建线索。
func (r *leadCommandRepository) Create(ctx context.Context, l *lead.Lead) error {
	model := toLeadModel(l)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	l.ID = model.ID
	l.CreatedAt = model.CreatedAt
	l.UpdatedAt = model.UpdatedAt
	return nil
}

// Update 更新线索。
func (r *leadCommandRepository) Update(ctx context.Context, l *lead.Lead) error {
	model := toLeadModel(l)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	l.UpdatedAt = model.UpdatedAt
	return nil
}

// Delete 删除线索。
func (r *leadCommandRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&LeadModel{}, id).Error
}
