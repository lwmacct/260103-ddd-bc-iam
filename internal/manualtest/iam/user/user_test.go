package user_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/application/user"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/transport/gin/manualtest"
)

// TestMain 在所有测试完成后清理测试数据。
func TestMain(m *testing.M) {
	code := m.Run()

	if os.Getenv("MANUAL") == "1" {
		cleanupTestUsers()
	}

	os.Exit(code)
}

// cleanupTestUsers 清理测试创建的用户。
func cleanupTestUsers() {
	c := manualtest.NewClient()
	if _, err := c.Login("admin", "admin123"); err != nil {
		return
	}

	// 测试用户名前缀列表（factory.go 创建的用户）
	userPrefixes := []string{"teammember_", "nonorgmember_", "orgmember_", "roletest_", "testuser_"}

	users, _, _ := manualtest.GetList[user.UserDTO](c, "/api/admin/users", map[string]string{"limit": "1000"})
	for _, u := range users {
		// 跳过系统用户
		if u.ID <= 2 {
			continue
		}
		for _, prefix := range userPrefixes {
			if len(u.Username) >= len(prefix) && u.Username[:len(prefix)] == prefix {
				_ = c.Delete(fmt.Sprintf("/api/admin/users/%d", u.ID))
				break
			}
		}
	}
}

// TestAdminUsersFlow 用户管理完整流程测试。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestAdminUsersFlow ./internal/integration/user/
func TestAdminUsersFlow(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 测试 1: 获取用户列表
	t.Log("\n测试 1: 获取用户列表")
	users, meta, err := manualtest.GetList[user.UserDTO](c, "/api/admin/users", map[string]string{
		"page":  "1",
		"limit": "10",
	})
	require.NoError(t, err, "获取用户列表失败")
	t.Logf("  用户数量: %d", len(users))
	if meta != nil {
		t.Logf("  总数: %d", meta.Total)
	}

	// 测试 2: 创建用户（使用工厂函数，返回清理控制）
	t.Log("\n测试 2: 创建用户")
	testUser, markDeleted := manualtest.CreateTestUserWithCleanupControl(t, c, "testuser")
	t.Logf("  创建成功! 用户 ID: %d", testUser.ID)

	// 验证创建的用户数据
	assert.NotEmpty(t, testUser.Username, "用户名不应为空")
	assert.NotEmpty(t, testUser.Email, "邮箱不应为空")

	// 测试 3: 获取用户详情
	t.Log("\n测试 3: 获取用户详情")
	userDetail, err := manualtest.Get[user.UserDTO](c, fmt.Sprintf("/api/admin/users/%d", testUser.ID), nil)
	require.NoError(t, err, "获取用户详情失败")
	t.Logf("  用户名: %s, 邮箱: %s", userDetail.Username, userDetail.Email)

	// 验证用户详情
	assert.Equal(t, testUser.ID, userDetail.ID, "用户 ID 不匹配")

	// 测试 4: 更新用户
	t.Log("\n测试 4: 更新用户")
	newFullName := "测试用户（已更新）"
	updateReq := user.UpdateDTO{
		RealName: &newFullName,
	}
	updatedUser, err := manualtest.Put[user.UserDTO](c, fmt.Sprintf("/api/admin/users/%d", testUser.ID), updateReq)
	require.NoError(t, err, "更新用户失败")
	t.Logf("  更新成功! 真实姓名: %s", updatedUser.RealName)

	// 验证更新后的字段
	assert.Equal(t, newFullName, updatedUser.RealName, "真实姓名未更新")

	// 测试 5: 删除用户
	t.Log("\n测试 5: 删除用户")
	err = c.Delete(fmt.Sprintf("/api/admin/users/%d", testUser.ID))
	require.NoError(t, err, "删除用户失败")
	t.Log("  删除成功!")

	// 标记已删除，避免 t.Cleanup 重复删除
	markDeleted()

	t.Log("\n用户管理流程测试完成!")
}

// TestListUsers 测试获取用户列表。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestListUsers ./internal/integration/user/
func TestListUsers(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	t.Log("获取用户列表...")
	users, meta, err := manualtest.GetList[user.UserDTO](c, "/api/admin/users", map[string]string{
		"page":  "1",
		"limit": "10",
	})
	require.NoError(t, err, "获取用户列表失败")

	t.Logf("用户数量: %d", len(users))
	if meta != nil {
		t.Logf("总数: %d, 总页数: %d", meta.Total, meta.TotalPages)
	}

	for _, u := range users {
		t.Logf("  - [%d] %s <%s> 状态: %s", u.ID, u.Username, u.Email, u.Status)
	}
}

// TestAssignRoles 测试分配用户角色。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestAssignRoles ./internal/integration/user/
func TestAssignRoles(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 创建测试用户
	t.Log("\n步骤 1: 创建测试用户")
	testUser := manualtest.CreateTestUser(t, c, "roletest")
	t.Logf("  创建成功! 用户 ID: %d", testUser.ID)

	// 分配角色（使用 user 角色 ID=2）
	t.Log("\n步骤 2: 分配角色")
	assignReq := user.AssignRolesDTO{
		RoleIDs: []uint{2}, // user 角色
	}
	t.Logf("  分配角色 IDs: %v", assignReq.RoleIDs)

	assignResp, err := manualtest.Put[user.UserWithRolesDTO](c, fmt.Sprintf("/api/admin/users/%d/roles", testUser.ID), assignReq)
	require.NoError(t, err, "分配角色失败")

	t.Logf("  分配成功! 用户现有角色数: %d", len(assignResp.Roles))
	for _, r := range assignResp.Roles {
		t.Logf("    - [%d] %s (%s)", r.ID, r.DisplayName, r.Name)
	}

	// 验证角色已分配
	require.NotEmpty(t, assignResp.Roles, "角色分配失败，用户没有角色")

	// 使用 assert.Contains 验证是否包含指定的角色 ID
	roleIDs := manualtest.ExtractIDs(assignResp.Roles, func(r user.RoleDTO) uint { return r.ID })
	assert.Contains(t, roleIDs, uint(2), "未找到预期的角色 ID=2")

	t.Log("\n角色分配测试完成!")
}

// TestBatchCreateUsers 测试批量创建用户。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestBatchCreateUsers ./internal/integration/user/
func TestBatchCreateUsers(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	timestamp := time.Now().Unix()
	username1 := fmt.Sprintf("batch1_%d", timestamp)
	username2 := fmt.Sprintf("batch2_%d", timestamp)

	// 确保测试结束时清理资源
	t.Cleanup(func() {
		// 使用搜索找到并删除测试用户（因为用户可能在列表末尾）
		for _, username := range []string{username1, username2} {
			searchParams := map[string]string{"search": username}
			users, _, _ := manualtest.GetList[user.UserWithRolesDTO](c, "/api/admin/users", searchParams)
			for _, u := range users {
				if u.Username == username {
					_ = c.Delete(fmt.Sprintf("/api/admin/users/%d", u.ID))
				}
			}
		}
	})

	// 步骤 1: 批量创建用户（2个成功 + 1个重复失败）
	t.Log("\n步骤 1: 批量创建用户")
	t.Logf("  用户1: %s", username1)
	t.Logf("  用户2: %s", username2)
	t.Logf("  用户3: %s (重复，应失败)", username1)

	batchReq := user.BatchCreateDTO{
		Users: []user.BatchItemDTO{
			{
				Username: username1,
				Email:    username1 + "@example.com",
				Password: "test123456",
				RealName: "批量用户1",
			},
			{
				Username: username2,
				Email:    username2 + "@example.com",
				Password: "test123456",
				RealName: "批量用户2",
			},
			{
				Username: username1, // 重复用户名
				Email:    "dup_" + username1 + "@example.com",
				Password: "test123456",
				RealName: "重复用户",
			},
		},
	}

	result, err := manualtest.Post[user.BatchCreateResultDTO](c, "/api/admin/users/batch", batchReq)
	require.NoError(t, err, "批量创建请求失败")

	t.Logf("\n批量创建结果:")
	t.Logf("  总数: %d", result.Total)
	t.Logf("  成功: %d", result.Success)
	t.Logf("  失败: %d", result.Failed)

	// 步骤 2: 验证结果
	t.Log("\n步骤 2: 验证结果")
	assert.Equal(t, 3, result.Total, "总数应为 3")
	assert.Equal(t, 2, result.Success, "成功数应为 2")
	assert.Equal(t, 1, result.Failed, "失败数应为 1")

	// 验证错误详情
	if len(result.Errors) > 0 {
		t.Log("  错误详情:")
		for _, e := range result.Errors {
			t.Logf("    - [%d] %s: %s", e.Index, e.Username, e.Error)
		}
	}

	// 步骤 3: 验证用户已创建
	// 注意：由于数据库有大量用户，新创建的用户 ID 最大，在列表末尾。
	// 使用搜索来验证用户是否创建成功，而不是依赖分页列表。
	t.Log("\n步骤 3: 验证用户已创建（通过搜索）")

	// 使用搜索端点查找创建的用户
	searchParams1 := map[string]string{"search": username1}
	users1, _, _ := manualtest.GetList[user.UserWithRolesDTO](c, "/api/admin/users", searchParams1)
	assert.NotEmpty(t, users1, "应该能搜索到用户1")
	assert.Equal(t, username1, users1[0].Username, "用户1用户名匹配")

	searchParams2 := map[string]string{"search": username2}
	users2, _, _ := manualtest.GetList[user.UserWithRolesDTO](c, "/api/admin/users", searchParams2)
	assert.NotEmpty(t, users2, "应该能搜索到用户2")
	assert.Equal(t, username2, users2[0].Username, "用户2用户名匹配")

	t.Log("\n批量创建用户测试完成!")
}

// TestSystemUserProtection 测试系统用户保护机制。
//
// 系统用户（root、admin）不可删除，部分字段不可修改。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestSystemUserProtection ./internal/integration/user/
func TestSystemUserProtection(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 获取系统用户
	t.Log("\n步骤 1: 获取系统用户")
	users, _, err := manualtest.GetList[user.UserDTO](c, "/api/admin/users", nil)
	require.NoError(t, err, "获取用户列表失败")

	var rootUser, adminUser *user.UserDTO
	for i := range users {
		if users[i].Username == "root" {
			rootUser = &users[i]
		}
		if users[i].Username == "admin" {
			adminUser = &users[i]
		}
	}
	require.NotNil(t, rootUser, "未找到 root 用户")
	require.NotNil(t, adminUser, "未找到 admin 用户")

	t.Logf("  root 用户: ID=%d, Type=%s", rootUser.ID, rootUser.Type)
	t.Logf("  admin 用户: ID=%d, Type=%s", adminUser.ID, adminUser.Type)

	// 验证系统用户标记
	assert.Equal(t, "system", rootUser.Type, "root 应为 system 类型")
	assert.Equal(t, "system", adminUser.Type, "admin 应为 system 类型")

	// 测试 2: 尝试删除 root 用户（应失败）
	t.Log("\n步骤 2: 尝试删除 root 用户（应失败）")
	err = c.Delete(fmt.Sprintf("/api/admin/users/%d", rootUser.ID))
	require.Error(t, err, "删除 root 用户应该失败")
	t.Logf("  预期失败: %v", err)

	// 测试 3: 尝试删除 admin 用户（应失败）
	t.Log("\n步骤 3: 尝试删除 admin 用户（应失败）")
	err = c.Delete(fmt.Sprintf("/api/admin/users/%d", adminUser.ID))
	require.Error(t, err, "删除 admin 用户应该失败")
	t.Logf("  预期失败: %v", err)

	// 测试 4: 尝试修改 root 用户名（应失败）
	t.Log("\n步骤 4: 尝试修改 root 用户名（应失败）")
	newUsername := "root_renamed"
	updateReq := user.UpdateDTO{
		Username: &newUsername,
	}
	_, err = manualtest.Put[user.UserDTO](c, fmt.Sprintf("/api/admin/users/%d", rootUser.ID), updateReq)
	require.Error(t, err, "修改 root 用户名应该失败")
	t.Logf("  预期失败: %v", err)

	// 测试 5: 尝试修改 root 状态（应失败）
	t.Log("\n步骤 5: 尝试修改 root 状态（应失败）")
	inactiveStatus := "inactive"
	statusReq := user.UpdateDTO{
		Status: &inactiveStatus,
	}
	_, err = manualtest.Put[user.UserDTO](c, fmt.Sprintf("/api/admin/users/%d", rootUser.ID), statusReq)
	require.Error(t, err, "修改 root 状态应该失败")
	t.Logf("  预期失败: %v", err)

	// 测试 6: 修改 admin 邮箱（应成功，非保护字段）
	t.Log("\n步骤 6: 修改 admin 邮箱（应成功）")
	newEmail := "admin_updated@example.com"
	emailReq := user.UpdateDTO{
		Email: &newEmail,
	}
	updatedAdmin, err := manualtest.Put[user.UserDTO](c, fmt.Sprintf("/api/admin/users/%d", adminUser.ID), emailReq)
	require.NoError(t, err, "修改 admin 邮箱应该成功")
	t.Logf("  修改成功! 新邮箱: %s", updatedAdmin.Email)

	// 恢复 admin 邮箱
	originalEmail := "admin@example.com"
	restoreReq := user.UpdateDTO{
		Email: &originalEmail,
	}
	_, _ = manualtest.Put[user.UserDTO](c, fmt.Sprintf("/api/admin/users/%d", adminUser.ID), restoreReq)

	t.Log("\n系统用户保护测试完成!")
}

// TestRootRoleProtection 测试 root 用户角色保护。
//
// root 用户的角色不可修改（始终拥有 *:*:* 权限）。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestRootRoleProtection ./internal/integration/user/
func TestRootRoleProtection(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 获取 root 用户
	t.Log("\n步骤 1: 获取 root 用户")
	users, _, err := manualtest.GetList[user.UserDTO](c, "/api/admin/users", nil)
	require.NoError(t, err, "获取用户列表失败")

	var rootUser *user.UserDTO
	for i := range users {
		if users[i].Username == "root" {
			rootUser = &users[i]
			break
		}
	}
	require.NotNil(t, rootUser, "未找到 root 用户")
	t.Logf("  root 用户 ID: %d", rootUser.ID)

	// 测试: 尝试修改 root 用户角色（应失败）
	t.Log("\n步骤 2: 尝试修改 root 用户角色（应失败）")
	assignReq := user.AssignRolesDTO{
		RoleIDs: []uint{2}, // user 角色
	}
	_, err = manualtest.Put[user.UserWithRolesDTO](c, fmt.Sprintf("/api/admin/users/%d/roles", rootUser.ID), assignReq)
	require.Error(t, err, "修改 root 用户角色应该失败")
	t.Logf("  预期失败: %v", err)

	t.Log("\nroot 用户角色保护测试完成!")
}
