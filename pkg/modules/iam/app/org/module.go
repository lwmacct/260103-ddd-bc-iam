package org

import (
	"go.uber.org/fx"
)

// OrgUseCases 组织管理用例处理器聚合
type OrgUseCases struct {
	Create *CreateHandler
	Update *UpdateHandler
	Delete *DeleteHandler
	Get    *GetHandler
	List   *ListHandler
}

// OrgMemberUseCases 组织成员管理用例处理器聚合
type OrgMemberUseCases struct {
	Add        *MemberAddHandler
	Remove     *MemberRemoveHandler
	UpdateRole *MemberUpdateRoleHandler
	List       *MemberListHandler
}

// TeamUseCases 团队管理用例处理器聚合
type TeamUseCases struct {
	Create *TeamCreateHandler
	Update *TeamUpdateHandler
	Delete *TeamDeleteHandler
	Get    *TeamGetHandler
	List   *TeamListHandler
}

// TeamMemberUseCases 团队成员管理用例处理器聚合
type TeamMemberUseCases struct {
	Add    *TeamMemberAddHandler
	Remove *TeamMemberRemoveHandler
	List   *TeamMemberListHandler
}

// UserOrgUseCases 用户视角的组织/团队查询用例处理器聚合
type UserOrgUseCases struct {
	ListOrgs  *UserOrgsHandler
	ListTeams *UserTeamsHandler
}

// Module 注册 Org 子模块所有依赖
var Module = fx.Module("iam.org",
	fx.Provide(
		// Org
		NewCreateHandler,
		NewUpdateHandler,
		NewDeleteHandler,
		NewGetHandler,
		NewListHandler,

		// OrgMember
		NewMemberAddHandler,
		NewMemberRemoveHandler,
		NewMemberUpdateRoleHandler,
		NewMemberListHandler,

		// Team
		NewTeamCreateHandler,
		NewTeamUpdateHandler,
		NewTeamDeleteHandler,
		NewTeamGetHandler,
		NewTeamListHandler,

		// TeamMember
		NewTeamMemberAddHandler,
		NewTeamMemberRemoveHandler,
		NewTeamMemberListHandler,

		// UserOrg
		NewUserOrgsHandler,
		NewUserTeamsHandler,

		// UseCases
		newOrgUseCases,
		newOrgMemberUseCases,
		newTeamUseCases,
		newTeamMemberUseCases,
		newUserOrgUseCases,
	),
)

// --- OrgUseCases ---

type orgUseCasesParams struct {
	fx.In

	Create *CreateHandler
	Update *UpdateHandler
	Delete *DeleteHandler
	Get    *GetHandler
	List   *ListHandler
}

func newOrgUseCases(p orgUseCasesParams) *OrgUseCases {
	return &OrgUseCases{
		Create: p.Create,
		Update: p.Update,
		Delete: p.Delete,
		Get:    p.Get,
		List:   p.List,
	}
}

// --- OrgMemberUseCases ---

type orgMemberUseCasesParams struct {
	fx.In

	Add        *MemberAddHandler
	Remove     *MemberRemoveHandler
	UpdateRole *MemberUpdateRoleHandler
	List       *MemberListHandler
}

func newOrgMemberUseCases(p orgMemberUseCasesParams) *OrgMemberUseCases {
	return &OrgMemberUseCases{
		Add:        p.Add,
		Remove:     p.Remove,
		UpdateRole: p.UpdateRole,
		List:       p.List,
	}
}

// --- TeamUseCases ---

type teamUseCasesParams struct {
	fx.In

	Create *TeamCreateHandler
	Update *TeamUpdateHandler
	Delete *TeamDeleteHandler
	Get    *TeamGetHandler
	List   *TeamListHandler
}

func newTeamUseCases(p teamUseCasesParams) *TeamUseCases {
	return &TeamUseCases{
		Create: p.Create,
		Update: p.Update,
		Delete: p.Delete,
		Get:    p.Get,
		List:   p.List,
	}
}

// --- TeamMemberUseCases ---

type teamMemberUseCasesParams struct {
	fx.In

	Add    *TeamMemberAddHandler
	Remove *TeamMemberRemoveHandler
	List   *TeamMemberListHandler
}

func newTeamMemberUseCases(p teamMemberUseCasesParams) *TeamMemberUseCases {
	return &TeamMemberUseCases{
		Add:    p.Add,
		Remove: p.Remove,
		List:   p.List,
	}
}

// --- UserOrgUseCases ---

type userOrgUseCasesParams struct {
	fx.In

	ListOrgs  *UserOrgsHandler
	ListTeams *UserTeamsHandler
}

func newUserOrgUseCases(p userOrgUseCasesParams) *UserOrgUseCases {
	return &UserOrgUseCases{
		ListOrgs:  p.ListOrgs,
		ListTeams: p.ListTeams,
	}
}
