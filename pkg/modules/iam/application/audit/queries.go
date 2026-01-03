package audit

import "time"

// GetQuery 获取审计日志查询
type GetQuery struct {
	LogID uint
}

// ListQuery 获取审计日志列表查询
type ListQuery struct {
	Page      int
	Limit     int
	UserID    *uint
	Action    string
	Resource  string
	Status    string
	StartDate *time.Time
	EndDate   *time.Time
}
