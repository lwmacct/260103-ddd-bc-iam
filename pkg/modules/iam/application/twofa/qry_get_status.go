package twofa

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/twofa"
)

// GetStatusHandler 获取 2FA 状态查询处理器
type GetStatusHandler struct {
	twofaService twofa.Service
}

// NewGetStatusHandler 创建获取 2FA 状态查询处理器
func NewGetStatusHandler(twofaService twofa.Service) *GetStatusHandler {
	return &GetStatusHandler{
		twofaService: twofaService,
	}
}

// Handle 处理获取 2FA 状态查询
func (h *GetStatusHandler) Handle(ctx context.Context, query GetStatusQuery) (*StatusResultDTO, error) {
	enabled, count, err := h.twofaService.GetStatus(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	return &StatusResultDTO{
		Enabled:            enabled,
		RecoveryCodesCount: count,
	}, nil
}
