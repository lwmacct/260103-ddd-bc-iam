package health

import "context"

// Status 表示组件的健康状态
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusDegraded  Status = "degraded"
)

// CheckResult 单个组件的检查结果
type CheckResult struct {
	Status Status
	Error  string
	Stats  map[string]any
}

// HealthReport 系统健康报告
type HealthReport struct {
	Status Status
	Checks map[string]CheckResult
}

// Checker 定义健康检查服务接口
type Checker interface {
	// Check 执行健康检查并返回报告
	Check(ctx context.Context) *HealthReport
}
