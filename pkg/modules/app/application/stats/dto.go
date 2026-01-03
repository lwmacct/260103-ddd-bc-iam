package stats

import "time"

// StatsDTO 系统统计响应 DTO
type StatsDTO struct {
	TotalUsers       int64                `json:"total_users"`
	ActiveUsers      int64                `json:"active_users"`
	InactiveUsers    int64                `json:"inactive_users"`
	BannedUsers      int64                `json:"banned_users"`
	TotalRoles       int64                `json:"total_roles"`
	TotalPermissions int64                `json:"total_permissions"`
	RecentAuditLogs  []AuditLogSummaryDTO `json:"recent_audit_logs,omitempty"`
}

// AuditLogSummaryDTO 审计日志摘要 DTO
type AuditLogSummaryDTO struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
