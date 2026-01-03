package user

import (
	"context"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/user"
)

// UpdateHandler 更新用户命令处理器
type UpdateHandler struct {
	userCommandRepo user.CommandRepository
	userQueryRepo   user.QueryRepository
}

// NewUpdateHandler 创建更新用户命令处理器
func NewUpdateHandler(
	userCommandRepo user.CommandRepository,
	userQueryRepo user.QueryRepository,
) *UpdateHandler {
	return &UpdateHandler{
		userCommandRepo: userCommandRepo,
		userQueryRepo:   userQueryRepo,
	}
}

// Handle 处理更新用户命令
func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*UpdateResultDTO, error) {
	// 1. 获取用户
	u, err := h.userQueryRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, user.ErrUserNotFound
	}

	// 2. 系统用户保护检查
	if cmd.Username != nil && *cmd.Username != u.Username {
		// 系统用户不可修改用户名
		if !u.CanModifyUsername() {
			return nil, user.ErrCannotModifySystemUsername
		}
		// 检查用户名是否已存在
		exists, err := h.userQueryRepo.ExistsByUsername(ctx, *cmd.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to check username existence: %w", err)
		}
		if exists {
			return nil, user.ErrUsernameAlreadyExists
		}
		u.Username = *cmd.Username
	}
	if cmd.Email != nil {
		newEmail := stringPtrValue(cmd.Email)
		oldEmail := stringPtrValue(u.Email)

		if err := h.validateEmailChange(ctx, newEmail, oldEmail); err != nil {
			return nil, err
		}

		u.Email = parseEmailPtr(cmd.Email, newEmail)
	}
	if cmd.RealName != nil {
		u.RealName = *cmd.RealName
	}
	if cmd.Nickname != nil {
		u.Nickname = *cmd.Nickname
	}
	if cmd.Phone != nil {
		// Phone 直接赋值指针（允许设置为空字符串 nil）
		u.Phone = cmd.Phone
	}
	if cmd.Signature != nil {
		u.Signature = *cmd.Signature
	}
	if cmd.Avatar != nil {
		u.Avatar = *cmd.Avatar
	}
	if cmd.Bio != nil {
		u.Bio = *cmd.Bio
	}
	if cmd.Status != nil {
		// root 用户状态不可修改
		if !u.CanModifyStatus() {
			return nil, user.ErrCannotModifyRootStatus
		}
		// 使用领域模型方法
		switch *cmd.Status {
		case "active":
			u.Activate()
		case "inactive":
			u.Deactivate()
		case "banned":
			u.Ban()
		default:
			return nil, user.ErrInvalidUserStatus
		}
	}

	// 3. 保存更新
	if err := h.userCommandRepo.Update(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &UpdateResultDTO{
		UserID: u.ID,
	}, nil
}

// validateEmailChange 验证邮箱变更是否合法。
// 如果新邮箱非空且与旧邮箱不同，检查是否已被其他用户使用。
func (h *UpdateHandler) validateEmailChange(ctx context.Context, newEmail, oldEmail string) error {
	if newEmail == "" || newEmail == oldEmail {
		return nil
	}
	exists, err := h.userQueryRepo.ExistsByEmail(ctx, newEmail)
	if err != nil {
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return user.ErrEmailAlreadyExists
	}
	return nil
}

// parseEmailPtr 解析邮箱指针：空字符串返回 nil，否则返回原指针。
func parseEmailPtr(emailPtr *string, value string) *string {
	if value == "" {
		return nil
	}
	return emailPtr
}
