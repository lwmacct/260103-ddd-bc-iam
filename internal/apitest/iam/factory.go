package iam

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/lwmacct/260103-ddd-shared/pkg/shared/apitest"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/role"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/user"
)

// CreateTestUser 创建测试用户并自动注册清理。
// 测试结束时会自动删除创建的用户。
// username 使用完整 UUID 确保唯一性，避免并行测试冲突。
func CreateTestUser(t *testing.T, c *Client, prefix string) *user.UserWithRolesDTO {
	t.Helper()

	username := fmt.Sprintf("%s_%s", prefix, uuid.New().String())
	req := user.CreateDTO{
		Username:  username,
		Email:     username + "@test.local",
		Password:  "test123456",
		RealName:  "测试用户",
		Nickname:  "测试",
		Phone:     "13800138000",
		Signature: "这是我的个性签名",
	}

	resp, err := apitest.Post[user.UserWithRolesDTO](&c.Client, "/api/admin/users", req)
	require.NoError(t, err, "创建测试用户失败: %s", username)

	userID := resp.ID
	t.Cleanup(func() {
		if userID > 0 {
			_ = c.Delete(fmt.Sprintf("/api/admin/users/%d", userID))
		}
	})

	return resp
}

// CreateTestRole 创建测试角色并自动注册清理。
// 测试结束时会自动删除创建的角色。
// rolename 使用完整 UUID 确保唯一性，避免并行测试冲突。
func CreateTestRole(t *testing.T, c *Client, prefix string) *role.CreateResultDTO {
	t.Helper()

	roleName := fmt.Sprintf("%s_%s", prefix, uuid.New().String())
	req := role.CreateDTO{
		Name:        roleName,
		DisplayName: "测试角色",
		Description: "自动创建的测试角色",
	}

	resp, err := apitest.Post[role.CreateResultDTO](&c.Client, "/api/admin/roles", req)
	require.NoError(t, err, "创建测试角色失败: %s", roleName)

	roleID := resp.RoleID
	t.Cleanup(func() {
		if roleID > 0 {
			_ = c.Delete(fmt.Sprintf("/api/admin/roles/%d", roleID))
		}
	})

	return resp
}

// CreateTestUserWithCleanupControl 创建测试用户，返回清理控制函数。
// 当测试本身需要删除用户时使用，删除成功后调用返回的 markDeleted 函数。
// username 使用完整 UUID 确保唯一性，避免并行测试冲突。
func CreateTestUserWithCleanupControl(t *testing.T, c *Client, prefix string) (*user.UserWithRolesDTO, func()) {
	t.Helper()

	username := fmt.Sprintf("%s_%s", prefix, uuid.New().String())
	req := user.CreateDTO{
		Username:  username,
		Email:     username + "@test.local",
		Password:  "test123456",
		RealName:  "测试用户",
		Nickname:  "测试",
		Phone:     "13800138000",
		Signature: "这是我的个性签名",
	}

	result, err := apitest.Post[user.UserWithRolesDTO](&c.Client, "/api/admin/users", req)
	require.NoError(t, err, "创建测试用户失败: %s", username)

	userID := result.ID
	deleted := false

	t.Cleanup(func() {
		if !deleted && userID > 0 {
			_ = c.Delete(fmt.Sprintf("/api/admin/users/%d", userID))
		}
	})

	return result, func() { deleted = true }
}

// CreateTestRoleWithCleanupControl 创建测试角色，返回清理控制函数。
// 当测试本身需要删除角色时使用，删除成功后调用返回的 markDeleted 函数。
// rolename 使用完整 UUID 确保唯一性，避免并行测试冲突。
func CreateTestRoleWithCleanupControl(t *testing.T, c *Client, prefix string) (*role.CreateResultDTO, func()) {
	t.Helper()

	roleName := fmt.Sprintf("%s_%s", prefix, uuid.New().String())
	req := role.CreateDTO{
		Name:        roleName,
		DisplayName: "测试角色",
		Description: "自动创建的测试角色",
	}

	result, err := apitest.Post[role.CreateResultDTO](&c.Client, "/api/admin/roles", req)
	require.NoError(t, err, "创建测试角色失败: %s", roleName)

	roleID := result.RoleID
	deleted := false

	t.Cleanup(func() {
		if !deleted && roleID > 0 {
			_ = c.Delete(fmt.Sprintf("/api/admin/roles/%d", roleID))
		}
	})

	return result, func() { deleted = true }
}
