// Package health 提供系统健康检查的基础设施实现。
//
// 本包实现 [domain/health.Checker] 接口。
//
// # 健康检查器
//
// [SystemChecker] 检查系统各组件的健康状态：
//   - 数据库连接检查（PostgreSQL/SQLite）
//   - Redis 连接检查
//   - 其他外部依赖检查（可扩展）
//
// 使用示例：
//
//	checker := health.NewSystemChecker(db, redisClient)
//	report := checker.Check(ctx)
//	if report.Status == domain.StatusHealthy {
//	    // 系统健康
//	}
//
// # 检查结果
//
// 检查返回 [domain/health.HealthReport]：
//   - Status: 整体健康状态（Healthy/Unhealthy/Degraded）
//   - Checks: 各组件检查结果列表
//   - Timestamp: 检查时间戳
//
// 单个组件检查结果 [domain/health.CheckResult]：
//   - Name: 组件名称（如 "database"、"redis"）
//   - Status: 组件状态
//   - Message: 状态描述或错误信息
//   - Duration: 检查耗时
//
// # 健康状态
//
//   - Healthy: 所有组件正常
//   - Degraded: 部分组件异常但核心功能可用
//   - Unhealthy: 核心组件异常，服务不可用
//
// # HTTP 端点
//
// 健康检查通常通过 /health 端点暴露：
//   - 200 OK: 健康
//   - 503 Service Unavailable: 不健康
//
// Kubernetes 探针：
//   - Liveness: /health/live
//   - Readiness: /health/ready
package health
