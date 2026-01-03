package system_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application/stats"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/platform/manualtest"
)

// TestHealthCheck 测试健康检查端点。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestHealthCheck ./internal/integration/system/
func TestHealthCheck(t *testing.T) {
	manualtest.SkipIfNotManual(t)

	c := manualtest.NewClient()

	t.Log("检查健康状态...")
	resp, err := c.R().Get("/health")
	require.NoError(t, err, "健康检查请求失败")
	require.False(t, resp.IsError(), "健康检查失败: 状态码 %d", resp.StatusCode())

	t.Logf("健康检查通过!")
	t.Logf("  状态码: %d", resp.StatusCode())
	t.Logf("  响应: %s", string(resp.Body()))
}

// TestSystemStats 测试系统统计端点。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestSystemStats ./internal/integration/system/
func TestSystemStats(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	t.Log("\n获取系统统计...")
	statsResult, err := manualtest.Get[stats.StatsDTO](c, "/api/admin/overview/stats", nil)
	require.NoError(t, err, "获取系统统计失败")

	// 验证返回的数据
	assert.GreaterOrEqual(t, statsResult.TotalUsers, int64(0), "总用户数不应为负数")
	assert.GreaterOrEqual(t, statsResult.ActiveUsers, int64(0), "活跃用户数不应为负数")
	assert.GreaterOrEqual(t, statsResult.TotalRoles, int64(0), "总角色数不应为负数")

	t.Logf("系统统计:")
	t.Logf("  总用户数: %d", statsResult.TotalUsers)
	t.Logf("  活跃用户数: %d", statsResult.ActiveUsers)
	t.Logf("  非活跃用户数: %d", statsResult.InactiveUsers)
	t.Logf("  封禁用户数: %d", statsResult.BannedUsers)
	t.Logf("  总角色数: %d", statsResult.TotalRoles)
	t.Logf("  总权限数: %d", statsResult.TotalPermissions)
	if len(statsResult.RecentAuditLogs) > 0 {
		t.Logf("  最近审计日志: %d 条", len(statsResult.RecentAuditLogs))
	}
}

// TestSwaggerDocs 测试 Swagger 文档端点。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestSwaggerDocs ./internal/integration/system/
func TestSwaggerDocs(t *testing.T) {
	manualtest.SkipIfNotManual(t)

	c := manualtest.NewClient()

	t.Log("检查 Swagger 文档...")
	resp, err := c.R().Get("/swagger/index.html")
	require.NoError(t, err, "Swagger 请求失败")

	if resp.IsError() {
		t.Logf("Swagger 文档不可用: 状态码 %d", resp.StatusCode())
	} else {
		t.Logf("Swagger 文档可用!")
		t.Logf("  状态码: %d", resp.StatusCode())
		t.Logf("  内容长度: %d bytes", len(resp.Body()))
	}
}
