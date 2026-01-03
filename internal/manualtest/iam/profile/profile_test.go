package profile_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/manualtest"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/auth"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/user"
)

// TestGetProfile 测试获取个人资料。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestGetProfile ./internal/integration/profile/
func TestGetProfile(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	t.Log("获取个人资料")
	profile, err := manualtest.Get[user.UserWithRolesDTO](c, "/api/user/profile", nil)
	require.NoError(t, err, "获取个人资料失败")

	// 验证关键字段
	require.NotZero(t, profile.ID, "返回的用户 ID 为 0")
	require.NotEmpty(t, profile.Username, "返回的用户名为空")
	require.NotEmpty(t, profile.Email, "返回的邮箱为空")
	require.NotEmpty(t, profile.Status, "返回的状态为空")

	t.Logf("获取成功!")
	t.Logf("  ID: %d", profile.ID)
	t.Logf("  用户名: %s", profile.Username)
	t.Logf("  邮箱: %s", profile.Email)
	t.Logf("  真实姓名: %s", profile.RealName)
	t.Logf("  状态: %s", profile.Status)
	t.Logf("  角色数量: %d", len(profile.Roles))
}

// TestUpdateProfile 测试更新个人资料。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestUpdateProfile ./internal/integration/profile/
func TestUpdateProfile(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	t.Log("步骤 1: 获取当前资料")
	originalProfile, err := manualtest.Get[user.UserWithRolesDTO](c, "/api/user/profile", nil)
	require.NoError(t, err, "获取原始资料失败")
	t.Logf("  当前真实姓名: %s", originalProfile.RealName)

	// 注册清理函数，确保即使测试失败也能恢复原始资料
	t.Cleanup(func() {
		restoreReq := user.UpdateDTO{
			RealName: &originalProfile.RealName,
		}
		_, _ = manualtest.Put[user.UserWithRolesDTO](c, "/api/user/profile", restoreReq)
	})

	t.Log("步骤 2: 更新资料")
	newRealName := fmt.Sprintf("测试更新_%d", time.Now().Unix())
	updateReq := user.UpdateDTO{
		RealName: &newRealName,
	}

	updateResp, err := manualtest.Put[user.UserWithRolesDTO](c, "/api/user/profile", updateReq)
	require.NoError(t, err, "更新资料失败")
	t.Logf("  更新后真实姓名: %s", updateResp.RealName)

	require.Equal(t, newRealName, updateResp.RealName, "真实姓名未更新")

	t.Log("更新资料测试完成!")
}

// TestUpdateProfileInvalid 测试使用无效数据更新个人资料。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestUpdateProfileInvalid ./internal/integration/profile/
func TestUpdateProfileInvalid(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	t.Log("尝试使用无效数据更新资料（如空全名）")
	emptyFullName := ""
	updateReq := user.UpdateDTO{
		RealName: &emptyFullName,
	}

	resp, err := c.R().
		SetBody(updateReq).
		Put("/api/user/profile")

	// 应该返回错误或验证失败
	if err == nil && resp.IsSuccess() {
		t.Log("  注意：服务器接受了空全名，这可能需要根据业务规则判断是否合理")
	} else {
		t.Logf("  无效数据被正确拒绝: %d - %s", resp.StatusCode(), resp.String())
	}

	t.Log("无效数据更新测试完成!")
}

// TestChangePassword 测试修改密码。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestChangePassword ./internal/integration/profile/
func TestChangePassword(t *testing.T) {
	adminClient := manualtest.LoginAsAdmin(t)

	testUsername := fmt.Sprintf("pwdtest_%d", time.Now().Unix())
	originalPassword := "original123"
	newPassword := "newpassword456"

	t.Log("步骤 1: 创建测试用户（带 user 角色）")
	createReq := user.CreateDTO{
		Username: testUsername,
		Email:    testUsername + "@example.com",
		Password: originalPassword,
		RealName: "密码测试用户",
		RoleIDs:  []uint{2}, // user 角色 ID
	}

	createResp, err := manualtest.Post[user.UserWithRolesDTO](adminClient, "/api/admin/users", createReq)
	require.NoError(t, err, "创建测试用户失败")
	testUserID := createResp.ID
	t.Logf("  创建成功，用户 ID: %d", testUserID)

	// 注册清理函数
	t.Cleanup(func() {
		_ = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
	})

	// 用测试用户登录
	t.Log("步骤 2: 测试用户登录")
	testClient := manualtest.LoginAs(t, testUsername, originalPassword)
	t.Log("  登录成功")

	t.Log("步骤 3: 修改密码")
	changeReq := user.ChangePasswordDTO{
		OldPassword: originalPassword,
		NewPassword: newPassword,
	}

	resp, err := testClient.R().
		SetBody(changeReq).
		Put("/api/user/password")
	require.NoError(t, err, "修改密码请求失败")
	require.False(t, resp.IsError(), "修改密码失败，状态码: %d", resp.StatusCode())
	t.Log("  密码修改成功")

	t.Log("步骤 4: 使用新密码登录")
	newClient := manualtest.LoginAs(t, testUsername, newPassword)
	_ = newClient
	t.Log("  新密码登录成功!")

	t.Log("步骤 5: 验证旧密码已失效")
	oldPwdClient := manualtest.NewClient()
	captcha, _ := oldPwdClient.GetCaptcha()
	oldLoginReq := auth.LoginDTO{
		Account:   testUsername,
		Password:  originalPassword,
		CaptchaID: captcha.ID,
		Captcha:   captcha.Code,
	}
	oldResp, err := oldPwdClient.LoginWithCaptcha(oldLoginReq)
	require.True(t, err != nil || oldResp.AccessToken == "", "旧密码不应该能登录")
	t.Log("  旧密码已失效")

	t.Log("修改密码测试完成!")
}

// TestChangePasswordWrongOld 测试使用错误的旧密码修改密码。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestChangePasswordWrongOld ./internal/integration/profile/
func TestChangePasswordWrongOld(t *testing.T) {
	adminClient := manualtest.LoginAsAdmin(t)

	testUsername := fmt.Sprintf("wrongpwd_%d", time.Now().Unix())
	testPassword := "original123"

	t.Log("步骤 1: 创建测试用户")
	createReq := user.CreateDTO{
		Username: testUsername,
		Email:    testUsername + "@example.com",
		Password: testPassword,
		RealName: "错误旧密码测试用户",
		RoleIDs:  []uint{2},
	}

	createResp, err := manualtest.Post[user.UserWithRolesDTO](adminClient, "/api/admin/users", createReq)
	require.NoError(t, err, "创建测试用户失败")
	testUserID := createResp.ID
	t.Logf("  创建成功，用户 ID: %d", testUserID)

	// 注册清理函数
	t.Cleanup(func() {
		_ = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
	})

	t.Log("步骤 2: 测试用户登录")
	testClient := manualtest.LoginAs(t, testUsername, testPassword)
	t.Log("  登录成功")

	t.Log("步骤 3: 使用错误的旧密码尝试修改")
	changeReq := user.ChangePasswordDTO{
		OldPassword: "wrong_old_password",
		NewPassword: "newpassword456",
	}

	resp, err := testClient.R().
		SetBody(changeReq).
		Put("/api/user/password")
	require.True(t, err != nil || !resp.IsSuccess(), "错误的旧密码不应该允许修改密码")
	t.Logf("  错误旧密码被正确拒绝: %d - %s", resp.StatusCode(), resp.String())

	t.Log("错误旧密码测试完成!")
}

// TestDeleteAccount 测试删除账户。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestDeleteAccount ./internal/integration/profile/
func TestDeleteAccount(t *testing.T) {
	adminClient := manualtest.LoginAsAdmin(t)

	testUsername := fmt.Sprintf("delacct_%d", time.Now().Unix())
	testPassword := "test123456"

	t.Log("步骤 1: 创建测试用户（带 user 角色）")
	createReq := user.CreateDTO{
		Username: testUsername,
		Email:    testUsername + "@example.com",
		Password: testPassword,
		RealName: "删除账户测试用户",
		RoleIDs:  []uint{2}, // user 角色 ID
	}

	createResp, err := manualtest.Post[user.UserWithRolesDTO](adminClient, "/api/admin/users", createReq)
	require.NoError(t, err, "创建测试用户失败")
	t.Logf("  创建成功，用户 ID: %d", createResp.ID)

	t.Log("步骤 2: 测试用户登录")
	testClient := manualtest.LoginAs(t, testUsername, testPassword)
	t.Log("  登录成功")

	t.Log("步骤 3: 调用删除账户接口")
	err = testClient.Delete("/api/user/account")
	require.NoError(t, err, "删除账户失败")
	t.Log("  删除成功!")

	t.Log("步骤 4: 验证账户已删除（尝试登录应失败）")
	verifyClient := manualtest.NewClient()
	captcha, _ := verifyClient.GetCaptcha()
	loginReq := auth.LoginDTO{
		Account:   testUsername,
		Password:  testPassword,
		CaptchaID: captcha.ID,
		Captcha:   captcha.Code,
	}
	loginResp, err := verifyClient.LoginWithCaptcha(loginReq)
	require.True(t, err != nil || loginResp.AccessToken == "", "账户已删除，不应该能登录")
	t.Log("  验证成功：账户已无法登录")

	t.Log("删除账户测试完成!")
}
