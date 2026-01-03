package captcha

// GenerateResultDTO 验证码生成结果
type GenerateResultDTO struct {
	// ID 验证码ID
	ID string `json:"id"`
	// Image Base64编码的图片
	Image string `json:"image"`
	// ExpireAt 过期时间戳（秒）
	ExpireAt int64 `json:"expire_at"`
	// Code 验证码值（仅开发模式返回）
	Code string `json:"code,omitempty"`
}
