package pat

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/pat"
)

// DeleteHandler 删除 Token 命令处理器
type DeleteHandler struct {
	patCommandRepo pat.CommandRepository
	patQueryRepo   pat.QueryRepository
}

// NewDeleteHandler 创建 DeleteHandler 实例
func NewDeleteHandler(
	patCommandRepo pat.CommandRepository,
	patQueryRepo pat.QueryRepository,
) *DeleteHandler {
	return &DeleteHandler{
		patCommandRepo: patCommandRepo,
		patQueryRepo:   patQueryRepo,
	}
}

// Handle 处理删除 Token 命令
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	token, err := h.patQueryRepo.FindByID(ctx, cmd.TokenID)
	if err != nil || token == nil {
		return errors.New("token not found")
	}

	if token.UserID != cmd.UserID {
		return errors.New("token does not belong to this user")
	}

	if err := h.patCommandRepo.Delete(ctx, cmd.TokenID); err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}
