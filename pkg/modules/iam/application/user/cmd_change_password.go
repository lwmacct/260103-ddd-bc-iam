package user

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/auth"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/user"
)

// ChangePasswordHandler 处理修改密码命令
type ChangePasswordHandler struct {
	userCommandRepo user.CommandRepository
	userQueryRepo   user.QueryRepository
	authService     auth.Service
}

// NewChangePasswordHandler 创建新的 ChangePasswordHandler
func NewChangePasswordHandler(
	userCommandRepo user.CommandRepository,
	userQueryRepo user.QueryRepository,
	authService auth.Service,
) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		userCommandRepo: userCommandRepo,
		userQueryRepo:   userQueryRepo,
		authService:     authService,
	}
}

// Handle 执行修改密码逻辑
func (h *ChangePasswordHandler) Handle(ctx context.Context, cmd ChangePasswordCommand) error {
	// 查询用户
	u, err := h.userQueryRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if err = h.authService.VerifyPassword(ctx, u.Password, cmd.OldPassword); err != nil {
		return user.ErrInvalidPassword
	}

	// 验证策略
	if err = h.authService.ValidatePasswordPolicy(ctx, cmd.NewPassword); err != nil {
		return err
	}

	// 生成新密码
	hashedPassword, err := h.authService.GeneratePasswordHash(ctx, cmd.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 更新密码
	if err := h.userCommandRepo.UpdatePassword(ctx, cmd.UserID, hashedPassword); err != nil {
		return err
	}

	return nil
}
