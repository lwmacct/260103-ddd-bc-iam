package application

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application/setting"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application/stats"
	domain_stats "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/stats"
	corepersistence "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/infrastructure/persistence"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/application/audit"
	app_captcha "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/application/captcha"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/application/org"
	iampersistence "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/infrastructure/persistence"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/platform/validation"
	infra_captcha "github.com/lwmacct/260101-go-pkg-ddd/pkg/shared/captcha"
)

// --- 用例模块结构体 ---

// AuditUseCases 审计日志用例处理器
type AuditUseCases struct {
	CreateLog *audit.CreateHandler
	Get       *audit.GetHandler
	List      *audit.ListHandler
}

// SettingUseCases 设置管理用例处理器
type SettingUseCases struct {
	Create         *setting.CreateHandler
	Update         *setting.UpdateHandler
	Delete         *setting.DeleteHandler
	BatchUpdate    *setting.BatchUpdateHandler
	Get            *setting.GetHandler
	List           *setting.ListHandler
	ListSettings   *setting.ListSettingsHandler
	CreateCategory *setting.CreateCategoryHandler
	UpdateCategory *setting.UpdateCategoryHandler
	DeleteCategory *setting.DeleteCategoryHandler
	GetCategory    *setting.GetCategoryHandler
	ListCategories *setting.ListCategoriesHandler
}

// UserSettingUseCases 用户设置用例处理器
type UserSettingUseCases struct {
	Set            *setting.UserSetHandler
	BatchSet       *setting.UserBatchSetHandler
	Reset          *setting.UserResetHandler
	ResetAll       *setting.UserResetAllHandler
	Get            *setting.UserGetHandler
	List           *setting.UserListHandler
	ListSettings   *setting.UserListSettingsHandler
	ListCategories *setting.UserListCategoriesHandler
}

// StatsUseCases 统计查询用例处理器
type StatsUseCases struct {
	GetStats *stats.GetStatsHandler
}

// CaptchaUseCases 验证码用例处理器
type CaptchaUseCases struct {
	Generate *app_captcha.GenerateHandler
}

// OrganizationUseCases 组织相关用例处理器
type OrganizationUseCases struct {
	// Organization
	Create *org.CreateHandler
	Update *org.UpdateHandler
	Delete *org.DeleteHandler
	Get    *org.GetHandler
	List   *org.ListHandler

	// Member
	MemberAdd        *org.MemberAddHandler
	MemberRemove     *org.MemberRemoveHandler
	MemberUpdateRole *org.MemberUpdateRoleHandler
	MemberList       *org.MemberListHandler

	// Team
	TeamCreate *org.TeamCreateHandler
	TeamUpdate *org.TeamUpdateHandler
	TeamDelete *org.TeamDeleteHandler
	TeamGet    *org.TeamGetHandler
	TeamList   *org.TeamListHandler

	// Team Member
	TeamMemberAdd    *org.TeamMemberAddHandler
	TeamMemberRemove *org.TeamMemberRemoveHandler
	TeamMemberList   *org.TeamMemberListHandler

	// User View
	UserOrgs  *org.UserOrgsHandler
	UserTeams *org.UserTeamsHandler
}

// --- Fx 模块 ---

// UseCaseModule 提供按领域组织的 App 模块用例处理器。
var UseCaseModule = fx.Module("app.usecase",
	fx.Provide(
		newAuditUseCases,
		newSettingUseCases,
		newUserSettingUseCases,
		newStatsUseCases,
		newCaptchaUseCases,
		newOrganizationUseCases,
	),
)

// --- 构造函数 ---

func newAuditUseCases(repos iampersistence.AuditRepositories) *AuditUseCases {
	return &AuditUseCases{
		CreateLog: audit.NewCreateHandler(repos.Command),
		Get:       audit.NewGetHandler(repos.Query),
		List:      audit.NewListHandler(repos.Query),
	}
}

// settingUseCasesParams 聚合 Setting 用例所需的依赖。
type settingUseCasesParams struct {
	fx.In

	SettingRepos  corepersistence.SettingRepositories
	SettingsCache setting.SettingsCacheService
}

func newSettingUseCases(p settingUseCasesParams) *SettingUseCases {
	validator := validation.NewJSONLogicValidator()

	return &SettingUseCases{
		Create:         setting.NewCreateHandler(p.SettingRepos.Command, p.SettingRepos.Query, p.SettingsCache),
		Update:         setting.NewUpdateHandler(p.SettingRepos.Command, p.SettingRepos.Query, validator, p.SettingsCache),
		Delete:         setting.NewDeleteHandler(p.SettingRepos.Command, p.SettingRepos.Query, p.SettingsCache),
		BatchUpdate:    setting.NewBatchUpdateHandler(p.SettingRepos.Command, p.SettingRepos.Query, validator, p.SettingsCache),
		Get:            setting.NewGetHandler(p.SettingRepos.Query),
		List:           setting.NewListHandler(p.SettingRepos.Query),
		ListSettings:   setting.NewListSettingsHandler(p.SettingRepos.Query, p.SettingRepos.CategoryQuery, p.SettingsCache),
		CreateCategory: setting.NewCreateCategoryHandler(p.SettingRepos.CategoryCommand, p.SettingRepos.CategoryQuery, p.SettingsCache),
		UpdateCategory: setting.NewUpdateCategoryHandler(p.SettingRepos.CategoryCommand, p.SettingRepos.CategoryQuery, p.SettingsCache),
		DeleteCategory: setting.NewDeleteCategoryHandler(p.SettingRepos.CategoryCommand, p.SettingRepos.CategoryQuery, p.SettingRepos.Query, p.SettingsCache),
		GetCategory:    setting.NewGetCategoryHandler(p.SettingRepos.CategoryQuery),
		ListCategories: setting.NewListCategoriesHandler(p.SettingRepos.CategoryQuery),
	}
}

// userSettingUseCasesParams 聚合 UserSetting 用例所需的依赖。
//
// 注意：UserSettingRepositories 来自 IAM 模块（跨模块依赖）。
type userSettingUseCasesParams struct {
	fx.In

	SettingRepos     corepersistence.SettingRepositories
	UserSettingRepos iampersistence.UserSettingRepositories
	SettingsCache    setting.SettingsCacheService
}

func newUserSettingUseCases(p userSettingUseCasesParams) *UserSettingUseCases {
	validator := validation.NewJSONLogicValidator()

	return &UserSettingUseCases{
		Set:            setting.NewUserSetHandler(p.SettingRepos.Query, p.UserSettingRepos.Command, validator),
		BatchSet:       setting.NewUserBatchSetHandler(p.SettingRepos.Query, p.UserSettingRepos.Command, validator),
		Reset:          setting.NewUserResetHandler(p.UserSettingRepos.Command),
		ResetAll:       setting.NewUserResetAllHandler(p.UserSettingRepos.Command),
		Get:            setting.NewUserGetHandler(p.SettingRepos.Query, p.UserSettingRepos.Query),
		List:           setting.NewUserListHandler(p.SettingRepos.Query, p.UserSettingRepos.Query),
		ListSettings:   setting.NewUserListSettingsHandler(p.SettingRepos.Query, p.UserSettingRepos.Query, p.SettingRepos.CategoryQuery, p.SettingsCache),
		ListCategories: setting.NewUserListCategoriesHandler(p.SettingRepos.Query, p.SettingRepos.CategoryQuery, p.SettingsCache),
	}
}

func newStatsUseCases(statsQuery domain_stats.QueryRepository) *StatsUseCases {
	return &StatsUseCases{
		GetStats: stats.NewGetStatsHandler(statsQuery),
	}
}

func newCaptchaUseCases(
	captchaCommand infra_captcha.CommandRepository,
	captchaSvc infra_captcha.Service,
) *CaptchaUseCases {
	return &CaptchaUseCases{
		Generate: app_captcha.NewGenerateHandler(captchaCommand, captchaSvc),
	}
}

// organizationUseCasesParams 聚合 Organization 用例所需的依赖。
type organizationUseCasesParams struct {
	fx.In

	OrgRepos iampersistence.OrganizationRepositories
}

func newOrganizationUseCases(p organizationUseCasesParams) *OrganizationUseCases {
	return &OrganizationUseCases{
		// Organization
		Create: org.NewCreateHandler(p.OrgRepos.Command, p.OrgRepos.Query, p.OrgRepos.MemberCommand),
		Update: org.NewUpdateHandler(p.OrgRepos.Command, p.OrgRepos.Query),
		Delete: org.NewDeleteHandler(
			p.OrgRepos.Command,
			p.OrgRepos.Query,
			p.OrgRepos.MemberQuery,
			p.OrgRepos.MemberCommand,
			p.OrgRepos.TeamQuery,
			p.OrgRepos.TeamCommand,
			p.OrgRepos.TeamMemberQuery,
			p.OrgRepos.TeamMemberCommand,
		),
		Get:  org.NewGetHandler(p.OrgRepos.Query),
		List: org.NewListHandler(p.OrgRepos.Query),

		// Member
		MemberAdd:        org.NewMemberAddHandler(p.OrgRepos.MemberCommand, p.OrgRepos.MemberQuery, p.OrgRepos.Query),
		MemberRemove:     org.NewMemberRemoveHandler(p.OrgRepos.MemberCommand, p.OrgRepos.MemberQuery),
		MemberUpdateRole: org.NewMemberUpdateRoleHandler(p.OrgRepos.MemberCommand, p.OrgRepos.MemberQuery),
		MemberList:       org.NewMemberListHandler(p.OrgRepos.MemberQuery),

		// Team
		TeamCreate: org.NewTeamCreateHandler(p.OrgRepos.TeamCommand, p.OrgRepos.TeamQuery, p.OrgRepos.Query, p.OrgRepos.TeamMemberCommand),
		TeamUpdate: org.NewTeamUpdateHandler(p.OrgRepos.TeamCommand, p.OrgRepos.TeamQuery),
		TeamDelete: org.NewTeamDeleteHandler(
			p.OrgRepos.TeamCommand,
			p.OrgRepos.TeamQuery,
			p.OrgRepos.TeamMemberQuery,
			p.OrgRepos.TeamMemberCommand,
		),
		TeamGet:  org.NewTeamGetHandler(p.OrgRepos.TeamQuery),
		TeamList: org.NewTeamListHandler(p.OrgRepos.TeamQuery),

		// Team Member
		TeamMemberAdd:    org.NewTeamMemberAddHandler(p.OrgRepos.TeamMemberCommand, p.OrgRepos.TeamMemberQuery, p.OrgRepos.TeamQuery, p.OrgRepos.MemberQuery),
		TeamMemberRemove: org.NewTeamMemberRemoveHandler(p.OrgRepos.TeamMemberCommand, p.OrgRepos.TeamQuery),
		TeamMemberList:   org.NewTeamMemberListHandler(p.OrgRepos.TeamMemberQuery, p.OrgRepos.TeamQuery),

		// User View
		UserOrgs:  org.NewUserOrgsHandler(p.OrgRepos.MemberQuery, p.OrgRepos.Query),
		UserTeams: org.NewUserTeamsHandler(p.OrgRepos.TeamMemberQuery, p.OrgRepos.TeamQuery, p.OrgRepos.Query),
	}
}
