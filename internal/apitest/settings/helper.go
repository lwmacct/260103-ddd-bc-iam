package settings

import (
	"os"
	"testing"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/260103-ddd-bc-iam/internal/apitest/iam"
	internalConfig "github.com/lwmacct/260103-ddd-bc-iam/internal/config"
	"github.com/lwmacct/260103-ddd-shared/pkg/shared/apitest"
)

// Client Settings 模块的测试客户端。
// 基于 apitest.Client，通过 IAM 的登录辅助获取认证。
type Client struct {
	apitest.Client
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

	baseURL := cfg.GetBaseUrl(false)
	return &Client{
		Client: *apitest.NewClient(baseURL),
	}
}

// HTTPClient 返回嵌入的 apitest.Client，用于泛型 HTTP 方法。
func (c *Client) HTTPClient() *apitest.Client {
	return &c.Client
}

// SkipIfNotAPITest 如果 API_TEST 环境变量未设置则跳过测试。
func SkipIfNotAPITest(t *testing.T) {
	t.Helper()
	if os.Getenv("API_TEST") == "" {
		t.SkipNow()
	}
}

// LoginAsAdmin 登录管理员账户，返回已认证的客户端。
// 委托给 iam.LoginAsAdmin 实现。
func LoginAsAdmin(t *testing.T) *Client {
	t.Helper()
	c := iam.LoginAsAdmin(t)
	return &Client{
		Client: c.Client,
	}
}

// LoginAsAdminForced 强制重新登录管理员账户（不使用缓存），返回已认证的客户端。
// 委托给 iam.LoginAsAdminForced 实现。
func LoginAsAdminForced(t *testing.T) *Client {
	t.Helper()
	c := iam.LoginAsAdminForced(t)
	return &Client{
		Client: c.Client,
	}
}

// LoginAs 使用指定账户登录，返回已认证的客户端。
// 委托给 iam.LoginAs 实现。
func LoginAs(t *testing.T, account, password string) *Client {
	t.Helper()
	c := iam.LoginAs(t, account, password)
	return &Client{
		Client: c.Client,
	}
}

// LoginAsForced 强制重新登录（不使用缓存），返回已认证的客户端。
// 委托给 iam.LoginAsForced 实现。
func LoginAsForced(t *testing.T, account, password string) *Client {
	t.Helper()
	c := iam.LoginAsForced(t, account, password)
	return &Client{
		Client: c.Client,
	}
}
