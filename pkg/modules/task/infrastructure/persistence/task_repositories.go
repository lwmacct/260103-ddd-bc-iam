package persistence

import (
	taskdomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/domain"
	"gorm.io/gorm"
)

// TaskRepositories 聚合任务读写仓储。
type TaskRepositories struct {
	Command taskdomain.CommandRepository
	Query   taskdomain.QueryRepository
}

// NewTaskRepositories 创建任务仓储聚合实例。
func NewTaskRepositories(db *gorm.DB) TaskRepositories {
	return TaskRepositories{
		Command: NewTaskCommandRepository(db),
		Query:   NewTaskQueryRepository(db),
	}
}
