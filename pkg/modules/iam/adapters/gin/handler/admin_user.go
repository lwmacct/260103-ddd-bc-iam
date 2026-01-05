package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/user"
	userDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/user"
)

// ListUsersQuery 用户列表查询参数
type ListUsersQuery struct {
	response.PaginationQueryDTO

	// Search 搜索关键词（用户名或邮箱）
	Search string `form:"search" json:"search" binding:"omitempty"`
}

// ToQuery 转换为 Application 层 Query 对象
func (q *ListUsersQuery) ToQuery() user.ListQuery {
	return user.ListQuery{
		Page:   q.GetPage(),
		Limit:  q.GetLimit(),
		Search: q.Search,
	}
}

// AdminUserHandler handles admin user management operations
type AdminUserHandler struct {
	createUserHandler      *user.CreateHandler
	updateUserHandler      *user.UpdateHandler
	deleteUserHandler      *user.DeleteHandler
	assignRolesHandler     *user.AssignRolesHandler
	batchCreateUserHandler *user.BatchCreateHandler
	getUserHandler         *user.GetHandler
	listUsersHandler       *user.ListHandler
}

// NewAdminUserHandler creates a new AdminUserHandler instance
func NewAdminUserHandler(
	createUserHandler *user.CreateHandler,
	updateUserHandler *user.UpdateHandler,
	deleteUserHandler *user.DeleteHandler,
	assignRolesHandler *user.AssignRolesHandler,
	batchCreateUserHandler *user.BatchCreateHandler,
	getUserHandler *user.GetHandler,
	listUsersHandler *user.ListHandler,
) *AdminUserHandler {
	return &AdminUserHandler{
		createUserHandler:      createUserHandler,
		updateUserHandler:      updateUserHandler,
		deleteUserHandler:      deleteUserHandler,
		assignRolesHandler:     assignRolesHandler,
		batchCreateUserHandler: batchCreateUserHandler,
		getUserHandler:         getUserHandler,
		listUsersHandler:       listUsersHandler,
	}
}

// CreateUser creates a new user (admin only)
//
//	@Summary		创建用户
//	@Description	管理员创建新用户账号，可同时分配角色
//	@Tags			admin-user
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		user.CreateDTO									true	"用户信息"
//	@Success		201		{object}	response.DataResponse[user.UserWithRolesDTO]	"用户创建成功"
//	@Failure		400		{object}	response.ErrorResponse							"参数错误或用户名/邮箱已存在"
//	@Failure		401		{object}	response.ErrorResponse							"未授权"
//	@Failure		403		{object}	response.ErrorResponse							"权限不足"
//	@Failure		500		{object}	response.ErrorResponse							"服务器内部错误"
//	@Router			/api/admin/users [post]
func (h *AdminUserHandler) CreateUser(c *gin.Context) {
	var dto user.CreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.createUserHandler.Handle(c.Request.Context(), user.CreateCommand(dto))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	createdUser, err := h.getUserHandler.Handle(c.Request.Context(), user.GetQuery{
		UserID:    result.UserID,
		WithRoles: true,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, createdUser)
}

// ListUsers lists all users with pagination (admin only)
//
//	@Summary		用户列表
//	@Description	分页获取所有用户列表（包含角色信息）
//	@Tags			admin-user
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			params	query		handler.ListUsersQuery							false	"查询参数"
//	@Success		200		{object}	response.PagedResponse[user.UserWithRolesDTO]	"用户列表"
//	@Failure		401		{object}	response.ErrorResponse							"未授权"
//	@Failure		403		{object}	response.ErrorResponse							"权限不足"
//	@Failure		500		{object}	response.ErrorResponse							"服务器内部错误"
//	@Router			/api/admin/users [get]
func (h *AdminUserHandler) ListUsers(c *gin.Context) {
	var q ListUsersQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.listUsersHandler.Handle(c.Request.Context(), q.ToQuery())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), q.GetPage(), q.GetLimit())
	response.List(c, result.Users, meta)
}

// GetUser gets a user by ID (admin only)
//
//	@Summary		用户详情
//	@Description	根据用户ID获取用户详细信息（包含角色信息）
//	@Tags			admin-user
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int												true	"用户ID"	minimum(1)
//	@Success		200	{object}	response.DataResponse[user.UserWithRolesDTO]	"用户详情"
//	@Failure		400	{object}	response.ErrorResponse							"无效的用户ID"
//	@Failure		401	{object}	response.ErrorResponse							"未授权"
//	@Failure		403	{object}	response.ErrorResponse							"权限不足"
//	@Failure		404	{object}	response.ErrorResponse							"用户不存在"
//	@Router			/api/admin/users/{id} [get]
func (h *AdminUserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, userDomain.ErrInvalidUserID.Error())
		return
	}

	userResp, err := h.getUserHandler.Handle(c.Request.Context(), user.GetQuery{
		UserID:    uint(id),
		WithRoles: true,
	})
	if err != nil {
		if errors.Is(err, userDomain.ErrUserNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, userResp)
}

// UpdateUser updates a user (admin only)
//
//	@Summary		更新用户
//	@Description	管理员更新用户的基本信息和状态
//	@Tags			admin-user
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int												true	"用户ID"	minimum(1)
//	@Param			request	body		user.UpdateDTO									true	"更新信息"
//	@Success		200		{object}	response.DataResponse[user.UserWithRolesDTO]	"用户更新成功"
//	@Failure		400		{object}	response.ErrorResponse							"无效的用户ID或参数错误"
//	@Failure		401		{object}	response.ErrorResponse							"未授权"
//	@Failure		403		{object}	response.ErrorResponse							"权限不足"
//	@Failure		404		{object}	response.ErrorResponse							"用户不存在"
//	@Failure		500		{object}	response.ErrorResponse							"服务器内部错误"
//	@Router			/api/admin/users/{id} [put]
func (h *AdminUserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, userDomain.ErrInvalidUserID.Error())
		return
	}

	var dto user.UpdateDTO
	if err = c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	_, err = h.updateUserHandler.Handle(c.Request.Context(), user.UpdateCommand{
		UserID:    uint(id),
		Username:  dto.Username,
		Email:     dto.Email,
		RealName:  dto.RealName,
		Nickname:  dto.Nickname,
		Phone:     dto.Phone,
		Signature: dto.Signature,
		Avatar:    dto.Avatar,
		Bio:       dto.Bio,
		Status:    dto.Status,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	updatedUser, err := h.getUserHandler.Handle(c.Request.Context(), user.GetQuery{
		UserID:    uint(id),
		WithRoles: true,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, updatedUser)
}

// DeleteUser deletes a user (admin only)
//
//	@Summary		删除用户
//	@Description	管理员删除指定用户（物理删除或软删除）
//	@Tags			admin-user
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int							true	"用户ID"	minimum(1)
//	@Success		200	{object}	response.MessageResponse	"用户删除成功"
//	@Failure		400	{object}	response.ErrorResponse		"无效的用户ID"
//	@Failure		401	{object}	response.ErrorResponse		"未授权"
//	@Failure		403	{object}	response.ErrorResponse		"权限不足"
//	@Failure		404	{object}	response.ErrorResponse		"用户不存在"
//	@Failure		500	{object}	response.ErrorResponse		"服务器内部错误"
//	@Router			/api/admin/users/{id} [delete]
func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, userDomain.ErrInvalidUserID.Error())
		return
	}

	if err := h.deleteUserHandler.Handle(c.Request.Context(), user.DeleteCommand{
		UserID: uint(id),
	}); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil)
}

// AssignRoles assigns roles to a user (admin only)
//
//	@Summary		分配角色
//	@Description	管理员为指定用户分配角色（会覆盖现有角色）
//	@Tags			admin-user
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int												true	"用户ID"	minimum(1)
//	@Param			request	body		user.AssignRolesDTO								true	"角色ID列表"
//	@Success		200		{object}	response.DataResponse[user.UserWithRolesDTO]	"角色分配成功"
//	@Failure		400		{object}	response.ErrorResponse							"无效的用户ID或参数错误"
//	@Failure		401		{object}	response.ErrorResponse							"未授权"
//	@Failure		403		{object}	response.ErrorResponse							"权限不足"
//	@Failure		404		{object}	response.ErrorResponse							"用户不存在"
//	@Failure		500		{object}	response.ErrorResponse							"服务器内部错误"
//	@Router			/api/admin/users/{id}/roles [put]
func (h *AdminUserHandler) AssignRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, userDomain.ErrInvalidUserID.Error())
		return
	}

	var req user.AssignRolesDTO
	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err = h.assignRolesHandler.Handle(c.Request.Context(), user.AssignRolesCommand{
		UserID:  uint(id),
		RoleIDs: req.RoleIDs,
	}); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 获取更新后的用户信息（包含角色）
	updatedUser, err := h.getUserHandler.Handle(c.Request.Context(), user.GetQuery{
		UserID:    uint(id),
		WithRoles: true,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, updatedUser)
}

// BatchCreateUsers creates multiple users at once (admin only)
//
//	@Summary		批量创建用户
//	@Description	管理员从 CSV 等来源批量创建用户，支持部分失败（单个失败不影响其他用户）
//	@Tags			admin-user
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		user.BatchCreateDTO									true	"用户列表（最多 100 个）"
//	@Success		200		{object}	response.DataResponse[user.BatchCreateResultDTO]	"批量创建结果"
//	@Failure		400		{object}	response.ErrorResponse								"参数错误"
//	@Failure		401		{object}	response.ErrorResponse								"未授权"
//	@Failure		403		{object}	response.ErrorResponse								"权限不足"
//	@Failure		500		{object}	response.ErrorResponse								"服务器内部错误"
//	@Router			/api/admin/users/batch [post]
func (h *AdminUserHandler) BatchCreateUsers(c *gin.Context) {
	var dto user.BatchCreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 构建 Command
	users := make([]user.BatchItemDTO, len(dto.Users))
	copy(users, dto.Users)

	result, err := h.batchCreateUserHandler.Handle(c.Request.Context(), user.BatchCreateCommand{
		Users: users,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 构建响应
	resp := user.BatchCreateResultDTO{
		Total:   result.Total,
		Success: result.Success,
		Failed:  result.Failed,
		Errors:  make([]user.BatchCreateErrorDTO, len(result.Errors)),
	}
	copy(resp.Errors, result.Errors)

	response.OK(c, resp)
}
