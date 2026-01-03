package app

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/audit"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/auth"
	captchaApp "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/captcha"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/org"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/pat"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/role"
	twofaApp "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/twofa"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/user"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/usersetting"
	authDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/auth"
	twofaDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/twofa"
	authInfra "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infra/auth"
	persistence "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infra/persistence"
	captchaInfra "github.com/lwmacct/260103-ddd-shared/pkg/shared/captcha"
	"github.com/lwmacct/260103-ddd-shared/pkg/shared/event"
)

// --- 用例模块结构体 ---

// AuditUseCases 审计日志用例处理器
type AuditUseCases struct {
	CreateLog *audit.CreateHandler
	Get       *audit.GetHandler
	List      *audit.ListHandler
}

// CaptchaUseCases 验证码用例处理器
type CaptchaUseCases struct {
	Generate *captchaApp.GenerateHandler
}

// AuthUseCases 认证用例处理器
type AuthUseCases struct {
	Login        *auth.LoginHandler
	Login2FA     *auth.Login2FAHandler
	Register     *auth.RegisterHandler
	RefreshToken *auth.RefreshTokenHandler
}

// UserUseCases 用户管理用例处理器
type UserUseCases struct {
	Create         *user.CreateHandler
	Update         *user.UpdateHandler
	Delete         *user.DeleteHandler
	AssignRoles    *user.AssignRolesHandler
	ChangePassword *user.ChangePasswordHandler
	BatchCreate    *user.BatchCreateHandler
	Get            *user.GetHandler
	List           *user.ListHandler
}

// RoleUseCases 角色管理用例处理器
type RoleUseCases struct {
	Create         *role.CreateHandler
	Update         *role.UpdateHandler
	Delete         *role.DeleteHandler
	SetPermissions *role.SetPermissionsHandler
	Get            *role.GetHandler
	List           *role.ListHandler
}

// PATUseCases 个人访问令牌用例处理器
type PATUseCases struct {
	Create  *pat.CreateHandler
	Delete  *pat.DeleteHandler
	Disable *pat.DisableHandler
	Enable  *pat.EnableHandler
	Get     *pat.GetHandler
	List    *pat.ListHandler
}

// TwoFAUseCases 双因素认证用例处理器
type TwoFAUseCases struct {
	Setup        *twofaApp.SetupHandler
	VerifyEnable *twofaApp.VerifyEnableHandler
	Disable      *twofaApp.DisableHandler
	GetStatus    *twofaApp.GetStatusHandler
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

// UserSettingUseCases 用户设置用例处理器（从 usersetting 包导入）
type UserSettingUseCases = usersetting.UserSettingUseCases

// --- Fx 模块 ---

// UseCaseModule 提供按领域组织的 IAM 模块用例处理器。
var UseCaseModule = fx.Module("iam.usecase",
	fx.Provide(
		newAuditUseCases,
		newCaptchaUseCases,
		newAuthUseCases,
		newUserUseCases,
		newRoleUseCases,
		newPATUseCases,
		newTwoFAUseCases,
		newOrganizationUseCases,
		// UserSettingUseCases 由 usersetting.UseCaseModule 提供
	),
)

// --- 构造函数 ---

func newAuditUseCases(repos persistence.AuditRepositories) *AuditUseCases {
	return &AuditUseCases{
		CreateLog: audit.NewCreateHandler(repos.Command),
		Get:       audit.NewGetHandler(repos.Query),
		List:      audit.NewListHandler(repos.Query),
	}
}

type captchaUseCasesParams struct {
	fx.In

	CaptchaCommand captchaInfra.CommandRepository
	CaptchaSvc     captchaInfra.Service
}

func newCaptchaUseCases(p captchaUseCasesParams) *CaptchaUseCases {
	return &CaptchaUseCases{
		Generate: captchaApp.NewGenerateHandler(p.CaptchaCommand, p.CaptchaSvc),
	}
}

// authUseCasesParams 聚合 Auth 用例所需的依赖。
type authUseCasesParams struct {
	fx.In

	UserRepos      persistence.UserRepositories
	CaptchaCommand captchaInfra.CommandRepository
	TwoFARepos     persistence.TwoFARepositories
	AuthSvc        authDomain.Service
	LoginSession   authDomain.SessionService
	TwoFASvc       twofaDomain.Service
	Audit          *AuditUseCases // IAM 内部依赖
}

func newAuthUseCases(p authUseCasesParams) *AuthUseCases {
	return &AuthUseCases{
		Login:        auth.NewLoginHandler(p.UserRepos.Query, p.CaptchaCommand, p.TwoFARepos.Query, p.AuthSvc, p.LoginSession, p.Audit.CreateLog),
		Login2FA:     auth.NewLogin2FAHandler(p.UserRepos.Query, p.AuthSvc, p.LoginSession, p.TwoFASvc, p.Audit.CreateLog),
		Register:     auth.NewRegisterHandler(p.UserRepos.Command, p.UserRepos.Query, p.AuthSvc),
		RefreshToken: auth.NewRefreshTokenHandler(p.UserRepos.Query, p.AuthSvc, p.Audit.CreateLog),
	}
}

// userUseCasesParams 聚合 User 用例所需的依赖。
type userUseCasesParams struct {
	fx.In

	UserRepos persistence.UserRepositories
	AuthSvc   authDomain.Service
	EventBus  event.EventBus
}

func newUserUseCases(p userUseCasesParams) *UserUseCases {
	return &UserUseCases{
		Create:         user.NewCreateHandler(p.UserRepos.Command, p.UserRepos.Query, p.AuthSvc),
		Update:         user.NewUpdateHandler(p.UserRepos.Command, p.UserRepos.Query),
		Delete:         user.NewDeleteHandler(p.UserRepos.Command, p.UserRepos.Query, p.EventBus),
		AssignRoles:    user.NewAssignRolesHandler(p.UserRepos.Command, p.UserRepos.Query, p.EventBus),
		ChangePassword: user.NewChangePasswordHandler(p.UserRepos.Command, p.UserRepos.Query, p.AuthSvc),
		BatchCreate:    user.NewBatchCreateHandler(p.UserRepos.Command, p.UserRepos.Query, p.AuthSvc),
		Get:            user.NewGetHandler(p.UserRepos.Query),
		List:           user.NewListHandler(p.UserRepos.Query),
	}
}

// roleUseCasesParams 聚合 Role 用例所需的依赖。
type roleUseCasesParams struct {
	fx.In

	RoleRepos persistence.RoleRepositories
	EventBus  event.EventBus
}

func newRoleUseCases(p roleUseCasesParams) *RoleUseCases {
	return &RoleUseCases{
		Create:         role.NewCreateHandler(p.RoleRepos.Command, p.RoleRepos.Query),
		Update:         role.NewUpdateHandler(p.RoleRepos.Command, p.RoleRepos.Query),
		Delete:         role.NewDeleteHandler(p.RoleRepos.Command, p.RoleRepos.Query),
		SetPermissions: role.NewSetPermissionsHandler(p.RoleRepos.Command, p.RoleRepos.Query, p.EventBus),
		Get:            role.NewGetHandler(p.RoleRepos.Query),
		List:           role.NewListHandler(p.RoleRepos.Query),
	}
}

// patUseCasesParams 聚合 PAT 用例所需的依赖。
type patUseCasesParams struct {
	fx.In

	PATRepos  persistence.PATRepositories
	UserRepos persistence.UserRepositories
	TokenGen  *authInfra.TokenGenerator
}

func newPATUseCases(p patUseCasesParams) *PATUseCases {
	return &PATUseCases{
		Create:  pat.NewCreateHandler(p.PATRepos.Command, p.UserRepos.Query, p.TokenGen),
		Delete:  pat.NewDeleteHandler(p.PATRepos.Command, p.PATRepos.Query),
		Disable: pat.NewDisableHandler(p.PATRepos.Command, p.PATRepos.Query),
		Enable:  pat.NewEnableHandler(p.PATRepos.Command, p.PATRepos.Query),
		Get:     pat.NewGetHandler(p.PATRepos.Query),
		List:    pat.NewListHandler(p.PATRepos.Query),
	}
}

func newTwoFAUseCases(twofaSvc twofaDomain.Service) *TwoFAUseCases {
	return &TwoFAUseCases{
		Setup:        twofaApp.NewSetupHandler(twofaSvc),
		VerifyEnable: twofaApp.NewVerifyEnableHandler(twofaSvc),
		Disable:      twofaApp.NewDisableHandler(twofaSvc),
		GetStatus:    twofaApp.NewGetStatusHandler(twofaSvc),
	}
}

// organizationUseCasesParams 聚合 Organization 用例所需的依赖。
type organizationUseCasesParams struct {
	fx.In

	OrgRepos persistence.OrganizationRepositories
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
