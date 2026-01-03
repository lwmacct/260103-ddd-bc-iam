// Package health 定义系统健康检查的领域模型和接口。
//
// 本包提供系统健康状态监控能力，定义了：
//   - [Status]: 健康状态枚举（healthy、unhealthy、degraded）
//   - [CheckResult]: 单个组件的检查结果
//   - [HealthReport]: 系统整体健康报告
//   - [Checker]: 健康检查服务接口
//
// 健康状态：
//   - [StatusHealthy]: 组件运行正常
//   - [StatusUnhealthy]: 组件不可用
//   - [StatusDegraded]: 组件部分可用或性能下降
//
// 使用场景：
// 通常用于 /health 或 /readiness 端点，
// 支持 Kubernetes 存活探针 (liveness) 和就绪探针 (readiness)。
//
// 依赖倒置：
// 本包仅定义接口，实现位于 infrastructure/health 包。
package health
