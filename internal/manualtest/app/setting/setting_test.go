package setting_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application/setting"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/platform/manualtest"
)

// TestMain 在所有测试完成后清理测试数据。
func TestMain(m *testing.M) {
	code := m.Run()

	if os.Getenv("MANUAL") == "1" {
		cleanupTestSettings()
	}

	os.Exit(code)
}

// cleanupTestSettings 清理测试创建的配置。
func cleanupTestSettings() {
	c := manualtest.NewClient()
	if _, err := c.Login("admin", "admin123"); err != nil {
		return
	}

	// 测试配置键前缀列表
	settingPrefixes := []string{"test_setting_"}

	settings, _, _ := manualtest.GetList[setting.SettingDTO](c, "/api/admin/settings", map[string]string{"limit": "1000"})
	for _, s := range settings {
		for _, prefix := range settingPrefixes {
			if len(s.Key) >= len(prefix) && s.Key[:len(prefix)] == prefix {
				_ = c.Delete("/api/admin/settings/" + s.Key)
				break
			}
		}
	}
}

// 测试配置前缀，用于隔离测试数据
const settingTestPrefix = "test_setting_"

// TestSettingsFlow 系统配置完整流程测试。
//
// 测试 CRUD 完整流程：创建 → 获取 → 更新 → 验证 Schema 一致性
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestSettingsFlow ./internal/integration/setting/
func TestSettingsFlow(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 测试 1: 创建配置
	t.Log("\n测试 1: 创建配置")
	settingKey := fmt.Sprintf("%s%d", settingTestPrefix, time.Now().Unix())
	createReq := map[string]any{
		"key":           settingKey,
		"default_value": "测试值",
		"category_id":   1, // general 分类
		"group":         "basic",
		"value_type":    "string",
		"label":         "测试配置",
	}

	created, err := manualtest.Post[setting.SettingDTO](c, "/api/admin/settings", createReq)
	require.NoError(t, err, "创建配置失败")
	require.NotZero(t, created.ID, "创建的配置 ID 为 0")
	assert.Equal(t, settingKey, created.Key)
	t.Logf("  ✓ 创建成功: ID=%d, Key=%s", created.ID, created.Key)

	t.Cleanup(func() {
		_ = c.Delete("/api/admin/settings/" + settingKey)
	})

	// 测试 2: 获取单个配置
	t.Log("\n测试 2: 获取单个配置")
	detail, err := manualtest.Get[setting.SettingDTO](c, "/api/admin/settings/"+settingKey, nil)
	require.NoError(t, err, "获取配置失败")
	assert.Equal(t, settingKey, detail.Key)
	assert.Equal(t, "测试配置", detail.Label)
	t.Logf("  ✓ Key=%s, Label=%s, Value=%v", detail.Key, detail.Label, detail.DefaultValue)

	// 测试 3: 更新配置
	t.Log("\n测试 3: 更新配置")
	updateReq := map[string]any{
		"default_value": "更新后的值",
		"label":         "更新后的标签",
	}
	updated, err := manualtest.Put[setting.SettingDTO](c, "/api/admin/settings/"+settingKey, updateReq)
	require.NoError(t, err, "更新配置失败")
	assert.Equal(t, "更新后的标签", updated.Label)
	t.Logf("  ✓ 更新成功: Label=%s, Value=%v", updated.Label, updated.DefaultValue)

	// 测试 4: 验证 Schema 包含新创建的配置（缓存一致性验证）
	t.Log("\n测试 4: 验证 Schema 缓存一致性")
	schema, err := manualtest.Get[[]setting.SettingsCategoryDTO](c, "/api/admin/settings", nil)
	require.NoError(t, err, "获取 Schema 失败")

	found := false
	for _, cat := range *schema {
		for _, group := range cat.Groups {
			for _, s := range group.Settings {
				if s.Key == settingKey {
					found = true
					t.Logf("  ✓ 配置存在于 Schema: 分类=%s, 分组=%s", cat.Category, group.Name)
					break
				}
			}
		}
	}
	assert.True(t, found, "创建的配置不在 Schema 中（缓存不一致）")

	t.Log("\n系统配置流程测试完成!")
}

// TestGetSettingsWithFilters 测试配置列表查询（Table-Driven）。
//
// 覆盖场景：获取所有配置、按 general/security 类别筛选
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestGetSettingsWithFilters ./internal/integration/setting/
func TestGetSettingsWithFilters(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 注意：GET /api/admin/settings 返回 []SettingsCategoryDTO（层级结构）
	// 使用 ?category=xxx 按分类 Key 筛选（不是 category_id）
	cases := []struct {
		name     string
		query    map[string]string
		validate func(t *testing.T, schema []setting.SettingsCategoryDTO)
	}{
		{
			name:  "获取所有配置",
			query: nil,
			validate: func(t *testing.T, schema []setting.SettingsCategoryDTO) {
				t.Helper()
				assert.NotEmpty(t, schema, "Schema 为空")
				t.Logf("分类总数: %d", len(schema))
				for _, cat := range schema {
					settingCount := 0
					for _, g := range cat.Groups {
						settingCount += len(g.Settings)
					}
					t.Logf("  [%s] %s: %d 分组, %d 配置", cat.Category, cat.Label, len(cat.Groups), settingCount)
				}
			},
		},
		{
			name:  "按 category=general 筛选",
			query: map[string]string{"category": "general"},
			validate: func(t *testing.T, schema []setting.SettingsCategoryDTO) {
				t.Helper()
				require.Len(t, schema, 1, "应该只返回 1 个分类")
				assert.Equal(t, "general", schema[0].Category, "分类 Key 不匹配")
				t.Logf("general 分类: %d 分组", len(schema[0].Groups))
				for _, g := range schema[0].Groups {
					t.Logf("  - %s: %d 配置", g.Name, len(g.Settings))
				}
			},
		},
		{
			name:  "按 category=security 筛选",
			query: map[string]string{"category": "security"},
			validate: func(t *testing.T, schema []setting.SettingsCategoryDTO) {
				t.Helper()
				require.Len(t, schema, 1, "应该只返回 1 个分类")
				assert.Equal(t, "security", schema[0].Category, "分类 Key 不匹配")
				t.Logf("security 分类: %d 分组", len(schema[0].Groups))
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			schema, err := manualtest.Get[[]setting.SettingsCategoryDTO](c, "/api/admin/settings", tc.query)
			require.NoError(t, err, "获取配置失败")
			tc.validate(t, *schema)
		})
	}
}

// TestBatchUpdateSettings 测试批量更新配置。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestBatchUpdateSettings ./internal/integration/setting/
func TestBatchUpdateSettings(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 先创建两个测试配置
	timestamp := time.Now().Unix()
	key1 := fmt.Sprintf("%sbatch1_%d", settingTestPrefix, timestamp)
	key2 := fmt.Sprintf("%sbatch2_%d", settingTestPrefix, timestamp)

	t.Log("\n准备: 创建两个测试配置...")
	for _, key := range []string{key1, key2} {
		createReq := map[string]any{
			"key":           key,
			"default_value": "初始值",
			"category_id":   1, // 使用 general 分类（ID=1）
			"group":         "batch",
			"value_type":    "string",
			"label":         "批量测试",
		}
		_, createErr := manualtest.Post[setting.SettingDTO](c, "/api/admin/settings", createReq)
		require.NoError(t, createErr, "创建配置 %s 失败", key)
		t.Logf("  创建配置: %s", key)
	}

	// 确保清理
	t.Cleanup(func() {
		for _, key := range []string{key1, key2} {
			if deleteErr := c.Delete("/api/admin/settings/" + key); deleteErr != nil {
				t.Logf("清理配置 %s 失败: %v", key, deleteErr)
			}
		}
	})

	// 批量更新
	t.Log("\n测试: 批量更新配置...")
	batchReq := map[string]any{
		"settings": []map[string]any{
			{"key": key1, "value": "批量更新值1"},
			{"key": key2, "value": "批量更新值2"},
		},
	}

	resp, err := c.R().
		SetBody(batchReq).
		Post("/api/admin/settings/batch")
	require.NoError(t, err, "批量更新请求失败")
	require.False(t, resp.IsError(), "批量更新失败: 状态码 %d", resp.StatusCode())
	t.Log("  批量更新成功!")

	// 验证更新结果
	t.Log("\n验证更新结果...")
	for i, key := range []string{key1, key2} {
		detail, getErr := manualtest.Get[setting.SettingDTO](c, "/api/admin/settings/"+key, nil)
		require.NoError(t, getErr, "获取配置 %s 失败", key)
		expected := fmt.Sprintf("批量更新值%d", i+1)
		assert.Equal(t, expected, detail.DefaultValue, "配置 %s 值不匹配", key)
		t.Logf("  ✓ %s = %v", key, detail.DefaultValue)
	}
}

// TestDeleteSetting 测试删除配置。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestDeleteSetting ./internal/integration/setting/
func TestDeleteSetting(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 使用 helper 创建配置（带清理控制）
	created, markDeleted := manualtest.CreateTestSettingWithCleanupControl(t, c, settingTestPrefix+"delete")
	t.Logf("创建测试配置: %s", created.Key)

	// 删除配置
	err := c.Delete("/api/admin/settings/" + created.Key)
	require.NoError(t, err, "删除配置失败")
	markDeleted() // 标记已删除，阻止 Cleanup 重复删除
	t.Log("  ✓ 删除成功")

	// 验证删除
	_, err = manualtest.Get[setting.SettingDTO](c, "/api/admin/settings/"+created.Key, nil)
	require.Error(t, err, "配置应该已被删除")
	t.Log("  ✓ 配置已确认删除")
}
