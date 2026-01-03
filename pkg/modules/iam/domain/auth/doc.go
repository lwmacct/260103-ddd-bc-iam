// Package auth 定义认证领域服务接口。
//
// 本包是认证系统的领域层核心，定义了：
//   - [Service]: 认证领域服务接口（密码管理、Token 生成与验证）
//   - [PasswordPolicy]: 密码策略值对象
//   - [TokenClaims]: JWT Token 声明结构
//   - 认证相关错误（见 errors.go）
//   - 认证相关常量（见 constants.go）
//
// 认证模式：
// 系统支持两种认证方式，由 [Service] 接口统一管理：
//   - JWT Token: 用于 Web 端登录，支持访问令牌和刷新令牌
//   - PAT (Personal Access Token): 用于 API 调用和自动化脚本
//
// 安全设计：
//   - 密码使用 BCrypt 哈希存储
//   - JWT 采用 HS256 签名算法
//   - Token 仅存储 user_id，权限信息从缓存实时查询（支持权限即时生效）
//
// 依赖倒置：
// 本包仅定义接口，实现位于 infrastructure/auth 包。
package auth
