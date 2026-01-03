package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/auth"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/user"
)

// RegisterHandler 注册命令处理器
type RegisterHandler struct {
	userCommandRepo user.CommandRepository
	userQueryRepo   user.QueryRepository
	authService     auth.Service
}

// NewRegisterHandler 创建注册命令处理器
func NewRegisterHandler(
	userCommandRepo user.CommandRepository,
	userQueryRepo user.QueryRepository,
	authService auth.Service,
) *RegisterHandler {
	return &RegisterHandler{
		userCommandRepo: userCommandRepo,
		userQueryRepo:   userQueryRepo,
		authService:     authService,
	}
}

// Handle 处理注册命令
func (h *RegisterHandler) Handle(ctx context.Context, cmd RegisterCommand) (*RegisterResultDTO, error) {
	// 1. 验证密码策略
	if err := h.authService.ValidatePasswordPolicy(ctx, cmd.Password); err != nil {
		return nil, err
	}

	// 2. 检查用户名是否已存在
	exists, err := h.userQueryRepo.ExistsByUsername(ctx, cmd.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return nil, user.ErrUsernameAlreadyExists
	}

	// 3. 检查邮箱是否已存在
	exists, err = h.userQueryRepo.ExistsByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, user.ErrEmailAlreadyExists
	}

	// 4. 生成密码哈希
	hashedPassword, err := h.authService.GeneratePasswordHash(ctx, cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 5. 创建用户
	// 将字符串转换为指针类型（Email 和 Phone 在 Domain 层是 nullable）
	var emailPtr *string
	if cmd.Email != "" {
		emailPtr = &cmd.Email
	}
	var phonePtr *string
	if cmd.Phone != "" {
		phonePtr = &cmd.Phone
	}

	newUser := &user.User{
		Username:  cmd.Username,
		Email:     emailPtr,
		Password:  hashedPassword,
		RealName:  cmd.RealName,
		Nickname:  cmd.Nickname,
		Phone:     phonePtr,
		Signature: cmd.Signature,
		Status:    "active",
	}

	if err = h.userCommandRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 6. 生成访问令牌（新架构：不传递 roles，权限从缓存查询）
	accessToken, expiresAt, err := h.authService.GenerateAccessToken(ctx, newUser.ID, newUser.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, _, err := h.authService.GenerateRefreshToken(ctx, newUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &RegisterResultDTO{
		UserID:       newUser.ID,
		Username:     newUser.Username,
		Email:        stringPtrValueAuth(newUser.Email),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(time.Until(expiresAt).Seconds()),
	}, nil
}

// stringPtrValueAuth 将 *string 转换为 string，nil 返回空字符串
func stringPtrValueAuth(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
