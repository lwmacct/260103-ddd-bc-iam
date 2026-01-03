package opportunity

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/opportunity"
)

// DeleteHandler 删除商机处理器。
type DeleteHandler struct {
	cmdRepo   opportunity.CommandRepository
	queryRepo opportunity.QueryRepository
}

// NewDeleteHandler 创建处理器实例。
func NewDeleteHandler(cmdRepo opportunity.CommandRepository, queryRepo opportunity.QueryRepository) *DeleteHandler {
	return &DeleteHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 执行删除商机命令。
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	// 验证商机存在
	if _, err := h.queryRepo.GetByID(ctx, cmd.ID); err != nil {
		return err
	}

	return h.cmdRepo.Delete(ctx, cmd.ID)
}
