package container

import (
	iamPersistence "github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/infra/persistence"
	userSettingsPersistence "github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/infra/persistence"
	settingsPersistence "github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/infra/persistence"
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
		// ========== IAM Models ==========
		// 用户和角色
		&iamPersistence.UserModel{},
		&iamPersistence.RoleModel{},
		// 认证和授权
		&iamPersistence.PersonalAccessTokenModel{},
		&iamPersistence.TwoFAModel{},
		// 组织和团队
		&iamPersistence.OrgModel{},
		&iamPersistence.TeamModel{},
		&iamPersistence.OrgMemberModel{},
		&iamPersistence.TeamMemberModel{},
		// 审计日志
		&iamPersistence.AuditModel{},

		// ========== Settings BC Models ==========
		&userSettingsPersistence.UserSettingModel{},
		&userSettingsPersistence.OrgSettingModel{},
		&userSettingsPersistence.TeamSettingModel{},

		// ========== Settings Models ==========
		&settingsPersistence.SettingModel{},
		&settingsPersistence.SettingCategoryModel{},
	}
}
