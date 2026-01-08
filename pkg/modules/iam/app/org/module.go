package org

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/infra/persistence"
)

// OrgUseCases 组织管理用例处理器
type OrgUseCases struct {
	Create *CreateHandler
	Update *UpdateHandler
	Delete *DeleteHandler
	Get    *GetHandler
	List   *ListHandler
}

// OrgMemberUseCases 组织成员管理用例处理器
type OrgMemberUseCases struct {
	Add        *MemberAddHandler
	Remove     *MemberRemoveHandler
	UpdateRole *MemberUpdateRoleHandler
	List       *MemberListHandler
}

// TeamUseCases 团队管理用例处理器
type TeamUseCases struct {
	Create *TeamCreateHandler
	Update *TeamUpdateHandler
	Delete *TeamDeleteHandler
	Get    *TeamGetHandler
	List   *TeamListHandler
}

// TeamMemberUseCases 团队成员管理用例处理器
type TeamMemberUseCases struct {
	Add    *TeamMemberAddHandler
	Remove *TeamMemberRemoveHandler
	List   *TeamMemberListHandler
}

// UserOrgUseCases 用户视角的组织/团队查询用例处理器
type UserOrgUseCases struct {
	ListOrgs  *UserOrgsHandler
	ListTeams *UserTeamsHandler
}

// Module 注册 Org 子模块依赖
var Module = fx.Module("iam.org",
	fx.Provide(
		newOrgUseCases,
		newOrgMemberUseCases,
		newTeamUseCases,
		newTeamMemberUseCases,
		newUserOrgUseCases,
	),
)

type orgUseCasesParams struct {
	fx.In

	OrgRepos persistence.OrganizationRepositories
}

func newOrgUseCases(p orgUseCasesParams) *OrgUseCases {
	return &OrgUseCases{
		Create: NewCreateHandler(p.OrgRepos.Command, p.OrgRepos.Query, p.OrgRepos.MemberCommand),
		Update: NewUpdateHandler(p.OrgRepos.Command, p.OrgRepos.Query),
		Delete: NewDeleteHandler(
			p.OrgRepos.Command,
			p.OrgRepos.Query,
			p.OrgRepos.MemberQuery,
			p.OrgRepos.MemberCommand,
			p.OrgRepos.TeamQuery,
			p.OrgRepos.TeamCommand,
			p.OrgRepos.TeamMemberQuery,
			p.OrgRepos.TeamMemberCommand,
		),
		Get:  NewGetHandler(p.OrgRepos.Query),
		List: NewListHandler(p.OrgRepos.Query),
	}
}

func newOrgMemberUseCases(p orgUseCasesParams) *OrgMemberUseCases {
	return &OrgMemberUseCases{
		Add:        NewMemberAddHandler(p.OrgRepos.MemberCommand, p.OrgRepos.MemberQuery, p.OrgRepos.Query),
		Remove:     NewMemberRemoveHandler(p.OrgRepos.MemberCommand, p.OrgRepos.MemberQuery),
		UpdateRole: NewMemberUpdateRoleHandler(p.OrgRepos.MemberCommand, p.OrgRepos.MemberQuery),
		List:       NewMemberListHandler(p.OrgRepos.MemberQuery),
	}
}

func newTeamUseCases(p orgUseCasesParams) *TeamUseCases {
	return &TeamUseCases{
		Create: NewTeamCreateHandler(p.OrgRepos.TeamCommand, p.OrgRepos.TeamQuery, p.OrgRepos.Query, p.OrgRepos.TeamMemberCommand),
		Update: NewTeamUpdateHandler(p.OrgRepos.TeamCommand, p.OrgRepos.TeamQuery),
		Delete: NewTeamDeleteHandler(
			p.OrgRepos.TeamCommand,
			p.OrgRepos.TeamQuery,
			p.OrgRepos.TeamMemberQuery,
			p.OrgRepos.TeamMemberCommand,
		),
		Get:  NewTeamGetHandler(p.OrgRepos.TeamQuery),
		List: NewTeamListHandler(p.OrgRepos.TeamQuery),
	}
}

func newTeamMemberUseCases(p orgUseCasesParams) *TeamMemberUseCases {
	return &TeamMemberUseCases{
		Add:    NewTeamMemberAddHandler(p.OrgRepos.TeamMemberCommand, p.OrgRepos.TeamMemberQuery, p.OrgRepos.TeamQuery, p.OrgRepos.MemberQuery),
		Remove: NewTeamMemberRemoveHandler(p.OrgRepos.TeamMemberCommand, p.OrgRepos.TeamQuery),
		List:   NewTeamMemberListHandler(p.OrgRepos.TeamMemberQuery, p.OrgRepos.TeamQuery),
	}
}

func newUserOrgUseCases(p orgUseCasesParams) *UserOrgUseCases {
	return &UserOrgUseCases{
		ListOrgs:  NewUserOrgsHandler(p.OrgRepos.MemberQuery, p.OrgRepos.Query),
		ListTeams: NewUserTeamsHandler(p.OrgRepos.TeamMemberQuery, p.OrgRepos.TeamQuery, p.OrgRepos.Query),
	}
}
