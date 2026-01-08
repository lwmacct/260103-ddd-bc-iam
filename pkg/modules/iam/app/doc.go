// Package app 提供 IAM 模块应用层的类型别名和模块聚合。
//
// # Overview
//
// 本包是 IAM 应用层的顶层入口，提供：
//   - 类型别名：为 Handler 层提供便捷的类型访问
//   - 模块聚合：通过 Fx 聚合所有应用层子模块
//
// # 导出的类型别名
//
// 认证相关：
//   - [AuthUseCases]: 认证用例（登录、注册、Token 刷新）
//   - [TwoFAUseCases]: 双因素认证用例
//   - [CaptchaUseCases]: 验证码用例
//
// 用户管理：
//   - [UserUseCases]: 用户管理用例
//
// 角色权限：
//   - [RoleUseCases]: 角色权限用例
//
// 个人访问令牌：
//   - [PATUseCases]: PAT 用例
//
// 组织团队：
//   - [OrgUseCases]: 组织管理用例
//   - [OrgMemberUseCases]: 组织成员管理用例
//   - [TeamUseCases]: 团队管理用例
//   - [TeamMemberUseCases]: 团队成员管理用例
//   - [UserOrgUseCases]: 用户视角组织查询用例
//
// 审计日志：
//   - [AuditUseCases]: 审计日志用例
//
// # UseCaseModule
//
// UseCaseModule 聚合所有应用层子模块，通过 Fx 依赖注入自动注册：
//
//	fx.New(
//	    app.UseCaseModule,  // ← 导入所有应用层 Handler
//	    // ... 其他模块
//	)
//
// # 模块列表
//
// UseCaseModule 包含以下子模块：
//   - audit.Module: 审计用例模块
//   - captcha.Module: 验证码模块
//   - auth.Module: 认证模块
//   - user.Module: 用户管理模块
//   - role.Module: 角色权限模块
//   - pat.Module: PAT 模块
//   - twofa.Module: 双因素认证模块
//   - org.Module: 组织团队模块
//
// # 依赖关系
//
// 本包不包含具体实现，仅做类型导出和模块聚合。
// 具体实现在各子包中（auth、user、role、org 等）。
package app
