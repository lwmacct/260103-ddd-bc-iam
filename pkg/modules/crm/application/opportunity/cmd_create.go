package opportunity

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/opportunity"
)

// CreateHandler 创建商机处理器。
type CreateHandler struct {
	cmdRepo opportunity.CommandRepository
}

// NewCreateHandler 创建处理器实例。
func NewCreateHandler(cmdRepo opportunity.CommandRepository) *CreateHandler {
	return &CreateHandler{cmdRepo: cmdRepo}
}

// Handle 执行创建商机命令。
func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*OpportunityDTO, error) {
	opp := &opportunity.Opportunity{
		Name:          cmd.Name,
		ContactID:     cmd.ContactID,
		CompanyID:     cmd.CompanyID,
		LeadID:        cmd.LeadID,
		Stage:         opportunity.StageProspecting, // 新建商机默认为初步接触阶段
		Amount:        cmd.Amount,
		Probability:   cmd.Probability,
		ExpectedClose: cmd.ExpectedClose,
		OwnerID:       cmd.OwnerID,
		Notes:         cmd.Notes,
	}

	if err := h.cmdRepo.Create(ctx, opp); err != nil {
		return nil, err
	}

	return ToOpportunityDTO(opp), nil
}
