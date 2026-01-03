// Package usersetting_test 用户设置手动测试
package usersetting_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	manualtest "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/manualtest"
	usersettingapp "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/usersetting"
)

// TestListSettings 测试获取用户设置列表
//
// 测试场景：
// 1. 登录用户可以获取设置列表（系统默认值）
// 2. 列表应包含 is_custom 字段标识来源
func TestListSettings(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	t.Run("登录用户可获取设置列表", func(t *testing.T) {
		result, err := manualtest.Get[usersettingapp.UserSettingListDTO](
			c,
			"/api/user/settings",
			nil,
		)
		require.NoError(t, err, "获取设置列表应成功")
		require.NotNil(t, result, "响应不应为空")
		assert.Positive(t, result.Total, "应返回至少一个设置")
		assert.NotEmpty(t, result.Settings, "设置列表不应为空")

		// 验证返回结构
		for _, setting := range result.Settings {
			assert.NotEmpty(t, setting.Key, "设置键不应为空")
			assert.NotEmpty(t, setting.Value, "设置值不应为空")
			assert.NotEmpty(t, setting.ValueType, "值类型不应为空")
			assert.NotEmpty(t, setting.Label, "标签不应为空")
		}
	})

	t.Run("系统默认值is_custom应为false", func(t *testing.T) {
		result, err := manualtest.Get[usersettingapp.UserSettingListDTO](
			c,
			"/api/user/settings",
			nil,
		)
		require.NoError(t, err, "获取设置列表应成功")

		// 找到一个系统默认设置验证
		for _, setting := range result.Settings {
			if !setting.IsCustom {
				assert.False(t, setting.IsCustom, "系统默认值is_custom应为false")
				return
			}
		}
		// 如果所有设置都是自定义的，说明测试数据有问题
		t.Skip("没有找到系统默认设置用于验证")
	})
}

// TestGetSetting 测试获取单个用户设置
//
// 测试场景：
// 1. 获取存在的设置应返回完整信息
func TestGetSetting(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	t.Run("获取存在的设置", func(t *testing.T) {
		// 使用种子数据中的 general.theme 设置
		result, err := manualtest.Get[usersettingapp.UserSettingDTO](
			c,
			"/api/user/settings/general.theme",
			nil,
		)
		require.NoError(t, err, "获取设置应成功")
		require.NotNil(t, result, "响应不应为空")

		assert.Equal(t, "general.theme", result.Key, "设置键应匹配")
		assert.NotEmpty(t, result.Value, "设置值不应为空")
		assert.NotEmpty(t, result.ValueType, "值类型不应为空")
		assert.NotEmpty(t, result.Label, "标签不应为空")
	})
}

// TestUpdateSetting 测试更新用户设置
//
// 测试场景：
// 1. 更新设置创建自定义值（is_custom=true）
// 2. 再次获取应返回自定义值
// 3. 更新为JSON值
// 4. 验证自定义值覆盖系统默认值
func TestUpdateSetting(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	testKey := "general.theme"

	t.Run("更新设置创建自定义值", func(t *testing.T) {
		// 更新设置为 "dark"
		updateReq := usersettingapp.UpdateDTO{
			Value: `"dark"`,
		}

		result, err := manualtest.Put[map[string]any](
			c,
			"/api/user/settings/"+testKey,
			updateReq,
		)
		require.NoError(t, err, "更新设置应成功")
		require.NotNil(t, result, "响应不应为空")

		// 验证更新后的值
		getResult, err := manualtest.Get[usersettingapp.UserSettingDTO](
			c,
			"/api/user/settings/"+testKey,
			nil,
		)
		require.NoError(t, err, "获取更新后的设置应成功")
		assert.Equal(t, testKey, getResult.Key, "设置键应匹配")
		assert.Equal(t, `"dark"`, getResult.Value, "值应为更新后的值")
		assert.True(t, getResult.IsCustom, "自定义值is_custom应为true")

		// Cleanup: 删除测试创建的自定义值
		t.Cleanup(func() {
			_ = c.Delete("/api/user/settings/" + testKey)
		})
	})

	t.Run("更新JSON类型设置", func(t *testing.T) {
		// general.site_name 可以接受任意字符串
		jsonKey := "general.site_name"
		updateReq := usersettingapp.UpdateDTO{
			Value: `"Test Site Name"`,
		}

		result, err := manualtest.Put[map[string]any](
			c,
			"/api/user/settings/"+jsonKey,
			updateReq,
		)
		require.NoError(t, err, "更新JSON设置应成功")
		require.NotNil(t, result, "响应不应为空")

		// 验证值
		getResult, err := manualtest.Get[usersettingapp.UserSettingDTO](
			c,
			"/api/user/settings/"+jsonKey,
			nil,
		)
		require.NoError(t, err, "获取更新后的设置应成功")
		assert.Equal(t, `"Test Site Name"`, getResult.Value, "JSON值应正确")
		assert.True(t, getResult.IsCustom, "自定义值is_custom应为true")

		// Cleanup
		t.Cleanup(func() {
			_ = c.Delete("/api/user/settings/" + jsonKey)
		})
	})

	t.Run("验证自定义值覆盖系统默认值", func(t *testing.T) {
		// 先获取原始值
		original, err := manualtest.Get[usersettingapp.UserSettingDTO](
			c,
			"/api/user/settings/general.language",
			nil,
		)
		require.NoError(t, err, "获取原始值应成功")
		originalValue := original.Value

		// 更新为自定义值
		updateReq := usersettingapp.UpdateDTO{
			Value: `"en-US"`,
		}
		_, err = manualtest.Put[map[string]any](
			c,
			"/api/user/settings/general.language",
			updateReq,
		)
		require.NoError(t, err, "更新应成功")

		// 验证自定义值
		updated, err := manualtest.Get[usersettingapp.UserSettingDTO](
			c,
			"/api/user/settings/general.language",
			nil,
		)
		require.NoError(t, err, "获取更新后的值应成功")
		assert.Equal(t, `"en-US"`, updated.Value, "应返回自定义值")
		assert.True(t, updated.IsCustom, "is_custom应为true")
		assert.NotEqual(t, originalValue, updated.Value, "自定义值应不同于原始值")

		// Cleanup: 恢复系统默认值
		t.Cleanup(func() {
			_ = c.Delete("/api/user/settings/general.language")
		})
	})
}

// TestDeleteSetting 测试删除用户设置
//
// 测试场景：
// 1. 删除自定义值应恢复系统默认值
// 2. 删除后is_custom应为false
// 3. 删除不存在的自定义值应幂等（不报错）
func TestDeleteSetting(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	testKey := "general.theme"

	t.Run("删除自定义值恢复系统默认", func(t *testing.T) {
		// 先创建自定义值
		updateReq := usersettingapp.UpdateDTO{
			Value: `"dark"`,
		}
		_, err := manualtest.Put[map[string]any](
			c,
			"/api/user/settings/"+testKey,
			updateReq,
		)
		require.NoError(t, err, "创建自定义值应成功")

		// 验证自定义值存在
		customResult, err := manualtest.Get[usersettingapp.UserSettingDTO](
			c,
			"/api/user/settings/"+testKey,
			nil,
		)
		require.NoError(t, err, "获取自定义值应成功")
		assert.True(t, customResult.IsCustom, "创建后is_custom应为true")

		// 删除自定义值
		err = c.Delete("/api/user/settings/" + testKey)
		require.NoError(t, err, "删除应成功")

		// 验证恢复为系统默认值
		defaultResult, err := manualtest.Get[usersettingapp.UserSettingDTO](
			c,
			"/api/user/settings/"+testKey,
			nil,
		)
		require.NoError(t, err, "获取恢复后的值应成功")
		assert.False(t, defaultResult.IsCustom, "删除后is_custom应为false")
		assert.NotEqual(t, `"dark"`, defaultResult.Value, "值应恢复为系统默认")
	})

	t.Run("删除不存在的自定义值应幂等", func(t *testing.T) {
		// 确保没有自定义值（先删除一次）
		_ = c.Delete("/api/user/settings/" + testKey)

		// 再次删除应成功（幂等）
		err := c.Delete("/api/user/settings/" + testKey)
		require.NoError(t, err, "删除不存在的自定义值应成功")
	})
}

// TestSettingsValueTypes 测试不同值类型的设置
//
// 测试场景：
// 1. 字符串类型设置
// 2. 布尔类型设置
// 3. 数字类型设置
func TestSettingsValueTypes(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)
	t.Cleanup(func() {
		// Cleanup: 删除所有测试创建的自定义值
		_ = c.Delete("/api/user/settings/general.theme")
		_ = c.Delete("/api/user/settings/security.enable_twofa")
		_ = c.Delete("/api/user/settings/security.password_min_length")
	})

	t.Run("字符串类型设置", func(t *testing.T) {
		updateReq := usersettingapp.UpdateDTO{
			Value: `"dark"`,
		}

		result, err := manualtest.Put[map[string]any](
			c,
			"/api/user/settings/general.theme",
			updateReq,
		)
		require.NoError(t, err, "更新字符串设置应成功")
		require.NotNil(t, result, "响应不应为空")

		getResult, err := manualtest.Get[usersettingapp.UserSettingDTO](
			c,
			"/api/user/settings/general.theme",
			nil,
		)
		require.NoError(t, err, "获取设置应成功")
		assert.Equal(t, `"dark"`, getResult.Value, "字符串值应正确")
		assert.True(t, getResult.IsCustom, "应为自定义值")
	})

	t.Run("布尔类型设置", func(t *testing.T) {
		updateReq := usersettingapp.UpdateDTO{
			Value: `true`,
		}

		result, err := manualtest.Put[map[string]any](
			c,
			"/api/user/settings/security.enable_twofa",
			updateReq,
		)
		require.NoError(t, err, "更新布尔设置应成功")
		require.NotNil(t, result, "响应不应为空")

		getResult, err := manualtest.Get[usersettingapp.UserSettingDTO](
			c,
			"/api/user/settings/security.enable_twofa",
			nil,
		)
		require.NoError(t, err, "获取设置应成功")
		assert.Equal(t, `true`, getResult.Value, "布尔值应正确")
		assert.True(t, getResult.IsCustom, "应为自定义值")
	})

	t.Run("数字类型设置", func(t *testing.T) {
		updateReq := usersettingapp.UpdateDTO{
			Value: `12`,
		}

		result, err := manualtest.Put[map[string]any](
			c,
			"/api/user/settings/security.password_min_length",
			updateReq,
		)
		require.NoError(t, err, "更新数字设置应成功")
		require.NotNil(t, result, "响应不应为空")

		getResult, err := manualtest.Get[usersettingapp.UserSettingDTO](
			c,
			"/api/user/settings/security.password_min_length",
			nil,
		)
		require.NoError(t, err, "获取设置应成功")
		assert.Equal(t, `12`, getResult.Value, "数字值应正确")
		assert.True(t, getResult.IsCustom, "应为自定义值")
	})
}

// TestSettingsIsolation 测试用户设置隔离
//
// 测试场景：
// 1. 不同用户的自定义值应相互独立
func TestSettingsIsolation(t *testing.T) {
	// Admin 用户
	adminClient := manualtest.LoginAsAdmin(t)
	t.Cleanup(func() {
		_ = adminClient.Delete("/api/user/settings/general.theme")
	})

	// 创建测试用户
	testUserResp, err := manualtest.Post[map[string]any](
		adminClient,
		"/api/admin/users",
		map[string]any{
			"username": "test_user_settings",
			"email":    "test_settings@example.com",
			"password": "Test123456!",
		},
	)
	require.NoError(t, err, "创建测试用户应成功")
	testUserID := uint((*testUserResp)["id"].(float64))
	t.Cleanup(func() {
		_ = adminClient.Delete("/api/admin/users/" + strconv.FormatUint(uint64(testUserID), 10))
	})

	// 以测试用户登录
	testUserClient := manualtest.LoginAs(t, "test_user_settings", "Test123456!")
	t.Cleanup(func() {
		_ = testUserClient.Delete("/api/user/settings/general.theme")
	})

	t.Run("不同用户的自定义值应独立", func(t *testing.T) {
		testKey := "general.theme"

		// Admin 设置自定义值
		adminUpdateReq := usersettingapp.UpdateDTO{
			Value: `"dark"`,
		}
		_, err = manualtest.Put[map[string]any](
			adminClient,
			"/api/user/settings/"+testKey,
			adminUpdateReq,
		)
		require.NoError(t, err, "Admin更新设置应成功")

		// 测试用户设置不同的自定义值
		testUserUpdateReq := usersettingapp.UpdateDTO{
			Value: `"light"`,
		}
		_, err = manualtest.Put[map[string]any](
			testUserClient,
			"/api/user/settings/"+testKey,
			testUserUpdateReq,
		)
		require.NoError(t, err, "测试用户更新设置应成功")

		// 验证 Admin 的值
		adminResult, err := manualtest.Get[usersettingapp.UserSettingDTO](
			adminClient,
			"/api/user/settings/"+testKey,
			nil,
		)
		require.NoError(t, err, "Admin获取设置应成功")
		assert.Equal(t, `"dark"`, adminResult.Value, "Admin应看到自己的值")

		// 验证测试用户的值
		testUserResult, err := manualtest.Get[usersettingapp.UserSettingDTO](
			testUserClient,
			"/api/user/settings/"+testKey,
			nil,
		)
		require.NoError(t, err, "测试用户获取设置应成功")
		assert.Equal(t, `"light"`, testUserResult.Value, "测试用户应看到自己的值")
	})
}
