package settings

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/lwmacct/260103-ddd-shared/pkg/shared/apitest"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/app/user"
)

// CreateTestUserSetting 创建测试用户设置值并自动注册清理。
// key 会自动添加 UUID 后缀以避免并行测试冲突。
func CreateTestUserSetting(t *testing.T, c *Client, key string, value any) *user.SettingsItemDTO {
	t.Helper()

	// 添加 UUID 后缀避免冲突
	uniqueKey := fmt.Sprintf("%s_%s", key, uuid.New().String()[:8])

	req := map[string]any{
		"value": value,
	}

	resp, err := apitest.Put[user.SettingsItemDTO](&c.Client, "/api/user/settings/"+uniqueKey, req)
	require.NoError(t, err, "创建测试设置失败: %s", uniqueKey)

	t.Cleanup(func() {
		_ = c.Delete("/api/user/settings/" + uniqueKey)
	})

	return resp
}

// CreateTestUserSettingWithCleanupControl 创建测试用户设置，返回清理控制函数。
// key 会自动添加 UUID 后缀以避免并行测试冲突。
func CreateTestUserSettingWithCleanupControl(t *testing.T, c *Client, key string, value any) (*user.SettingsItemDTO, func()) {
	t.Helper()

	// 添加 UUID 后缀避免冲突
	uniqueKey := fmt.Sprintf("%s_%s", key, uuid.New().String()[:8])

	req := map[string]any{
		"value": value,
	}

	result, err := apitest.Put[user.SettingsItemDTO](&c.Client, "/api/user/settings/"+uniqueKey, req)
	require.NoError(t, err, "创建测试设置失败: %s", uniqueKey)

	deleted := false

	t.Cleanup(func() {
		if !deleted {
			_ = c.Delete("/api/user/settings/" + uniqueKey)
		}
	})

	return result, func() { deleted = true }
}
