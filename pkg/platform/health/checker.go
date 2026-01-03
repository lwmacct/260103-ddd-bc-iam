package health

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	cacheinfra "github.com/lwmacct/260101-go-pkg-ddd/pkg/platform/cache"
	database "github.com/lwmacct/260101-go-pkg-ddd/pkg/platform/db"
	domainHealth "github.com/lwmacct/260101-go-pkg-ddd/pkg/shared/health"
)

// SystemChecker 系统健康检查器实现
type SystemChecker struct {
	db          *gorm.DB
	redisClient *redis.Client
}

// 确保 SystemChecker 实现了 domainHealth.Checker 接口
var _ domainHealth.Checker = (*SystemChecker)(nil)

// NewSystemChecker 创建系统健康检查器
func NewSystemChecker(db *gorm.DB, redisClient *redis.Client) *SystemChecker {
	return &SystemChecker{
		db:          db,
		redisClient: redisClient,
	}
}

// Check 执行系统健康检查
func (c *SystemChecker) Check(ctx context.Context) *domainHealth.HealthReport {
	report := &domainHealth.HealthReport{
		Status: domainHealth.StatusHealthy,
		Checks: make(map[string]domainHealth.CheckResult),
	}

	allHealthy := true

	// 检查数据库连接
	report.Checks["database"] = c.checkDatabase(ctx)
	if report.Checks["database"].Status != domainHealth.StatusHealthy {
		allHealthy = false
	}

	// 检查 Redis 连接
	report.Checks["redis"] = c.checkRedis(ctx)
	if report.Checks["redis"].Status != domainHealth.StatusHealthy {
		allHealthy = false
	}

	if !allHealthy {
		report.Status = domainHealth.StatusDegraded
	}

	return report
}

// checkDatabase 检查数据库健康状态
func (c *SystemChecker) checkDatabase(ctx context.Context) domainHealth.CheckResult {
	if err := database.HealthCheck(ctx, c.db); err != nil {
		return domainHealth.CheckResult{
			Status: domainHealth.StatusUnhealthy,
			Error:  err.Error(),
		}
	}

	// 获取数据库连接池统计
	stats, _ := database.GetStats(c.db)
	return domainHealth.CheckResult{
		Status: domainHealth.StatusHealthy,
		Stats:  stats,
	}
}

// checkRedis 检查 Redis 健康状态
func (c *SystemChecker) checkRedis(ctx context.Context) domainHealth.CheckResult {
	if err := cacheinfra.HealthCheck(ctx, c.redisClient); err != nil {
		return domainHealth.CheckResult{
			Status: domainHealth.StatusUnhealthy,
			Error:  err.Error(),
		}
	}

	return domainHealth.CheckResult{
		Status: domainHealth.StatusHealthy,
	}
}
