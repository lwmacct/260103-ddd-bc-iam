package persistence

import (
	"context"
	"errors"

	taskdomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/domain"
	"gorm.io/gorm"
)

type taskQueryRepository struct {
	db *gorm.DB
}

// NewTaskQueryRepository 创建任务读仓储实例。
func NewTaskQueryRepository(db *gorm.DB) taskdomain.QueryRepository {
	return &taskQueryRepository{db: db}
}

func (r *taskQueryRepository) GetByID(ctx context.Context, id uint) (*taskdomain.Task, error) {
	var model TaskModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, taskdomain.ErrTaskNotFound
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *taskQueryRepository) GetByIDAndTeam(ctx context.Context, id, orgID, teamID uint) (*taskdomain.Task, error) {
	var model TaskModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND org_id = ? AND team_id = ?", id, orgID, teamID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, taskdomain.ErrTaskNotFound
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *taskQueryRepository) ListByTeam(ctx context.Context, orgID, teamID uint, offset, limit int) ([]*taskdomain.Task, error) {
	var models []TaskModel
	err := r.db.WithContext(ctx).
		Where("org_id = ? AND team_id = ?", orgID, teamID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return mapTaskModelsToEntities(models), nil
}

func (r *taskQueryRepository) CountByTeam(ctx context.Context, orgID, teamID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&TaskModel{}).
		Where("org_id = ? AND team_id = ?", orgID, teamID).
		Count(&count).Error
	return count, err
}
