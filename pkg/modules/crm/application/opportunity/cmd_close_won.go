package opportunity

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/opportunity"
)

// CloseWonHandler 成交处理器。
type CloseWonHandler struct {
	cmdRepo   opportunity.CommandRepository
	queryRepo opportunity.QueryRepository
}

// NewCloseWonHandler 创建处理器实例。
func NewCloseWonHandler(cmdRepo opportunity.CommandRepository, queryRepo opportunity.QueryRepository) *CloseWonHandler {
	return &CloseWonHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 执行成交命令。
func (h *CloseWonHandler) Handle(ctx context.Context, cmd CloseWonCommand) (*OpportunityDTO, error) {
	opp, err := h.queryRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	if err := opp.CloseWon(); err != nil {
		return nil, err
	}

	if err := h.cmdRepo.Update(ctx, opp); err != nil {
		return nil, err
	}

	return ToOpportunityDTO(opp), nil
}
