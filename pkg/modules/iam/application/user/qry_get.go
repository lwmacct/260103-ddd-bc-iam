package user

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/user"
)

// GetHandler 获取用户查询处理器
type GetHandler struct {
	userQueryRepo user.QueryRepository
}

// NewGetHandler 创建获取用户查询处理器
func NewGetHandler(userQueryRepo user.QueryRepository) *GetHandler {
	return &GetHandler{
		userQueryRepo: userQueryRepo,
	}
}

// Handle 处理获取用户查询
func (h *GetHandler) Handle(ctx context.Context, query GetQuery) (*UserWithRolesDTO, error) {
	var u *user.User
	var err error

	if query.WithRoles {
		u, err = h.userQueryRepo.GetByIDWithRoles(ctx, query.UserID)
	} else {
		u, err = h.userQueryRepo.GetByID(ctx, query.UserID)
	}

	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	response := &UserWithRolesDTO{
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

	if query.WithRoles {
		roles := make([]RoleDTO, 0, len(u.Roles))
		for _, r := range u.Roles {
			roles = append(roles, RoleDTO{
				ID:          r.ID,
				Name:        r.Name,
				DisplayName: r.DisplayName,
				Description: r.Description,
			})
		}
		response.Roles = roles
	}

	return response, nil
}
