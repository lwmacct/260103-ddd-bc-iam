package lead

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/lead"
)

// DeleteHandler 删除线索处理器。
type DeleteHandler struct {
	cmdRepo   lead.CommandRepository
	queryRepo lead.QueryRepository
}

// NewDeleteHandler 创建 DeleteHandler。
func NewDeleteHandler(cmdRepo lead.CommandRepository, queryRepo lead.QueryRepository) *DeleteHandler {
	return &DeleteHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 处理删除线索命令。
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	// 验证线索存在
	_, err := h.queryRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	return h.cmdRepo.Delete(ctx, cmd.ID)
}
