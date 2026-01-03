package task_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	taskapplication "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/application"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/platform/manualtest"
)

// 种子数据: acme org (ID=1), engineering team (ID=1), admin 是 owner 和 team lead
const (
	testOrgID  uint = 1
	testTeamID uint = 1
)

// taskBasePath 返回任务 API 基础路径
func taskBasePath() string {
	return fmt.Sprintf("/api/org/%d/teams/%d/tasks", testOrgID, testTeamID)
}

// taskPath 返回单个任务 API 路径
func taskPath(id uint) string {
	return fmt.Sprintf("%s/%d", taskBasePath(), id)
}

// TestTaskCRUD 任务 CRUD 完整流程测试。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestTaskCRUD ./internal/manualtest/task/
func TestTaskCRUD(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 测试 1: 创建任务
	t.Log("\n测试 1: 创建任务")
	createReq := taskapplication.CreateTaskDTO{
		Title:       "测试任务",
		Description: "这是一个测试任务",
	}
	createdTask, err := manualtest.Post[taskapplication.TaskDTO](c, taskBasePath(), createReq)
	require.NoError(t, err, "创建任务失败")
	t.Logf("  创建成功! 任务 ID: %d, 标题: %s", createdTask.ID, createdTask.Title)

	taskID := createdTask.ID
	t.Cleanup(func() {
		if taskID > 0 {
			_ = c.Delete(taskPath(taskID))
		}
	})

	// 验证创建结果
	assert.Equal(t, "测试任务", createdTask.Title)
	assert.Equal(t, "这是一个测试任务", createdTask.Description)
	assert.Equal(t, "pending", createdTask.Status)
	assert.Equal(t, testOrgID, createdTask.OrgID)
	assert.Equal(t, testTeamID, createdTask.TeamID)

	// 测试 2: 获取任务详情
	t.Log("\n测试 2: 获取任务详情")
	taskDetail, err := manualtest.Get[taskapplication.TaskDTO](c, taskPath(taskID), nil)
	require.NoError(t, err, "获取任务详情失败")
	t.Logf("  标题: %s, 状态: %s", taskDetail.Title, taskDetail.Status)

	assert.Equal(t, taskID, taskDetail.ID)
	assert.Equal(t, "测试任务", taskDetail.Title)

	// 测试 3: 更新任务
	t.Log("\n测试 3: 更新任务标题")
	newTitle := "更新后的任务标题"
	updateReq := taskapplication.UpdateTaskDTO{
		Title: &newTitle,
	}
	updatedTask, err := manualtest.Put[taskapplication.TaskDTO](c, taskPath(taskID), updateReq)
	require.NoError(t, err, "更新任务失败")
	t.Logf("  更新成功! 新标题: %s", updatedTask.Title)

	assert.Equal(t, newTitle, updatedTask.Title)

	// 测试 4: 获取任务列表
	t.Log("\n测试 4: 获取任务列表")
	tasks, meta, err := manualtest.GetList[taskapplication.TaskDTO](c, taskBasePath(), map[string]string{
		"page":  "1",
		"limit": "10",
	})
	require.NoError(t, err, "获取任务列表失败")
	t.Logf("  任务数量: %d", len(tasks))
	if meta != nil {
		t.Logf("  总数: %d", meta.Total)
	}

	// 验证列表中包含创建的任务
	taskIDs := make([]uint, len(tasks))
	for i, tsk := range tasks {
		taskIDs[i] = tsk.ID
	}
	assert.Contains(t, taskIDs, taskID, "列表中应包含创建的任务")

	// 测试 5: 删除任务
	t.Log("\n测试 5: 删除任务")
	err = c.Delete(taskPath(taskID))
	require.NoError(t, err, "删除任务失败")
	t.Log("  删除成功!")

	// 标记已删除，避免 Cleanup 重复删除
	taskID = 0

	t.Log("\n任务 CRUD 测试完成!")
}

// TestTaskStatusTransition 任务状态流转测试。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestTaskStatusTransition ./internal/manualtest/task/
func TestTaskStatusTransition(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 创建测试任务
	t.Log("\n步骤 1: 创建测试任务")
	createReq := taskapplication.CreateTaskDTO{
		Title:       "状态测试任务",
		Description: "用于测试状态流转",
	}
	createdTask, err := manualtest.Post[taskapplication.TaskDTO](c, taskBasePath(), createReq)
	require.NoError(t, err, "创建任务失败")
	t.Logf("  任务 ID: %d, 初始状态: %s", createdTask.ID, createdTask.Status)

	taskID := createdTask.ID
	t.Cleanup(func() {
		if taskID > 0 {
			_ = c.Delete(taskPath(taskID))
		}
	})

	assert.Equal(t, "pending", createdTask.Status, "初始状态应为 pending")

	// 步骤 2: pending → in_progress
	t.Log("\n步骤 2: 状态从 pending → in_progress")
	inProgressStatus := "in_progress"
	updateReq := taskapplication.UpdateTaskDTO{
		Status: &inProgressStatus,
	}
	updatedTask, err := manualtest.Put[taskapplication.TaskDTO](c, taskPath(taskID), updateReq)
	require.NoError(t, err, "更新状态失败")
	t.Logf("  新状态: %s", updatedTask.Status)

	assert.Equal(t, "in_progress", updatedTask.Status)

	// 步骤 3: in_progress → completed
	t.Log("\n步骤 3: 状态从 in_progress → completed")
	completedStatus := "completed"
	updateReq = taskapplication.UpdateTaskDTO{
		Status: &completedStatus,
	}
	updatedTask, err = manualtest.Put[taskapplication.TaskDTO](c, taskPath(taskID), updateReq)
	require.NoError(t, err, "更新状态失败")
	t.Logf("  新状态: %s", updatedTask.Status)

	assert.Equal(t, "completed", updatedTask.Status)

	t.Log("\n任务状态流转测试完成!")
}

// TestTaskInvalidStatusTransition 测试无效状态转换。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestTaskInvalidStatusTransition ./internal/manualtest/task/
func TestTaskInvalidStatusTransition(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 创建测试任务
	t.Log("\n步骤 1: 创建测试任务")
	createReq := taskapplication.CreateTaskDTO{
		Title: "无效状态转换测试",
	}
	createdTask, err := manualtest.Post[taskapplication.TaskDTO](c, taskBasePath(), createReq)
	require.NoError(t, err, "创建任务失败")
	t.Logf("  任务 ID: %d, 状态: %s", createdTask.ID, createdTask.Status)

	taskID := createdTask.ID
	t.Cleanup(func() {
		if taskID > 0 {
			_ = c.Delete(taskPath(taskID))
		}
	})

	// 先完成任务
	t.Log("\n步骤 2: 完成任务（pending → completed 是允许的）")
	completedStatus := "completed"
	updateReq := taskapplication.UpdateTaskDTO{
		Status: &completedStatus,
	}
	updatedTask, err := manualtest.Put[taskapplication.TaskDTO](c, taskPath(taskID), updateReq)
	require.NoError(t, err, "完成任务失败")
	t.Logf("  新状态: %s", updatedTask.Status)
	assert.Equal(t, "completed", updatedTask.Status)

	// 步骤 3: 尝试 completed → pending（应失败）
	t.Log("\n步骤 3: 尝试 completed → pending（应失败）")
	pendingStatus := "pending"
	updateReq = taskapplication.UpdateTaskDTO{
		Status: &pendingStatus,
	}
	_, err = manualtest.Put[taskapplication.TaskDTO](c, taskPath(taskID), updateReq)
	require.Error(t, err, "从 completed 到 pending 应该失败")
	t.Logf("  预期失败: %v", err)

	t.Log("\n无效状态转换测试完成!")
}

// TestTaskWithAssignee 测试任务指派功能。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestTaskWithAssignee ./internal/manualtest/task/
func TestTaskWithAssignee(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 步骤 1: 创建带指派人的任务
	t.Log("\n步骤 1: 创建带指派人的任务")
	adminUserID := uint(2) // admin 用户 ID
	createReq := taskapplication.CreateTaskDTO{
		Title:       "指派任务",
		Description: "指派给 admin 用户",
		AssigneeID:  &adminUserID,
	}
	createdTask, err := manualtest.Post[taskapplication.TaskDTO](c, taskBasePath(), createReq)
	require.NoError(t, err, "创建任务失败")
	t.Logf("  任务 ID: %d, 指派给: %v", createdTask.ID, createdTask.AssigneeID)

	taskID := createdTask.ID
	t.Cleanup(func() {
		if taskID > 0 {
			_ = c.Delete(taskPath(taskID))
		}
	})

	require.NotNil(t, createdTask.AssigneeID, "AssigneeID 不应为空")
	assert.Equal(t, adminUserID, *createdTask.AssigneeID)

	// 步骤 2: 更新指派人为空（取消指派）
	t.Log("\n步骤 2: 取消指派")
	var nilAssignee *uint = nil
	updateReq := taskapplication.UpdateTaskDTO{
		AssigneeID: nilAssignee,
	}
	updatedTask, err := manualtest.Put[taskapplication.TaskDTO](c, taskPath(taskID), updateReq)
	require.NoError(t, err, "取消指派失败")
	t.Logf("  更新后指派人: %v", updatedTask.AssigneeID)

	// 注意: 如果 DTO 不传 AssigneeID 字段，后端可能不会更新
	// 这里验证 API 行为

	t.Log("\n任务指派测试完成!")
}

// TestTaskOrgTeamIsolation 测试组织/团队数据隔离。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestTaskOrgTeamIsolation ./internal/manualtest/task/
func TestTaskOrgTeamIsolation(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 创建一个任务
	t.Log("\n步骤 1: 创建任务")
	createReq := taskapplication.CreateTaskDTO{
		Title: "隔离测试任务",
	}
	createdTask, err := manualtest.Post[taskapplication.TaskDTO](c, taskBasePath(), createReq)
	require.NoError(t, err, "创建任务失败")
	t.Logf("  任务 ID: %d, OrgID: %d, TeamID: %d", createdTask.ID, createdTask.OrgID, createdTask.TeamID)

	taskID := createdTask.ID
	t.Cleanup(func() {
		if taskID > 0 {
			_ = c.Delete(taskPath(taskID))
		}
	})

	// 验证任务绑定到正确的 org 和 team
	assert.Equal(t, testOrgID, createdTask.OrgID, "OrgID 应匹配")
	assert.Equal(t, testTeamID, createdTask.TeamID, "TeamID 应匹配")

	// 尝试用不存在的 org/team 访问（应该失败，因为用户不是成员）
	t.Log("\n步骤 2: 尝试访问不存在的组织（应失败）")
	invalidPath := fmt.Sprintf("/api/org/9999/teams/%d/tasks", testTeamID)
	_, _, err = manualtest.GetList[taskapplication.TaskDTO](c, invalidPath, nil)
	require.Error(t, err, "访问不存在的组织应失败")
	t.Logf("  预期失败: %v", err)

	t.Log("\n组织/团队数据隔离测试完成!")
}

// TestListTasks 测试任务列表分页。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestListTasks ./internal/manualtest/task/
func TestListTasks(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	t.Log("获取任务列表...")
	tasks, meta, err := manualtest.GetList[taskapplication.TaskDTO](c, taskBasePath(), map[string]string{
		"page":  "1",
		"limit": "10",
	})
	require.NoError(t, err, "获取任务列表失败")

	t.Logf("任务数量: %d", len(tasks))
	if meta != nil {
		t.Logf("总数: %d, 总页数: %d", meta.Total, meta.TotalPages)
	}

	for _, tsk := range tasks {
		t.Logf("  - [%d] %s (状态: %s, 指派: %v)", tsk.ID, tsk.Title, tsk.Status, tsk.AssigneeID)
	}
}
