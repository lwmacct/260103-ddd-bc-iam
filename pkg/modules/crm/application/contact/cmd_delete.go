package contact

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/contact"
)

// DeleteHandler 删除联系人处理器。
type DeleteHandler struct {
	cmdRepo   contact.CommandRepository
	queryRepo contact.QueryRepository
}

// NewDeleteHandler 创建 DeleteHandler。
func NewDeleteHandler(cmdRepo contact.CommandRepository, queryRepo contact.QueryRepository) *DeleteHandler {
	return &DeleteHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 处理删除联系人命令。
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	// 验证联系人存在
	_, err := h.queryRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	return h.cmdRepo.Delete(ctx, cmd.ID)
}
