package iam

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	internalConfig "github.com/lwmacct/260103-ddd-bc-iam/internal/config"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/config"
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

// SkipIfNotAPITest 如果 API_TEST 环境变量未设置则跳过测试。
func SkipIfNotAPITest(t *testing.T) {
	t.Helper()
	if os.Getenv("API_TEST") == "" {
		t.SkipNow()
	}
}

// NewClientFromConfig 从配置创建测试客户端。
func NewClientFromConfig() *Client {
	cfg, err := cfgm.Load(
		internalConfig.DefaultConfig(),
		cfgm.WithCallerSkip(2),
	)
	if err != nil {
		panic("加载配置失败: " + err.Error())
	}

	// 转换为 IAM 配置
	iamCfg := toIAMConfig(cfg)

	baseURL := cfg.GetBaseUrl(false)
	return NewClient(baseURL, iamCfg.Auth.DevSecret)
}

// toIAMConfig 从通用配置中提取 IAM 模块配置（测试用）。
func toIAMConfig(cfg *internalConfig.Config) config.Config {
	return config.Config{
		Auth: config.Auth{
			DevSecret: cfg.Auth.DevSecret,
		},
	}
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
	SkipIfNotAPITest(t)

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
	SkipIfNotAPITest(t)

	cacheKey := account + ":" + password

	// 检查缓存
	if cached, ok := sessionCache.Load(cacheKey); ok {
		if session, ok := cached.(*cachedSession); ok && time.Now().Before(session.expiresAt) {
			c := NewClientFromConfig()
			c.SetToken(session.token)
			return c
		}
		// token 过期或类型错误，删除缓存
		sessionCache.Delete(cacheKey)
	}

	// 执行真实登录
	c := NewClientFromConfig()
	resp, err := c.Login(account, password)
	require.NoError(t, err, "登录失败: account=%s", account)

	// 缓存 session（token 有效期 30 分钟，缓存 25 分钟）
	sessionCache.Store(cacheKey, &cachedSession{
		token:     resp.AccessToken,
		expiresAt: time.Now().Add(25 * time.Minute),
	})

	return c
}
