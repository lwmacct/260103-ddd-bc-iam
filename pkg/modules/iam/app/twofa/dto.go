// Package twofa 提供两步验证应用层 DTO。
package twofa

// SetupDTO 2FA 设置响应 DTO。
type SetupDTO struct {
	Secret    string `json:"secret"`     // TOTP 密钥（用户可手动输入）
	QRCodeURL string `json:"qrcode_url"` // 二维码 URL
	QRCodeImg string `json:"qrcode_img"` // Base64 编码的二维码图片
}

// EnableDTO 2FA 启用响应 DTO。
type EnableDTO struct {
	RecoveryCodes []string `json:"recovery_codes"` // 恢复码列表
	Message       string   `json:"message"`        // 提示消息
}

// StatusDTO 2FA 状态响应 DTO。
type StatusDTO struct {
	Enabled            bool `json:"enabled"`              // 是否启用
	RecoveryCodesCount int  `json:"recovery_codes_count"` // 剩余恢复码数量
}

// VerifyDTO 验证 TOTP 代码请求 DTO。
type VerifyDTO struct {
	Code string `json:"code" binding:"required" example:"123456"` // TOTP 验证码
}

// SetupResultDTO 2FA 设置结果 DTO（Handler 返回类型）
type SetupResultDTO struct {
	Secret    string `json:"secret"`
	QRCodeURL string `json:"qrcode_url"`
	QRCodeImg string `json:"qrcode_img"`
}

// EnableResultDTO 验证并启用 2FA 结果 DTO（Handler 返回类型）
type EnableResultDTO struct {
	RecoveryCodes []string `json:"recovery_codes"`
}

// StatusResultDTO 获取 2FA 状态结果 DTO（Handler 返回类型）
type StatusResultDTO struct {
	Enabled            bool `json:"enabled"`
	RecoveryCodesCount int  `json:"recovery_codes_count"`
}
