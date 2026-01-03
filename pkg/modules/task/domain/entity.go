package task

import "time"

// Status 任务状态。
type Status string

const (
	// StatusPending 待处理。
	StatusPending Status = "pending"
	// StatusInProgress 进行中。
	StatusInProgress Status = "in_progress"
	// StatusCompleted 已完成。
	StatusCompleted Status = "completed"
)

// Task 团队任务实体。
//
// 任务隶属于特定的组织和团队，支持指派给团队成员。
type Task struct {
	ID          uint      `json:"id"`
	OrgID       uint      `json:"org_id"`
	TeamID      uint      `json:"team_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	AssigneeID  *uint     `json:"assignee_id,omitempty"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// IsCompleted 报告任务是否已完成。
func (t *Task) IsCompleted() bool {
	return t.Status == StatusCompleted
}

// CanTransitionTo 报告是否可以转换到目标状态。
func (t *Task) CanTransitionTo(target Status) bool {
	switch t.Status {
	case StatusPending:
		return target == StatusInProgress || target == StatusCompleted
	case StatusInProgress:
		return target == StatusCompleted || target == StatusPending
	case StatusCompleted:
		return target == StatusInProgress // 允许重新打开
	}
	return false
}

// Start 将任务状态变更为进行中。
func (t *Task) Start() {
	t.Status = StatusInProgress
}

// Complete 将任务标记为完成。
func (t *Task) Complete() {
	t.Status = StatusCompleted
}
