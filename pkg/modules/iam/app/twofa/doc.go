// Package twofa 实现双因素认证（Two-Factor Authentication）的应用层用例。
//
// 本包提供基于 TOTP（RFC 6238）的双因素认证功能：
//
// # Command（写操作）
//
//   - [SetupHandler]: 初始化双因素认证，生成密钥和二维码
//   - [VerifyEnableHandler]: 验证 TOTP 码并启用双因素认证
//   - [DisableHandler]: 禁用双因素认证
//
// # Query（读操作）
//
//   - [GetStatusHandler]: 获取用户双因素认证状态
//
// # DTO 与映射
//
// 请求 DTO：
//   - [SetupCommand]: 设置请求（用户 ID）
//   - [VerifyEnableCommand]: 验证并启用请求（TOTP 码）
//   - [DisableCommand]: 禁用请求（用户 ID）
//
// 响应 DTO：
//   - [SetupResultDTO]: 设置结果（密钥、二维码 Base64、恢复码）
//   - [VerifyEnableResultDTO]: 验证结果（恢复码）
//   - [StatusDTO]: 认证状态（是否启用）
//
// # 使用流程
//
// 1. 调用 Setup() → 生成密钥和二维码
// 2. 用户使用认证器 APP 扫描二维码
// 3. 调用 VerifyEnable() → 验证 TOTP 码并返回恢复码
// 4. 后续登录时调用 Verify() → 验证 TOTP 或恢复码
//
// # 安全设计
//
//   - 密钥：80 位随机数
//   - 恢复码：一次性，验证后删除
//   - 二维码：Base64 PNG 输出
//
// # Thread Safety
//
// 所有 Handler 都是无状态的，仅依赖注入的 Repository（通过 Fx 管理）。
// Handler 本身是并发安全的。
//
// # 依赖关系
//
// 本包依赖 Domain 层的 [TwoFA](github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/twofa) 实体和 Repository 接口。
package twofa
