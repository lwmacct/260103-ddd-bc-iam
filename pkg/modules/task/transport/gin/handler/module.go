package handler

import (
	"go.uber.org/fx"

	taskapplication "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/application"
)

// HandlersResult 批量返回 Handler。
type HandlersResult struct {
	fx.Out

	Task *TaskHandler
}

// NewAllHandlers 创建所有 Handler 实例。
func NewAllHandlers(useCases *taskapplication.TaskUseCases) HandlersResult {
	return HandlersResult{
		Task: NewTaskHandler(
			useCases.Create,
			useCases.Update,
			useCases.Delete,
			useCases.Get,
			useCases.List,
		),
	}
}
