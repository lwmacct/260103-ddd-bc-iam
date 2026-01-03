package audit

import (
	"context"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/audit"
)

// CreateHandler 创建审计日志命令处理器
type CreateHandler struct {
	auditCommandRepo audit.CommandRepository
}

// NewCreateHandler 创建处理器实例
func NewCreateHandler(repo audit.CommandRepository) *CreateHandler {
	return &CreateHandler{
		auditCommandRepo: repo,
	}
}

// Handle 处理创建审计日志命令
func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) error {
	log := &audit.Audit{
		UserID:      cmd.UserID,
		Username:    cmd.Username,
		Action:      cmd.Action,
		Resource:    cmd.Resource,
		ResourceID:  cmd.ResourceID,
		IPAddress:   cmd.IPAddress,
		UserAgent:   cmd.UserAgent,
		Details:     cmd.Details,
		Status:      cmd.Status,
		RequestID:   cmd.RequestID,
		OperationID: cmd.OperationID,
	}

	return h.auditCommandRepo.Create(ctx, log)
}
