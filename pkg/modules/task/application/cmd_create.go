package task

import (
	"context"

	taskdomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/domain"
)

// CreateHandler 创建任务处理器。
type CreateHandler struct {
	commandRepo taskdomain.CommandRepository
}

// NewCreateHandler 创建 CreateHandler 实例。
func NewCreateHandler(commandRepo taskdomain.CommandRepository) *CreateHandler {
	return &CreateHandler{
		commandRepo: commandRepo,
	}
}

// Handle 处理创建任务命令。
func (h *CreateHandler) Handle(ctx context.Context, cmd CreateTaskCommand) (*TaskDTO, error) {
	t := &taskdomain.Task{
		OrgID:       cmd.OrgID,
		TeamID:      cmd.TeamID,
		Title:       cmd.Title,
		Description: cmd.Description,
		Status:      taskdomain.StatusPending,
		AssigneeID:  cmd.AssigneeID,
		CreatedBy:   cmd.CreatedBy,
	}

	if err := h.commandRepo.Create(ctx, t); err != nil {
		return nil, err
	}

	return ToTaskDTO(t), nil
}
