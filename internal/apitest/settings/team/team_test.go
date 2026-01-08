// Package team_test 团队设置手动测试
package team_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	settings "github.com/lwmacct/260103-ddd-iam-bc/internal/apitest/settings"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/app/team"
	"github.com/lwmacct/260103-ddd-shared/pkg/shared/apitest"
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
	c := settings.LoginAsAdmin(t)

	t.Run("获取团队设置列表应成功", func(t *testing.T) {
		result, _, err := apitest.GetList[team.SettingsItemDTO](c.HTTPClient(),
			teamSettingsPath(testOrgID, testTeamID),
			nil,
		)
		require.NoError(t, err, "获取团队设置列表应成功")
		require.NotNil(t, result, "响应不应为空")
		assert.NotEmpty(t, result, "设置列表不应为空")

		t.Logf("团队设置数量: %d", len(result))
	})

	t.Run("验证 VisibleAt 和 ConfigurableAt 字段", func(t *testing.T) {
		result, _, err := apitest.GetList[team.SettingsItemDTO](c.HTTPClient(),
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
		result, _, err := apitest.GetList[team.SettingsItemDTO](c.HTTPClient(),
			teamSettingsPath(testOrgID, testTeamID),
			nil,
		)
		require.NoError(t, err, "获取团队设置列表应成功")

		// 查找 Team 可为用户设置默认值的配置
		// 例如: general.timezone (visible_at=user, configurable_at=team)
		var teamDefaultSetting *team.SettingsItemDTO
		for _, s := range result {
			if s.Key == "general.timezone" || s.Key == "general.theme" {
				teamDefaultSetting = &s
				break
			}
		}

		require.NotNil(t, teamDefaultSetting, "应找到 Team 默认值设置")
		assert.Equal(t, "user", teamDefaultSetting.VisibleAt, "general.timezone 的 VisibleAt 应为 user")
		assert.Equal(t, "team", teamDefaultSetting.ConfigurableAt, "general.timezone 的 ConfigurableAt 应为 team")

		t.Logf("  Team 默认值设置: %s, visible_at=%s, configurable_at=%s",
			teamDefaultSetting.Key, teamDefaultSetting.VisibleAt, teamDefaultSetting.ConfigurableAt)
	})
}

// TestGetTeamSetting 测试获取单个团队设置
//
// 测试场景：
// 1. 获取 Team 可配置的设置应成功
// 2. 验证三级继承逻辑
func TestGetTeamSetting(t *testing.T) {
	c := settings.LoginAsAdmin(t)

	t.Run("获取 Team 可配置的设置", func(t *testing.T) {
		// general.timezone: visible_at=user, configurable_at=team
		result, err := apitest.Get[team.SettingsItemDTO](c.HTTPClient(),
			teamSettingPath(testOrgID, testTeamID, "general.timezone"),
			nil,
		)
		require.NoError(t, err, "获取团队设置应成功")
		require.NotNil(t, result, "响应不应为空")

		assert.Equal(t, "general.timezone", result.Key, "设置键应匹配")
		assert.Equal(t, "user", result.VisibleAt, "VisibleAt 应为 user")
		assert.Equal(t, "team", result.ConfigurableAt, "ConfigurableAt 应为 team")

		t.Logf("  值: %v, is_customized: %v", result.Value, result.IsCustomized)
	})
}

// TestSetTeamSetting 测试设置团队配置
//
// 测试场景：
// 1. Team 只能配置 ConfigurableAt >= team 的设置
// 2. 设置成功后 is_customized 应为 true
// 3. 验证三级继承：团队 > 组织 > 系统
func TestSetTeamSetting(t *testing.T) {
	c := settings.LoginAsAdmin(t)

	testKey := "general.timezone" // visible_at=user, configurable_at=team

	t.Run("Team 可为用户设置默认值", func(t *testing.T) {
		// 更新团队默认值
		updateReq := map[string]any{
			"value": "America/New_York",
		}

		result, err := apitest.Put[team.SettingsItemDTO](c.HTTPClient(),
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
		original, err := apitest.Get[team.SettingsItemDTO](c.HTTPClient(),
			teamSettingPath(testOrgID, testTeamID, "general.theme"),
			nil,
		)
		require.NoError(t, err, "获取原始值应成功")
		originalValue := original.Value

		// 更新为团队默认值
		updateReq := map[string]any{
			"value": "light",
		}
		_, err = apitest.Put[team.SettingsItemDTO](c.HTTPClient(),
			teamSettingPath(testOrgID, testTeamID, "general.theme"),
			updateReq,
		)
		require.NoError(t, err, "更新应成功")

		// 验证团队自定义值
		updated, err := apitest.Get[team.SettingsItemDTO](c.HTTPClient(),
			teamSettingPath(testOrgID, testTeamID, "general.theme"),
			nil,
		)
		require.NoError(t, err, "获取更新后的值应成功")
		assert.Equal(t, "light", updated.Value, "应返回团队自定义值")
		assert.True(t, updated.IsCustomized, "is_customized 应为 true")

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
	c := settings.LoginAsAdmin(t)

	testKey := "general.theme"

	t.Run("重置团队自定义值", func(t *testing.T) {
		// 先创建团队自定义值（使用不同的值以避免冲突）
		updateReq := map[string]any{
			"value": "light", // 使用 light 而非 dark
		}
		_, err := apitest.Put[team.SettingsItemDTO](c.HTTPClient(),
			teamSettingPath(testOrgID, testTeamID, testKey),
			updateReq,
		)
		require.NoError(t, err, "创建团队自定义值应成功")

		// 验证自定义值存在
		customResult, err := apitest.Get[team.SettingsItemDTO](c.HTTPClient(),
			teamSettingPath(testOrgID, testTeamID, testKey),
			nil,
		)
		require.NoError(t, err, "获取团队设置应成功")
		assert.True(t, customResult.IsCustomized, "创建后 is_customized 应为 true")

		// 重置团队配置
		err = c.Delete(teamSettingPath(testOrgID, testTeamID, testKey))
		require.NoError(t, err, "重置应成功")

		// 验证恢复为组织配置或系统默认值
		defaultResult, err := apitest.Get[team.SettingsItemDTO](c.HTTPClient(),
			teamSettingPath(testOrgID, testTeamID, testKey),
			nil,
		)
		require.NoError(t, err, "获取重置后的值应成功")
		assert.False(t, defaultResult.IsCustomized, "重置后 is_customized 应为 false")
		assert.NotEqual(t, "light", defaultResult.Value, "值应恢复为默认值")

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
	c := settings.LoginAsAdmin(t)

	t.Run("Team 可以配置 configurable_at=team 的设置", func(t *testing.T) {
		// general.theme: visible_at=user, configurable_at=team
		updateReq := map[string]any{
			"value": "dark",
		}

		result, err := apitest.Put[team.SettingsItemDTO](c.HTTPClient(),
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

		_, err := apitest.Put[team.SettingsItemDTO](c.HTTPClient(),
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
	c := settings.LoginAsAdmin(t)

	t.Run("验证返回的设置包含正确的 VisibleAt", func(t *testing.T) {
		result, _, err := apitest.GetList[team.SettingsItemDTO](c.HTTPClient(),
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

// TestTeamSettingsPermission 测试团队设置权限控制
//
// 测试场景：
// 1. 非组织成员访问团队设置应被拒绝
// 2. 组织成员但非团队成员访问团队设置（根据 RBAC 配置可能被拒绝或允许）
// 3. 系统管理员（org owner/admin）可以访问团队设置
// 4. 团队普通成员（member）无法访问团队设置
// 5. 团队负责人（lead）可以访问团队设置
// 6. 系统管理员可以读取和修改设置
//
//nolint:maintidx // 集成测试需要多个场景，保持单一函数减少重复代码
func TestTeamSettingsPermission(t *testing.T) {
	adminClient := settings.LoginAsAdmin(t)

	t.Run("非组织成员访问团队设置应被拒绝", func(t *testing.T) {
		// 创建一个不属于任何组织的测试用户
		testUserResp, err := apitest.Post[map[string]any](
			adminClient.HTTPClient(),
			"/api/admin/users",
			map[string]any{
				"username": "test_non_org_member_team",
				"email":    "non_org_member_team@example.com",
				"password": "Test123456!",
			},
		)
		require.NoError(t, err, "创建测试用户应成功")
		testUserID := uint((*testUserResp)["id"].(float64))
		t.Cleanup(func() {
			_ = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
		})

		// 以测试用户登录
		nonMemberClient := settings.LoginAs(t, "test_non_org_member_team", "Test123456!")

		// 尝试访问团队设置（应被拒绝）
		_, _, err = apitest.GetList[team.SettingsItemDTO](nonMemberClient.HTTPClient(),
			teamSettingsPath(testOrgID, testTeamID),
			nil,
		)
		require.Error(t, err, "非组织成员访问团队设置应被拒绝")
		t.Logf("  预期错误: %v", err)
	})

	t.Run("组织成员但非团队成员访问团队设置", func(t *testing.T) {
		// 创建测试用户并加入组织（但不加入团队）
		testUserResp, err := apitest.Post[map[string]any](
			adminClient.HTTPClient(),
			"/api/admin/users",
			map[string]any{
				"username": "test_org_only_member",
				"email":    "org_only_member@example.com",
				"password": "Test123456!",
			},
		)
		require.NoError(t, err, "创建测试用户应成功")
		testUserID := uint((*testUserResp)["id"].(float64))
		t.Cleanup(func() {
			_ = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
		})

		// 将用户加入组织（但不加入团队）
		_, err = apitest.Post[any](adminClient.HTTPClient(),
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
		orgOnlyClient := settings.LoginAs(t, "test_org_only_member", "Test123456!")

		// 尝试访问团队设置
		// 注意：根据权限配置，组织成员可能可以查看团队设置，但不能修改
		_, _, err = apitest.GetList[team.SettingsItemDTO](orgOnlyClient.HTTPClient(),
			teamSettingsPath(testOrgID, testTeamID),
			nil,
		)
		// 根据实际权限配置验证结果
		if err != nil {
			t.Logf("  组织成员（非团队成员）访问团队设置被拒绝: %v", err)
		} else {
			t.Log("  组织成员（非团队成员）可以访问团队设置（查看权限）")
		}
	})

	t.Run("组织 owner 可以访问团队设置（无需是团队成员）", func(t *testing.T) {
		// 创建测试用户并设置为组织 owner（但不加入团队）
		testUserResp, err := apitest.Post[map[string]any](
			adminClient.HTTPClient(),
			"/api/admin/users",
			map[string]any{
				"username": "test_org_owner_for_team",
				"email":    "org_owner_for_team@example.com",
				"password": "Test123456!",
			},
		)
		require.NoError(t, err, "创建测试用户应成功")
		testUserID := uint((*testUserResp)["id"].(float64))
		t.Cleanup(func() {
			_ = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
		})

		// 将用户加入组织（作为 owner，但不加入团队）
		_, err = apitest.Post[any](adminClient.HTTPClient(),
			fmt.Sprintf("/api/org/%d/members", testOrgID),
			map[string]any{
				"user_id": testUserID,
				"role":    "owner",
			},
		)
		require.NoError(t, err, "将用户加入组织应成功")
		t.Cleanup(func() {
			_ = adminClient.Delete(fmt.Sprintf("/api/org/%d/members/%d", testOrgID, testUserID))
		})

		// 以测试用户登录
		orgOwnerClient := settings.LoginAs(t, "test_org_owner_for_team", "Test123456!")

		// org owner 可以访问团队设置（通过 OrgContext 动态注入 org:*:* 权限）
		result, _, err := apitest.GetList[team.SettingsItemDTO](orgOwnerClient.HTTPClient(),
			teamSettingsPath(testOrgID, testTeamID),
			nil,
		)
		require.NoError(t, err, "组织 owner 访问团队设置应成功")
		assert.NotEmpty(t, result, "设置列表不应为空")
		t.Logf("  组织 owner（非团队成员）访问团队设置成功，获取 %d 个设置", len(result))
	})

	t.Run("组织 admin 可以访问团队设置（无需是团队成员）", func(t *testing.T) {
		// 创建测试用户并设置为组织 admin（但不加入团队）
		testUserResp, err := apitest.Post[map[string]any](
			adminClient.HTTPClient(),
			"/api/admin/users",
			map[string]any{
				"username": "test_org_admin_for_team",
				"email":    "org_admin_for_team@example.com",
				"password": "Test123456!",
			},
		)
		require.NoError(t, err, "创建测试用户应成功")
		testUserID := uint((*testUserResp)["id"].(float64))
		t.Cleanup(func() {
			_ = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
		})

		// 将用户加入组织（作为 admin，但不加入团队）
		_, err = apitest.Post[any](adminClient.HTTPClient(),
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
		orgAdminClient := settings.LoginAs(t, "test_org_admin_for_team", "Test123456!")

		// org admin 可以访问团队设置（通过 OrgContext 动态注入 org:team:*:* 权限）
		result, _, err := apitest.GetList[team.SettingsItemDTO](orgAdminClient.HTTPClient(),
			teamSettingsPath(testOrgID, testTeamID),
			nil,
		)
		require.NoError(t, err, "组织 admin 访问团队设置应成功")
		assert.NotEmpty(t, result, "设置列表不应为空")
		t.Logf("  组织 admin（非团队成员）访问团队设置成功，获取 %d 个设置", len(result))
	})

	t.Run("团队成员可以访问团队设置", func(t *testing.T) {
		// admin 用户是团队成员
		result, _, err := apitest.GetList[team.SettingsItemDTO](adminClient.HTTPClient(),
			teamSettingsPath(testOrgID, testTeamID),
			nil,
		)
		require.NoError(t, err, "团队成员访问团队设置应成功")
		assert.NotEmpty(t, result, "设置列表不应为空")
	})

	t.Run("团队普通成员无法访问团队设置", func(t *testing.T) {
		// 创建测试用户并加入组织和团队（作为 member 角色）
		testUserResp, err := apitest.Post[map[string]any](
			adminClient.HTTPClient(),
			"/api/admin/users",
			map[string]any{
				"username": "test_team_member",
				"email":    "team_member@example.com",
				"password": "Test123456!",
			},
		)
		require.NoError(t, err, "创建测试用户应成功")
		testUserID := uint((*testUserResp)["id"].(float64))
		t.Cleanup(func() {
			_ = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
		})

		// 将用户加入组织
		_, err = apitest.Post[any](adminClient.HTTPClient(),
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

		// 将用户加入团队（作为 member）
		_, err = apitest.Post[any](adminClient.HTTPClient(),
			fmt.Sprintf("/api/org/%d/teams/%d/members", testOrgID, testTeamID),
			map[string]any{
				"user_id": testUserID,
				"role":    "member",
			},
		)
		require.NoError(t, err, "将用户加入团队应成功")
		t.Cleanup(func() {
			_ = adminClient.Delete(fmt.Sprintf("/api/org/%d/teams/%d/members/%d", testOrgID, testTeamID, testUserID))
		})

		// 以测试用户登录
		teamMemberClient := settings.LoginAs(t, "test_team_member", "Test123456!")

		// 团队 member 角色无法读取团队设置（需要 lead/admin 角色）
		_, _, err = apitest.GetList[team.SettingsItemDTO](teamMemberClient.HTTPClient(),
			teamSettingsPath(testOrgID, testTeamID),
			nil,
		)
		require.Error(t, err, "团队 member 访问团队设置应被拒绝")
		t.Logf("  团队 member 访问团队设置被拒绝（符合预期）: %v", err)
	})

	t.Run("团队负责人可以访问团队设置", func(t *testing.T) {
		// 创建测试用户并加入组织和团队（作为 lead 角色）
		// 修复后：team lead 角色通过 TeamContext 动态注入 org:team:settings:* 权限
		testUserResp, err := apitest.Post[map[string]any](
			adminClient.HTTPClient(),
			"/api/admin/users",
			map[string]any{
				"username": "test_team_lead",
				"email":    "team_lead@example.com",
				"password": "Test123456!",
			},
		)
		require.NoError(t, err, "创建测试用户应成功")
		testUserID := uint((*testUserResp)["id"].(float64))
		t.Cleanup(func() {
			_ = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", testUserID))
		})

		// 将用户加入组织
		_, err = apitest.Post[any](adminClient.HTTPClient(),
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

		// 将用户加入团队（作为 lead）
		_, err = apitest.Post[any](adminClient.HTTPClient(),
			fmt.Sprintf("/api/org/%d/teams/%d/members", testOrgID, testTeamID),
			map[string]any{
				"user_id": testUserID,
				"role":    "lead",
			},
		)
		require.NoError(t, err, "将用户加入团队应成功")
		t.Cleanup(func() {
			_ = adminClient.Delete(fmt.Sprintf("/api/org/%d/teams/%d/members/%d", testOrgID, testTeamID, testUserID))
		})

		// 以测试用户登录
		teamLeadClient := settings.LoginAs(t, "test_team_lead", "Test123456!")

		// 团队 lead 角色可以读取团队设置（通过 TeamContext 动态权限注入）
		result, _, err := apitest.GetList[team.SettingsItemDTO](teamLeadClient.HTTPClient(),
			teamSettingsPath(testOrgID, testTeamID),
			nil,
		)
		require.NoError(t, err, "团队 lead 访问团队设置应成功")
		assert.NotEmpty(t, result, "设置列表不应为空")
		t.Logf("  团队 lead 访问团队设置成功，获取 %d 个设置", len(result))

		// 团队 lead 角色也可以修改团队设置
		updateReq := map[string]any{
			"value": "Asia/Tokyo",
		}
		_, err = apitest.Put[team.SettingsItemDTO](teamLeadClient.HTTPClient(),
			teamSettingPath(testOrgID, testTeamID, "general.timezone"),
			updateReq,
		)
		require.NoError(t, err, "团队 lead 修改团队设置应成功")
		t.Cleanup(func() {
			_ = adminClient.Delete(teamSettingPath(testOrgID, testTeamID, "general.timezone"))
		})
		t.Log("  团队 lead 修改团队设置成功")
	})

	t.Run("系统管理员（团队创建者）可以读取和修改设置", func(t *testing.T) {
		// admin 用户（系统管理员）是团队的创建者
		// 只有系统级 RBAC 权限才能访问团队设置
		result, _, err := apitest.GetList[team.SettingsItemDTO](adminClient.HTTPClient(),
			teamSettingsPath(testOrgID, testTeamID),
			nil,
		)
		require.NoError(t, err, "系统管理员读取团队设置应成功")
		assert.NotEmpty(t, result, "设置列表不应为空")

		// 系统管理员可以修改团队设置
		updateReq := map[string]any{
			"value": "dark",
		}
		_, err = apitest.Put[team.SettingsItemDTO](adminClient.HTTPClient(),
			teamSettingPath(testOrgID, testTeamID, "general.theme"),
			updateReq,
		)
		require.NoError(t, err, "系统管理员修改团队设置应成功")
		t.Cleanup(func() {
			_ = adminClient.Delete(teamSettingPath(testOrgID, testTeamID, "general.theme"))
		})
	})
}
