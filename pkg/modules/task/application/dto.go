package task

import "time"

// ============================================================================
// Request DTOs
// ============================================================================

// CreateTaskDTO 创建任务请求 DTO。
type CreateTaskDTO struct {
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Description string `json:"description" binding:"max=2000"`
	AssigneeID  *uint  `json:"assignee_id" binding:"omitempty,gt=0"`
}

// UpdateTaskDTO 更新任务请求 DTO。
type UpdateTaskDTO struct {
	Title       *string `json:"title" binding:"omitempty,min=1,max=200"`
	Description *string `json:"description" binding:"omitempty,max=2000"`
	Status      *string `json:"status" binding:"omitempty,oneof=pending in_progress completed"`
	AssigneeID  *uint   `json:"assignee_id" binding:"omitempty,gt=0"`
}

// ============================================================================
// Response DTOs
// ============================================================================

// TaskDTO 任务响应 DTO。
type TaskDTO struct {
	ID          uint      `json:"id"`
	OrgID       uint      `json:"org_id"`
	TeamID      uint      `json:"team_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	AssigneeID  *uint     `json:"assignee_id,omitempty"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
