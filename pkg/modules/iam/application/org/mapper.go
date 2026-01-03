package org

import "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/org"

// ============================================================================
// Organization Mappers
// ============================================================================

// ToOrgDTO 将组织实体转换为 DTO
func ToOrgDTO(org *org.Org) *OrgDTO {
	if org == nil {
		return nil
	}
	return &OrgDTO{
		ID:          org.ID,
		Name:        org.Name,
		DisplayName: org.DisplayName,
		Description: org.Description,
		Avatar:      org.Avatar,
		Status:      org.Status,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}
}

// ToOrgDTOs 将组织实体列表转换为 DTO 列表
func ToOrgDTOs(orgs []*org.Org) []*OrgDTO {
	if len(orgs) == 0 {
		return []*OrgDTO{}
	}
	dtos := make([]*OrgDTO, 0, len(orgs))
	for _, org := range orgs {
		if dto := ToOrgDTO(org); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}

// ============================================================================
// Team Mappers
// ============================================================================

// ToTeamDTO 将团队实体转换为 DTO
func ToTeamDTO(team *org.Team) *TeamDTO {
	if team == nil {
		return nil
	}
	return &TeamDTO{
		ID:          team.ID,
		OrgID:       team.OrgID,
		Name:        team.Name,
		DisplayName: team.DisplayName,
		Description: team.Description,
		Avatar:      team.Avatar,
		CreatedAt:   team.CreatedAt,
		UpdatedAt:   team.UpdatedAt,
	}
}

// ToTeamDTOs 将团队实体列表转换为 DTO 列表
func ToTeamDTOs(teams []*org.Team) []*TeamDTO {
	if len(teams) == 0 {
		return []*TeamDTO{}
	}
	dtos := make([]*TeamDTO, 0, len(teams))
	for _, team := range teams {
		if dto := ToTeamDTO(team); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}

// ============================================================================
// Member Mappers
// ============================================================================

// ToMemberDTO 将成员实体转换为 DTO
func ToMemberDTO(member *org.Member) *MemberDTO {
	if member == nil {
		return nil
	}
	return &MemberDTO{
		ID:       member.ID,
		OrgID:    member.OrgID,
		UserID:   member.UserID,
		Role:     string(member.Role),
		JoinedAt: member.JoinedAt,
	}
}

// ToMemberDTOs 将成员实体列表转换为 DTO 列表
func ToMemberDTOs(members []*org.Member) []*MemberDTO {
	if len(members) == 0 {
		return []*MemberDTO{}
	}
	dtos := make([]*MemberDTO, 0, len(members))
	for _, member := range members {
		if dto := ToMemberDTO(member); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}

// ToMemberWithUserDTO 将 MemberWithUser 值对象转换为 DTO（包含用户信息）
func ToMemberWithUserDTO(mwu *org.MemberWithUser) *MemberDTO {
	if mwu == nil {
		return nil
	}
	return &MemberDTO{
		ID:       mwu.ID,
		OrgID:    mwu.OrgID,
		UserID:   mwu.UserID,
		Role:     string(mwu.Role),
		JoinedAt: mwu.JoinedAt,
		// 用户信息
		Username: mwu.Username,
		Email:    mwu.Email,
		FullName: mwu.FullName,
		Avatar:   mwu.Avatar,
	}
}

// ToMemberWithUserDTOs 将 MemberWithUser 值对象列表转换为 DTO 列表
func ToMemberWithUserDTOs(mwuList []*org.MemberWithUser) []*MemberDTO {
	if len(mwuList) == 0 {
		return []*MemberDTO{}
	}
	dtos := make([]*MemberDTO, 0, len(mwuList))
	for _, mwu := range mwuList {
		if dto := ToMemberWithUserDTO(mwu); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}

// ============================================================================
// Team Member Mappers
// ============================================================================

// ToTeamMemberDTO 将团队成员实体转换为 DTO
func ToTeamMemberDTO(member *org.TeamMember) *TeamMemberDTO {
	if member == nil {
		return nil
	}
	return &TeamMemberDTO{
		ID:       member.ID,
		TeamID:   member.TeamID,
		UserID:   member.UserID,
		Role:     string(member.Role),
		JoinedAt: member.JoinedAt,
	}
}

// ToTeamMemberDTOs 将团队成员实体列表转换为 DTO 列表
func ToTeamMemberDTOs(members []*org.TeamMember) []*TeamMemberDTO {
	if len(members) == 0 {
		return []*TeamMemberDTO{}
	}
	dtos := make([]*TeamMemberDTO, 0, len(members))
	for _, member := range members {
		if dto := ToTeamMemberDTO(member); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}
