package auth

// LoginCommand 登录命令
type LoginCommand struct {
	Account   string // 用户名或邮箱
	Password  string
	CaptchaID string
	Captcha   string
	ClientIP  string // 客户端 IP（用于审计日志）
	UserAgent string // 用户代理（用于审计日志）
}

// Login2FACommand 二次认证命令
type Login2FACommand struct {
	SessionToken  string
	TwoFactorCode string
	ClientIP      string // 客户端 IP（用于审计日志）
	UserAgent     string // 用户代理（用于审计日志）
}

// RegisterCommand 注册命令
type RegisterCommand struct {
	Username  string
	Email     string
	Password  string
	RealName  string
	Nickname  string
	Phone     string
	Signature string
}

// RefreshTokenCommand 刷新令牌命令
type RefreshTokenCommand struct {
	RefreshToken string
	ClientIP     string // 客户端 IP（用于审计日志）
	UserAgent    string // 用户代理（用于审计日志）
}
