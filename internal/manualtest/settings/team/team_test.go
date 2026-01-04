// Package team_test 团队设置手动测试
package team_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	manualtest "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/adapters/gin/manualtest"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/app/team"
)

// 种子数据: acme org (ID=1), engineering team (ID=1)
const (
	testOrgID  uint = 1
	testTeamID uint = 1
)

// teamSettingsPath 返回团队设置 API 基础路径
func teamSettingsPath(_ /*orgID*/, _ /*teamID*/ uint) string {
	return fmt.Sprintf("/api/org/%d/teams/%d/settings", testOrgID, testTeamID)
}

// teamSettingPath 返回单个团队设置 API 路径
func teamSettingPath(_ /*orgID*/, _ /*teamID*/ uint, key string) string {
	return fmt.Sprintf("%s/%s", teamSettingsPath(testOrgID, testTeamID), key)
}

// TestListTeamSettings 测试获取团队设置列表
//
// 测试场景：
// 1. 团队设置列表应包含所有 Team 可配置的设置
// 2. 验证 VisibleAt/ConfigurableAt 字段存在
// 3. 验证 Team 默认值设置（VisibleAt=user, ConfigurableAt=team）
func TestListTeamSettings(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	t.Run("获取团队设置列表应成功", func(t *testing.T) {
		result, _, err := manualtest.GetList[team.TeamSettingDTO](
			c,
			teamSettingsPath(testOrgID, testTeamID),
			nil,
		)
		require.NoError(t, err, "获取团队设置列表应成功")
		require.NotNil(t, result, "响应不应为空")
		assert.NotEmpty(t, result, "设置列表不应为空")

		t.Logf("团队设置数量: %d", len(result))
	})

	t.Run("验证 VisibleAt 和 ConfigurableAt 字段", func(t *testing.T) {
		result, _, err := manualtest.GetList[team.TeamSettingDTO](
			c,
			teamSettingsPath(testOrgID, testTeamID),
			nil,
		)
		require.NoError(t, err, "获取团队设置列表应成功")

		// 验证每个设置都有 VisibleAt 和 ConfigurableAt 字段
		for _, setting := range result {
			assert.NotEmpty(t, setting.VisibleAt, "VisibleAt 不应为空")
			assert.NotEmpty(t, setting.ConfigurableAt, "ConfigurableAt 不应空")
			t.Logf("  %s: visible_at=%s, configurable_at=%s",
				setting.Key, setting.VisibleAt, setting.ConfigurableAt)
		}
	})

	t.Run("验证 Team 默认值设置（VisibleAt=user, ConfigurableAt=team）", func(t *testing.T) {
		result, _, err := manualtest.GetList[team.TeamSettingDTO](
			c,
			teamSettingsPath(testOrgID, testTeamID),
			nil,
		)
		require.NoError(t, err, "获取团队设置列表应成功")

		// 查找 Team 可为用户设置默认值的配置
		// 例如: general.timezone (visible_at=user, configurable_at=team)
		var teamDefaultSetting *team.TeamSettingDTO
		for _, s := range result {
			if s.Key == "general.timezone" || s.Key == "general.theme" {
				teamDefaultSetting = &s
				break
			}
		}

		require.NotNil(t, teamDefaultSetting, "应找到 Team 默认值设置")
		assert.Equal(t, "user", teamDefaultSetting.VisibleAt, "general.timezone 的 VisibleAt 应为 user")
		assert.Equal(t, "team", teamDefaultSetting.ConfigurableAt, "general.timezone 的 ConfigurableAt 应为 team")
		assert.True(t, teamDefaultSetting.IsTeamDefault, "IsTeamDefault 应为 true")

		t.Logf("  Team 默认值设置: %s, visible_at=%s, configurable_at=%s, is_team_default=%v",
			teamDefaultSetting.Key, teamDefaultSetting.VisibleAt, teamDefaultSetting.ConfigurableAt, teamDefaultSetting.IsTeamDefault)
	})
}

// TestGetTeamSetting 测试获取单个团队设置
//
// 测试场景：
// 1. 获取 Team 可配置的设置应成功
// 2. 验证三级继承逻辑
func TestGetTeamSetting(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	t.Run("获取 Team 可配置的设置", func(t *testing.T) {
		// general.timezone: visible_at=user, configurable_at=team
		result, err := manualtest.Get[team.TeamSettingDTO](
			c,
			teamSettingPath(testOrgID, testTeamID, "general.timezone"),
			nil,
		)
		require.NoError(t, err, "获取团队设置应成功")
		require.NotNil(t, result, "响应不应为空")

		assert.Equal(t, "general.timezone", result.Key, "设置键应匹配")
		assert.Equal(t, "user", result.VisibleAt, "VisibleAt 应为 user")
		assert.Equal(t, "team", result.ConfigurableAt, "ConfigurableAt 应为 team")
		assert.True(t, result.IsTeamDefault, "IsTeamDefault 应为 true")

		t.Logf("  继承来源: %s, 值: %v", result.InheritedFrom, result.Value)
	})
}

// TestSetTeamSetting 测试设置团队配置
//
// 测试场景：
// 1. Team 只能配置 ConfigurableAt >= team 的设置
// 2. 设置成功后 is_customized 应为 true
// 3. 验证三级继承：团队 > 组织 > 系统
func TestSetTeamSetting(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	testKey := "general.timezone" // visible_at=user, configurable_at=team

	t.Run("Team 可为用户设置默认值", func(t *testing.T) {
		// 更新团队默认值
		updateReq := map[string]any{
			"value": "America/New_York",
		}

		result, err := manualtest.Put[team.TeamSettingDTO](
			c,
			teamSettingPath(testOrgID, testTeamID, testKey),
			updateReq,
		)
		require.NoError(t, err, "更新团队设置应成功")
		require.NotNil(t, result, "响应不应为空")

		assert.Equal(t, testKey, result.Key, "设置键应匹配")
		assert.True(t, result.IsCustomized, "设置后 is_customized 应为 true")
		assert.Equal(t, "America/New_York", result.Value, "值应为更新后的值")

		// Cleanup: 恢复默认值
		t.Cleanup(func() {
			_ = c.Delete(teamSettingPath(testOrgID, testTeamID, testKey))
		})
	})

	t.Run("验证自定义值覆盖系统默认值", func(t *testing.T) {
		// 先获取当前值
		original, err := manualtest.Get[team.TeamSettingDTO](
			c,
			teamSettingPath(testOrgID, testTeamID, "general.theme"),
			nil,
		)
		require.NoError(t, err, "获取原始值应成功")
		originalValue := original.Value

		// 更新为团队默认值
		updateReq := map[string]any{
			"value": "light",
		}
		_, err = manualtest.Put[team.TeamSettingDTO](
			c,
			teamSettingPath(testOrgID, testTeamID, "general.theme"),
			updateReq,
		)
		require.NoError(t, err, "更新应成功")

		// 验证团队自定义值
		updated, err := manualtest.Get[team.TeamSettingDTO](
			c,
			teamSettingPath(testOrgID, testTeamID, "general.theme"),
			nil,
		)
		require.NoError(t, err, "获取更新后的值应成功")
		assert.Equal(t, "light", updated.Value, "应返回团队自定义值")
		assert.True(t, updated.IsCustomized, "is_customized 应为 true")
		assert.Equal(t, "team", updated.InheritedFrom, "inherited_from 应为 team")

		// Cleanup: 恢复原始值
		t.Cleanup(func() {
			_ = c.Delete(teamSettingPath(testOrgID, testTeamID, "general.theme"))
		})

		// 恢复原始值（如果有）
		if original.IsCustomized {
			_ = c.Delete(teamSettingPath(testOrgID, testTeamID, "general.theme"))
			if originalValue != nil && originalValue != "system" {
				_ = c.Delete(teamSettingPath(testOrgID, testTeamID, "general.theme"))
			}
		}
	})
}

// TestResetTeamSetting 测试重置团队配置
//
// 测试场景：
// 1. 重置后应恢复组织配置或系统默认值
// 2. 重置后 is_customized 应为 false
func TestResetTeamSetting(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	testKey := "general.theme"

	t.Run("重置团队自定义值", func(t *testing.T) {
		// 先创建团队自定义值（使用不同的值以避免冲突）
		updateReq := map[string]any{
			"value": "light", // 使用 light 而非 dark
		}
		_, err := manualtest.Put[team.TeamSettingDTO](
			c,
			teamSettingPath(testOrgID, testTeamID, testKey),
			updateReq,
		)
		require.NoError(t, err, "创建团队自定义值应成功")

		// 验证自定义值存在
		customResult, err := manualtest.Get[team.TeamSettingDTO](
			c,
			teamSettingPath(testOrgID, testTeamID, testKey),
			nil,
		)
		require.NoError(t, err, "获取团队设置应成功")
		assert.True(t, customResult.IsCustomized, "创建后 is_customized 应为 true")
		assert.Equal(t, "team", customResult.InheritedFrom, "inherited_from 应为 team")

		// 重置团队配置
		err = c.Delete(teamSettingPath(testOrgID, testTeamID, testKey))
		require.NoError(t, err, "重置应成功")

		// 验证恢复为组织配置或系统默认值
		defaultResult, err := manualtest.Get[team.TeamSettingDTO](
			c,
			teamSettingPath(testOrgID, testTeamID, testKey),
			nil,
		)
		require.NoError(t, err, "获取重置后的值应成功")
		assert.False(t, defaultResult.IsCustomized, "重置后 is_customized 应为 false")
		assert.NotEqual(t, "light", defaultResult.Value, "值应恢复为默认值")
		assert.NotEqual(t, "team", defaultResult.InheritedFrom, "inherited_from 不应为 team")

		// Cleanup: 确保没有残留的自定义值
		t.Cleanup(func() {
			_ = c.Delete(teamSettingPath(testOrgID, testTeamID, testKey))
		})
	})
}

// TestTeamSettingsConfigurableOnly 测试 Team 只能配置 ConfigurableAt >= team 的设置
//
// 测试场景：
// 1. Team 可以配置 configurable_at=team 的设置
// 2. Team 不能配置 configurable_at=system 的设置
func TestTeamSettingsConfigurableOnly(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	t.Run("Team 可以配置 configurable_at=team 的设置", func(t *testing.T) {
		// general.theme: visible_at=user, configurable_at=team
		updateReq := map[string]any{
			"value": "dark",
		}

		result, err := manualtest.Put[team.TeamSettingDTO](
			c,
			teamSettingPath(testOrgID, testTeamID, "general.theme"),
			updateReq,
		)
		require.NoError(t, err, "Team 应能配置 configurable_at=team 的设置")
		require.NotNil(t, result, "响应不应为空")

		// Cleanup
		t.Cleanup(func() {
			_ = c.Delete(teamSettingPath(testOrgID, testTeamID, "general.theme"))
		})
	})

	t.Run("Team 不能配置 system 级别设置", func(t *testing.T) {
		// general.site_name: visible_at=system, configurable_at=system
		updateReq := map[string]any{
			"value": "Test Site",
		}

		_, err := manualtest.Put[team.TeamSettingDTO](
			c,
			teamSettingPath(testOrgID, testTeamID, "general.site_name"),
			updateReq,
		)
		// 应该返回错误（设置不允许在 team 级别配置）
		assert.Error(t, err, "Team 不应能配置 system 级别设置")
	})
}

// TestTeamSettingsVisibility 测试设置可见性
//
// 测试场景：
// 1. Team 设置列表只包含 Team 可配置的设置（configurable_at >= team）
// 2. 验证不同 VisibleAt 级别的设置
func TestTeamSettingsVisibility(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	t.Run("验证返回的设置包含正确的 VisibleAt", func(t *testing.T) {
		result, _, err := manualtest.GetList[team.TeamSettingDTO](
			c,
			teamSettingsPath(testOrgID, testTeamID),
			nil,
		)
		require.NoError(t, err, "获取团队设置列表应成功")

		// 验证不同 VisibleAt 级别的设置
		visibleAtLevels := make(map[string]int)
		for _, s := range result {
			visibleAtLevels[s.VisibleAt]++
		}

		t.Logf("  VisibleAt 分布: %v", visibleAtLevels)

		// Team 设置列表应只包含 Team 可配置的设置
		// 当前种子数据中，只有 user 级别的设置允许 Team 配置
		// （如 general.theme, general.timezone 的 configurable_at=team）
		assert.Contains(t, visibleAtLevels, "user", "应包含 user 级别设置（Team 可配置的默认值）")

		// system 级别设置不应出现在列表中，因为 Team 不可配置
		assert.NotContains(t, visibleAtLevels, "system", "不应包含 system 级别设置（Team 不可配置）")
	})
}
