package persistence

import (
	"time"

	taskdomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/domain"
	"gorm.io/gorm"
)

// TaskModel 定义任务的 GORM 持久化模型。
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type TaskModel struct {
	ID          uint   `gorm:"primaryKey"`
	OrgID       uint   `gorm:"index:idx_tasks_org_team;not null"`
	TeamID      uint   `gorm:"index:idx_tasks_org_team;not null"`
	Title       string `gorm:"size:200;not null"`
	Description string `gorm:"type:text"`
	Status      string `gorm:"size:20;default:'pending';not null;index"`
	AssigneeID  *uint  `gorm:"index"`
	CreatedBy   uint   `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 指定任务表名。
func (TaskModel) TableName() string {
	return "tasks"
}

func newTaskModelFromEntity(entity *taskdomain.Task) *TaskModel {
	if entity == nil {
		return nil
	}

	return &TaskModel{
		ID:          entity.ID,
		OrgID:       entity.OrgID,
		TeamID:      entity.TeamID,
		Title:       entity.Title,
		Description: entity.Description,
		Status:      string(entity.Status),
		AssigneeID:  entity.AssigneeID,
		CreatedBy:   entity.CreatedBy,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

// ToEntity 将 GORM Model 转换为 Domain Entity。
func (m *TaskModel) ToEntity() *taskdomain.Task {
	if m == nil {
		return nil
	}

	return &taskdomain.Task{
		ID:          m.ID,
		OrgID:       m.OrgID,
		TeamID:      m.TeamID,
		Title:       m.Title,
		Description: m.Description,
		Status:      taskdomain.Status(m.Status),
		AssigneeID:  m.AssigneeID,
		CreatedBy:   m.CreatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func mapTaskModelsToEntities(models []TaskModel) []*taskdomain.Task {
	if len(models) == 0 {
		return nil
	}

	tasks := make([]*taskdomain.Task, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			tasks = append(tasks, entity)
		}
	}
	return tasks
}
