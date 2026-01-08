package pat_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260103-ddd-iam-bc/internal/apitest/iam"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/pat"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/role"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/user"
	"github.com/lwmacct/260103-ddd-shared/pkg/shared/apitest"
)

// TestMain 在所有测试完成后清理测试数据。
func TestMain(m *testing.M) {
	code := m.Run()

	if os.Getenv("MANUAL") == "1" {
		cleanupTestData()
	}

	os.Exit(code)
}

// cleanupTestData 清理测试创建的 PAT、角色和用户。
func cleanupTestData() {
	c := iam.NewClientFromConfig()
	if _, err := c.Login("admin", "admin123"); err != nil {
		return
	}

	// 清理测试 PAT tokens
	patPrefixes := []string{"test_pat_", "limited_pat_", "scope_test_pat_", "full_scope_pat_", "self_scope_pat_"}
	tokens, _, _ := apitest.GetList[pat.TokenDTO](c.HTTPClient(), "/api/user/tokens", map[string]string{"limit": "1000"})
	for _, token := range tokens {
		for _, prefix := range patPrefixes {
			if len(token.Name) >= len(prefix) && token.Name[:len(prefix)] == prefix {
				_ = c.Delete(fmt.Sprintf("/api/user/tokens/%d", token.ID))
				break
			}
		}
	}

	// 清理测试角色
	rolePrefixes := []string{"pat_scope_test_role_"}
	roles, _, _ := apitest.GetList[role.RoleDTO](c.HTTPClient(), "/api/admin/roles", map[string]string{"limit": "1000"})
	for _, r := range roles {
		for _, prefix := range rolePrefixes {
			if len(r.Name) >= len(prefix) && r.Name[:len(prefix)] == prefix {
				_ = c.Delete(fmt.Sprintf("/api/admin/roles/%d", r.ID))
				break
			}
		}
	}

	// 清理测试用户
	userPrefixes := []string{"pat_scope_user_"}
	users, _, _ := apitest.GetList[user.UserDTO](c.HTTPClient(), "/api/admin/users", map[string]string{"limit": "1000"})
	for _, u := range users {
		for _, prefix := range userPrefixes {
			if len(u.Username) >= len(prefix) && u.Username[:len(prefix)] == prefix {
				_ = c.Delete(fmt.Sprintf("/api/admin/users/%d", u.ID))
				break
			}
		}
	}
}

// TestPATFlow PAT 令牌完整流程测试。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestPATFlow ./internal/integration/pat/
func TestPATFlow(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	// 测试 1: 获取 PAT 列表
	t.Log("\n测试 1: 获取 PAT 列表")
	tokens, _, err := apitest.GetList[pat.TokenDTO](c.HTTPClient(), "/api/user/tokens", nil)
	require.NoError(t, err, "获取 PAT 列表失败")
	t.Logf("  现有 PAT 数量: %d", len(tokens))

	// 测试 2: 创建 PAT
	t.Log("\n测试 2: 创建 PAT")
	tokenName := fmt.Sprintf("test_pat_%d", time.Now().Unix())
	expiresIn := 30 // 30 天
	createReq := pat.CreateDTO{
		Name:        tokenName,
		Scopes:      []string{"self"}, // 仅 self 域权限
		ExpiresIn:   &expiresIn,
		Description: "测试用 PAT",
	}
	t.Logf("  创建 PAT: %s", tokenName)

	created, err := apitest.Post[pat.CreateResultDTO](c.HTTPClient(), "/api/user/tokens", createReq)
	require.NoError(t, err, "创建 PAT 失败")
	require.NotZero(t, created.Token.ID, "创建的 PAT ID 为 0")
	require.NotEmpty(t, created.PlainToken, "未返回明文令牌")

	tokenID := created.Token.ID

	// 注册清理：测试结束时删除 PAT（无论成功或失败）
	deleted := false
	t.Cleanup(func() {
		if !deleted {
			if delErr := c.Delete(fmt.Sprintf("/api/user/tokens/%d", tokenID)); delErr != nil {
				t.Logf("清理 PAT 失败: %v", delErr)
			}
		}
	})

	assert.Equal(t, tokenName, created.Token.Name, "Token 名称不匹配")
	assert.Equal(t, "active", created.Token.Status, "期望初始状态为 active")
	t.Logf("  创建成功! PAT ID: %d", created.Token.ID)
	t.Logf("  明文令牌: %s... (仅显示一次)", created.PlainToken[:20])
	t.Logf("  状态: %s", created.Token.Status)

	// 测试 3: 获取 PAT 详情
	t.Log("\n测试 3: 获取 PAT 详情")
	detail, err := apitest.Get[pat.TokenDTO](c.HTTPClient(), fmt.Sprintf("/api/user/tokens/%d", tokenID), nil)
	require.NoError(t, err, "获取 PAT 详情失败")
	assert.Equal(t, tokenID, detail.ID, "Token ID 不匹配")
	assert.Equal(t, tokenName, detail.Name, "Token 名称不匹配")
	assert.NotEmpty(t, detail.Scopes, "Scopes 列表为空")
	t.Logf("  名称: %s", detail.Name)
	t.Logf("  前缀: %s", detail.TokenPrefix)
	t.Logf("  Scopes: %v", detail.Scopes)
	t.Logf("  状态: %s", detail.Status)
	if detail.ExpiresAt != nil {
		t.Logf("  过期时间: %s", detail.ExpiresAt.Format("2006-01-02 15:04:05"))
	}

	// 测试 4: 禁用 PAT
	t.Log("\n测试 4: 禁用 PAT")
	resp, err := c.R().Patch(fmt.Sprintf("/api/user/tokens/%d/disable", tokenID))
	require.NoError(t, err, "禁用 PAT 失败")
	require.False(t, resp.IsError(), "禁用 PAT 失败: 状态码 %d", resp.StatusCode())
	t.Log("  禁用成功!")

	// 验证状态
	disabled, err := apitest.Get[pat.TokenDTO](c.HTTPClient(), fmt.Sprintf("/api/user/tokens/%d", tokenID), nil)
	require.NoError(t, err, "获取 PAT 详情失败")
	assert.Equal(t, "disabled", disabled.Status, "期望状态为 disabled")
	t.Logf("  当前状态: %s", disabled.Status)

	// 测试 5: 启用 PAT
	t.Log("\n测试 5: 启用 PAT")
	resp, err = c.R().Patch(fmt.Sprintf("/api/user/tokens/%d/enable", tokenID))
	require.NoError(t, err, "启用 PAT 失败")
	require.False(t, resp.IsError(), "启用 PAT 失败: 状态码 %d", resp.StatusCode())
	t.Log("  启用成功!")

	// 验证状态
	enabled, err := apitest.Get[pat.TokenDTO](c.HTTPClient(), fmt.Sprintf("/api/user/tokens/%d", tokenID), nil)
	require.NoError(t, err, "获取 PAT 详情失败")
	assert.Equal(t, "active", enabled.Status, "期望状态为 active")
	t.Logf("  当前状态: %s", enabled.Status)

	// 测试 6: 删除 PAT
	t.Log("\n测试 6: 删除 PAT")
	err = c.Delete(fmt.Sprintf("/api/user/tokens/%d", tokenID))
	require.NoError(t, err, "删除 PAT 失败")
	deleted = true // 标记已删除，Cleanup 不再重复删除
	t.Log("  删除成功!")

	t.Log("\nPAT 令牌流程测试完成!")
}

// TestListPATs 测试获取 PAT 列表。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestListPATs ./internal/integration/pat/
func TestListPATs(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	t.Log("获取 PAT 列表...")
	tokens, meta, err := apitest.GetList[pat.TokenDTO](c.HTTPClient(), "/api/user/tokens", nil)
	require.NoError(t, err, "获取 PAT 列表失败")

	t.Logf("PAT 数量: %d", len(tokens))
	if meta != nil {
		t.Logf("总数: %d", meta.Total)
	}

	for _, token := range tokens {
		statusIcon := "✓"
		if token.Status != "active" {
			statusIcon = "✗"
		}
		t.Logf("  [%s] [%d] %s (%s)", statusIcon, token.ID, token.Name, token.TokenPrefix)
	}
}

// TestPATWithScopes 测试创建带特定 Scope 的 PAT。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestPATWithScopes ./internal/integration/pat/
func TestPATWithScopes(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	// 创建带限制 Scope 的 PAT
	t.Log("\n创建带限制 Scope 的 PAT...")
	tokenName := fmt.Sprintf("limited_pat_%d", time.Now().Unix())
	createReq := pat.CreateDTO{
		Name: tokenName,
		Scopes: []string{
			"self",
			"sys",
		},
		Description: "仅限 self 和 sys 域操作",
	}

	created, err := apitest.Post[pat.CreateResultDTO](c.HTTPClient(), "/api/user/tokens", createReq)
	require.NoError(t, err, "创建 PAT 失败")
	require.NotZero(t, created.Token.ID, "创建的 PAT ID 为 0")
	assert.Len(t, created.Token.Scopes, 2, "期望 Scopes 数量为 2 (self, sys)")
	assert.ElementsMatch(t, []string{"self", "sys"}, created.Token.Scopes, "Scope 内容应匹配")
	t.Logf("  创建成功! PAT ID: %d", created.Token.ID)
	t.Logf("  Scopes: %v", created.Token.Scopes)

	// 确保清理
	t.Cleanup(func() {
		if err := c.Delete(fmt.Sprintf("/api/user/tokens/%d", created.Token.ID)); err != nil {
			t.Logf("清理 PAT 失败: %v", err)
		}
	})
}

// TestListPATScopes 测试获取可用 Scope 列表。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestListPATScopes ./internal/integration/pat/
func TestListPATScopes(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	t.Log("获取 PAT Scope 列表...")

	// 使用 patdomain.ScopeInfo 类型
	type ScopeInfo struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		Description string `json:"description"`
	}

	scopes, _, err := apitest.GetList[ScopeInfo](c.HTTPClient(), "/api/user/tokens/scopes", nil)
	require.NoError(t, err, "获取 Scope 列表失败")
	require.NotEmpty(t, scopes, "Scope 列表不应为空")

	t.Logf("可用 Scope 数量: %d", len(scopes))
	for _, scope := range scopes {
		t.Logf("  - %s (%s): %s", scope.Name, scope.DisplayName, scope.Description)
	}

	// 验证预期的 Scope
	scopeNames := make([]string, len(scopes))
	for i, s := range scopes {
		scopeNames[i] = s.Name
	}
	assert.Contains(t, scopeNames, "full", "应包含 full scope")
	assert.Contains(t, scopeNames, "self", "应包含 self scope")
	assert.Contains(t, scopeNames, "sys", "应包含 sys scope")
}

// TestPATScopeEnforcement 测试 PAT Scope 权限过滤是否生效。
//
// 隐性角色机制：所有已认证用户自动拥有 "user" 角色及其权限 (self:*:*)。
// 因此 Admin 用户实际权限为 *:*:* (admin 角色) + self:*:* (隐性 user 角色)。
//
// self scope PAT 过滤结果：
//   - *:*:* → 被过滤（不匹配 self: 前缀）
//   - self:*:* → 保留
//
// 测试验证：
//   - self scope PAT 可以访问 self 域 API（200）
//   - self scope PAT 不能访问 sys 域 API（403）
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestPATScopeEnforcement ./internal/integration/pat/
func TestPATScopeEnforcement(t *testing.T) {
	// 1. 使用 JWT 登录创建 PAT
	jwtClient := iam.LoginAsAdmin(t)

	t.Log("\n===== PAT Scope 权限过滤测试 =====")
	t.Log("注意: Admin 用户有隐性 'user' 角色，self scope PAT 保留 self:*:* 权限")

	// 2. 创建仅 self scope 的 PAT
	t.Log("\n步骤 1: 创建仅 self scope 的 PAT")
	tokenName := fmt.Sprintf("scope_test_pat_%d", time.Now().Unix())
	createReq := pat.CreateDTO{
		Name:        tokenName,
		Scopes:      []string{"self"}, // 仅 self 域权限
		Description: "Scope 过滤测试用 PAT",
	}

	created, err := apitest.Post[pat.CreateResultDTO](jwtClient.HTTPClient(), "/api/user/tokens", createReq)
	require.NoError(t, err, "创建 PAT 失败")
	require.NotEmpty(t, created.PlainToken, "未返回明文令牌")
	t.Logf("  PAT ID: %d", created.Token.ID)
	t.Logf("  Scopes: %v", created.Token.Scopes)
	t.Logf("  PlainToken: %s...", created.PlainToken[:20])

	// 注册清理
	t.Cleanup(func() {
		if delErr := jwtClient.Delete(fmt.Sprintf("/api/user/tokens/%d", created.Token.ID)); delErr != nil {
			t.Logf("清理 PAT 失败: %v", delErr)
		}
	})

	// 3. 创建使用 PAT 的客户端
	patClient := iam.NewClientFromConfig()
	patClient.SetToken(created.PlainToken)

	// 4. 测试 self 域 API（应成功，因为隐性 user 角色有 self:*:* 权限）
	t.Log("\n步骤 2: 使用 self scope PAT 访问 self 域 API")
	t.Log("  预期: 200（隐性 user 角色的 self:*:* 权限被保留）")

	// 访问 /api/user/profile（self:profile:get）
	resp, err := patClient.R().Get("/api/user/profile")
	require.NoError(t, err, "请求失败")
	t.Logf("  GET /api/user/profile -> 状态码: %d", resp.StatusCode())
	// 隐性 user 角色有 self:*:* 权限，匹配 self scope，可以访问
	assert.Equal(t, 200, resp.StatusCode(), "self scope PAT 应能访问 self 域（隐性 user 角色权限）")

	// 5. 测试 sys 域 API（应被拒绝，因为 *:*:* 被过滤）
	t.Log("\n步骤 3: 使用 self scope PAT 访问 sys 域 API")
	t.Log("  预期: 403（*:*:* 被过滤，无 sys 域权限）")

	// 访问 /api/admin/users（admin:users:list）- 正确路径
	resp, err = patClient.R().Get("/api/admin/users")
	require.NoError(t, err, "请求失败")
	t.Logf("  GET /api/admin/users -> 状态码: %d", resp.StatusCode())
	assert.Equal(t, 403, resp.StatusCode(), "self scope PAT 不应能访问 sys 域（*:*:* 已被过滤）")

	t.Log("\n===== PAT Scope 权限过滤测试完成 =====")
	t.Log("结论: PAT Scope 正确过滤权限，保留匹配前缀的权限（包括隐性角色权限）")
}

// TestPATFullScopeAccess 测试 full scope PAT 可以访问所有 API。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestPATFullScopeAccess ./internal/integration/pat/
func TestPATFullScopeAccess(t *testing.T) {
	jwtClient := iam.LoginAsAdmin(t)

	t.Log("\n===== PAT Full Scope 测试 =====")

	// 创建 full scope PAT
	t.Log("\n步骤 1: 创建 full scope PAT")
	tokenName := fmt.Sprintf("full_scope_pat_%d", time.Now().Unix())
	createReq := pat.CreateDTO{
		Name:        tokenName,
		Scopes:      []string{"full"}, // 完整权限
		Description: "Full scope 测试用 PAT",
	}

	created, err := apitest.Post[pat.CreateResultDTO](jwtClient.HTTPClient(), "/api/user/tokens", createReq)
	require.NoError(t, err, "创建 PAT 失败")
	t.Logf("  PAT ID: %d, Scopes: %v", created.Token.ID, created.Token.Scopes)

	t.Cleanup(func() {
		if delErr := jwtClient.Delete(fmt.Sprintf("/api/user/tokens/%d", created.Token.ID)); delErr != nil {
			t.Logf("清理 PAT 失败: %v", delErr)
		}
	})

	// 使用 PAT 客户端
	patClient := iam.NewClientFromConfig()
	patClient.SetToken(created.PlainToken)

	// 测试 self 域 API
	t.Log("\n步骤 2: 使用 full scope PAT 访问 self 域 API")
	resp, err := patClient.R().Get("/api/user/profile")
	require.NoError(t, err)
	t.Logf("  GET /api/user/profile -> 状态码: %d", resp.StatusCode())
	assert.Equal(t, 200, resp.StatusCode(), "full scope 应能访问 self 域")

	// 测试 sys 域 API（正确路径）
	t.Log("\n步骤 3: 使用 full scope PAT 访问 sys 域 API")
	resp, err = patClient.R().Get("/api/admin/users")
	require.NoError(t, err)
	t.Logf("  GET /api/admin/users -> 状态码: %d", resp.StatusCode())
	assert.Equal(t, 200, resp.StatusCode(), "full scope 应能访问 sys 域")

	t.Log("\n===== PAT Full Scope 测试完成 =====")
}

// TestPATScopeWithRegularUser 测试普通用户的 PAT Scope 权限过滤。
//
// 这是关键测试：验证有具体域权限的用户使用限制性 Scope 时，
// PAT 权限被正确限制在 Scope 范围内。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestPATScopeWithRegularUser ./internal/integration/pat/
func TestPATScopeWithRegularUser(t *testing.T) {
	adminClient := iam.LoginAsAdmin(t)

	t.Log("\n===== 普通用户 PAT Scope 权限过滤测试 =====")
	t.Log("目标: 验证有具体域权限的用户，PAT Scope 过滤是否正确工作")

	// 步骤 1: 创建测试角色（包含 self 和 sys 域权限）
	t.Log("\n步骤 1: 创建测试角色（包含 self + sys 域权限）")
	roleName := "pat_scope_test_role_" + uuid.New().String()[:8]
	createRoleReq := role.CreateDTO{
		Name:        roleName,
		DisplayName: "PAT Scope 测试角色",
		Description: "用于测试 PAT Scope 过滤的角色",
	}

	createdRole, err := apitest.Post[role.CreateResultDTO](adminClient.HTTPClient(), "/api/admin/roles", createRoleReq)
	require.NoError(t, err, "创建测试角色失败")
	t.Logf("  角色 ID: %d, 名称: %s", createdRole.RoleID, createdRole.Name)

	// 注册角色清理
	t.Cleanup(func() {
		_ = adminClient.Delete(fmt.Sprintf("/api/admin/roles/%d", createdRole.RoleID))
	})

	// 步骤 2: 为角色设置权限（self 域 + sys 域）
	t.Log("\n步骤 2: 设置角色权限")
	permissions := []role.PermissionInputDTO{
		// self 域权限
		{OperationPattern: "self:profile:get", ResourcePattern: "*"},
		{OperationPattern: "self:profile:update", ResourcePattern: "*"},
		{OperationPattern: "self:tokens:*", ResourcePattern: "*"}, // PAT 管理权限
		// sys 域权限
		{OperationPattern: "admin:users:list", ResourcePattern: "*"},
		{OperationPattern: "admin:users:get", ResourcePattern: "*"},
	}
	setPermReq := role.SetPermissionsDTO{Permissions: permissions}

	resp, err := adminClient.R().
		SetBody(setPermReq).
		Put(fmt.Sprintf("/api/admin/roles/%d/permissions", createdRole.RoleID))
	require.NoError(t, err, "设置角色权限请求失败")
	require.False(t, resp.IsError(), "设置角色权限失败，状态码: %d", resp.StatusCode())
	t.Logf("  设置了 %d 个权限: self:profile:*, self:tokens:*, admin:users:list/get", len(permissions))

	// 步骤 3: 创建测试用户并分配角色
	t.Log("\n步骤 3: 创建测试用户并分配角色")
	username := "pat_scope_user_" + uuid.New().String()[:8]
	password := "testpass123"
	createUserReq := user.CreateDTO{
		Username: username,
		Email:    username + "@test.local",
		Password: password,
		RealName: "PAT Scope 测试用户",
		RoleIDs:  []uint{createdRole.RoleID},
	}

	createdUser, err := apitest.Post[user.UserWithRolesDTO](adminClient.HTTPClient(), "/api/admin/users", createUserReq)
	require.NoError(t, err, "创建测试用户失败")
	t.Logf("  用户 ID: %d, 用户名: %s", createdUser.ID, createdUser.Username)
	t.Logf("  分配角色: %v", createdUser.Roles)

	// 注册用户清理
	t.Cleanup(func() {
		_ = adminClient.Delete(fmt.Sprintf("/api/admin/users/%d", createdUser.ID))
	})

	// 步骤 4: 使用测试用户登录
	t.Log("\n步骤 4: 测试用户登录")
	userClient := iam.NewClientFromConfig()
	_, err = userClient.Login(username, password)
	require.NoError(t, err, "测试用户登录失败")
	t.Log("  登录成功!")

	// 步骤 5: 创建仅 self scope 的 PAT
	t.Log("\n步骤 5: 创建仅 self scope 的 PAT")
	patName := "self_scope_pat_" + uuid.New().String()[:8]
	createPATReq := pat.CreateDTO{
		Name:        patName,
		Scopes:      []string{"self"}, // 仅 self 域
		Description: "仅 self scope 的测试 PAT",
	}

	createdPAT, err := apitest.Post[pat.CreateResultDTO](userClient.HTTPClient(), "/api/user/tokens", createPATReq)
	require.NoError(t, err, "创建 PAT 失败")
	require.NotEmpty(t, createdPAT.PlainToken, "未返回明文令牌")
	t.Logf("  PAT ID: %d, Scopes: %v", createdPAT.Token.ID, createdPAT.Token.Scopes)

	// 注册 PAT 清理
	t.Cleanup(func() {
		_ = userClient.Delete(fmt.Sprintf("/api/user/tokens/%d", createdPAT.Token.ID))
	})

	// 步骤 6: 使用 PAT 测试权限
	t.Log("\n步骤 6: 使用 self scope PAT 测试 API 访问")
	patClient := iam.NewClientFromConfig()
	patClient.SetToken(createdPAT.PlainToken)

	// 6a. 测试 self 域 API（应该成功）
	t.Log("\n  6a. 访问 self 域 API（期望: 200）")
	resp, err = patClient.R().Get("/api/user/profile")
	require.NoError(t, err, "请求失败")
	t.Logf("    GET /api/user/profile -> 状态码: %d", resp.StatusCode())
	assert.Equal(t, 200, resp.StatusCode(), "self scope PAT 应能访问 self:profile:get")

	// 6b. 测试 sys 域 API（应该被拒绝，因为 PAT scope 仅限 self）
	t.Log("\n  6b. 访问 sys 域 API（期望: 403，因为 PAT scope 仅限 self）")
	resp, err = patClient.R().Get("/api/admin/users")
	require.NoError(t, err, "请求失败")
	t.Logf("    GET /api/admin/users -> 状态码: %d", resp.StatusCode())
	assert.Equal(t, 403, resp.StatusCode(), "self scope PAT 不应能访问 sys 域")

	// 步骤 7: 创建 full scope PAT 验证用户完整权限
	t.Log("\n步骤 7: 创建 full scope PAT 验证用户完整权限")
	fullPatName := "full_scope_pat_" + uuid.New().String()[:8]
	createFullPATReq := pat.CreateDTO{
		Name:        fullPatName,
		Scopes:      []string{"full"}, // 完整权限
		Description: "完整 scope 的测试 PAT",
	}

	createdFullPAT, err := apitest.Post[pat.CreateResultDTO](userClient.HTTPClient(), "/api/user/tokens", createFullPATReq)
	require.NoError(t, err, "创建 full scope PAT 失败")
	t.Logf("  Full PAT ID: %d, Scopes: %v", createdFullPAT.Token.ID, createdFullPAT.Token.Scopes)

	t.Cleanup(func() {
		_ = userClient.Delete(fmt.Sprintf("/api/user/tokens/%d", createdFullPAT.Token.ID))
	})

	fullPatClient := iam.NewClientFromConfig()
	fullPatClient.SetToken(createdFullPAT.PlainToken)

	// 7a. full scope 应能访问 self 域
	t.Log("\n  7a. Full scope PAT 访问 self 域 API（期望: 200）")
	resp, err = fullPatClient.R().Get("/api/user/profile")
	require.NoError(t, err, "请求失败")
	t.Logf("    GET /api/user/profile -> 状态码: %d", resp.StatusCode())
	assert.Equal(t, 200, resp.StatusCode(), "full scope PAT 应能访问 self 域")

	// 7b. full scope 应能访问 sys 域（用户有 admin:users:list 权限）
	t.Log("\n  7b. Full scope PAT 访问 sys 域 API（期望: 200）")
	resp, err = fullPatClient.R().Get("/api/admin/users")
	require.NoError(t, err, "请求失败")
	t.Logf("    GET /api/admin/users -> 状态码: %d", resp.StatusCode())
	assert.Equal(t, 200, resp.StatusCode(), "full scope PAT 应能访问 sys 域")

	// 7c. 验证 PAT 不能超过用户权限（用户没有 admin:users:create 权限）
	t.Log("\n  7c. Full scope PAT 创建用户（期望: 403，用户无此权限）")
	createAttempt := user.CreateDTO{
		Username: "should_fail_user",
		Email:    "shouldfail@test.local",
		Password: "test123456",
	}
	resp, err = fullPatClient.R().SetBody(createAttempt).Post("/api/admin/users")
	require.NoError(t, err, "请求失败")
	t.Logf("    POST /api/admin/users -> 状态码: %d", resp.StatusCode())
	assert.Equal(t, 403, resp.StatusCode(), "PAT 不应能超过用户本身权限")

	t.Log("\n===== 普通用户 PAT Scope 权限过滤测试完成 =====")
	t.Log("结论:")
	t.Log("  - self scope PAT: 仅能访问 self 域 API ✓")
	t.Log("  - full scope PAT: 能访问用户权限范围内的所有 API ✓")
	t.Log("  - PAT 永远不能超过用户本身权限 ✓")
}
