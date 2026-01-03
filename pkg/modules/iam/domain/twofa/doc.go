// Package twofa 定义双因素认证（Two-Factor Authentication）领域模型。
//
// 本包实现基于 TOTP（时间同步一次性密码）的双因素认证，定义了：
//   - [TwoFA]: 用户 2FA 配置实体
//   - [RecoveryCodes]: 恢复码值对象（见 value_objects.go）
//   - [CommandRepository]: 写仓储接口
//   - [QueryRepository]: 读仓储接口
//   - 2FA 领域错误（见 errors.go）
//
// 兼容的 Authenticator 应用：
//   - Google Authenticator
//   - Microsoft Authenticator
//   - 其他标准 TOTP 应用
//
// 核心功能：
//   - TOTP 密钥管理：[TwoFA.Secret] 存储 Base32 编码的密钥
//   - 恢复码：[TwoFA.RecoveryCodes] 用于设备丢失时的账户恢复
//   - 状态管理：[TwoFA.Enable] / [TwoFA.Disable]
//
// 安全设计：
//   - Secret 字段不在 JSON 中暴露
//   - 恢复码为一次性使用（[TwoFA.UseRecoveryCode]）
//   - 恢复码格式：xxxx-xxxx（8 位数字）
//
// 依赖倒置：
// 本包仅定义接口，实现位于 infrastructure/persistence 和 infrastructure/twofa 包。
package twofa
