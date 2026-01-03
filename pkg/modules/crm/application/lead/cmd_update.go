package lead

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/lead"
)

// UpdateHandler 更新线索处理器。
type UpdateHandler struct {
	cmdRepo   lead.CommandRepository
	queryRepo lead.QueryRepository
}

// NewUpdateHandler 创建 UpdateHandler。
func NewUpdateHandler(cmdRepo lead.CommandRepository, queryRepo lead.QueryRepository) *UpdateHandler {
	return &UpdateHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 处理更新线索命令。
func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*LeadDTO, error) {
	l, err := h.queryRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	// 已关闭的线索不能修改
	if l.IsClosed() {
		return nil, lead.ErrAlreadyClosed
	}

	if cmd.Title != nil {
		l.Title = *cmd.Title
	}
	if cmd.ContactID != nil {
		l.ContactID = cmd.ContactID
	}
	if cmd.CompanyName != nil {
		l.CompanyName = *cmd.CompanyName
	}
	if cmd.Source != nil {
		l.Source = *cmd.Source
	}
	if cmd.Score != nil {
		l.Score = *cmd.Score
	}
	if cmd.Notes != nil {
		l.Notes = *cmd.Notes
	}

	if err := h.cmdRepo.Update(ctx, l); err != nil {
		return nil, err
	}

	return ToLeadDTO(l), nil
}
