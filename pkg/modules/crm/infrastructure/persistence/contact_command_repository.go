package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/contact"
)

// contactCommandRepository 联系人写仓储实现。
type contactCommandRepository struct {
	db *gorm.DB
}

// NewContactCommandRepository 创建联系人写仓储。
func NewContactCommandRepository(db *gorm.DB) contact.CommandRepository {
	return &contactCommandRepository{db: db}
}

// Create 创建联系人。
func (r *contactCommandRepository) Create(ctx context.Context, entity *contact.Contact) error {
	model := newContactModelFromEntity(entity)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	// 回写生成的字段
	entity.ID = model.ID
	entity.CreatedAt = model.CreatedAt
	entity.UpdatedAt = model.UpdatedAt
	return nil
}

// Update 更新联系人。
func (r *contactCommandRepository) Update(ctx context.Context, entity *contact.Contact) error {
	model := newContactModelFromEntity(entity)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	entity.UpdatedAt = model.UpdatedAt
	return nil
}

// Delete 删除联系人。
func (r *contactCommandRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&ContactModel{}, id).Error
}
