package stats

import "time"

// SystemStats 系统统计信息值对象
type SystemStats struct {
	TotalUsers       int64
	ActiveUsers      int64
	InactiveUsers    int64
	BannedUsers      int64
	TotalRoles       int64
	TotalPermissions int64
	RecentAuditLogs  []AuditLogSummary
}

// AuditLogSummary 审计日志摘要
type AuditLogSummary struct {
	ID        uint
	UserID    uint
	Username  string
	Action    string
	Resource  string
	Status    string
	CreatedAt time.Time
}

// QueryRepository 定义统计查询仓储接口
type QueryRepository interface {
	// GetSystemStats 获取系统统计信息
	GetSystemStats(recentLogsLimit int) (*SystemStats, error)

	// GetUserCountByStatus 按状态统计用户数量
	GetUserCountByStatus(status string) (int64, error)

	// GetTotalUsers 获取用户总数
	GetTotalUsers() (int64, error)

	// GetTotalRoles 获取角色总数
	GetTotalRoles() (int64, error)

	// GetTotalPermissions 获取权限总数
	GetTotalPermissions() (int64, error)

	// GetRecentAuditLogs 获取最近的审计日志
	GetRecentAuditLogs(limit int) ([]AuditLogSummary, error)
}
