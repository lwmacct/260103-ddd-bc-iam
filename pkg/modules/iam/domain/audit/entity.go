package audit

import "time"

// Audit 审计日志实体，用于追踪用户操作
type Audit struct {
	ID          uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	UserID      uint
	Username    string
	Action      string
	Resource    string
	ResourceID  string
	IPAddress   string
	UserAgent   string
	Details     string
	Status      string
	RequestID   string // 请求追踪 ID（复用 OTel Trace ID）
	OperationID string // API 操作标识符
}

// FilterOptions 审计日志过滤条件
type FilterOptions struct {
	UserID    *uint
	Action    string
	Resource  string
	Status    string
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	Limit     int
}

// 操作状态常量
const (
	StatusSuccess = "success"
	StatusFailed  = "failed"
	StatusPending = "pending"
)

// IsSuccess 检查操作是否成功
func (a *Audit) IsSuccess() bool {
	return a.Status == StatusSuccess
}

// IsFailed 检查操作是否失败
func (a *Audit) IsFailed() bool {
	return a.Status == StatusFailed
}

// IsPending 检查操作是否待定
func (a *Audit) IsPending() bool {
	return a.Status == StatusPending
}

// GetResourceIdentifier 获取完整资源标识符
func (a *Audit) GetResourceIdentifier() string {
	if a.ResourceID == "" {
		return a.Resource
	}
	return a.Resource + ":" + a.ResourceID
}

// IsRecentlyCreated 检查是否在指定时间范围内创建
func (a *Audit) IsRecentlyCreated(duration time.Duration) bool {
	return time.Since(a.CreatedAt) <= duration
}

// HasDetails 检查是否有详情信息
func (a *Audit) HasDetails() bool {
	return a.Details != ""
}

// IsUserAction 检查是否为用户操作（有用户 ID）
func (a *Audit) IsUserAction() bool {
	return a.UserID > 0
}

// IsSystemAction 检查是否为系统操作（无用户 ID）
func (a *Audit) IsSystemAction() bool {
	return a.UserID == 0
}

// MatchesFilter 检查日志是否匹配过滤条件
func (a *Audit) MatchesFilter(filter FilterOptions) bool {
	if filter.UserID != nil && a.UserID != *filter.UserID {
		return false
	}
	if filter.Action != "" && a.Action != filter.Action {
		return false
	}
	if filter.Resource != "" && a.Resource != filter.Resource {
		return false
	}
	if filter.Status != "" && a.Status != filter.Status {
		return false
	}
	if filter.StartDate != nil && a.CreatedAt.Before(*filter.StartDate) {
		return false
	}
	if filter.EndDate != nil && a.CreatedAt.After(*filter.EndDate) {
		return false
	}
	return true
}

// IsValidFilter 检查过滤条件是否有效
func (f *FilterOptions) IsValidFilter() bool {
	if f.Page < 0 || f.Limit < 0 {
		return false
	}
	if f.StartDate != nil && f.EndDate != nil && f.StartDate.After(*f.EndDate) {
		return false
	}
	return true
}

// SetDefaults 设置默认分页值
func (f *FilterOptions) SetDefaults() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 20
	}
	if f.Limit > 100 {
		f.Limit = 100
	}
}
