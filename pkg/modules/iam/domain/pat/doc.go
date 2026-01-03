// Package pat 定义个人访问令牌（Personal Access Token）领域模型。
//
// PAT 是一种长期有效的 API 认证凭证，本包定义了：
//   - [PersonalAccessToken]: PAT 实体
//   - [TokenListItem]: PAT 列表项（不含敏感信息）
//   - [PermissionList]: 权限列表值对象（见 value_objects.go）
//   - [StringList]: 字符串列表值对象（IP 白名单等）
//   - [CommandRepository]: 写仓储接口
//   - [QueryRepository]: 读仓储接口
//   - PAT 领域错误（见 errors.go）
//
// 适用场景：
//   - CI/CD 自动化脚本
//   - 第三方应用集成
//   - 命令行工具认证
//
// 安全特性：
//   - Token 以哈希形式存储，原始值仅在创建时返回一次
//   - 支持细粒度权限控制（[PersonalAccessToken.Permissions]）
//   - 可配置过期时间（[PersonalAccessToken.ExpiresAt]）
//   - 支持 IP 白名单（[PersonalAccessToken.IPWhitelist]）
//   - 提供 TokenPrefix 用于识别（不暴露完整 Token）
//
// Token 状态：
//   - active: 活跃可用
//   - disabled: 已禁用
//   - expired: 已过期
//
// 依赖倒置：
// 本包仅定义接口，实现位于 infrastructure/persistence 包。
package pat
