package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/captcha"
)

// CaptchaHandler 验证码处理器
type CaptchaHandler struct {
	generateHandler *captcha.GenerateHandler
	devSecret       string
}

// NewCaptchaHandler 创建验证码处理器
func NewCaptchaHandler(
	generateHandler *captcha.GenerateHandler,
	devSecret string,
) *CaptchaHandler {
	return &CaptchaHandler{
		generateHandler: generateHandler,
		devSecret:       devSecret,
	}
}

// GetCaptcha 获取验证码
//
//	@Summary		获取验证码
//	@Description	生成图形验证码用于登录。支持开发模式（通过code和secret参数指定验证码值）
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			code	query		string												false	"开发模式：指定验证码值"
//	@Param			secret	query		string												false	"开发模式：密钥"
//	@Success		200		{object}	response.DataResponse[captcha.GenerateResultDTO]	"验证码生成成功"
//	@Failure		500		{object}	response.ErrorResponse								"生成失败"
//	@Router			/api/auth/captcha [get]
func (h *CaptchaHandler) GetCaptcha(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析开发模式参数
	code := c.Query("code")
	secret := c.Query("secret")
	isDevMode := code != "" && secret != "" && secret == h.devSecret

	// 构建命令
	cmd := captcha.GenerateCommand{
		DevMode:    isDevMode,
		CustomCode: code,
	}

	// 调用 Application Handler
	result, err := h.generateHandler.Handle(ctx, cmd)
	if err != nil {
		response.InternalError(c, "failed to generate captcha", err.Error())
		return
	}

	response.OK(c, result)
}
