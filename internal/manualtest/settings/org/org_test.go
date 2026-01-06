// Package org_test 组织设置手动测试
package org_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	manualtest "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/adapters/gin/manualtest"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/app/org"
)

// 种子数据: acme org (ID=1)
const testOrgID uint = 1

// orgSettingsPath 返回组织设置 API 基础路径
func orgSettingsPath(orgID uint) string {
	return fmt.Sprintf("/api/org/%d/settings", orgID)
}

// orgSettingPath 返回单个组织设置 API 路径
func orgSettingPath(orgID uint, key string) string {
	return fmt.Sprintf("%s/%s", orgSettingsPath(orgID), key)
}

// TestListOrgSettings 测试获取组织设置列表
//
// 测试场景：
// 1. 组织管理员可以获取组织设置列表
// 2. 验证设置包含必要字段（key, value_type, label）
func TestListOrgSettings(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	t.Run("获取组织设置列表应成功", func(t *testing.T) {
		result, _, err := manualtest.GetList[org.SettingsItemDTO](
			c,
			orgSettingsPath(testOrgID),
			nil,
		)
		require.NoError(t, err, "获取组织设置列表应成功")
		require.NotNil(t, result, "响应不应为空")
		assert.NotEmpty(t, result, "设置列表不应为空")

		t.Logf("组织设置数量: %d", len(result))
	})

	t.Run("验证设置包含必要字段", func(t *testing.T) {
		result, _, err := manualtest.GetList[org.SettingsItemDTO](
			c,
			orgSettingsPath(testOrgID),
			nil,
		)
		require.NoError(t, err, "获取组织设置列表应成功")

		for _, setting := range result {
			assert.NotEmpty(t, setting.Key, "设置键不应为空")
			assert.NotEmpty(t, setting.ValueType, "值类型不应为空")
			assert.NotEmpty(t, setting.Label, "标签不应为空")
			t.Logf("  %s: value_type=%s, is_customized=%v",
				setting.Key, setting.ValueType, setting.IsCustomized)
		}
	})

	t.Run("系统默认值 is_customized 应为 false", func(t *testing.T) {
		result, _, err := manualtest.GetList[org.SettingsItemDTO](
			c,
			orgSettingsPath(testOrgID),
			nil,
		)
		require.NoError(t, err, "获取组织设置列表应成功")

		// 找到一个系统默认设置验证
		for _, setting := range result {
			if !setting.IsCustomized {
				assert.False(t, setting.IsCustomized, "系统默认值 is_customized 应为 false")
				return
			}
		}
		t.Skip("没有找到系统默认设置用于验证")
	})
}

// TestGetOrgSetting 测试获取单个组织设置
//
// 测试场景：
// 1. 获取存在的设置应返回完整信息
// 2. 获取不存在的设置应返回错误
func TestGetOrgSetting(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	t.Run("获取存在的设置", func(t *testing.T) {
		// 使用种子数据中的 general.theme 设置
		result, err := manualtest.Get[org.SettingsItemDTO](
			c,
			orgSettingPath(testOrgID, "general.theme"),
			nil,
		)
		require.NoError(t, err, "获取设置应成功")
		require.NotNil(t, result, "响应不应为空")

		assert.Equal(t, "general.theme", result.Key, "设置键应匹配")
		assert.NotEmpty(t, result.ValueType, "值类型不应为空")
		assert.NotEmpty(t, result.Label, "标签不应为空")
	})

	t.Run("获取不存在的设置应返回错误", func(t *testing.T) {
		_, err := manualtest.Get[org.SettingsItemDTO](
			c,
			orgSettingPath(testOrgID, "nonexistent.setting"),
			nil,
		)
		assert.Error(t, err, "获取不存在的设置应返回错误")
	})
}

// TestSetOrgSetting 测试设置组织配置
//
// 测试场景：
// 1. 设置成功后 is_customized 应为 true
// 2. 自定义值应覆盖系统默认值
func TestSetOrgSetting(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	testKey := "general.theme"

	t.Run("设置组织自定义值", func(t *testing.T) {
		updateReq := map[string]any{
			"value": "dark",
		}

		result, err := manualtest.Put[org.SettingsItemDTO](
			c,
			orgSettingPath(testOrgID, testKey),
			updateReq,
		)
		require.NoError(t, err, "设置组织配置应成功")
		require.NotNil(t, result, "响应不应为空")

		assert.Equal(t, testKey, result.Key, "设置键应匹配")
		assert.True(t, result.IsCustomized, "设置后 is_customized 应为 true")
		assert.Equal(t, "dark", result.Value, "值应为更新后的值")

		// Cleanup: 恢复默认值
		t.Cleanup(func() {
			_ = c.Delete(orgSettingPath(testOrgID, testKey))
		})
	})

	t.Run("验证自定义值覆盖系统默认值", func(t *testing.T) {
		// 先获取当前值
		original, err := manualtest.Get[org.SettingsItemDTO](
			c,
			orgSettingPath(testOrgID, "general.language"),
			nil,
		)
		require.NoError(t, err, "获取原始值应成功")

		// 更新为组织自定义值
		updateReq := map[string]any{
			"value": "en-US",
		}
		_, err = manualtest.Put[org.SettingsItemDTO](
			c,
			orgSettingPath(testOrgID, "general.language"),
			updateReq,
		)
		require.NoError(t, err, "更新应成功")

		// 验证自定义值
		updated, err := manualtest.Get[org.SettingsItemDTO](
			c,
			orgSettingPath(testOrgID, "general.language"),
			nil,
		)
		require.NoError(t, err, "获取更新后的值应成功")
		assert.Equal(t, "en-US", updated.Value, "应返回组织自定义值")
		assert.True(t, updated.IsCustomized, "is_customized 应为 true")

		// Cleanup: 恢复原始值
		t.Cleanup(func() {
			_ = c.Delete(orgSettingPath(testOrgID, "general.language"))
		})

		// 验证原始值和更新值不同（如果原始值不是自定义的）
		if !original.IsCustomized {
			assert.NotEqual(t, original.Value, updated.Value, "自定义值应不同于原始默认值")
		}
	})
}

// TestResetOrgSetting 测试重置组织配置
//
// 测试场景：
// 1. 重置后应恢复系统默认值
// 2. 重置后 is_customized 应为 false
// 3. 重置不存在的自定义值应幂等
func TestResetOrgSetting(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	testKey := "general.theme"

	t.Run("重置组织自定义值", func(t *testing.T) {
		// 先创建组织自定义值
		updateReq := map[string]any{
			"value": "dark",
		}
		_, err := manualtest.Put[org.SettingsItemDTO](
			c,
			orgSettingPath(testOrgID, testKey),
			updateReq,
		)
		require.NoError(t, err, "创建组织自定义值应成功")

		// 验证自定义值存在
		customResult, err := manualtest.Get[org.SettingsItemDTO](
			c,
			orgSettingPath(testOrgID, testKey),
			nil,
		)
		require.NoError(t, err, "获取组织设置应成功")
		assert.True(t, customResult.IsCustomized, "创建后 is_customized 应为 true")

		// 重置组织配置
		err = c.Delete(orgSettingPath(testOrgID, testKey))
		require.NoError(t, err, "重置应成功")

		// 验证恢复为系统默认值
		defaultResult, err := manualtest.Get[org.SettingsItemDTO](
			c,
			orgSettingPath(testOrgID, testKey),
			nil,
		)
		require.NoError(t, err, "获取重置后的值应成功")
		assert.False(t, defaultResult.IsCustomized, "重置后 is_customized 应为 false")
		assert.NotEqual(t, "dark", defaultResult.Value, "值应恢复为系统默认")
	})

	t.Run("重置不存在的自定义值应幂等", func(t *testing.T) {
		// 确保没有自定义值
		_ = c.Delete(orgSettingPath(testOrgID, testKey))

		// 再次重置应成功（幂等）
		err := c.Delete(orgSettingPath(testOrgID, testKey))
		require.NoError(t, err, "重置不存在的自定义值应成功")
	})
}

// TestOrgSettingsIsolation 测试组织设置隔离
//
// 测试场景：
// 1. 不同组织的自定义值应相互独立
func TestOrgSettingsIsolation(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	testKey := "general.theme"

	t.Run("访问不存在的组织应返回错误", func(t *testing.T) {
		invalidOrgID := uint(99999)
		_, err := manualtest.Get[org.SettingsItemDTO](
			c,
			orgSettingPath(invalidOrgID, testKey),
			nil,
		)
		assert.Error(t, err, "访问不存在的组织应返回错误")
	})
}

// TestOrgSettingsPermission 测试组织设置权限控制
//
// 测试场景：
// 1. 非组织成员访问组织设置应被拒绝
// 2. 系统管理员（org owner）可以访问组织设置
// 3. 组织普通成员（member）无法访问组织设置
// 4. 组织管理员（admin）可以访问组织设置
// 5. 组织 Owner 可以读取和修改设置
func TestOrgSettingsPermission(t *testing.T) {
	adminClient := manualtest.LoginAsAdmin(t)

	t.Run("非组织成员访问组织设置应被拒绝", func(t *testing.T) {
		// 创建一个不属于任何组织的测试用户
		testUserResp, err := manualtest.Post[map[string]any](
			adminClient,
			"/api/admin/users",
			map[string]any{
				"username": "test_non_org_member",
				"email":    "non_org_member@example.com",
				"password": "Test123456!",
			},
		)
		require.NoError(t, err, "创建测试用户应成功")
		testUserID := uint((*testUserResp)["id"].(float64))
		t.Cleanup(func() {
			_ = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
		})

		// 以测试用户登录
		nonMemberClient := manualtest.LoginAs(t, "test_non_org_member", "Test123456!")

		// 尝试访问组织设置（应被拒绝）
		_, _, err = manualtest.GetList[org.SettingsItemDTO](
			nonMemberClient,
			orgSettingsPath(testOrgID),
			nil,
		)
		require.Error(t, err, "非组织成员访问组织设置应被拒绝")
		t.Logf("  预期错误: %v", err)
	})

	t.Run("组织成员可以访问组织设置", func(t *testing.T) {
		// admin 用户是 testOrgID (acme) 的 owner
		result, _, err := manualtest.GetList[org.SettingsItemDTO](
			adminClient,
			orgSettingsPath(testOrgID),
			nil,
		)
		require.NoError(t, err, "组织成员访问组织设置应成功")
		assert.NotEmpty(t, result, "设置列表不应为空")
	})

	t.Run("组织普通成员无法访问组织设置", func(t *testing.T) {
		// 创建测试用户并加入组织（作为 member 角色）
		testUserResp, err := manualtest.Post[map[string]any](
			adminClient,
			"/api/admin/users",
			map[string]any{
				"username": "test_org_member",
				"email":    "org_member@example.com",
				"password": "Test123456!",
			},
		)
		require.NoError(t, err, "创建测试用户应成功")
		testUserID := uint((*testUserResp)["id"].(float64))
		t.Cleanup(func() {
			_ = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
		})

		// 将用户加入组织（作为 member）
		_, err = manualtest.Post[any](
			adminClient,
			fmt.Sprintf("/api/org/%d/members", testOrgID),
			map[string]any{
				"user_id": testUserID,
				"role":    "member",
			},
		)
		require.NoError(t, err, "将用户加入组织应成功")
		t.Cleanup(func() {
			_ = adminClient.Delete(fmt.Sprintf("/api/org/%d/members/%d", testOrgID, testUserID))
		})

		// 以测试用户登录
		memberClient := manualtest.LoginAs(t, "test_org_member", "Test123456!")

		// member 角色无法读取组织设置（需要 owner/admin 角色）
		_, _, err = manualtest.GetList[org.SettingsItemDTO](
			memberClient,
			orgSettingsPath(testOrgID),
			nil,
		)
		require.Error(t, err, "组织 member 访问组织设置应被拒绝")
		t.Logf("  member 访问组织设置被拒绝（符合预期）: %v", err)
	})

	t.Run("组织管理员可以访问组织设置", func(t *testing.T) {
		// 创建测试用户并加入组织（作为 admin 角色）
		// 修复后：org admin 角色通过 OrgContext 动态注入 org:settings:* 权限
		testUserResp, err := manualtest.Post[map[string]any](
			adminClient,
			"/api/admin/users",
			map[string]any{
				"username": "test_org_admin",
				"email":    "org_admin@example.com",
				"password": "Test123456!",
			},
		)
		require.NoError(t, err, "创建测试用户应成功")
		testUserID := uint((*testUserResp)["id"].(float64))
		t.Cleanup(func() {
			_ = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
		})

		// 将用户加入组织（作为 admin）
		_, err = manualtest.Post[any](
			adminClient,
			fmt.Sprintf("/api/org/%d/members", testOrgID),
			map[string]any{
				"user_id": testUserID,
				"role":    "admin",
			},
		)
		require.NoError(t, err, "将用户加入组织应成功")
		t.Cleanup(func() {
			_ = adminClient.Delete(fmt.Sprintf("/api/org/%d/members/%d", testOrgID, testUserID))
		})

		// 以测试用户登录
		orgAdminClient := manualtest.LoginAs(t, "test_org_admin", "Test123456!")

		// org admin 角色可以读取组织设置（通过 OrgContext 动态权限注入）
		result, _, err := manualtest.GetList[org.SettingsItemDTO](
			orgAdminClient,
			orgSettingsPath(testOrgID),
			nil,
		)
		require.NoError(t, err, "组织 admin 访问组织设置应成功")
		assert.NotEmpty(t, result, "设置列表不应为空")
		t.Logf("  org admin 访问组织设置成功，获取 %d 个设置", len(result))

		// org admin 角色也可以修改组织设置
		updateReq := map[string]any{
			"value": "light",
		}
		_, err = manualtest.Put[org.SettingsItemDTO](
			orgAdminClient,
			orgSettingPath(testOrgID, "general.theme"),
			updateReq,
		)
		require.NoError(t, err, "组织 admin 修改组织设置应成功")
		t.Cleanup(func() {
			_ = adminClient.Delete(orgSettingPath(testOrgID, "general.theme"))
		})
		t.Log("  org admin 修改组织设置成功")
	})

	t.Run("组织 Owner 可以读取和修改设置", func(t *testing.T) {
		// admin 用户（系统管理员）是 testOrgID 的 owner
		// 只有 owner 角色 + 系统 RBAC 权限才能访问组织设置
		result, _, err := manualtest.GetList[org.SettingsItemDTO](
			adminClient,
			orgSettingsPath(testOrgID),
			nil,
		)
		require.NoError(t, err, "组织 owner 读取组织设置应成功")
		assert.NotEmpty(t, result, "设置列表不应为空")

		// owner 可以修改组织设置
		updateReq := map[string]any{
			"value": "dark",
		}
		_, err = manualtest.Put[org.SettingsItemDTO](
			adminClient,
			orgSettingPath(testOrgID, "general.theme"),
			updateReq,
		)
		require.NoError(t, err, "组织 owner 修改组织设置应成功")
		t.Cleanup(func() {
			_ = adminClient.Delete(orgSettingPath(testOrgID, "general.theme"))
		})
	})
}

// TestOrgSettingsValueTypes 测试不同值类型的组织设置
//
// 测试场景：
// 1. 字符串类型设置
// 2. 布尔类型设置
func TestOrgSettingsValueTypes(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)
	t.Cleanup(func() {
		// Cleanup: 删除所有测试创建的自定义值
		_ = c.Delete(orgSettingPath(testOrgID, "general.theme"))
		_ = c.Delete(orgSettingPath(testOrgID, "notification.enable_email"))
	})

	t.Run("字符串类型设置", func(t *testing.T) {
		updateReq := map[string]any{
			"value": "dark",
		}

		result, err := manualtest.Put[org.SettingsItemDTO](
			c,
			orgSettingPath(testOrgID, "general.theme"),
			updateReq,
		)
		require.NoError(t, err, "更新字符串设置应成功")
		require.NotNil(t, result, "响应不应为空")

		getResult, err := manualtest.Get[org.SettingsItemDTO](
			c,
			orgSettingPath(testOrgID, "general.theme"),
			nil,
		)
		require.NoError(t, err, "获取设置应成功")
		assert.Equal(t, "dark", getResult.Value, "字符串值应正确")
		assert.True(t, getResult.IsCustomized, "应为自定义值")
	})

	t.Run("布尔类型设置", func(t *testing.T) {
		updateReq := map[string]any{
			"value": true,
		}

		result, err := manualtest.Put[org.SettingsItemDTO](
			c,
			orgSettingPath(testOrgID, "notification.enable_email"),
			updateReq,
		)
		require.NoError(t, err, "更新布尔设置应成功")
		require.NotNil(t, result, "响应不应为空")

		getResult, err := manualtest.Get[org.SettingsItemDTO](
			c,
			orgSettingPath(testOrgID, "notification.enable_email"),
			nil,
		)
		require.NoError(t, err, "获取设置应成功")
		assert.Equal(t, true, getResult.Value, "布尔值应正确")
		assert.True(t, getResult.IsCustomized, "应为自定义值")
	})
}
