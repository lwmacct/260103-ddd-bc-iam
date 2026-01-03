package user

import "slices"

// ============================================================================
// 用户类型
// ============================================================================

// UserType 用户类型枚举。
type UserType string

const (
	// UserTypeHuman 人类用户。
	// 通过用户名/邮箱 + 密码登录，支持 2FA。
	UserTypeHuman UserType = "human"

	// UserTypeService 服务账户。
	// 无密码，仅通过 PAT (Personal Access Token) 认证。
	// 用于 CI/CD、API 集成、自动化脚本等场景。
	UserTypeService UserType = "service"

	// UserTypeSystem 系统用户。
	// 系统预置用户（如 root、admin），不可删除。
	// 用于系统初始化和管理。
	UserTypeSystem UserType = "system"
)

// ValidUserTypes 有效的用户类型列表。
var ValidUserTypes = []UserType{UserTypeHuman, UserTypeService, UserTypeSystem}

// IsValidUserType 检查用户类型是否有效。
func IsValidUserType(t UserType) bool {
	return slices.Contains(ValidUserTypes, t)
}

// ============================================================================
// 系统预置用户
// ============================================================================

// RootUsername 超级管理员用户名（硬编码）。
//
// root 用户拥有系统最高权限，无需数据库角色配置：
//   - 自动拥有 *:*:* 权限（全域、全操作、全资源）
//   - 绕过所有权限检查
//   - 通常由 Seeder 在系统初始化时创建
//
// 安全注意：
//   - 仅用于系统维护和紧急操作
//   - 生产环境应限制 root 账户的使用
const RootUsername = "root"

// AdminUsername 系统管理员用户名。
const AdminUsername = "admin"

// SystemUsernames 系统预置用户名列表。
// 这些用户标记为 IsSystem=true，不可删除。
var SystemUsernames = []string{RootUsername, AdminUsername}

// IsSystemUsername 检查是否为系统预置用户名。
func IsSystemUsername(username string) bool {
	return slices.Contains(SystemUsernames, username)
}
