package twofa

// SetupCommand 设置 2FA 命令
type SetupCommand struct {
	UserID uint
}

// VerifyEnableCommand 验证并启用 2FA 命令
type VerifyEnableCommand struct {
	UserID uint
	Code   string
}

// DisableCommand 禁用 2FA 命令
type DisableCommand struct {
	UserID uint
}
