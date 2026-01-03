package lead

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/lead"
)

// ConvertHandler 转化为商机处理器。
type ConvertHandler struct {
	cmdRepo   lead.CommandRepository
	queryRepo lead.QueryRepository
}

// NewConvertHandler 创建 ConvertHandler。
func NewConvertHandler(cmdRepo lead.CommandRepository, queryRepo lead.QueryRepository) *ConvertHandler {
	return &ConvertHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 处理转化为商机命令。
func (h *ConvertHandler) Handle(ctx context.Context, cmd ConvertCommand) (*LeadDTO, error) {
	l, err := h.queryRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	if err := l.Convert(); err != nil {
		return nil, err
	}

	if err := h.cmdRepo.Update(ctx, l); err != nil {
		return nil, err
	}

	return ToLeadDTO(l), nil
}
