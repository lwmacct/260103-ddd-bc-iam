package pat

import (
	"context"
	"errors"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/pat"
)

// GetHandler 获取 Token 查询处理器
type GetHandler struct {
	patQueryRepo pat.QueryRepository
}

// NewGetHandler 创建 GetHandler 实例
func NewGetHandler(patQueryRepo pat.QueryRepository) *GetHandler {
	return &GetHandler{
		patQueryRepo: patQueryRepo,
	}
}

// Handle 处理获取 Token 查询
func (h *GetHandler) Handle(ctx context.Context, query GetQuery) (*TokenInfoDTO, error) {
	// 1. 查询 Token
	token, err := h.patQueryRepo.FindByID(ctx, query.TokenID)
	if err != nil || token == nil {
		return nil, errors.New("token not found")
	}

	// 2. 验证所有权
	if token.UserID != query.UserID {
		return nil, errors.New("token does not belong to this user")
	}

	return ToTokenInfoDTO(token), nil
}
