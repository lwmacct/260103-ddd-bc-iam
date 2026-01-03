package pat

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/auth"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/pat"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/user"
)

// InternalCreateTokenResult Handler 内部返回类型（包含领域实体）
// 注意：这是内部类型，不应该直接序列化为 HTTP 响应
type InternalCreateTokenResult struct {
	Token      *pat.PersonalAccessToken
	PlainToken string
}

// CreateHandler 创建 Token 命令处理器
type CreateHandler struct {
	patCommandRepo pat.CommandRepository
	userQueryRepo  user.QueryRepository
	tokenGenerator auth.TokenGenerator
}

// NewCreateHandler 创建 CreateHandler 实例
func NewCreateHandler(
	patCommandRepo pat.CommandRepository,
	userQueryRepo user.QueryRepository,
	tokenGenerator auth.TokenGenerator,
) *CreateHandler {
	return &CreateHandler{
		patCommandRepo: patCommandRepo,
		userQueryRepo:  userQueryRepo,
		tokenGenerator: tokenGenerator,
	}
}

// Handle 处理创建 Token 命令
func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*InternalCreateTokenResult, error) {
	if cmd.UserID == 0 {
		return nil, errors.New("user ID is required")
	}

	// 1. 验证用户存在
	u, err := h.userQueryRepo.GetByIDWithRoles(ctx, cmd.UserID)
	if err != nil || u == nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 2. 处理 Scopes：默认 ["full"]
	scopes := cmd.Scopes
	if len(scopes) == 0 {
		scopes = []string{string(pat.ScopeFull)}
	}

	// 3. 验证 Scopes 有效性
	if invalid := pat.ValidateScopes(scopes); len(invalid) > 0 {
		return nil, fmt.Errorf("invalid scopes: %s", strings.Join(invalid, ", "))
	}

	// 4. 生成 Token
	plainToken, hashedToken, prefix, err := h.tokenGenerator.GeneratePAT()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 5. 创建 PAT 实体
	patEntity := &pat.PersonalAccessToken{
		UserID:      cmd.UserID,
		Name:        cmd.Name,
		Token:       hashedToken,
		TokenPrefix: prefix,
		Scopes:      scopes,
		ExpiresAt:   cmd.ExpiresAt,
		LastUsedAt:  nil,
		Status:      pat.StatusActive,
		IPWhitelist: cmd.IPWhitelist,
		Description: cmd.Description,
	}

	// 6. 保存 Token
	if err := h.patCommandRepo.Create(ctx, patEntity); err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return &InternalCreateTokenResult{
		Token:      patEntity,
		PlainToken: plainToken, // 返回明文 token（仅此一次）
	}, nil
}
