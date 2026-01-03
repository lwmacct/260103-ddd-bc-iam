package persistence

import (
	"context"

	taskdomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/domain"
	"gorm.io/gorm"
)

type taskCommandRepository struct {
	db *gorm.DB
}

// NewTaskCommandRepository 创建任务写仓储实例。
func NewTaskCommandRepository(db *gorm.DB) taskdomain.CommandRepository {
	return &taskCommandRepository{db: db}
}

func (r *taskCommandRepository) Create(ctx context.Context, t *taskdomain.Task) error {
	model := newTaskModelFromEntity(t)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	t.ID = model.ID
	t.CreatedAt = model.CreatedAt
	t.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *taskCommandRepository) Update(ctx context.Context, t *taskdomain.Task) error {
	model := newTaskModelFromEntity(t)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *taskCommandRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&TaskModel{}, id).Error
}
