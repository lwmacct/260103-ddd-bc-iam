package twofa

import (
	"context"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/twofa"
)

// DisableHandler 禁用 2FA 命令处理器
type DisableHandler struct {
	twofaService twofa.Service
}

// NewDisableHandler 创建禁用 2FA 命令处理器
func NewDisableHandler(twofaService twofa.Service) *DisableHandler {
	return &DisableHandler{
		twofaService: twofaService,
	}
}

// Handle 处理禁用 2FA 命令
func (h *DisableHandler) Handle(ctx context.Context, cmd DisableCommand) error {
	return h.twofaService.Disable(ctx, cmd.UserID)
}
