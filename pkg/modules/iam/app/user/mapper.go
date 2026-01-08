package user

import (
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/user"
)

// stringPtrValue 将 *string 转换为 string，nil 返回空字符串
func stringPtrValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// ToUserDTO 将领域模型 User 转换为应用层 UserDTO
func ToUserDTO(u *user.User) *UserDTO {
	if u == nil {
		return nil
	}

	return &UserDTO{
		ID:        u.ID,
		Username:  u.Username,
		Email:     stringPtrValue(u.Email),
		RealName:  u.RealName,
		Nickname:  u.Nickname,
		Phone:     stringPtrValue(u.Phone),
		Signature: u.Signature,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		Status:    u.Status,
		Type:      string(u.Type),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ToUserWithRolesDTO 将领域模型 User 转换为应用层 UserWithRolesDTO（包含角色信息）
func ToUserWithRolesDTO(u *user.User) *UserWithRolesDTO {
	if u == nil {
		return nil
	}

	// 转换角色信息（nil 或空切片返回 nil，避免 JSON 中出现 []）
	var roles []RoleDTO
	if len(u.Roles) > 0 {
		roles = make([]RoleDTO, 0, len(u.Roles))
		for _, role := range u.Roles {
			roles = append(roles, RoleDTO{
				ID:          role.ID,
				Name:        role.Name,
				DisplayName: role.DisplayName,
				Description: role.Description,
			})
		}
	}

	return &UserWithRolesDTO{
		ID:        u.ID,
		Username:  u.Username,
		Email:     stringPtrValue(u.Email),
		RealName:  u.RealName,
		Nickname:  u.Nickname,
		Phone:     stringPtrValue(u.Phone),
		Signature: u.Signature,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		Status:    u.Status,
		Type:      string(u.Type),
		Roles:     roles,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
