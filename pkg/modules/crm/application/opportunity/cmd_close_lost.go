package opportunity

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/opportunity"
)

// CloseLostHandler 丢单处理器。
type CloseLostHandler struct {
	cmdRepo   opportunity.CommandRepository
	queryRepo opportunity.QueryRepository
}

// NewCloseLostHandler 创建处理器实例。
func NewCloseLostHandler(cmdRepo opportunity.CommandRepository, queryRepo opportunity.QueryRepository) *CloseLostHandler {
	return &CloseLostHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 执行丢单命令。
func (h *CloseLostHandler) Handle(ctx context.Context, cmd CloseLostCommand) (*OpportunityDTO, error) {
	opp, err := h.queryRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	if err := opp.CloseLost(); err != nil {
		return nil, err
	}

	if err := h.cmdRepo.Update(ctx, opp); err != nil {
		return nil, err
	}

	return ToOpportunityDTO(opp), nil
}
