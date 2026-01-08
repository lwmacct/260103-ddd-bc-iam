// Package manualtest 提供 IAM 模块的 HTTP 集成测试辅助工具。
//
// # Overview
//
// 本包提供 HTTP 集成测试的工具函数，用于手动测试 IAM 模块的 API 端点：
//   - [Client]: HTTP 测试客户端，支持认证请求
//   - Factory: 测试资源工厂函数（用户、角色、配置等）
//   - Helper: 测试辅助函数（登录、缓存）
//   - Assert: 测试断言辅助函数
//
// # Usage
//
//	// 创建测试客户端并登录
//	c := manualtest.NewClient()
//	admin, err := c.Login("admin", "admin123")
//	assert.NoError(t, err)
//
//	// 创建测试资源（自动注册清理）
//	user := manualtest.CreateTestUser(t, c, "testuser")
//
//	// 使用断言辅助函数
//	manualtest.AssertUserHasRole(t, user, "editor")
//
// # Thread Safety
//
// 测试工具函数设计用于测试环境，不保证生产环境使用的并发安全性。
// Client 是并发安全的，但每个测试应使用独立的 Client 实例以避免状态干扰。
//
// # 运行测试
//
//	# 运行所有 manualtest
//	MANUAL=1 go test -v -count=1 ./internal/manualtest/...
//
//	# 运行特定模块测试
//	MANUAL=1 go test -v -count=1 -run TestLoginScenarios ./internal/manualtest/iam/auth/
package manualtest
