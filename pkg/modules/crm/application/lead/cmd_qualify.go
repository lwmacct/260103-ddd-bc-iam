package lead

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/lead"
)

// QualifyHandler 转换到已确认状态处理器。
type QualifyHandler struct {
	cmdRepo   lead.CommandRepository
	queryRepo lead.QueryRepository
}

// NewQualifyHandler 创建 QualifyHandler。
func NewQualifyHandler(cmdRepo lead.CommandRepository, queryRepo lead.QueryRepository) *QualifyHandler {
	return &QualifyHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 处理转换到已确认状态命令。
func (h *QualifyHandler) Handle(ctx context.Context, cmd QualifyCommand) (*LeadDTO, error) {
	l, err := h.queryRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	if err := l.Qualify(); err != nil {
		return nil, err
	}

	if err := h.cmdRepo.Update(ctx, l); err != nil {
		return nil, err
	}

	return ToLeadDTO(l), nil
}
