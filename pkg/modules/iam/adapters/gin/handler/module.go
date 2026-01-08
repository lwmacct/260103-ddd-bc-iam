package handler

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/config"
)

// Handlers 聚合 IAM 模块的所有 HTTP 处理器。
//
// 设计说明：使用聚合结构体而非 fx.Out 导出单独 Handler，
// 符合 Fx 规范（禁止 fx.Out 导出 12+ 字段），维护成本更低。
type Handlers struct {
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
	fx.Provide(NewHandlers),
)

// HandlersParams 聚合创建 Handler 所需的依赖。
type HandlersParams struct {
	fx.In

	IAMConfig *config.Config

	// IAM 模块用例（12 个清晰的领域聚合）
	Auth       *app.AuthUseCases
	User       *app.UserUseCases
	Role       *app.RoleUseCases
	PAT        *app.PATUseCases
	TwoFA      *app.TwoFAUseCases
	Audit      *app.AuditUseCases
	Captcha    *app.CaptchaUseCases
	Org        *app.OrgUseCases
	OrgMember  *app.OrgMemberUseCases
	Team       *app.TeamUseCases
	TeamMember *app.TeamMemberUseCases
	UserOrg    *app.UserOrgUseCases
}

// NewHandlers 创建所有 IAM HTTP 处理器。
func NewHandlers(p HandlersParams) *Handlers {
	return &Handlers{
		Auth:        NewAuthHandler(p.Auth),
		UserProfile: NewUserProfileHandler(p.User),
		AdminUser:   NewAdminUserHandler(p.User),
		Role:        NewRoleHandler(p.Role),
		PAT:         NewPATHandler(p.PAT),
		TwoFA:       NewTwoFAHandler(p.TwoFA),
		UserOrg:     NewUserOrgHandler(p.UserOrg),
		Audit:       NewAuditHandler(p.Audit),
		Captcha:     NewCaptchaHandler(p.Captcha, p.IAMConfig.Auth.DevSecret),
		Org:         NewOrgHandler(p.Org),
		OrgMember:   NewOrgMemberHandler(p.OrgMember),
		Team:        NewTeamHandler(p.Team),
		TeamMember:  NewTeamMemberHandler(p.TeamMember),
	}
}
