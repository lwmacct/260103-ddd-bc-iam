// Package manualtest 提供 Settings 模块的 HTTP 集成测试辅助工具。
//
// 测试工具包括：
//   - Client: HTTP 测试客户端，支持认证请求
//   - Factory: 测试资源工厂函数（设置值）
//   - Helper: 测试辅助函数（复用 IAM 登录）
//   - Assert: 测试断言辅助函数
//
// 使用方式：
//
//	c := manualtest.NewClient()
//	c.SetToken(iamToken) // 使用 IAM 登录获取 token
//	result, err := manualtest.Get[UserSettingDTO](c, "/api/user/settings", nil)
package manualtest
