package captcha

// GenerateCommand 生成验证码命令
type GenerateCommand struct {
	// DevMode 是否为开发模式
	DevMode bool
	// CustomCode 自定义验证码（开发模式）
	CustomCode string
}
