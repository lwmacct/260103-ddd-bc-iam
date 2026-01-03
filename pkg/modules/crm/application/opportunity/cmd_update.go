package opportunity

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/opportunity"
)

// UpdateHandler 更新商机处理器。
type UpdateHandler struct {
	cmdRepo   opportunity.CommandRepository
	queryRepo opportunity.QueryRepository
}

// NewUpdateHandler 创建处理器实例。
func NewUpdateHandler(cmdRepo opportunity.CommandRepository, queryRepo opportunity.QueryRepository) *UpdateHandler {
	return &UpdateHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 执行更新商机命令。
func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*OpportunityDTO, error) {
	opp, err := h.queryRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	// 已关闭的商机不能更新
	if opp.IsClosed() {
		return nil, opportunity.ErrAlreadyClosed
	}

	// 应用部分更新
	if cmd.Name != nil {
		opp.Name = *cmd.Name
	}
	if cmd.ContactID != nil {
		opp.ContactID = *cmd.ContactID
	}
	if cmd.CompanyID != nil {
		opp.CompanyID = cmd.CompanyID
	}
	if cmd.Amount != nil {
		opp.Amount = *cmd.Amount
	}
	if cmd.Probability != nil {
		opp.Probability = *cmd.Probability
	}
	if cmd.ExpectedClose != nil {
		opp.ExpectedClose = cmd.ExpectedClose
	}
	if cmd.Notes != nil {
		opp.Notes = *cmd.Notes
	}

	if err := h.cmdRepo.Update(ctx, opp); err != nil {
		return nil, err
	}

	return ToOpportunityDTO(opp), nil
}
