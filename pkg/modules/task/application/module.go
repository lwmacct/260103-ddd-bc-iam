package task

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/infrastructure/persistence"
)

// UseCaseModule 提供 task 领域的用例处理器。
var UseCaseModule = fx.Module("task.usecase",
	fx.Provide(
		NewCreateHandler,
		NewUpdateHandler,
		NewDeleteHandler,
		NewGetHandler,
		NewListHandler,
	),
)

// TaskUseCases 聚合任务用例处理器。
type TaskUseCases struct {
	Create *CreateHandler
	Update *UpdateHandler
	Delete *DeleteHandler
	Get    *GetHandler
	List   *ListHandler
}

// NewTaskUseCases 创建任务用例聚合实例。
func NewTaskUseCases(repos persistence.TaskRepositories) *TaskUseCases {
	return &TaskUseCases{
		Create: NewCreateHandler(repos.Command),
		Update: NewUpdateHandler(repos.Command, repos.Query),
		Delete: NewDeleteHandler(repos.Command, repos.Query),
		Get:    NewGetHandler(repos.Query),
		List:   NewListHandler(repos.Query),
	}
}
