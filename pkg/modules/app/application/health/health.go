// Package health 重新导出健康检查领域类型供 Adapters 层使用
// 遵循 DDD 依赖方向: Adapters → Application → Domain
package health

import "github.com/lwmacct/260101-go-pkg-ddd/pkg/shared/health"

// 重新导出领域类型
type (
	Status       = health.Status
	CheckResult  = health.CheckResult
	HealthReport = health.HealthReport
	Checker      = health.Checker
)

// 重新导出状态常量
const (
	StatusHealthy   = health.StatusHealthy
	StatusUnhealthy = health.StatusUnhealthy
	StatusDegraded  = health.StatusDegraded
)
