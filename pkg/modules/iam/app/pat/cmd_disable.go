package pat

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/pat"
)

// DisableHandler 禁用 Token 命令处理器
type DisableHandler struct {
	patCommandRepo pat.CommandRepository
	patQueryRepo   pat.QueryRepository
}

// NewDisableHandler 创建 DisableHandler 实例
func NewDisableHandler(
	patCommandRepo pat.CommandRepository,
	patQueryRepo pat.QueryRepository,
) *DisableHandler {
	return &DisableHandler{
		patCommandRepo: patCommandRepo,
		patQueryRepo:   patQueryRepo,
	}
}

// Handle 处理禁用 Token 命令
func (h *DisableHandler) Handle(ctx context.Context, cmd DisableCommand) error {
	token, err := h.patQueryRepo.FindByID(ctx, cmd.TokenID)
	if err != nil || token == nil {
		return errors.New("token not found")
	}

	if token.UserID != cmd.UserID {
		return errors.New("token does not belong to this user")
	}

	if err := h.patCommandRepo.Disable(ctx, cmd.TokenID); err != nil {
		return fmt.Errorf("failed to disable token: %w", err)
	}

	return nil
}
