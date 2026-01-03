package lead

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/lead"
)

// ContactHandler 转换到已联系状态处理器。
type ContactHandler struct {
	cmdRepo   lead.CommandRepository
	queryRepo lead.QueryRepository
}

// NewContactHandler 创建 ContactHandler。
func NewContactHandler(cmdRepo lead.CommandRepository, queryRepo lead.QueryRepository) *ContactHandler {
	return &ContactHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 处理转换到已联系状态命令。
func (h *ContactHandler) Handle(ctx context.Context, cmd ContactCommand) (*LeadDTO, error) {
	l, err := h.queryRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	if err := l.Contact(); err != nil {
		return nil, err
	}

	if err := h.cmdRepo.Update(ctx, l); err != nil {
		return nil, err
	}

	return ToLeadDTO(l), nil
}
