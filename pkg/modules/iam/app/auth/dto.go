package auth

import "github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/auth"

// 重新导出领域错误供 Adapters 层使用（遵循 DDD 依赖方向）
var (
	ErrInvalidToken = auth.ErrInvalidToken
	ErrTokenExpired = auth.ErrTokenExpired
)

// LoginDTO 登录请求
type LoginDTO struct {
	Account   string `json:"account" binding:"required" example:"admin"`         // 手机号/用户名/邮箱
	Password  string `json:"password" binding:"required" example:"admin123"`     // 密码
	CaptchaID string `json:"captcha_id" binding:"required" example:"dev-123456"` // 验证码ID
	Captcha   string `json:"captcha" binding:"required" example:"9999"`          // 验证码
}

// Login2FADTO 二次认证请求
type Login2FADTO struct {
	SessionToken  string `json:"session_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIs..."` // 登录时返回的临时会话令牌
	TwoFactorCode string `json:"two_factor_code" binding:"required,len=6" example:"123456"`          // 6位TOTP验证码
}

// RegisterDTO 注册请求
type RegisterDTO struct {
	Username  string `json:"username" binding:"required,min=3,max=50" example:"john_doe"`
	Email     string `json:"email" binding:"required,email" example:"john@example.com"`
	Password  string `json:"password" binding:"required,min=6" example:"password123"`
	RealName  string `json:"real_name" binding:"max=100" example:"John Doe"`
	Nickname  string `json:"nickname" binding:"max=50" example:"Johnny"`
	Phone     string `json:"phone" binding:"omitempty,len=11" example:"13800138000"`
	Signature string `json:"signature" binding:"max=255" example:"Hello World"`
}

// RefreshTokenDTO 刷新令牌请求
type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// TokenDTO 令牌响应 DTO
type TokenDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// LoginResultDTO 登录结果 DTO（Handler 返回类型）
type LoginResultDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	Requires2FA  bool   `json:"requires_2fa"`
	SessionToken string `json:"session_token"`
	Roles        []LoginRoleDTO
}

// RefreshTokenResultDTO 刷新令牌结果 DTO（Handler 返回类型）
type RefreshTokenResultDTO struct {
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	TokenType    string         `json:"token_type"`
	ExpiresIn    int            `json:"expires_in"`
	UserID       uint           `json:"user_id"`
	Username     string         `json:"username"`
	Roles        []LoginRoleDTO `json:"roles"`
}

// ToRefreshTokenResponse 将 RefreshTokenResultDTO 转换为 HTTP 响应格式
// 复用 LoginResponseDTO 结构，保持登录和刷新响应格式一致
func (r *RefreshTokenResultDTO) ToRefreshTokenResponse() *LoginResponseDTO {
	return &LoginResponseDTO{
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
		TokenType:    r.TokenType,
		ExpiresIn:    r.ExpiresIn,
		User: UserBriefDTO{
			UserID:   r.UserID,
			Username: r.Username,
			Roles:    r.Roles,
		},
	}
}

// RegisterResultDTO 注册结果 DTO（Handler 返回类型）
type RegisterResultDTO struct {
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// LoginRoleDTO 登录响应中的角色信息 DTO
type LoginRoleDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// UserBriefDTO 用户简要信息响应 DTO
type UserBriefDTO struct {
	UserID   uint           `json:"user_id"`
	Username string         `json:"username"`
	Roles    []LoginRoleDTO `json:"roles"`
}

// TwoFARequiredDTO 需要二次认证响应 DTO
type TwoFARequiredDTO struct {
	Requires2FA  bool   `json:"requires_2fa"`
	SessionToken string `json:"session_token"`
}

// LoginResponseDTO 登录成功 HTTP 响应 DTO（与 HTTP API 响应格式匹配）
// 支持两种场景：正常登录（返回 token）或需要 2FA（返回 session_token）
type LoginResponseDTO struct {
	AccessToken  string       `json:"access_token,omitempty"`
	RefreshToken string       `json:"refresh_token,omitempty"`
	TokenType    string       `json:"token_type,omitempty"`
	ExpiresIn    int          `json:"expires_in,omitempty"`
	User         UserBriefDTO `json:"user,omitzero"`
	// 2FA 相关（当需要 2FA 时返回）
	Requires2FA  bool   `json:"requires_2fa,omitempty"`
	SessionToken string `json:"session_token,omitempty"`
}

// ToLoginResponse 将 LoginResultDTO 转换为 HTTP 响应格式
func (r *LoginResultDTO) ToLoginResponse() *LoginResponseDTO {
	return &LoginResponseDTO{
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
		TokenType:    r.TokenType,
		ExpiresIn:    r.ExpiresIn,
		User: UserBriefDTO{
			UserID:   r.UserID,
			Username: r.Username,
			Roles:    r.Roles,
		},
		Requires2FA:  r.Requires2FA,
		SessionToken: r.SessionToken,
	}
}
