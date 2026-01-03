package audit

import "time"

// AuditDTO 审计日志响应 DTO
type AuditDTO struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Details   string    `json:"details"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ListDTO 审计日志列表响应 DTO
type ListDTO struct {
	Logs  []*AuditDTO `json:"logs"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}

// AuditActionsResponseDTO 审计操作定义响应 DTO
// 供前端筛选器动态获取操作定义、分类和操作类型选项
type AuditActionsResponseDTO struct {
	Actions    []AuditActionDTO    `json:"actions"`    // 操作定义列表
	Categories []CategoryOptionDTO `json:"categories"` // 分类选项
	Operations []OperationTypeDTO  `json:"operations"` // 操作类型选项
}

// AuditActionDTO 审计操作定义 DTO
type AuditActionDTO struct {
	Action      string `json:"action"`      // "user.create"
	Operation   string `json:"operation"`   // "create"
	Category    string `json:"category"`    // "user"
	Label       string `json:"label"`       // "创建用户"
	Description string `json:"description"` // "Create new user"
}

// CategoryOptionDTO 分类选项 DTO
type CategoryOptionDTO struct {
	Value string `json:"value"` // "user"
	Label string `json:"label"` // "用户管理"
}

// OperationTypeDTO 操作类型选项 DTO
type OperationTypeDTO struct {
	Value string `json:"value"` // "create"
	Label string `json:"label"` // "创建"
}
