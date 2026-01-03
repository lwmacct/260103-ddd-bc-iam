package task

import (
	"context"

	taskdomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/domain"
)

// GetHandler 获取任务详情处理器。
type GetHandler struct {
	queryRepo taskdomain.QueryRepository
}

// NewGetHandler 创建 GetHandler 实例。
func NewGetHandler(queryRepo taskdomain.QueryRepository) *GetHandler {
	return &GetHandler{
		queryRepo: queryRepo,
	}
}

// Handle 处理获取任务详情查询。
func (h *GetHandler) Handle(ctx context.Context, query GetTaskQuery) (*TaskDTO, error) {
	t, err := h.queryRepo.GetByIDAndTeam(ctx, query.ID, query.OrgID, query.TeamID)
	if err != nil {
		return nil, err
	}
	return ToTaskDTO(t), nil
}
