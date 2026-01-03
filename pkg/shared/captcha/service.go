package captcha

// Service 定义验证码生成的领域服务接口
type Service interface {
	// GenerateRandomCode 生成随机验证码
	// 返回 (captchaID, imageBase64, code, error)
	GenerateRandomCode() (captchaID string, imageBase64 string, code string, err error)

	// GenerateCustomCodeImage 生成指定文本的验证码图片（用于开发模式）
	GenerateCustomCodeImage(text string) (imageBase64 string, err error)

	// GenerateDevCaptchaID 生成开发模式验证码ID
	GenerateDevCaptchaID() string

	// GetDefaultExpiration 获取默认过期时间
	GetDefaultExpiration() int64
}
