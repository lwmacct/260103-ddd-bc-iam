package twofa_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260103-ddd-iam-bc/internal/apitest/iam"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/auth"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/twofa"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/user"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/response"
	"github.com/lwmacct/260103-ddd-shared/pkg/shared/apitest"
)

// TestGetTwoFAStatus 测试获取 2FA 状态。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestGetTwoFAStatus ./internal/integration/twofa/
func TestGetTwoFAStatus(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	t.Log("获取 2FA 状态")
	status, err := apitest.Get[twofa.StatusDTO](c.HTTPClient(), "/api/auth/2fa/status", nil)
	require.NoError(t, err, "获取 2FA 状态失败")

	t.Logf("2FA 状态获取成功!")
	t.Logf("  启用状态: %v", status.Enabled)
	t.Logf("  剩余恢复码数量: %d", status.RecoveryCodesCount)
}

// TestTwoFAFlow 2FA 完整流程测试（设置 → 验证启用 → 状态检查 → 禁用）。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestTwoFAFlow ./internal/integration/twofa/
func TestTwoFAFlow(t *testing.T) {
	adminClient := iam.LoginAsAdmin(t)

	testUsername := fmt.Sprintf("twofa_test_%d", time.Now().Unix())
	testPassword := "password123"

	t.Log("步骤 1: 创建测试用户（带 user 角色）")
	createReq := user.CreateDTO{
		Username: testUsername,
		Email:    testUsername + "@example.com",
		Password: testPassword,
		RealName: "2FA 测试用户",
		RoleIDs:  []uint{2}, // user 角色 ID
	}

	createResp, err := apitest.Post[user.UserWithRolesDTO](adminClient.HTTPClient(), "/api/admin/users", createReq)
	require.NoError(t, err, "创建测试用户失败")
	testUserID := createResp.ID
	t.Logf("  创建成功，用户 ID: %d", testUserID)

	// 确保测试结束时清理资源
	t.Cleanup(func() {
		if delErr := adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID)); delErr != nil {
			t.Logf("清理测试用户失败: %v", delErr)
		}
	})

	// 用测试用户登录
	t.Log("步骤 2: 测试用户登录")
	testClient := iam.LoginAs(t, testUsername, testPassword)
	t.Log("  登录成功")

	// 设置 2FA
	t.Log("步骤 3: 设置 2FA")
	setup, err := apitest.Post[twofa.SetupDTO](testClient.HTTPClient(), "/api/auth/2fa/setup", nil)
	require.NoError(t, err, "设置 2FA 失败")
	require.NotEmpty(t, setup.Secret, "2FA 密钥为空")
	require.NotEmpty(t, setup.QRCodeURL, "二维码 URL 为空")
	require.NotEmpty(t, setup.QRCodeImg, "二维码图片为空")
	t.Logf("  密钥: %s", setup.Secret)
	t.Logf("  二维码 URL: %s...", setup.QRCodeURL[:50])
	t.Log("  二维码图片: [已生成]")

	// 使用密钥生成 TOTP 代码
	t.Log("步骤 4: 生成并验证 TOTP 代码")
	code, err := totp.GenerateCode(setup.Secret, time.Now())
	require.NoError(t, err, "生成 TOTP 代码失败")
	t.Logf("  生成的 TOTP 代码: %s", code)

	// 验证并启用 2FA
	verifyReq := map[string]string{"code": code}
	enableResp, err := apitest.Post[twofa.EnableDTO](testClient.HTTPClient(), "/api/auth/2fa/verify", verifyReq)
	require.NoError(t, err, "验证并启用 2FA 失败")
	assert.NotEmpty(t, enableResp.RecoveryCodes, "应返回恢复码，但列表为空")
	t.Logf("  2FA 启用成功!")
	t.Logf("  恢复码数量: %d", len(enableResp.RecoveryCodes))
	if len(enableResp.RecoveryCodes) > 0 {
		t.Logf("  第一个恢复码: %s", enableResp.RecoveryCodes[0])
	}

	// 检查 2FA 状态
	t.Log("步骤 5: 检查 2FA 状态")
	status, err := apitest.Get[twofa.StatusDTO](testClient.HTTPClient(), "/api/auth/2fa/status", nil)
	require.NoError(t, err, "获取 2FA 状态失败")
	require.True(t, status.Enabled, "2FA 应该已启用")
	t.Logf("  2FA 已启用，剩余恢复码: %d", status.RecoveryCodesCount)

	// 禁用 2FA
	t.Log("步骤 6: 禁用 2FA")
	resp, err := testClient.R().Post("/api/auth/2fa/disable")
	require.NoError(t, err, "禁用 2FA 请求失败")
	require.False(t, resp.IsError(), "禁用 2FA 失败，状态码: %d", resp.StatusCode())
	t.Log("  2FA 已禁用")

	// 验证 2FA 已禁用
	status2, err := apitest.Get[twofa.StatusDTO](testClient.HTTPClient(), "/api/auth/2fa/status", nil)
	require.NoError(t, err, "获取 2FA 状态失败")
	require.False(t, status2.Enabled, "2FA 应该已禁用")
	t.Log("  验证：2FA 已禁用")

	t.Log("2FA 完整流程测试完成!")
}

// TestSetup2FA 测试设置 2FA。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestSetup2FA ./internal/integration/twofa/
func TestSetup2FA(t *testing.T) {
	adminClient := iam.LoginAsAdmin(t)

	testUsername := fmt.Sprintf("setup2fa_%d", time.Now().Unix())

	t.Log("步骤 1: 创建测试用户")
	createReq := user.CreateDTO{
		Username: testUsername,
		Email:    testUsername + "@example.com",
		Password: "password123",
		RealName: "2FA Setup 测试用户",
		RoleIDs:  []uint{2},
	}

	createResp, err := apitest.Post[user.UserWithRolesDTO](adminClient.HTTPClient(), "/api/admin/users", createReq)
	require.NoError(t, err, "创建测试用户失败")
	testUserID := createResp.ID
	t.Logf("  创建成功，用户 ID: %d", testUserID)

	// 确保测试结束时清理资源
	t.Cleanup(func() {
		if delErr := adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID)); delErr != nil {
			t.Logf("清理测试用户失败: %v", delErr)
		}
	})

	// 用测试用户登录
	t.Log("步骤 2: 测试用户登录")
	c := iam.LoginAs(t, testUsername, "password123")
	t.Log("  登录成功")

	t.Log("步骤 3: 设置 2FA")
	setup, err := apitest.Post[twofa.SetupDTO](c.HTTPClient(), "/api/auth/2fa/setup", nil)
	require.NoError(t, err, "设置 2FA 失败")

	require.NotEmpty(t, setup.Secret, "2FA 密钥为空")
	require.NotEmpty(t, setup.QRCodeURL, "二维码 URL 为空")
	require.NotEmpty(t, setup.QRCodeImg, "二维码图片为空")

	t.Logf("2FA 设置成功!")
	t.Logf("  密钥: %s", setup.Secret)
	t.Logf("  二维码 URL 长度: %d", len(setup.QRCodeURL))
	t.Logf("  二维码图片大小: %d bytes", len(setup.QRCodeImg))
}

// TestDisable2FA 测试禁用 2FA。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestDisable2FA ./internal/integration/twofa/
func TestDisable2FA(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	t.Log("禁用 2FA")
	resp, err := c.R().Post("/api/auth/2fa/disable")
	require.NoError(t, err, "禁用 2FA 请求失败")

	// 即使 2FA 未启用，禁用也应该成功（幂等操作）
	if resp.IsError() {
		t.Logf("禁用 2FA 返回状态码: %d，响应: %s", resp.StatusCode(), resp.String())
	} else {
		t.Log("  禁用 2FA 成功")
	}

	// 检查状态
	t.Log("检查 2FA 状态")
	status, err := apitest.Get[twofa.StatusDTO](c.HTTPClient(), "/api/auth/2fa/status", nil)
	require.NoError(t, err, "获取 2FA 状态失败")

	t.Logf("  2FA 状态: 启用=%v, 恢复码=%d", status.Enabled, status.RecoveryCodesCount)
}

// TestLogin2FA 测试 2FA 登录流程（启用 2FA 后的完整登录）。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestLogin2FA ./internal/integration/twofa/
func TestLogin2FA(t *testing.T) {
	adminClient := iam.LoginAsAdmin(t)

	testUsername := fmt.Sprintf("login2fa_%d", time.Now().Unix())
	testPassword := "password123"

	t.Log("步骤 1: 创建测试用户（带 user 角色）")
	createReq := user.CreateDTO{
		Username: testUsername,
		Email:    testUsername + "@example.com",
		Password: testPassword,
		RealName: "2FA Login 测试用户",
		RoleIDs:  []uint{2}, // user 角色 ID
	}

	createResp, err := apitest.Post[user.UserWithRolesDTO](adminClient.HTTPClient(), "/api/admin/users", createReq)
	require.NoError(t, err, "创建测试用户失败")
	testUserID := createResp.ID
	t.Logf("  创建成功，用户 ID: %d", testUserID)

	// 确保测试结束时清理资源
	t.Cleanup(func() {
		_ = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
	})

	// 步骤 2: 测试用户登录并设置 2FA
	t.Log("步骤 2: 测试用户登录并设置 2FA")
	testClient := iam.LoginAs(t, testUsername, testPassword)
	t.Log("  登录成功")

	// 设置 2FA
	setup, err := apitest.Post[twofa.SetupDTO](testClient.HTTPClient(), "/api/auth/2fa/setup", nil)
	require.NoError(t, err, "设置 2FA 失败")
	t.Logf("  2FA 密钥: %s", setup.Secret)

	// 生成 TOTP 代码并启用 2FA
	code, err := totp.GenerateCode(setup.Secret, time.Now())
	require.NoError(t, err, "生成 TOTP 代码失败")
	verifyReq := map[string]string{"code": code}
	_, err = apitest.Post[twofa.EnableDTO](testClient.HTTPClient(), "/api/auth/2fa/verify", verifyReq)
	require.NoError(t, err, "启用 2FA 失败")
	t.Log("  2FA 已启用")

	// 步骤 3: 尝试再次登录（应返回 requires_2fa=true）
	t.Log("步骤 3: 再次登录（应触发 2FA 验证）")
	newClient := iam.NewClientFromConfig()
	loginResp, err := newClient.Login(testUsername, testPassword)
	if err != nil {
		// 如果登录返回错误但是因为需要 2FA，这是预期行为
		t.Logf("  登录返回: %v（这可能是预期的 2FA 挑战）", err)
	}

	require.NotNil(t, loginResp, "预期返回登录响应")
	require.True(t, loginResp.Requires2FA, "预期返回 requires_2fa=true")
	require.NotEmpty(t, loginResp.SessionToken, "预期返回 session_token")
	t.Logf("  收到 2FA 挑战")
	t.Logf("  Session Token: %s...", loginResp.SessionToken[:20])

	// 步骤 4: 使用 2FA 完成登录
	t.Log("步骤 4: 使用 2FA 完成登录")
	// 生成新的 TOTP 代码
	newCode, err := totp.GenerateCode(setup.Secret, time.Now())
	require.NoError(t, err, "生成 TOTP 代码失败")
	t.Logf("  TOTP 代码: %s", newCode)

	login2FAReq := auth.Login2FADTO{
		SessionToken:  loginResp.SessionToken,
		TwoFactorCode: newCode,
	}

	// 使用原始请求获取更详细的错误信息
	var finalResult response.DataResponse[auth.LoginResponseDTO]
	resp, err := newClient.R().
		SetBody(login2FAReq).
		SetResult(&finalResult).
		Post("/api/auth/login/2fa")
	require.NoError(t, err, "2FA 登录请求失败")
	require.False(t, resp.IsError(), "2FA 登录失败，状态码: %d, 响应: %s", resp.StatusCode(), resp.String())

	finalResp := &finalResult.Data
	require.NotEmpty(t, finalResp.AccessToken, "2FA 登录后未返回 access_token")
	t.Log("  2FA 登录成功!")
	t.Logf("  Access Token: %s...", finalResp.AccessToken[:30])
	t.Logf("  用户名: %s", finalResp.User.Username)

	// 步骤 5: 验证 token 可用
	t.Log("步骤 5: 验证 token 可用")
	newClient.SetToken(finalResp.AccessToken)
	profile, err := apitest.Get[user.UserWithRolesDTO](newClient.HTTPClient(), "/api/user/profile", nil)
	require.NoError(t, err, "使用新 token 获取资料失败")
	assert.Equal(t, testUsername, profile.Username, "用户名不匹配")
	t.Logf("  Token 验证成功，用户: %s", profile.Username)

	t.Log("\n2FA 登录测试完成!")
}
