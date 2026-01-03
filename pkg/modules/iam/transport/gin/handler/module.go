package handler

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/config"
	iamapplication "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/application"
)

// HandlersResult 使用 fx.Out 批量返回 IAM 模块的所有 HTTP 处理器。
type HandlersResult struct {
	fx.Out

	Auth        *AuthHandler
	UserProfile *UserProfileHandler
	AdminUser   *AdminUserHandler
	Role        *RoleHandler
	PAT         *PATHandler
	TwoFA       *TwoFAHandler
	UserOrg     *UserOrgHandler
	Audit       *AuditHandler
	Captcha     *CaptchaHandler
	Org         *OrgHandler
	OrgMember   *OrgMemberHandler
	Team        *TeamHandler
	TeamMember  *TeamMemberHandler
}

// HandlerModule 提供 IAM 模块的所有 HTTP 处理器。
var HandlerModule = fx.Module("iam.handler",
	fx.Provide(newAllHandlers),
)

// handlersParams 聚合创建 Handler 所需的依赖。
type handlersParams struct {
	fx.In

	Config *config.Config

	// IAM 模块用例
	Auth         *iamapplication.AuthUseCases
	User         *iamapplication.UserUseCases
	Role         *iamapplication.RoleUseCases
	PAT          *iamapplication.PATUseCases
	TwoFA        *iamapplication.TwoFAUseCases
	Audit        *iamapplication.AuditUseCases
	Captcha      *iamapplication.CaptchaUseCases
	Organization *iamapplication.OrganizationUseCases
}

func newAllHandlers(p handlersParams) HandlersResult {
	return HandlersResult{
		Auth: NewAuthHandler(
			p.Auth.Login,
			p.Auth.Login2FA,
			p.Auth.Register,
			p.Auth.RefreshToken,
		),
		UserProfile: NewUserProfileHandler(
			p.User.Get,
			p.User.Update,
			p.User.ChangePassword,
			p.User.Delete,
		),
		AdminUser: NewAdminUserHandler(
			p.User.Create,
			p.User.Update,
			p.User.Delete,
			p.User.AssignRoles,
			p.User.BatchCreate,
			p.User.Get,
			p.User.List,
		),
		Role: NewRoleHandler(
			p.Role.Create,
			p.Role.Update,
			p.Role.Delete,
			p.Role.SetPermissions,
			p.Role.Get,
			p.Role.List,
		),
		PAT: NewPATHandler(
			p.PAT.Create,
			p.PAT.Delete,
			p.PAT.Disable,
			p.PAT.Enable,
			p.PAT.Get,
			p.PAT.List,
		),
		TwoFA: NewTwoFAHandler(
			p.TwoFA.Setup,
			p.TwoFA.VerifyEnable,
			p.TwoFA.Disable,
			p.TwoFA.GetStatus,
		),
		UserOrg: NewUserOrgHandler(
			p.Organization.UserOrgs,
			p.Organization.UserTeams,
		),
		Audit: NewAuditHandler(
			p.Audit.List,
			p.Audit.Get,
		),
		Captcha: NewCaptchaHandler(
			p.Captcha.Generate,
			p.Config.Auth.DevSecret,
		),
		Org: NewOrgHandler(
			p.Organization.Create,
			p.Organization.Update,
			p.Organization.Delete,
			p.Organization.Get,
			p.Organization.List,
		),
		OrgMember: NewOrgMemberHandler(
			p.Organization.MemberAdd,
			p.Organization.MemberRemove,
			p.Organization.MemberUpdateRole,
			p.Organization.MemberList,
		),
		Team: NewTeamHandler(
			p.Organization.TeamCreate,
			p.Organization.TeamUpdate,
			p.Organization.TeamDelete,
			p.Organization.TeamGet,
			p.Organization.TeamList,
		),
		TeamMember: NewTeamMemberHandler(
			p.Organization.TeamMemberAdd,
			p.Organization.TeamMemberRemove,
			p.Organization.TeamMemberList,
		),
	}
}
