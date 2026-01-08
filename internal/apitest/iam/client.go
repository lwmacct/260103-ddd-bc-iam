package iam

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/response"
	"github.com/lwmacct/260103-ddd-shared/pkg/shared/apitest"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/auth"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/captcha"
)

// Client IAM 模块的测试客户端。
// 嵌入 apitest.Client 提供基础 HTTP 功能，并添加 IAM 特定方法（登录、验证码）。
type Client struct {
	apitest.Client

	devSecret string
}

// NewClient 创建 IAM 测试客户端。
func NewClient(baseURL, devSecret string) *Client {
	return &Client{
		Client:    *apitest.NewClient(baseURL),
		devSecret: devSecret,
	}
}

// R 返回 resty.Request，用于自定义请求。
func (c *Client) R() *resty.Request {
	return c.Client.R()
}

// HTTPClient 返回嵌入的 apitest.Client，用于泛型 HTTP 方法。
func (c *Client) HTTPClient() *apitest.Client {
	return &c.Client
}

// GetCaptcha 获取验证码（开发模式）。
func (c *Client) GetCaptcha() (*captcha.GenerateResultDTO, error) {
	var result response.DataResponse[captcha.GenerateResultDTO]

	resp, err := c.R().
		SetQueryParams(map[string]string{
			"code":   "9999",
			"secret": c.devSecret,
		}).
		SetResult(&result).
		Get("/api/auth/captcha")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("状态码 %d: %s", resp.StatusCode(), result.Message)
	}

	return &result.Data, nil
}

// Login 执行登录（自动获取验证码）。
func (c *Client) Login(account, password string) (*auth.LoginResponseDTO, error) {
	captchaResp, err := c.GetCaptcha()
	if err != nil {
		return nil, fmt.Errorf("获取验证码失败: %w", err)
	}

	return c.LoginWithCaptcha(auth.LoginDTO{
		Account:   account,
		Password:  password,
		CaptchaID: captchaResp.ID,
		Captcha:   captchaResp.Code,
	})
}

// LoginWithCaptcha 使用指定验证码登录。
func (c *Client) LoginWithCaptcha(req auth.LoginDTO) (*auth.LoginResponseDTO, error) {
	var result response.DataResponse[auth.LoginResponseDTO]

	resp, err := c.R().
		SetBody(req).
		SetResult(&result).
		Post("/api/auth/login")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	if result.Data.AccessToken != "" {
		c.SetToken(result.Data.AccessToken)
	}

	if resp.IsError() {
		return &result.Data, fmt.Errorf("登录失败 [%d]: %s", resp.StatusCode(), result.Message)
	}

	return &result.Data, nil
}
