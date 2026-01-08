package org_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260103-ddd-iam-bc/internal/apitest/iam"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/org"
	"github.com/lwmacct/260103-ddd-shared/pkg/shared/apitest"
)

// 种子数据: acme org (ID=1), admin 是 owner
const testOrgID uint = 1

// teamBasePath 返回团队 API 基础路径
func teamBasePath(orgID uint) string {
	return fmt.Sprintf("/api/org/%d/teams", orgID)
}

// teamPath 返回单个团队 API 路径
func teamPath(orgID, teamID uint) string {
	return fmt.Sprintf("%s/%d", teamBasePath(orgID), teamID)
}

// teamMembersBasePath 返回团队成员 API 基础路径
func teamMembersBasePath(orgID, teamID uint) string {
	return fmt.Sprintf("/api/org/%d/teams/%d/members", orgID, teamID)
}

// teamMemberPath 返回单个团队成员 API 路径
func teamMemberPath(orgID, teamID, userID uint) string {
	return fmt.Sprintf("%s/%d", teamMembersBasePath(orgID, teamID), userID)
}

// TestTeamCRUD 测试团队完整 CRUD 流程。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestTeamCRUD ./internal/integration/org/
func TestTeamCRUD(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	// 测试 1: 创建团队
	t.Log("\n测试 1: 创建团队")
	createReq := org.CreateTeamDTO{
		Name:        "testteam_" + uuid.New().String()[:8],
		DisplayName: "测试团队",
		Description: "这是一个测试团队",
	}
	createdTeam, err := apitest.Post[org.TeamDTO](c.HTTPClient(), teamBasePath(testOrgID), createReq)
	require.NoError(t, err, "创建团队失败")
	teamID := createdTeam.ID
	t.Cleanup(func() {
		if teamID > 0 {
			_ = c.Delete(teamPath(testOrgID, teamID))
		}
	})
	t.Logf("  创建成功! 团队 ID: %d", teamID)
	t.Logf("  名称: %s, 显示名: %s", createdTeam.Name, createdTeam.DisplayName)

	// 验证字段
	assert.Equal(t, createReq.Name, createdTeam.Name, "名称不匹配")
	assert.Equal(t, createReq.DisplayName, createdTeam.DisplayName, "显示名不匹配")
	assert.Equal(t, testOrgID, createdTeam.OrgID, "组织 ID 不匹配")

	// 测试 2: 获取团队列表
	t.Log("\n测试 2: 获取团队列表")
	teams, meta, err := apitest.GetList[org.TeamDTO](c.HTTPClient(), teamBasePath(testOrgID), map[string]string{})
	require.NoError(t, err, "获取团队列表失败")
	t.Logf("  团队数量: %d", len(teams))
	if meta != nil {
		t.Logf("  总数: %d, 总页数: %d", meta.Total, meta.TotalPages)
	}

	// 验证列表中包含创建的团队
	teamIDs := apitest.ExtractIDs(teams, func(tm org.TeamDTO) uint { return tm.ID })
	assert.Contains(t, teamIDs, teamID, "列表中应包含新创建的团队")

	// 测试 3: 获取团队详情
	t.Log("\n测试 3: 获取团队详情")
	teamDetail, err := apitest.Get[org.TeamDTO](c.HTTPClient(), teamPath(testOrgID, teamID), nil)
	require.NoError(t, err, "获取团队详情失败")
	t.Logf("  详情: %s - %s", teamDetail.Name, teamDetail.Description)
	assert.Equal(t, teamID, teamDetail.ID, "团队 ID 不匹配")

	// 测试 4: 更新团队
	t.Log("\n测试 4: 更新团队")
	newDisplayName := "更新后的团队名称"
	newDescription := "更新后的团队描述"
	updateReq := org.UpdateTeamDTO{
		DisplayName: &newDisplayName,
		Description: &newDescription,
	}
	updatedTeam, err := apitest.Put[org.TeamDTO](c.HTTPClient(), teamPath(testOrgID, teamID), updateReq)
	require.NoError(t, err, "更新团队失败")
	t.Logf("  更新成功! 显示名: %s", updatedTeam.DisplayName)

	// 验证更新后的字段
	assert.Equal(t, newDisplayName, updatedTeam.DisplayName, "显示名未更新")
	assert.Equal(t, newDescription, updatedTeam.Description, "描述未更新")

	// 测试 5: 删除团队
	t.Log("\n测试 5: 删除团队")
	err = c.Delete(teamPath(testOrgID, teamID))
	require.NoError(t, err, "删除团队失败")
	t.Log("  删除成功!")

	// 标记已删除，避免 t.Cleanup 重复删除
	teamID = 0

	t.Log("\n团队 CRUD 流程测试完成!")
}

// TestTeamListByOrg 测试按组织获取团队。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestTeamListByOrg ./internal/integration/org/
func TestTeamListByOrg(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	// 步骤 1: 创建两个团队
	t.Log("\n步骤 1: 创建测试团队")
	team1Name := "team1_" + uuid.New().String()[:8]
	team1Req := org.CreateTeamDTO{
		Name:        team1Name,
		DisplayName: "团队 1",
	}
	team1, err := apitest.Post[org.TeamDTO](c.HTTPClient(), teamBasePath(testOrgID), team1Req)
	require.NoError(t, err, "创建团队 1 失败")
	team1ID := team1.ID
	t.Cleanup(func() {
		if team1ID > 0 {
			_ = c.Delete(teamPath(testOrgID, team1ID))
		}
	})

	team2Name := "team2_" + uuid.New().String()[:8]
	team2Req := org.CreateTeamDTO{
		Name:        team2Name,
		DisplayName: "团队 2",
	}
	team2, err := apitest.Post[org.TeamDTO](c.HTTPClient(), teamBasePath(testOrgID), team2Req)
	require.NoError(t, err, "创建团队 2 失败")
	team2ID := team2.ID
	t.Cleanup(func() {
		if team2ID > 0 {
			_ = c.Delete(teamPath(testOrgID, team2ID))
		}
	})
	t.Logf("  团队创建成功! ID1: %d, ID2: %d", team1ID, team2ID)

	// 步骤 2: 获取组织下的所有团队
	t.Log("\n步骤 2: 获取组织下的团队列表")
	// 使用较大的 limit 确保能获取到刚创建的团队
	teams, _, err := apitest.GetList[org.TeamDTO](c.HTTPClient(), teamBasePath(testOrgID), map[string]string{})
	require.NoError(t, err, "获取团队列表失败")
	t.Logf("  团队数量: %d", len(teams))

	// 验证两个团队都在列表中
	teamIDs := apitest.ExtractIDs(teams, func(tm org.TeamDTO) uint { return tm.ID })
	assert.Contains(t, teamIDs, team1ID, "列表应包含团队 1")
	assert.Contains(t, teamIDs, team2ID, "列表应包含团队 2")

	// 验证所有团队都属于同一组织
	for _, tm := range teams {
		assert.Equal(t, testOrgID, tm.OrgID, "团队应属于指定组织")
	}

	t.Log("\n按组织获取团队测试完成!")
}

// TestTeamOrgIsolation 测试组织隔离验证。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestTeamOrgIsolation ./internal/integration/org/
func TestTeamOrgIsolation(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	// 步骤 1: 在 testOrgID 创建团队
	t.Log("\n步骤 1: 创建测试团队")
	createReq := org.CreateTeamDTO{
		Name:        "isolated_" + uuid.New().String()[:8],
		DisplayName: "隔离测试团队",
	}
	createdTeam, err := apitest.Post[org.TeamDTO](c.HTTPClient(), teamBasePath(testOrgID), createReq)
	require.NoError(t, err, "创建团队失败")
	teamID := createdTeam.ID
	t.Cleanup(func() {
		if teamID > 0 {
			_ = c.Delete(teamPath(testOrgID, teamID))
		}
	})
	t.Logf("  团队创建成功! ID: %d", teamID)

	// 步骤 2: 尝试用错误的 org_id 访问（应失败）
	t.Log("\n步骤 2: 尝试用错误的 org_id 访问（应失败）")
	invalidOrgID := uint(99999)
	_, err = apitest.Get[org.TeamDTO](c.HTTPClient(), teamPath(invalidOrgID, teamID), nil)
	require.Error(t, err, "访问其他组织的团队应失败")
	t.Logf("  预期失败: %v", err)

	// 步骤 3: 尝试用错误的 org_id 更新（应失败）
	t.Log("\n步骤 3: 尝试用错误的 org_id 更新（应失败）")
	newDesc := "尝试更新"
	updateReq := org.UpdateTeamDTO{
		Description: &newDesc,
	}
	_, err = apitest.Put[org.TeamDTO](c.HTTPClient(), teamPath(invalidOrgID, teamID), updateReq)
	require.Error(t, err, "用错误的 org_id 更新应失败")
	t.Logf("  预期失败: %v", err)

	t.Log("\n组织隔离测试完成!")
}

// TestTeamMemberCRUD 测试团队成员完整流程。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestTeamMemberCRUD ./internal/integration/org/
func TestTeamMemberCRUD(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	// 步骤 1: 创建测试团队
	t.Log("\n步骤 1: 创建测试团队")
	createTeamReq := org.CreateTeamDTO{
		Name:        "teammember_" + uuid.New().String()[:8],
		DisplayName: "团队成员测试",
	}
	createdTeam, err := apitest.Post[org.TeamDTO](c.HTTPClient(), teamBasePath(testOrgID), createTeamReq)
	require.NoError(t, err, "创建团队失败")
	teamID := createdTeam.ID
	t.Cleanup(func() {
		if teamID > 0 {
			_ = c.Delete(teamPath(testOrgID, teamID))
		}
	})
	t.Logf("  团队创建成功! ID: %d", teamID)

	// 步骤 2: 创建测试用户并加入组织
	t.Log("\n步骤 2: 创建测试用户并加入组织")
	testUser := iam.CreateTestUser(t, c, "teammember")
	t.Logf("  用户创建成功! ID: %d", testUser.ID)

	// 先将用户加入组织
	addOrgMemberReq := org.AddMemberDTO{
		UserID: testUser.ID,
		Role:   "member",
	}
	_, err = apitest.Post[org.MemberDTO](c.HTTPClient(), fmt.Sprintf("/api/org/%d/members", testOrgID), addOrgMemberReq)
	require.NoError(t, err, "用户加入组织失败")
	t.Cleanup(func() {
		_ = c.Delete(fmt.Sprintf("/api/org/%d/members/%d", testOrgID, testUser.ID))
	})

	// 步骤 3: 添加团队成员
	t.Log("\n步骤 3: 添加团队成员")
	addTeamMemberReq := org.AddTeamMemberDTO{
		UserID: testUser.ID,
		Role:   "member",
	}
	teamMember, err := apitest.Post[org.TeamMemberDTO](c.HTTPClient(), teamMembersBasePath(testOrgID, teamID), addTeamMemberReq)
	require.NoError(t, err, "添加团队成员失败")
	t.Logf("  成员添加成功! 角色: %s", teamMember.Role)

	// 验证成员信息
	assert.Equal(t, testUser.ID, teamMember.UserID, "用户 ID 不匹配")
	assert.Equal(t, teamID, teamMember.TeamID, "团队 ID 不匹配")
	assert.Equal(t, "member", teamMember.Role, "角色不匹配")
	// TODO: 用户信息未加载，需要 Repository Join User 表
	// assert.NotEmpty(t, teamMember.Username, "用户名应加载")

	// 清理：移除团队成员
	t.Cleanup(func() {
		_ = c.Delete(teamMemberPath(testOrgID, teamID, testUser.ID))
	})

	// 步骤 4: 获取团队成员列表
	t.Log("\n步骤 4: 获取团队成员列表")
	members, meta, err := apitest.GetList[org.TeamMemberDTO](c.HTTPClient(), teamMembersBasePath(testOrgID, teamID), map[string]string{})
	require.NoError(t, err, "获取团队成员列表失败")
	t.Logf("  成员数量: %d", len(members))
	if meta != nil {
		t.Logf("  总数: %d", meta.Total)
	}

	// 验证列表中包含新成员
	userIDs := apitest.ExtractIDs(members, func(m org.TeamMemberDTO) uint { return m.UserID })
	assert.Contains(t, userIDs, testUser.ID, "成员列表应包含新添加的用户")

	// 步骤 5: 移除团队成员
	t.Log("\n步骤 5: 移除团队成员")
	err = c.Delete(teamMemberPath(testOrgID, teamID, testUser.ID))
	require.NoError(t, err, "移除团队成员失败")
	t.Logf("  成员移除成功!")

	// 验证成员已移除
	membersAfter, _, err := apitest.GetList[org.TeamMemberDTO](c.HTTPClient(), teamMembersBasePath(testOrgID, teamID), nil)
	require.NoError(t, err, "获取团队成员列表失败")
	userIDsAfter := apitest.ExtractIDs(membersAfter, func(m org.TeamMemberDTO) uint { return m.UserID })
	assert.NotContains(t, userIDsAfter, testUser.ID, "成员列表不应包含已移除的用户")

	t.Log("\n团队成员流程测试完成!")
}

// TestTeamMemberRequiresOrgMembership 测试成员必须是组织成员。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestTeamMemberRequiresOrgMembership ./internal/integration/org/
func TestTeamMemberRequiresOrgMembership(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	// 步骤 1: 创建测试团队
	t.Log("\n步骤 1: 创建测试团队")
	createTeamReq := org.CreateTeamDTO{
		Name:        "orgcheck_" + uuid.New().String()[:8],
		DisplayName: "组织成员检查测试",
	}
	createdTeam, err := apitest.Post[org.TeamDTO](c.HTTPClient(), teamBasePath(testOrgID), createTeamReq)
	require.NoError(t, err, "创建团队失败")
	teamID := createdTeam.ID
	t.Cleanup(func() {
		if teamID > 0 {
			_ = c.Delete(teamPath(testOrgID, teamID))
		}
	})
	t.Logf("  团队创建成功! ID: %d", teamID)

	// 步骤 2: 创建测试用户（但不加入组织）
	t.Log("\n步骤 2: 创建测试用户")
	testUser := iam.CreateTestUser(t, c, "nonorgmember")
	t.Logf("  用户创建成功! ID: %d", testUser.ID)

	// 步骤 3: 尝试将非组织成员加入团队（应失败）
	t.Log("\n步骤 3: 尝试将非组织成员加入团队（应失败）")
	addTeamMemberReq := org.AddTeamMemberDTO{
		UserID: testUser.ID,
		Role:   "member",
	}
	_, err = apitest.Post[org.TeamMemberDTO](c.HTTPClient(), teamMembersBasePath(testOrgID, teamID), addTeamMemberReq)
	require.Error(t, err, "非组织成员加入团队应该失败")
	t.Logf("  预期失败: %v", err)

	t.Log("\n组织成员要求测试完成!")
}

// TestUserTeams 测试用户视角的团队列表。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestUserTeams ./internal/integration/org/
func TestUserTeams(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	t.Log("\n获取用户加入的团队列表...")
	userTeams, _, err := apitest.GetList[org.UserTeamDTO](c.HTTPClient(), "/api/user/teams", nil)
	require.NoError(t, err, "获取用户团队列表失败")
	t.Logf("  用户加入的团队数量: %d", len(userTeams))

	for _, ut := range userTeams {
		t.Logf("  - [%d] %s (%s) - 组织: %s, 角色: %s",
			ut.ID, ut.DisplayName, ut.Name, ut.OrgName, ut.Role)
	}

	t.Log("\n用户团队列表测试完成!")
}

// TestTeamWithInvalidOrgID 测试使用无效组织 ID 创建团队。
//
// 手动运行:
//
//	API_TEST=1 go test -v -run TestTeamWithInvalidOrgID ./internal/integration/org/
func TestTeamWithInvalidOrgID(t *testing.T) {
	c := iam.LoginAsAdmin(t)

	t.Log("\n测试使用不存在的组织 ID 创建团队...")
	invalidOrgID := uint(99999)

	createReq := org.CreateTeamDTO{
		Name:        "invalidorg_" + uuid.New().String()[:8],
		DisplayName: "无效组织测试",
	}
	_, err := apitest.Post[org.TeamDTO](c.HTTPClient(), teamBasePath(invalidOrgID), createReq)
	require.Error(t, err, "使用不存在的组织 ID 创建团队应该失败")
	t.Logf("  预期失败: %v", err)

	t.Log("\n无效组织 ID 测试完成!")
}
