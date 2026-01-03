package opportunity

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/opportunity"
)

// AdvanceHandler 推进商机阶段处理器。
type AdvanceHandler struct {
	cmdRepo   opportunity.CommandRepository
	queryRepo opportunity.QueryRepository
}

// NewAdvanceHandler 创建处理器实例。
func NewAdvanceHandler(cmdRepo opportunity.CommandRepository, queryRepo opportunity.QueryRepository) *AdvanceHandler {
	return &AdvanceHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 执行推进阶段命令。
func (h *AdvanceHandler) Handle(ctx context.Context, cmd AdvanceCommand) (*OpportunityDTO, error) {
	opp, err := h.queryRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	if err := opp.AdvanceTo(cmd.Stage); err != nil {
		return nil, err
	}

	if err := h.cmdRepo.Update(ctx, opp); err != nil {
		return nil, err
	}

	return ToOpportunityDTO(opp), nil
}
