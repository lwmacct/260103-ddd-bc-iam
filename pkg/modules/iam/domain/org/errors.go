package org

import "errors"

// 组织相关错误
var (
	// ErrOrgNotFound 组织不存在
	ErrOrgNotFound = errors.New("组织不存在")

	// ErrOrgAlreadyExists 组织已存在
	ErrOrgAlreadyExists = errors.New("组织已存在")

	// ErrOrgNameAlreadyExists 组织名称已存在
	ErrOrgNameAlreadyExists = errors.New("组织名称已存在")

	// ErrOrgSuspended 组织已被暂停
	ErrOrgSuspended = errors.New("组织已被暂停")

	// ErrOrgHasMembers 组织还有成员，不能删除
	ErrOrgHasMembers = errors.New("组织还有成员，不能删除")

	// ErrOrgHasTeams 组织还有团队，不能删除
	ErrOrgHasTeams = errors.New("组织还有团队，不能删除")

	// ErrInvalidOrgID 无效的组织 ID
	ErrInvalidOrgID = errors.New("无效的组织 ID")
)

// 团队相关错误
var (
	// ErrTeamNotFound 团队不存在
	ErrTeamNotFound = errors.New("团队不存在")

	// ErrTeamAlreadyExists 团队已存在
	ErrTeamAlreadyExists = errors.New("团队已存在")

	// ErrTeamNameAlreadyExists 团队名称已存在（组织内）
	ErrTeamNameAlreadyExists = errors.New("团队名称已存在")

	// ErrTeamHasMembers 团队还有成员，不能删除
	ErrTeamHasMembers = errors.New("团队还有成员，不能删除")

	// ErrTeamNotInOrg 团队不属于该组织
	ErrTeamNotInOrg = errors.New("团队不属于该组织")

	// ErrInvalidTeamID 无效的团队 ID
	ErrInvalidTeamID = errors.New("无效的团队 ID")
)

// 成员相关错误
var (
	// ErrMemberNotFound 成员不存在
	ErrMemberNotFound = errors.New("成员不存在")

	// ErrMemberAlreadyExists 成员已存在
	ErrMemberAlreadyExists = errors.New("成员已存在")

	// ErrNotOrgMember 不是组织成员
	ErrNotOrgMember = errors.New("不是组织成员")

	// ErrNotTeamMember 不是团队成员
	ErrNotTeamMember = errors.New("不是团队成员")

	// ErrInvalidMemberRole 无效的成员角色
	ErrInvalidMemberRole = errors.New("无效的成员角色")

	// ErrInvalidTeamMemberRole 无效的团队成员角色
	ErrInvalidTeamMemberRole = errors.New("无效的团队成员角色")

	// ErrCannotRemoveLastOwner 不能移除最后一个所有者
	ErrCannotRemoveLastOwner = errors.New("不能移除最后一个所有者")

	// ErrCannotDemoteLastOwner 不能降级最后一个所有者
	ErrCannotDemoteLastOwner = errors.New("不能降级最后一个所有者")

	// ErrCannotRemoveSelf 不能移除自己
	ErrCannotRemoveSelf = errors.New("不能移除自己")

	// ErrMustBeOrgMemberFirst 必须先是组织成员才能加入团队
	ErrMustBeOrgMemberFirst = errors.New("必须先成为组织成员才能加入团队")

	// ErrTeamMemberAlreadyExists 团队成员已存在
	ErrTeamMemberAlreadyExists = errors.New("团队成员已存在")

	// ErrTeamMemberNotFound 团队成员不存在
	ErrTeamMemberNotFound = errors.New("团队成员不存在")
)

// 权限相关错误
var (
	// ErrNoPermission 没有权限
	ErrNoPermission = errors.New("没有权限执行此操作")

	// ErrOnlyOwnerCanTransfer 只有所有者可以转让所有权
	ErrOnlyOwnerCanTransfer = errors.New("只有所有者可以转让所有权")

	// ErrOnlyAdminCanManage 只有管理员可以管理成员
	ErrOnlyAdminCanManage = errors.New("只有管理员可以管理成员")
)
