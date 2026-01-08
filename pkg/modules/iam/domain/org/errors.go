package org

import "errors"

// 组织相关错误
var (
	// ErrOrgNotFound 组织不存在
	ErrOrgNotFound = errors.New("organization not found")

	// ErrOrgAlreadyExists 组织已存在
	ErrOrgAlreadyExists = errors.New("organization already exists")

	// ErrOrgNameAlreadyExists 组织名称已存在
	ErrOrgNameAlreadyExists = errors.New("organization name already exists")

	// ErrOrgSuspended 组织已被暂停
	ErrOrgSuspended = errors.New("organization has been suspended")

	// ErrOrgHasMembers 组织还有成员，不能删除
	ErrOrgHasMembers = errors.New("organization has members, cannot delete")

	// ErrOrgHasTeams 组织还有团队，不能删除
	ErrOrgHasTeams = errors.New("organization has teams, cannot delete")

	// ErrInvalidOrgID 无效的组织 ID
	ErrInvalidOrgID = errors.New("invalid organization ID")
)

// 团队相关错误
var (
	// ErrTeamNotFound 团队不存在
	ErrTeamNotFound = errors.New("team not found")

	// ErrTeamAlreadyExists 团队已存在
	ErrTeamAlreadyExists = errors.New("team already exists")

	// ErrTeamNameAlreadyExists 团队名称已存在（组织内）
	ErrTeamNameAlreadyExists = errors.New("team name already exists")

	// ErrTeamHasMembers 团队还有成员，不能删除
	ErrTeamHasMembers = errors.New("team has members, cannot delete")

	// ErrTeamNotInOrg 团队不属于该组织
	ErrTeamNotInOrg = errors.New("team does not belong to this organization")

	// ErrInvalidTeamID 无效的团队 ID
	ErrInvalidTeamID = errors.New("invalid team ID")
)

// 成员相关错误
var (
	// ErrMemberNotFound 成员不存在
	ErrMemberNotFound = errors.New("member not found")

	// ErrMemberAlreadyExists 成员已存在
	ErrMemberAlreadyExists = errors.New("member already exists")

	// ErrNotOrgMember 不是组织成员
	ErrNotOrgMember = errors.New("not an organization member")

	// ErrNotTeamMember 不是团队成员
	ErrNotTeamMember = errors.New("not a team member")

	// ErrInvalidMemberRole 无效的成员角色
	ErrInvalidMemberRole = errors.New("invalid member role")

	// ErrInvalidTeamMemberRole 无效的团队成员角色
	ErrInvalidTeamMemberRole = errors.New("invalid team member role")

	// ErrCannotRemoveLastOwner 不能移除最后一个所有者
	ErrCannotRemoveLastOwner = errors.New("cannot remove the last owner")

	// ErrCannotDemoteLastOwner 不能降级最后一个所有者
	ErrCannotDemoteLastOwner = errors.New("cannot demote the last owner")

	// ErrCannotRemoveSelf 不能移除自己
	ErrCannotRemoveSelf = errors.New("cannot remove yourself")

	// ErrMustBeOrgMemberFirst 必须先是组织成员才能加入团队
	ErrMustBeOrgMemberFirst = errors.New("must be an organization member before joining a team")

	// ErrTeamMemberAlreadyExists 团队成员已存在
	ErrTeamMemberAlreadyExists = errors.New("team member already exists")

	// ErrTeamMemberNotFound 团队成员不存在
	ErrTeamMemberNotFound = errors.New("team member not found")
)

// 权限相关错误
var (
	// ErrNoPermission 没有权限
	ErrNoPermission = errors.New("no permission to perform this action")

	// ErrOnlyOwnerCanTransfer 只有所有者可以转让所有权
	ErrOnlyOwnerCanTransfer = errors.New("only owner can transfer ownership")

	// ErrOnlyAdminCanManage 只有管理员可以管理成员
	ErrOnlyAdminCanManage = errors.New("only admin can manage members")
)
