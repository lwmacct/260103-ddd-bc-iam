// Package iam 提供 IAM 模块的集成测试工具。
//
// 本包基于 shared/apitest.Client，提供 IAM 特定的测试辅助功能：
//   - [Client]: BC 特定客户端（嵌入 apitest.Client），支持登录、验证码
//   - Factory: 测试资源工厂函数（用户、角色等）
//   - Helper: 登录辅助函数（带 session 缓存）
//   - Assert: 测试断言辅助函数
//
// # 使用方式
//
//	// 创建测试客户端并登录
//	c := iam.NewClientFromConfig()
//	c.Login("admin", "admin123")
//
//	// 使用登录辅助（带 session 缓存）
//	c := iam.LoginAs(t, "admin", "admin123")
//
//	// 创建测试资源（自动注册清理）
//	user := iam.CreateTestUser(t, c, "testuser")
//
//	# 运行测试
//	API_TEST=1 go test -v -count=1 ./internal/manualtest/iam/...
package iam
