package task

import (
	"context"

	taskdomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/domain"
)

// ListHandler 任务列表处理器。
type ListHandler struct {
	queryRepo taskdomain.QueryRepository
}

// NewListHandler 创建 ListHandler 实例。
func NewListHandler(queryRepo taskdomain.QueryRepository) *ListHandler {
	return &ListHandler{
		queryRepo: queryRepo,
	}
}

// ListResult 列表查询结果。
type ListResult struct {
	Items []*TaskDTO
	Total int64
}

// Handle 处理任务列表查询。
func (h *ListHandler) Handle(ctx context.Context, query ListTasksQuery) (*ListResult, error) {
	total, err := h.queryRepo.CountByTeam(ctx, query.OrgID, query.TeamID)
	if err != nil {
		return nil, err
	}

	tasks, err := h.queryRepo.ListByTeam(ctx, query.OrgID, query.TeamID, query.Offset, query.Limit)
	if err != nil {
		return nil, err
	}

	return &ListResult{
		Items: ToTaskDTOs(tasks),
		Total: total,
	}, nil
}
