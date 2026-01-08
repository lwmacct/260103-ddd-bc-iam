// Package apitest 提供手动集成测试辅助工具。
//
// 本包提供 HTTP 客户端和工厂函数，用于针对运行中的服务进行 API 集成测试。
// 测试需要手动触发，通过 MANUAL 环境变量控制执行。
//
// 运行方式：
//
//	注意：必须包含 -count=1 参数以禁用测试结果缓存。
//	Go 1.10+ 默认缓存测试结果，apitest 依赖外部服务状态，
//	禁用缓存确保每次测试都真正执行，避免假阴性。
//
//	# 运行所有测试
//	MANUAL=1 go test -v -count=1 ./internal/apitest/... 2>&1 | grep -E "FAIL|PASS"
//
//	# 运行单个 BC 测试
//	MANUAL=1 go test -v -count=1 ./internal/apitest/iam/...
//	MANUAL=1 go test -v -count=1 ./internal/apitest/app/...
//
//	# 串行执行（服务端压力大时）
//	MANUAL=1 go test -v -count=1 -p 1 ./internal/apitest/...
//
//	# 运行单个测试函数
//	MANUAL=1 go test -v -count=1 -run TestLoginScenarios ./internal/apitest/iam/auth/
//
// 核心类型：
//   - [Client]: HTTP 测试客户端，封装 resty 库
//   - [Get], [Post], [Put], [Delete]: 泛型 HTTP 方法
//   - [CreateTestUser], [CreateTestRole]: 资源工厂函数
//   - [LoginAsAdmin], [LoginAs]: 登录辅助函数
//
// 测试目录按 Bounded Context 组织：
//   - iam/: IAM 域测试（auth, user, role, profile, pat, twofa, org, audit）
//   - app/: App 域测试（setting, cache, system）
//   - task/: Task 域测试（task）
package apitest
