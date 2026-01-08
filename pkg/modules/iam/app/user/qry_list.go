package user

import (
	"context"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/user"
)

// ListHandler 获取用户列表查询处理器
type ListHandler struct {
	userQueryRepo user.QueryRepository
}

// NewListHandler 创建获取用户列表查询处理器
func NewListHandler(userQueryRepo user.QueryRepository) *ListHandler {
	return &ListHandler{
		userQueryRepo: userQueryRepo,
	}
}

// Handle 处理获取用户列表查询
func (h *ListHandler) Handle(ctx context.Context, query ListQuery) (*UserListDTO, error) {
	// 根据是否有搜索关键词选择不同的查询方法
	users, total, err := h.fetchUsers(ctx, query)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO（使用统一的 mapper 函数）
	userResponses := make([]*UserDTO, 0, len(users))
	for _, u := range users {
		userResponses = append(userResponses, ToUserDTO(u))
	}

	return &UserListDTO{
		Users: userResponses,
		Total: total,
	}, nil
}

// fetchUsers 根据查询条件获取用户列表和总数
func (h *ListHandler) fetchUsers(ctx context.Context, query ListQuery) ([]*user.User, int64, error) {
	offset := query.GetOffset()
	if query.Search != "" {
		return h.searchUsers(ctx, query.Search, offset, query.Limit)
	}
	return h.listAllUsers(ctx, offset, query.Limit)
}

// searchUsers 搜索用户
func (h *ListHandler) searchUsers(ctx context.Context, keyword string, offset, limit int) ([]*user.User, int64, error) {
	users, err := h.userQueryRepo.Search(ctx, keyword, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	total, err := h.userQueryRepo.CountBySearch(ctx, keyword)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// listAllUsers 获取所有用户列表
func (h *ListHandler) listAllUsers(ctx context.Context, offset, limit int) ([]*user.User, int64, error) {
	users, err := h.userQueryRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	total, err := h.userQueryRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
