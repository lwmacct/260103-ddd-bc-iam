// Package captcha 实现验证码的应用层用例。
//
// 本包仅提供 Command Handler：
//
// # Command（写操作）
//
//   - [command.GenerateCaptchaHandler]: 生成图形验证码
//
// # 验证码流程
//
// 生成流程：
//  1. 调用 GenerateCaptchaHandler 生成验证码
//  2. 返回验证码 ID 和 Base64 编码的图片
//  3. 验证码存储在内存中（有过期时间）
//
// 验证流程：
// 验证码验证在 [domain/captcha.Service] 层处理，
// 通常在登录等敏感操作前调用。
//
// 安全特性：
//   - 验证码有过期时间
//   - 验证后立即失效（一次性使用）
//   - 支持大小写不敏感匹配
//
// 依赖：
//   - [domain/captcha.Service]: 验证码领域服务接口
//
// 依赖注入：所有 Handler 通过 [bootstrap.Container] 注册。
package captcha
