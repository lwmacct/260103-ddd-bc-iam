package container

import (
	persistence "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infra/persistence"
)

// ContainerOptions 容器初始化选项。
type ContainerOptions struct {
	AutoMigrate bool // 是否自动迁移数据库（仅建议在开发环境使用）
}

// DefaultOptions 返回默认的容器选项。
func DefaultOptions() *ContainerOptions {
	return &ContainerOptions{
		AutoMigrate: false, // 生产环境默认：不自动迁移
	}
}

// GetAllModels 返回所有需要迁移的领域模型。
// 新增领域模型时，需在此处注册。
func GetAllModels() []any {
	return []any{
		// 用户和角色
		&persistence.UserModel{},
		&persistence.RoleModel{},
		// 认证和授权
		&persistence.PersonalAccessTokenModel{},
		&persistence.TwoFAModel{},
		// 组织和团队
		&persistence.OrgModel{},
		&persistence.TeamModel{},
		&persistence.OrgMemberModel{},
		&persistence.TeamMemberModel{},
		// 审计日志
		&persistence.AuditModel{},
	}
}
