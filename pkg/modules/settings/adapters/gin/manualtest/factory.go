package manualtest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/app/user"
)

// CreateTestUserSetting 创建测试用户设置值并自动注册清理。
// 测试结束时会自动删除创建的设置值。
func CreateTestUserSetting(t *testing.T, c *Client, key string, value any) *user.SettingsItemDTO {
	t.Helper()

	req := map[string]any{
		"value": value,
	}

	resp, err := Put[user.SettingsItemDTO](c, "/api/user/settings/"+key, req)
	require.NoError(t, err, "创建测试设置失败: %s", key)

	t.Cleanup(func() {
		_ = c.Delete("/api/user/settings/" + key)
	})

	return resp
}

// CreateTestUserSettingWithCleanupControl 创建测试用户设置，返回清理控制函数。
// 当测试本身需要删除设置时使用，删除成功后调用返回的 markDeleted 函数。
func CreateTestUserSettingWithCleanupControl(t *testing.T, c *Client, key string, value any) (*user.SettingsItemDTO, func()) {
	t.Helper()

	req := map[string]any{
		"value": value,
	}

	result, err := Put[user.SettingsItemDTO](c, "/api/user/settings/"+key, req)
	require.NoError(t, err, "创建测试设置失败: %s", key)

	deleted := false

	t.Cleanup(func() {
		if !deleted {
			_ = c.Delete("/api/user/settings/" + key)
		}
	})

	return result, func() { deleted = true }
}
