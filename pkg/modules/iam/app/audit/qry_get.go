package audit

import (
	"context"
	"errors"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/audit"
)

// GetHandler 获取审计日志查询处理器
type GetHandler struct {
	auditQueryRepo audit.QueryRepository
}

// NewGetHandler 创建 GetHandler 实例
func NewGetHandler(auditQueryRepo audit.QueryRepository) *GetHandler {
	return &GetHandler{
		auditQueryRepo: auditQueryRepo,
	}
}

// Handle 处理获取审计日志查询
func (h *GetHandler) Handle(ctx context.Context, query GetQuery) (*AuditDTO, error) {
	log, err := h.auditQueryRepo.FindByID(ctx, query.LogID)
	if err != nil || log == nil {
		return nil, errors.New("audit log not found")
	}

	return ToAuditDTO(log), nil
}
