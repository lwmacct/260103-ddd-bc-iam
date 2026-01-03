package auth_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/application/auth"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/application/user"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/platform/manualtest"
)

// TestLoginScenarios 测试各种登录场景（Table-Driven）。
//
// 覆盖场景：正确凭证、错误密码、错误验证码、不存在用户
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestLoginScenarios ./internal/integration/auth/
func TestLoginScenarios(t *testing.T) {
	manualtest.SkipIfNotManual(t)

	cases := []struct {
		name        string
		account     string
		password    string
		captcha     string // "valid" = 使用正确验证码, 其他值 = 使用该值
		wantSuccess bool
	}{
		{
			name:        "正确凭证登录成功",
			account:     "admin",
			password:    "admin123",
			captcha:     "valid",
			wantSuccess: true,
		},
		{
			name:        "错误密码被拒绝",
			account:     "admin",
			password:    "wrong_password",
			captcha:     "valid",
			wantSuccess: false,
		},
		{
			name:        "错误验证码被拒绝",
			account:     "admin",
			password:    "admin123",
			captcha:     "0000",
			wantSuccess: false,
		},
		{
			name:        "不存在的用户被拒绝",
			account:     "nonexistent_user_12345",
			password:    "any_password",
			captcha:     "valid",
			wantSuccess: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := manualtest.NewClient()

			// 获取验证码
			captcha, err := c.GetCaptcha()
			require.NoError(t, err, "获取验证码失败")

			// 构造请求
			captchaCode := tc.captcha
			if captchaCode == "valid" {
				captchaCode = captcha.Code
			}

			req := auth.LoginDTO{
				Account:   tc.account,
				Password:  tc.password,
				CaptchaID: captcha.ID,
				Captcha:   captchaCode,
			}

			resp, err := c.LoginWithCaptcha(req)

			// 验证结果
			if tc.wantSuccess {
				require.NoError(t, err, "期望成功，但登录失败")
				require.NotEmpty(t, resp.AccessToken, "未返回 access_token")
				require.NotEmpty(t, resp.RefreshToken, "未返回 refresh_token")
				require.NotZero(t, resp.User.UserID, "未返回有效用户 ID")
				t.Logf("登录成功: token=%s...", resp.AccessToken[:20])
			} else {
				require.True(t, err != nil || resp.AccessToken == "", "期望失败，但登录成功了")
				t.Logf("登录被正确拒绝")
			}
		})
	}
}

// TestGetCaptcha 测试获取验证码。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestGetCaptcha ./internal/integration/auth/
func TestGetCaptcha(t *testing.T) {
	manualtest.SkipIfNotManual(t)

	c := manualtest.NewClient()

	t.Log("获取验证码（开发模式）...")
	captcha, err := c.GetCaptcha()
	require.NoError(t, err, "获取验证码失败")
	require.NotEmpty(t, captcha.ID, "验证码 ID 为空")
	require.NotEmpty(t, captcha.Code, "验证码答案为空（开发模式应返回）")

	t.Logf("验证码获取成功!")
	t.Logf("  ID: %s", captcha.ID)
	t.Logf("  Code: %s", captcha.Code)
}

// TestAuthFlow 完整认证流程测试。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestAuthFlow ./internal/integration/auth/
func TestAuthFlow(t *testing.T) {
	manualtest.SkipIfNotManual(t)

	c := manualtest.NewClient()

	t.Log("步骤 1: 获取验证码")
	captcha, err := c.GetCaptcha()
	require.NoError(t, err, "获取验证码失败")
	t.Logf("  验证码 ID: %s", captcha.ID)

	t.Log("步骤 2: 登录")
	loginResp, err := c.Login("admin", "admin123")
	require.NoError(t, err, "登录失败")
	require.NotEmpty(t, loginResp.AccessToken, "登录成功但未返回 token")
	t.Logf("  登录成功，获取到 token")

	t.Log("步骤 3: 访问用户列表（验证 token）")
	resp, err := c.R().
		SetQueryParams(map[string]string{"page": "1", "limit": "1"}).
		Get("/api/admin/users")
	require.NoError(t, err, "请求失败")
	require.False(t, resp.IsError(), "预期状态码 200，实际 %d", resp.StatusCode())
	t.Log("  Token 验证成功，可以访问受保护资源")

	t.Log("认证流程测试完成!")
}

// TestRegister 测试用户注册。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestRegister ./internal/integration/auth/
func TestRegister(t *testing.T) {
	manualtest.SkipIfNotManual(t)

	c := manualtest.NewClient()

	// 生成唯一用户名
	testUsername := fmt.Sprintf("reguser_%d", time.Now().Unix())
	testEmail := testUsername + "@example.com"

	t.Log("测试用户注册...")
	t.Logf("  用户名: %s", testUsername)
	t.Logf("  邮箱: %s", testEmail)

	registerReq := auth.RegisterDTO{
		Username: testUsername,
		Email:    testEmail,
		Password: "password123",
		RealName: "注册测试用户",
	}

	resp, err := manualtest.Post[auth.RegisterResultDTO](c, "/api/auth/register", registerReq)
	require.NoError(t, err, "注册失败")
	require.NotZero(t, resp.UserID, "注册成功但未返回 user_id")
	require.NotEmpty(t, resp.AccessToken, "注册成功但未返回 access_token")

	t.Logf("注册成功!")
	t.Logf("  User ID: %d", resp.UserID)
	t.Logf("  Username: %s", resp.Username)
	t.Logf("  Access Token: %s...", resp.AccessToken[:50])

	// 注册清理函数，确保即使后续操作失败也能删除测试用户
	adminClient := manualtest.NewClient()
	_, adminErr := adminClient.Login("admin", "admin123")
	if adminErr == nil {
		t.Cleanup(func() {
			_ = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", resp.UserID))
		})
	}
}

// TestRegisterDuplicate 测试注册重复用户名。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestRegisterDuplicate ./internal/integration/auth/
func TestRegisterDuplicate(t *testing.T) {
	manualtest.SkipIfNotManual(t)

	c := manualtest.NewClient()

	// 生成唯一用户名
	testUsername := fmt.Sprintf("dupuser_%d", time.Now().Unix())
	testEmail := testUsername + "@example.com"

	t.Log("步骤 1: 注册第一个用户")
	registerReq := auth.RegisterDTO{
		Username: testUsername,
		Email:    testEmail,
		Password: "password123",
		RealName: "重复测试用户",
	}

	firstResp, err := manualtest.Post[auth.RegisterResultDTO](c, "/api/auth/register", registerReq)
	require.NoError(t, err, "首次注册失败")
	t.Logf("  首次注册成功，用户 ID: %d", firstResp.UserID)

	// 注册清理函数，确保即使后续操作失败也能删除测试用户
	adminClient := manualtest.NewClient()
	_, adminErr := adminClient.Login("admin", "admin123")
	if adminErr == nil {
		t.Cleanup(func() {
			_ = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", firstResp.UserID))
		})
	}

	t.Log("步骤 2: 尝试注册同名用户")
	duplicateReq := auth.RegisterDTO{
		Username: testUsername, // 相同用户名
		Email:    "another@example.com",
		Password: "password456",
		RealName: "重复测试用户2",
	}

	_, err = manualtest.Post[auth.RegisterResultDTO](c, "/api/auth/register", duplicateReq)
	require.Error(t, err, "重复用户名应该返回错误")
	t.Logf("  重复用户名被正确拒绝: %v", err)
}

// TestRefreshTokenScenarios 测试 Token 刷新各种场景（Table-Driven）。
//
// 覆盖场景：有效 token、无效 token、空 token、格式错误的 JWT
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestRefreshTokenScenarios ./internal/integration/auth/
func TestRefreshTokenScenarios(t *testing.T) {
	manualtest.SkipIfNotManual(t)

	// 先登录获取有效的 refresh_token
	c := manualtest.NewClient()
	loginResp, err := c.Login("admin", "admin123")
	require.NoError(t, err, "登录失败")

	cases := []struct {
		name         string
		refreshToken string
		wantSuccess  bool
	}{
		{
			name:         "有效 token 刷新成功",
			refreshToken: loginResp.RefreshToken,
			wantSuccess:  true,
		},
		{
			name:         "无效 token 被拒绝",
			refreshToken: "invalid_token_string",
			wantSuccess:  false,
		},
		{
			name:         "空 token 被拒绝",
			refreshToken: "",
			wantSuccess:  false,
		},
		{
			name:         "格式错误的 JWT 被拒绝",
			refreshToken: "not.a.valid.jwt.format",
			wantSuccess:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			refreshReq := auth.RefreshTokenDTO{
				RefreshToken: tc.refreshToken,
			}

			resp, err := manualtest.Post[auth.RefreshTokenResultDTO](c, "/api/auth/refresh", refreshReq)

			if tc.wantSuccess {
				require.NoError(t, err, "期望成功，但刷新失败")
				require.NotEmpty(t, resp.AccessToken, "刷新成功但未返回新 access_token")
				t.Logf("Token 刷新成功: token=%s...", resp.AccessToken[:20])
			} else {
				require.Error(t, err, "期望失败，但刷新成功了")
				t.Logf("无效 token 被正确拒绝: %v", err)
			}
		})
	}
}

// TestGetCurrentUser 测试获取当前登录用户信息。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestGetCurrentUser ./internal/integration/auth/
func TestGetCurrentUser(t *testing.T) {
	manualtest.SkipIfNotManual(t)

	c := manualtest.NewClient()

	t.Log("步骤 1: 登录")
	_, err := c.Login("admin", "admin123")
	require.NoError(t, err, "登录失败")
	t.Log("  登录成功")

	t.Log("步骤 2: 获取当前用户信息")
	me, err := manualtest.Get[user.UserWithRolesDTO](c, "/api/user/profile", nil)
	require.NoError(t, err, "获取当前用户失败")
	require.NotZero(t, me.ID, "返回的用户 ID 为 0")
	require.NotEmpty(t, me.Username, "返回的用户名为空")

	t.Logf("获取成功!")
	t.Logf("  ID: %d", me.ID)
	t.Logf("  用户名: %s", me.Username)
	t.Logf("  邮箱: %s", me.Email)
	t.Logf("  角色数量: %d", len(me.Roles))
	for _, role := range me.Roles {
		t.Logf("    - %s (%s)", role.DisplayName, role.Name)
	}
}
