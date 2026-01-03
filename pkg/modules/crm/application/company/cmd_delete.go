package company

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/company"
)

// DeleteHandler 删除公司处理器。
type DeleteHandler struct {
	cmdRepo   company.CommandRepository
	queryRepo company.QueryRepository
}

// NewDeleteHandler 创建 DeleteHandler。
func NewDeleteHandler(cmdRepo company.CommandRepository, queryRepo company.QueryRepository) *DeleteHandler {
	return &DeleteHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 处理删除公司命令。
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	// 验证公司存在
	_, err := h.queryRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	return h.cmdRepo.Delete(ctx, cmd.ID)
}
