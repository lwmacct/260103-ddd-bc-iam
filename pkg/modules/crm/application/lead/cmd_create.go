package lead

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/lead"
)

// CreateHandler 创建线索处理器。
type CreateHandler struct {
	cmdRepo   lead.CommandRepository
	queryRepo lead.QueryRepository
}

// NewCreateHandler 创建 CreateHandler。
func NewCreateHandler(cmdRepo lead.CommandRepository, queryRepo lead.QueryRepository) *CreateHandler {
	return &CreateHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 处理创建线索命令。
func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*LeadDTO, error) {
	l := &lead.Lead{
		Title:       cmd.Title,
		ContactID:   cmd.ContactID,
		CompanyName: cmd.CompanyName,
		Source:      cmd.Source,
		Status:      lead.StatusNew,
		Score:       cmd.Score,
		OwnerID:     cmd.OwnerID,
		Notes:       cmd.Notes,
	}

	if err := h.cmdRepo.Create(ctx, l); err != nil {
		return nil, err
	}

	return ToLeadDTO(l), nil
}
