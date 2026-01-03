// Package twofa 提供双因素认证的基础设施实现。
//
// 本包基于 TOTP (RFC 6238) 标准实现 2FA 功能，兼容主流认证器应用：
//   - Google Authenticator
//   - Microsoft Authenticator
//   - Authy, 1Password 等
//
// # 核心组件
//
//   - [Service]: 2FA 服务，提供完整生命周期管理
//
// # 功能概览
//
//   - Setup: 生成 TOTP 密钥和二维码（Base64 PNG）
//   - VerifyAndEnable: 验证首次 TOTP 码并启用 2FA，同时生成恢复码
//   - Verify: 验证 TOTP 码或恢复码
//   - Disable: 禁用用户的 2FA
//   - GetStatus: 查询 2FA 启用状态和剩余恢复码数量
//
// # 安全设计
//
//   - TOTP 密钥使用 80 位（10 字节）随机数
//   - 恢复码为一次性使用，使用后自动删除
//   - 密钥存储在数据库中，不对外暴露
//
// # 使用示例
//
//	service := twofa.NewService(cmdRepo, qryRepo, userRepo, "MyApp")
//
//	// 1. 设置 2FA（返回二维码）
//	result, err := service.Setup(ctx, userID)
//	// result.QRCodeImg 可直接用于 <img src="...">
//
//	// 2. 用户扫码后，验证并启用
//	recoveryCodes, err := service.VerifyAndEnable(ctx, userID, totpCode)
//
//	// 3. 后续登录时验证
//	valid, err := service.Verify(ctx, userID, totpCode)
//
// # 依赖
//
//   - github.com/pquerna/otp/totp: TOTP 实现
//   - github.com/skip2/go-qrcode: 二维码生成
package twofa
