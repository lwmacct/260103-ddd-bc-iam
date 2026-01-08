// Package settings 提供 Settings 模块的集成测试工具。
//
// 本包基于 shared/apitest.Client，提供 Settings 特定的测试辅助功能：
//   - Factory: 测试资源工厂函数（用户设置值）
//   - Helper: 登录辅助函数（委托给 iam 包）
//   - Assert: 测试断言辅助函数
//
// # 使用方式
//
//	// 创建测试客户端并登录
//	c := settings.NewClientFromConfig()
//	c = settings.LoginAs(t, "admin", "admin123")
//
//	// 创建测试资源（自动注册清理）
//	setting := settings.CreateTestUserSetting(t, c, "theme", "dark")
//
//	# 运行测试
//	API_TEST=1 go test -v -count=1 ./internal/manualtest/settings/...
package settings
