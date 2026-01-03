package twofa

import "context"

// SetupResult 2FA 设置结果
type SetupResult struct {
	Secret    string // TOTP 密钥
	QRCodeURL string // 二维码 URL
	QRCodeImg string // Base64 编码的二维码图片
}

// Service 2FA 领域服务接口
// 定义双因素认证的核心能力
type Service interface {
	// Setup 设置 2FA（生成密钥和二维码）
	Setup(ctx context.Context, userID uint) (*SetupResult, error)

	// VerifyAndEnable 验证 TOTP 代码并启用 2FA
	// 返回恢复码列表
	VerifyAndEnable(ctx context.Context, userID uint, code string) ([]string, error)

	// Verify 验证 TOTP 代码或恢复码
	Verify(ctx context.Context, userID uint, code string) (bool, error)

	// Disable 禁用 2FA
	Disable(ctx context.Context, userID uint) error

	// GetStatus 获取 2FA 状态
	// 返回是否启用和剩余恢复码数量
	GetStatus(ctx context.Context, userID uint) (enabled bool, recoveryCodesCount int, err error)
}
