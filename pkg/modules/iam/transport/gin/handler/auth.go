package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/application/auth"
	authDomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/auth"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
)

// AuthHandler 认证处理器（新架构）
type AuthHandler struct {
	loginHandler        *auth.LoginHandler
	login2FAHandler     *auth.Login2FAHandler
	registerHandler     *auth.RegisterHandler
	refreshTokenHandler *auth.RefreshTokenHandler
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(
	loginHandler *auth.LoginHandler,
	login2FAHandler *auth.Login2FAHandler,
	registerHandler *auth.RegisterHandler,
	refreshTokenHandler *auth.RefreshTokenHandler,
) *AuthHandler {
	return &AuthHandler{
		loginHandler:        loginHandler,
		login2FAHandler:     login2FAHandler,
		registerHandler:     registerHandler,
		refreshTokenHandler: refreshTokenHandler,
	}
}

// Register 用户注册
//
//	@Summary		注册
//	@Description	创建新用户账号，注册成功后自动登录并返回访问令牌
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		auth.RegisterDTO								true	"注册信息"
//	@Success		201		{object}	response.DataResponse[auth.RegisterResultDTO]	"注册成功"
//	@Failure		400		{object}	response.ErrorResponse							"参数错误或用户名/邮箱已存在"
//	@Router			/api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req auth.RegisterDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	result, err := h.registerHandler.Handle(c.Request.Context(), auth.RegisterCommand(req))

	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, result)
}

// Login 用户登录
//
//	@Summary		登录
//	@Description	使用手机号/用户名/邮箱和密码登录系统，需要提供图形验证码。如果启用了2FA，返回session_token用于后续2FA验证
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		auth.LoginDTO									true	"登录凭证"
//	@Success		200		{object}	response.DataResponse[auth.LoginResponseDTO]	"登录成功或需要2FA验证"
//	@Failure		401		{object}	response.ErrorResponse							"登录失败：凭证无效、验证码错误或账户被禁用"
//	@Router			/api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用 Use Case Handler（传递 IP 和 UserAgent 用于审计日志）
	result, err := h.loginHandler.Handle(c.Request.Context(), auth.LoginCommand{
		Account:   req.Account,
		Password:  req.Password,
		CaptchaID: req.CaptchaID,
		Captcha:   req.Captcha,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	})

	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	// 检查是否需要 2FA
	if result.Requires2FA {
		response.OK(c, &auth.TwoFARequiredDTO{
			Requires2FA:  true,
			SessionToken: result.SessionToken,
		}, authDomain.MsgTwoFARequired)
		return
	}

	// 正常登录成功
	response.OK(c, result.ToLoginResponse())
}

// Login2FA 二次认证登录
//
//	@Summary		2FA 登录
//	@Description	使用session_token和2FA验证码完成登录（适用于启用了2FA的账户）
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		auth.Login2FADTO						true	"二次认证凭证"
//	@Success		200		{object}	response.DataResponse[auth.TokenDTO]	"登录成功"
//	@Failure		401		{object}	response.ErrorResponse					"验证失败：session_token无效或2FA验证码错误"
//	@Router			/api/auth/login/2fa [post]
func (h *AuthHandler) Login2FA(c *gin.Context) {
	var req auth.Login2FADTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用 Use Case Handler（传递 IP 和 UserAgent 用于审计日志）
	result, err := h.login2FAHandler.Handle(c.Request.Context(), auth.Login2FACommand{
		SessionToken:  req.SessionToken,
		TwoFactorCode: req.TwoFactorCode,
		ClientIP:      c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
	})

	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.OK(c, result.ToLoginResponse())
}

// RefreshToken 刷新访问令牌
//
//	@Summary		刷新令牌
//	@Description	使用refresh_token获取新的access_token和refresh_token，延长会话有效期
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		auth.RefreshTokenDTO							true	"刷新令牌"
//	@Success		200		{object}	response.DataResponse[auth.LoginResponseDTO]	"令牌刷新成功"
//	@Failure		401		{object}	response.ErrorResponse							"刷新令牌无效或已过期"
//	@Router			/api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req auth.RefreshTokenDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 构建命令（添加 IP 和 UserAgent 用于审计）
	cmd := auth.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
		ClientIP:     c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
	}

	result, err := h.refreshTokenHandler.Handle(c.Request.Context(), cmd)

	if err != nil {
		// 处理特定错误
		if errors.Is(err, authDomain.ErrInvalidToken) || errors.Is(err, authDomain.ErrTokenExpired) {
			response.Unauthorized(c, err.Error())
			return
		}
		response.Unauthorized(c, err.Error())
		return
	}

	response.OK(c, result.ToRefreshTokenResponse())
}
