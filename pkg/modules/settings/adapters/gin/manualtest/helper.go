package manualtest

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	internalConfig "github.com/lwmacct/260103-ddd-bc-iam/internal/config"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/auth"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/captcha"
	"github.com/stretchr/testify/require"
)

// cachedSession 缓存的登录会话。
type cachedSession struct {
	token     string
	expiresAt time.Time
}

// sessionCache 存储已认证的会话，避免重复登录。
// key: "account:password", value: *cachedSession
var sessionCache sync.Map

// SkipIfNotManual 如果 MANUAL 环境变量未设置则跳过测试。
func SkipIfNotManual(t *testing.T) {
	t.Helper()
	if os.Getenv("MANUAL") == "" {
		t.SkipNow()
	}
}

// NewClient 创建测试客户端，从环境变量读取配置。
func NewClient() *Client {
	cfg, err := cfgm.Load(
		internalConfig.DefaultConfig(),
		cfgm.WithCallerSkip(2),
	)
	if err != nil {
		panic("加载配置失败: " + err.Error())
	}

	baseURL := cfg.GetBaseUrl(false)
	return newClientWithSecret(baseURL, cfg.Auth.DevSecret)
}

// newClientWithSecret 创建带开发密钥的测试客户端。
func newClientWithSecret(baseURL, devSecret string) *Client {
	c := newClient(baseURL)
	c.devSecret = devSecret
	return c
}

// LoginAsAdmin 登录管理员账户，返回已认证的客户端。
// 登录失败会导致测试立即失败。
func LoginAsAdmin(t *testing.T) *Client {
	t.Helper()
	return LoginAs(t, "admin", "admin123")
}

// LoginAsAdminForced 强制重新登录管理员账户（不使用缓存），返回已认证的客户端。
// 用于需要最新权限的场景（如权限变更后的测试）。
func LoginAsAdminForced(t *testing.T) *Client {
	t.Helper()
	return LoginAsForced(t, "admin", "admin123")
}

// LoginAsForced 强制重新登录（不使用缓存），返回已认证的客户端。
func LoginAsForced(t *testing.T, account, password string) *Client {
	t.Helper()
	SkipIfNotManual(t)

	// 清除缓存强制重新登录
	cacheKey := account + ":" + password
	sessionCache.Delete(cacheKey)
	return LoginAs(t, account, password)
}

// LoginAs 使用指定账户登录，返回已认证的客户端。
// 会复用缓存的 session 避免重复登录。
// 登录失败会导致测试立即失败。
func LoginAs(t *testing.T, account, password string) *Client {
	t.Helper()
	SkipIfNotManual(t)

	cacheKey := account + ":" + password

	// 检查缓存
	if cached, ok := sessionCache.Load(cacheKey); ok {
		if session, ok := cached.(*cachedSession); ok && time.Now().Before(session.expiresAt) {
			c := NewClient()
			c.SetToken(session.token)
			return c
		}
		// token 过期或类型错误，删除缓存
		sessionCache.Delete(cacheKey)
	}

	// 执行真实登录
	c := NewClient()
	resp, err := c.Login(account, password)
	require.NoError(t, err, "登录失败: account=%s", account)

	// 缓存 session（token 有效期 30 分钟，缓存 25 分钟）
	sessionCache.Store(cacheKey, &cachedSession{
		token:     resp.AccessToken,
		expiresAt: time.Now().Add(25 * time.Minute),
	})

	return c
}

// GetCaptcha 获取验证码（开发模式）。
func (c *Client) GetCaptcha() (*captcha.GenerateResultDTO, error) {
	var result struct {
		Code    int                       `json:"code"`
		Message string                    `json:"message"`
		Data    captcha.GenerateResultDTO `json:"data"`
	}

	resp, err := c.resty.R().
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

	return c.LoginWithCaptcha(map[string]any{
		"account":    account,
		"password":   password,
		"captcha_id": captchaResp.ID,
		"captcha":    captchaResp.Code,
	})
}

// LoginWithCaptcha 使用指定验证码登录。
func (c *Client) LoginWithCaptcha(req map[string]any) (*auth.LoginResponseDTO, error) {
	var result struct {
		Code    int                   `json:"code"`
		Message string                `json:"message"`
		Data    auth.LoginResponseDTO `json:"data"`
	}

	resp, err := c.resty.R().
		SetBody(req).
		SetResult(&result).
		Post("/api/auth/login")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	if result.Data.AccessToken != "" {
		c.token = result.Data.AccessToken
		c.resty.SetAuthToken(c.token)
	}

	if resp.IsError() {
		return &result.Data, fmt.Errorf("登录失败 [%d]: %s", resp.StatusCode(), result.Message)
	}

	return &result.Data, nil
}
