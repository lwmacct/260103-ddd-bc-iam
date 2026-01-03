package task

import (
	"context"

	taskdomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/domain"
)

// UpdateHandler 更新任务处理器。
type UpdateHandler struct {
	commandRepo taskdomain.CommandRepository
	queryRepo   taskdomain.QueryRepository
}

// NewUpdateHandler 创建 UpdateHandler 实例。
func NewUpdateHandler(
	commandRepo taskdomain.CommandRepository,
	queryRepo taskdomain.QueryRepository,
) *UpdateHandler {
	return &UpdateHandler{
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

// Handle 处理更新任务命令。
func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateTaskCommand) (*TaskDTO, error) {
	// 获取现有任务（验证归属）
	t, err := h.queryRepo.GetByIDAndTeam(ctx, cmd.ID, cmd.OrgID, cmd.TeamID)
	if err != nil {
		return nil, err
	}

	// 应用更新
	if cmd.Title != nil {
		t.Title = *cmd.Title
	}
	if cmd.Description != nil {
		t.Description = *cmd.Description
	}
	if cmd.Status != nil {
		newStatus := taskdomain.Status(*cmd.Status)
		if !t.CanTransitionTo(newStatus) {
			return nil, taskdomain.ErrInvalidStatusTransition
		}
		t.Status = newStatus
	}
	if cmd.AssigneeID != nil {
		t.AssigneeID = cmd.AssigneeID
	}

	if err := h.commandRepo.Update(ctx, t); err != nil {
		return nil, err
	}

	return ToTaskDTO(t), nil
}
