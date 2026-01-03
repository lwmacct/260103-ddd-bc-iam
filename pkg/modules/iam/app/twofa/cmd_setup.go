package twofa

import (
	"context"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/twofa"
)

// SetupHandler 设置 2FA 命令处理器
type SetupHandler struct {
	twofaService twofa.Service
}

// NewSetupHandler 创建设置 2FA 命令处理器
func NewSetupHandler(twofaService twofa.Service) *SetupHandler {
	return &SetupHandler{
		twofaService: twofaService,
	}
}

// Handle 处理设置 2FA 命令
func (h *SetupHandler) Handle(ctx context.Context, cmd SetupCommand) (*SetupResultDTO, error) {
	result, err := h.twofaService.Setup(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	return &SetupResultDTO{
		Secret:    result.Secret,
		QRCodeURL: result.QRCodeURL,
		QRCodeImg: result.QRCodeImg,
	}, nil
}
