package lead

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/lead"
)

// LoseHandler 标记为丢失处理器。
type LoseHandler struct {
	cmdRepo   lead.CommandRepository
	queryRepo lead.QueryRepository
}

// NewLoseHandler 创建 LoseHandler。
func NewLoseHandler(cmdRepo lead.CommandRepository, queryRepo lead.QueryRepository) *LoseHandler {
	return &LoseHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 处理标记为丢失命令。
func (h *LoseHandler) Handle(ctx context.Context, cmd LoseCommand) (*LeadDTO, error) {
	l, err := h.queryRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	if err := l.Lose(); err != nil {
		return nil, err
	}

	if err := h.cmdRepo.Update(ctx, l); err != nil {
		return nil, err
	}

	return ToLeadDTO(l), nil
}
