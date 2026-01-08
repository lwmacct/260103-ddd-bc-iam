package pat

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/pat"
)

// EnableHandler 启用 Token 命令处理器
type EnableHandler struct {
	patCommandRepo pat.CommandRepository
	patQueryRepo   pat.QueryRepository
}

// NewEnableHandler 创建 EnableHandler 实例
func NewEnableHandler(
	patCommandRepo pat.CommandRepository,
	patQueryRepo pat.QueryRepository,
) *EnableHandler {
	return &EnableHandler{
		patCommandRepo: patCommandRepo,
		patQueryRepo:   patQueryRepo,
	}
}

// Handle 处理启用 Token 命令
func (h *EnableHandler) Handle(ctx context.Context, cmd EnableCommand) error {
	token, err := h.patQueryRepo.FindByID(ctx, cmd.TokenID)
	if err != nil || token == nil {
		return errors.New("token not found")
	}

	if token.UserID != cmd.UserID {
		return errors.New("token does not belong to this user")
	}

	if token.IsExpired() {
		return errors.New("token is expired and cannot be enabled")
	}

	if err := h.patCommandRepo.Enable(ctx, cmd.TokenID); err != nil {
		return fmt.Errorf("failed to enable token: %w", err)
	}

	return nil
}
