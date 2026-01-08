package app

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/audit"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/auth"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/captcha"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/org"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/pat"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/role"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/twofa"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/user"
)

// UseCases 类型别名：为 Handler 层提供便捷访问
type (
	AuditUseCases      = audit.AuditUseCases
	CaptchaUseCases    = captcha.CaptchaUseCases
	AuthUseCases       = auth.AuthUseCases
	UserUseCases       = user.UserUseCases
	RoleUseCases       = role.RoleUseCases
	PATUseCases        = pat.PATUseCases
	TwoFAUseCases      = twofa.TwoFAUseCases
	OrgUseCases        = org.OrgUseCases
	OrgMemberUseCases  = org.OrgMemberUseCases
	TeamUseCases       = org.TeamUseCases
	TeamMemberUseCases = org.TeamMemberUseCases
	UserOrgUseCases    = org.UserOrgUseCases
)

// UseCaseModule 聚合所有 IAM 子模块
var UseCaseModule = fx.Module("iam.usecase",
	audit.Module,
	captcha.Module,
	auth.Module,
	user.Module,
	role.Module,
	pat.Module,
	twofa.Module,
	org.Module,
)
