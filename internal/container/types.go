package container

import (
	corepersistence "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/infrastructure/persistence"
	crmpersistence "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/infrastructure/persistence"
	iampersistence "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/infrastructure/persistence"
	taskpersistence "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/infrastructure/persistence"
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
		&iampersistence.UserModel{},
		&iampersistence.RoleModel{},
		&iampersistence.PersonalAccessTokenModel{},
		&iampersistence.AuditModel{},
		&iampersistence.TwoFAModel{},
		&corepersistence.SettingModel{},
		&corepersistence.SettingCategoryModel{},
		&iampersistence.UserSettingModel{},
		// 组织和团队
		&iampersistence.OrgModel{},
		&iampersistence.TeamModel{},
		&iampersistence.OrgMemberModel{},
		&iampersistence.TeamMemberModel{},
		// 任务
		&taskpersistence.TaskModel{},
		// CRM
		&crmpersistence.ContactModel{},
		&crmpersistence.CompanyModel{},
		&crmpersistence.LeadModel{},
		&crmpersistence.OpportunityModel{},
	}
}
