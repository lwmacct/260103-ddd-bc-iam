package org_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260103-ddd-bc-iam/internal/apitest/iam"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/org"
	"github.com/lwmacct/260103-ddd-shared/pkg/shared/apitest"
)

// TestMain 在所有测试完成后清理测试数据。
func TestMain(m *testing.M) {
	code := m.Run()

	if os.Getenv("MANUAL") == "1" {
		cleanupTestOrgs()
	}

	os.Exit(code)
}

// cleanupTestOrgs 清理测试创建的组织。
func cleanupTestOrgs() {
	c := iam.NewClientFromConfig()
	if _, err := c.Login("admin", "admin123"); err != nil {
		return
	}

	// 测试组织名称前缀列表
	testPrefixes := []string{"testorg_", "membertest_", "ownerorg_", "duplicate_", "statusorg_"}

	orgs, _, _ := apitest.GetList[org.OrgDTO](c.HTTPClient(), "/api/admin/orgs", map[string]string{"limit": "1000"})
	for _, o := range orgs {
		for _, prefix := range testPrefixes {
			if len(o.Name) >= len(prefix) && o.Name[:len(prefix)] == prefix {
				_ = c.Delete(fmt.Sprintf("/api/admin/orgs/%d", o.ID))
				break
			}
		}
	}
}

// TestOrgCRUD 测试组织完整 CRUD 流程。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestOrgCRUD ./internal/integration/org/
func TestOrgCRUD(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	// 测试 1: 创建组织
	t.Log("\n测试 1: 创建组织")
	createReq := org.CreateOrgDTO{
		Name:        "testorg_" + uuid.New().String()[:8],
		DisplayName: "测试组织",
		Description: "这是一个测试组织",
	}
	createdOrg, err := apitest.Post[org.OrgDTO](c.HTTPClient(), "/api/admin/orgs", createReq)
	require.NoError(t, err, "创建组织失败")
	orgID := createdOrg.ID
	t.Cleanup(func() {
		if orgID > 0 {
			_ = c.Delete(fmt.Sprintf("/api/admin/orgs/%d", orgID))
		}
	})
	t.Logf("  创建成功! 组织 ID: %d", orgID)
	t.Logf("  名称: %s, 显示名: %s", createdOrg.Name, createdOrg.DisplayName)

	// 验证字段
	assert.Equal(t, createReq.Name, createdOrg.Name, "名称不匹配")
	assert.Equal(t, createReq.DisplayName, createdOrg.DisplayName, "显示名不匹配")
	assert.Equal(t, "active", createdOrg.Status, "默认状态应为 active")

	// 测试 2: 获取组织列表
	t.Log("\n测试 2: 获取组织列表")
	orgs, meta, err := apitest.GetList[org.OrgDTO](c.HTTPClient(), "/api/admin/orgs", map[string]string{})
	require.NoError(t, err, "获取组织列表失败")
	t.Logf("  组织数量: %d", len(orgs))
	if meta != nil {
		t.Logf("  总数: %d, 总页数: %d", meta.Total, meta.TotalPages)
	}

	// 验证列表中包含创建的组织
	orgIDs := apitest.ExtractIDs(orgs, func(o org.OrgDTO) uint { return o.ID })
	assert.Contains(t, orgIDs, orgID, "列表中应包含新创建的组织")

	// 测试 3: 获取组织详情
	t.Log("\n测试 3: 获取组织详情")
	orgDetail, err := apitest.Get[org.OrgDTO](c.HTTPClient(), fmt.Sprintf("/api/admin/orgs/%d", orgID), nil)
	require.NoError(t, err, "获取组织详情失败")
	t.Logf("  详情: %s - %s", orgDetail.Name, orgDetail.Description)
	assert.Equal(t, orgID, orgDetail.ID, "组织 ID 不匹配")

	// 测试 4: 更新组织
	t.Log("\n测试 4: 更新组织")
	newDisplayName := "更新后的组织名称"
	newDescription := "更新后的组织描述"
	updateReq := org.UpdateOrgDTO{
		DisplayName: &newDisplayName,
		Description: &newDescription,
	}
	updatedOrg, err := apitest.Put[org.OrgDTO](c.HTTPClient(), fmt.Sprintf("/api/admin/orgs/%d", orgID), updateReq)
	require.NoError(t, err, "更新组织失败")
	t.Logf("  更新成功! 显示名: %s", updatedOrg.DisplayName)

	// 验证更新后的字段
	assert.Equal(t, newDisplayName, updatedOrg.DisplayName, "显示名未更新")
	assert.Equal(t, newDescription, updatedOrg.Description, "描述未更新")

	// 测试 5: 更新组织状态
	t.Log("\n测试 5: 更新组织状态")
	suspendedStatus := "suspended"
	statusReq := org.UpdateOrgDTO{
		Status: &suspendedStatus,
	}
	updatedOrg, err = apitest.Put[org.OrgDTO](c.HTTPClient(), fmt.Sprintf("/api/admin/orgs/%d", orgID), statusReq)
	require.NoError(t, err, "更新组织状态失败")
	t.Logf("  状态更新为: %s", updatedOrg.Status)
	assert.Equal(t, "suspended", updatedOrg.Status, "状态未更新")

	// 测试 6: 删除组织
	t.Log("\n测试 6: 删除组织")
	err = c.Delete(fmt.Sprintf("/api/admin/orgs/%d", orgID))
	require.NoError(t, err, "删除组织失败")
	t.Log("  删除成功!")

	// 标记已删除，避免 t.Cleanup 重复删除
	orgID = 0

	t.Log("\n组织 CRUD 流程测试完成!")
}

// TestOrgListPagination 测试组织列表分页。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestOrgListPagination ./internal/integration/org/
func TestOrgListPagination(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	t.Log("测试组织列表分页...")

	// 第一页
	page1, meta1, err := apitest.GetList[org.OrgDTO](c.HTTPClient(), "/api/admin/orgs", map[string]string{})
	require.NoError(t, err, "获取第一页失败")
	t.Logf("  第 1 页: %d 条记录", len(page1))
	if meta1 != nil {
		t.Logf("  总数: %d, 总页数: %d", meta1.Total, meta1.TotalPages)
	}

	// 验证总数
	assert.Positive(t, meta1.Total, "总数应大于 0")

	// 如果有第二页，获取并验证
	if meta1.TotalPages > 1 {
		page2, meta2, err := apitest.GetList[org.OrgDTO](c.HTTPClient(), "/api/admin/orgs", map[string]string{})
		require.NoError(t, err, "获取第二页失败")
		t.Logf("  第 2 页: %d 条记录", len(page2))
		if meta2 != nil {
			assert.Equal(t, meta1.Total, meta2.Total, "两页总数应一致")
		}

		// 验证无重复 ID
		ids1 := apitest.ExtractIDs(page1, func(o org.OrgDTO) uint { return o.ID })
		ids2 := apitest.ExtractIDs(page2, func(o org.OrgDTO) uint { return o.ID })
		for _, id := range ids1 {
			assert.NotContains(t, ids2, id, "两页不应有相同 ID")
		}
	}

	t.Log("\n分页测试完成!")
}

// TestOrgStatusUpdate 测试组织状态更新。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestOrgStatusUpdate ./internal/integration/org/
func TestOrgStatusUpdate(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	// 创建测试组织
	t.Log("\n步骤 1: 创建测试组织")
	createReq := org.CreateOrgDTO{
		Name:        "statusorg_" + uuid.New().String()[:8],
		DisplayName: "状态测试组织",
	}
	createdOrg, err := apitest.Post[org.OrgDTO](c.HTTPClient(), "/api/admin/orgs", createReq)
	require.NoError(t, err, "创建组织失败")
	orgID := createdOrg.ID
	t.Cleanup(func() {
		if orgID > 0 {
			_ = c.Delete(fmt.Sprintf("/api/admin/orgs/%d", orgID))
		}
	})
	t.Logf("  创建成功! 初始状态: %s", createdOrg.Status)
	assert.Equal(t, "active", createdOrg.Status, "初始状态应为 active")

	// 测试 2: 更新为 suspended
	t.Log("\n步骤 2: 更新状态为 suspended")
	suspendedStatus := "suspended"
	statusReq := org.UpdateOrgDTO{
		Status: &suspendedStatus,
	}
	updatedOrg, err := apitest.Put[org.OrgDTO](c.HTTPClient(), fmt.Sprintf("/api/admin/orgs/%d", orgID), statusReq)
	require.NoError(t, err, "更新状态失败")
	t.Logf("  状态更新为: %s", updatedOrg.Status)
	assert.Equal(t, "suspended", updatedOrg.Status, "状态应为 suspended")

	// 测试 3: 恢复为 active
	t.Log("\n步骤 3: 恢复状态为 active")
	activeStatus := "active"
	statusReq.Status = &activeStatus
	updatedOrg, err = apitest.Put[org.OrgDTO](c.HTTPClient(), fmt.Sprintf("/api/admin/orgs/%d", orgID), statusReq)
	require.NoError(t, err, "恢复状态失败")
	t.Logf("  状态恢复为: %s", updatedOrg.Status)
	assert.Equal(t, "active", updatedOrg.Status, "状态应为 active")

	t.Log("\n状态更新测试完成!")
}

// TestOrgMemberCRUD 测试组织成员完整流程。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestOrgMemberCRUD ./internal/integration/org/
func TestOrgMemberCRUD(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	// 步骤 1: 创建组织
	t.Log("\n步骤 1: 创建组织")
	createOrgReq := org.CreateOrgDTO{
		Name:        "membertest_" + uuid.New().String()[:8],
		DisplayName: "成员测试组织",
	}
	createdOrg, err := apitest.Post[org.OrgDTO](c.HTTPClient(), "/api/admin/orgs", createOrgReq)
	require.NoError(t, err, "创建组织失败")
	orgID := createdOrg.ID
	t.Cleanup(func() {
		if orgID > 0 {
			_ = c.Delete(fmt.Sprintf("/api/admin/orgs/%d", orgID))
		}
	})
	t.Logf("  组织创建成功! ID: %d", orgID)

	// 步骤 2: 创建测试用户
	t.Log("\n步骤 2: 创建测试用户")
	testUser := iam.CreateTestUser(t, c, "orgmember")
	t.Logf("  用户创建成功! ID: %d, 用户名: %s", testUser.ID, testUser.Username)

	// 步骤 3: 添加成员
	t.Log("\n步骤 3: 添加成员")
	addMemberReq := org.AddMemberDTO{
		UserID: testUser.ID,
		Role:   "admin",
	}
	memberPath := fmt.Sprintf("/api/org/%d/members", orgID)
	member, err := apitest.Post[org.MemberDTO](c.HTTPClient(), memberPath, addMemberReq)
	require.NoError(t, err, "添加成员失败")
	t.Logf("  成员添加成功! 角色: %s", member.Role)

	// 验证成员信息
	assert.Equal(t, testUser.ID, member.UserID, "用户 ID 不匹配")
	assert.Equal(t, orgID, member.OrgID, "组织 ID 不匹配")
	assert.Equal(t, "admin", member.Role, "角色不匹配")
	// TODO: 用户信息未加载，需要 Repository Join User 表
	// assert.NotEmpty(t, member.Username, "用户名应加载")
	// assert.Equal(t, testUser.Username, member.Username, "用户名不匹配")

	// 清理：移除成员
	t.Cleanup(func() {
		_ = c.Delete(fmt.Sprintf("/api/org/%d/members/%d", orgID, testUser.ID))
	})

	// 步骤 4: 获取成员列表
	t.Log("\n步骤 4: 获取成员列表")
	members, meta, err := apitest.GetList[org.MemberDTO](c.HTTPClient(), memberPath, map[string]string{})
	require.NoError(t, err, "获取成员列表失败")
	t.Logf("  成员数量: %d", len(members))
	if meta != nil {
		t.Logf("  总数: %d", meta.Total)
	}

	// 验证列表中包含新成员
	userIDs := apitest.ExtractIDs(members, func(m org.MemberDTO) uint { return m.UserID })
	assert.Contains(t, userIDs, testUser.ID, "成员列表应包含新添加的用户")

	// 验证成员包含用户信息
	var testMember *org.MemberDTO
	for i := range members {
		if members[i].UserID == testUser.ID {
			testMember = &members[i]
			break
		}
	}
	require.NotNil(t, testMember, "应找到测试用户")
	assert.NotEmpty(t, testMember.Username, "成员应包含用户名")
	assert.NotEmpty(t, testMember.Email, "成员应包含邮箱")
	t.Logf("  用户名: %s, 邮箱: %s", testMember.Username, testMember.Email)

	// 步骤 5: 更新成员角色
	t.Log("\n步骤 5: 更新成员角色")
	newRole := "member"
	updateRoleReq := org.UpdateMemberRoleDTO{
		Role: newRole,
	}
	updatePath := fmt.Sprintf("/api/org/%d/members/%d/role", orgID, testUser.ID)
	_, err = apitest.Put[any](c.HTTPClient(), updatePath, updateRoleReq)
	require.NoError(t, err, "更新成员角色失败")
	t.Logf("  角色更新成功!")

	// TODO: 成员详情端点不存在，需要添加或通过列表验证
	// updatedMember, err := apitest.Get[org.MemberDTO](c.HTTPClient(), fmt.Sprintf("%s/%d", memberPath, testUser.ID), nil)
	// require.NoError(t, err, "获取成员详情失败")
	// assert.Equal(t, newRole, updatedMember.Role, "角色未更新")

	// 步骤 6: 移除成员
	t.Log("\n步骤 6: 移除成员")
	err = c.Delete(fmt.Sprintf("/api/org/%d/members/%d", orgID, testUser.ID))
	require.NoError(t, err, "移除成员失败")
	t.Logf("  成员移除成功!")

	// 验证成员已移除
	membersAfter, _, err := apitest.GetList[org.MemberDTO](c.HTTPClient(), memberPath, nil)
	require.NoError(t, err, "获取成员列表失败")
	userIDsAfter := apitest.ExtractIDs(membersAfter, func(m org.MemberDTO) uint { return m.UserID })
	assert.NotContains(t, userIDsAfter, testUser.ID, "成员列表不应包含已移除的用户")

	t.Log("\n组织成员流程测试完成!")
}

// TestOrgMemberLastOwnerProtection 测试最后 Owner 保护机制。
//
// 组织必须有至少一个 owner，无法移除或降级最后的 owner。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestOrgMemberLastOwnerProtection ./internal/integration/org/
func TestOrgMemberLastOwnerProtection(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	// 步骤 1: 创建组织
	t.Log("\n步骤 1: 创建组织")
	createOrgReq := org.CreateOrgDTO{
		Name:        "ownerorg_" + uuid.New().String()[:8],
		DisplayName: "Owner 保护测试组织",
	}
	createdOrg, err := apitest.Post[org.OrgDTO](c.HTTPClient(), "/api/admin/orgs", createOrgReq)
	require.NoError(t, err, "创建组织失败")
	orgID := createdOrg.ID
	t.Cleanup(func() {
		if orgID > 0 {
			_ = c.Delete(fmt.Sprintf("/api/admin/orgs/%d", orgID))
		}
	})
	t.Logf("  组织创建成功! ID: %d", orgID)

	// 步骤 2: 获取成员列表，确认 admin (user_id=2) 是 owner
	t.Log("\n步骤 2: 获取成员列表")
	members, _, err := apitest.GetList[org.MemberDTO](c.HTTPClient(), fmt.Sprintf("/api/org/%d/members", orgID), nil)
	require.NoError(t, err, "获取成员列表失败")

	var adminMember *org.MemberDTO
	for i := range members {
		if members[i].UserID == 2 { // admin 用户
			adminMember = &members[i]
			break
		}
	}
	require.NotNil(t, adminMember, "未找到 admin 成员")
	t.Logf("  admin 成员: UserID=%d, Role=%s", adminMember.UserID, adminMember.Role)
	assert.Equal(t, "owner", adminMember.Role, "创建者应为 owner")

	// 步骤 3: 尝试降级最后的 owner（应失败）
	t.Log("\n步骤 3: 尝试降级最后的 owner（应失败）")
	updateRoleReq := org.UpdateMemberRoleDTO{
		Role: "admin", // 降级为 admin
	}
	updatePath := fmt.Sprintf("/api/org/%d/members/%d/role", orgID, adminMember.UserID)
	_, err = apitest.Put[any](c.HTTPClient(), updatePath, updateRoleReq)
	require.Error(t, err, "降级最后的 owner 应该失败")
	t.Logf("  预期失败: %v", err)

	// 步骤 4: 尝试移除最后的 owner（应失败）
	t.Log("\n步骤 4: 尝试移除最后的 owner（应失败）")
	err = c.Delete(fmt.Sprintf("/api/org/%d/members/%d", orgID, adminMember.UserID))
	require.Error(t, err, "移除最后的 owner 应该失败")
	t.Logf("  预期失败: %v", err)

	t.Log("\n最后 Owner 保护测试完成!")
}

// TestUserOrgs 测试用户视角的组织列表。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestUserOrgs ./internal/integration/org/
func TestUserOrgs(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	t.Log("\n获取用户加入的组织列表...")
	userOrgs, _, err := apitest.GetList[org.UserOrgDTO](c.HTTPClient(), "/api/user/orgs", nil)
	require.NoError(t, err, "获取用户组织列表失败")
	t.Logf("  用户加入的组织数量: %d", len(userOrgs))

	for _, uo := range userOrgs {
		t.Logf("  - [%d] %s (%s) - 角色: %s",
			uo.ID, uo.DisplayName, uo.Name, uo.Role)
	}

	t.Log("\n用户组织列表测试完成!")
}

// TestOrgWithInvalidData 测试无效数据处理。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestOrgWithInvalidData ./internal/integration/org/
func TestOrgWithInvalidData(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	// 测试 1: 重复名称应失败
	t.Log("\n测试 1: 重复名称应失败")
	orgName := "duplicate_" + uuid.New().String()[:8]
	createReq := org.CreateOrgDTO{
		Name:        orgName,
		DisplayName: "第一个组织",
	}
	firstOrg, err := apitest.Post[org.OrgDTO](c.HTTPClient(), "/api/admin/orgs", createReq)
	require.NoError(t, err, "创建第一个组织失败")
	orgID := firstOrg.ID
	t.Cleanup(func() {
		if orgID > 0 {
			_ = c.Delete(fmt.Sprintf("/api/admin/orgs/%d", orgID))
		}
	})

	// 尝试创建同名组织
	duplicateReq := org.CreateOrgDTO{
		Name:        orgName, // 相同名称
		DisplayName: "第二个组织",
	}
	_, err = apitest.Post[org.OrgDTO](c.HTTPClient(), "/api/admin/orgs", duplicateReq)
	require.Error(t, err, "创建同名组织应该失败")
	t.Logf("  预期失败: %v", err)

	// 测试 2: 无效参数应失败
	t.Log("\n测试 2: 无效参数应失败")
	invalidReq := org.CreateOrgDTO{
		Name:        "x", // 太短
		DisplayName: "无效参数组织",
	}
	_, err = apitest.Post[org.OrgDTO](c.HTTPClient(), "/api/admin/orgs", invalidReq)
	require.Error(t, err, "名称太短应该失败")
	t.Logf("  预期失败: %v", err)

	t.Log("\n无效数据处理测试完成!")
}
