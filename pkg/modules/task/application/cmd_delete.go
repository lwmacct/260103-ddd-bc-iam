package task

import (
	"context"

	taskdomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/domain"
)

// DeleteHandler 删除任务处理器。
type DeleteHandler struct {
	commandRepo taskdomain.CommandRepository
	queryRepo   taskdomain.QueryRepository
}

// NewDeleteHandler 创建 DeleteHandler 实例。
func NewDeleteHandler(
	commandRepo taskdomain.CommandRepository,
	queryRepo taskdomain.QueryRepository,
) *DeleteHandler {
	return &DeleteHandler{
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

// Handle 处理删除任务命令。
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteTaskCommand) error {
	// 验证任务存在且归属正确
	_, err := h.queryRepo.GetByIDAndTeam(ctx, cmd.ID, cmd.OrgID, cmd.TeamID)
	if err != nil {
		return err
	}

	return h.commandRepo.Delete(ctx, cmd.ID)
}
