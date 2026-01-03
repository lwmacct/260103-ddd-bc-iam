package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/application/role"
	roleDomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/role"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
)

// ListRolesQuery 角色列表查询参数
type ListRolesQuery struct {
	response.PaginationQueryDTO
}

// ToQuery 转换为 Application 层 Query 对象
func (q *ListRolesQuery) ToQuery() role.ListQuery {
	return role.ListQuery{
		Page:  q.GetPage(),
		Limit: q.GetLimit(),
	}
}

// RoleHandler handles role management operations (DDD+CQRS Use Case Pattern)
type RoleHandler struct {
	// Command Handlers
	createRoleHandler     *role.CreateHandler
	updateRoleHandler     *role.UpdateHandler
	deleteRoleHandler     *role.DeleteHandler
	setPermissionsHandler *role.SetPermissionsHandler

	// Query Handlers
	getRoleHandler   *role.GetHandler
	listRolesHandler *role.ListHandler
}

// NewRoleHandler creates a new RoleHandler instance
func NewRoleHandler(
	createRoleHandler *role.CreateHandler,
	updateRoleHandler *role.UpdateHandler,
	deleteRoleHandler *role.DeleteHandler,
	setPermissionsHandler *role.SetPermissionsHandler,
	getRoleHandler *role.GetHandler,
	listRolesHandler *role.ListHandler,
) *RoleHandler {
	return &RoleHandler{
		createRoleHandler:     createRoleHandler,
		updateRoleHandler:     updateRoleHandler,
		deleteRoleHandler:     deleteRoleHandler,
		setPermissionsHandler: setPermissionsHandler,
		getRoleHandler:        getRoleHandler,
		listRolesHandler:      listRolesHandler,
	}
}

// CreateRole creates a new role
//
//	@Summary		创建角色
//	@Description	管理员创建新的系统角色
//	@Tags			Admin - Role Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		role.CreateDTO								true	"角色信息"
//	@Success		201		{object}	response.DataResponse[role.CreateResultDTO]	"角色创建成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误或角色名已存在"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		403		{object}	response.ErrorResponse						"权限不足"
//	@Failure		500		{object}	response.ErrorResponse						"服务器内部错误"
//	@Router			/api/admin/roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req role.CreateDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	result, err := h.createRoleHandler.Handle(c.Request.Context(), role.CreateCommand(req))

	if err != nil {
		if errors.Is(err, roleDomain.ErrRoleNameExists) {
			response.Conflict(c, err.Error())
			return
		}
		if errors.Is(err, roleDomain.ErrInvalidRoleName) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	// 转换为 DTO 响应
	resp := role.CreateResultDTO{
		RoleID:      result.RoleID,
		Name:        result.Name,
		DisplayName: result.DisplayName,
	}
	response.Created(c, resp)
}

// ListRoles lists all roles
//
//	@Summary		角色列表
//	@Description	分页获取所有系统角色
//	@Tags			Admin - Role Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			params	query		handler.ListRolesQuery					false	"查询参数"
//	@Success		200		{object}	response.PagedResponse[role.RoleDTO]	"角色列表"
//	@Failure		401		{object}	response.ErrorResponse					"未授权"
//	@Failure		403		{object}	response.ErrorResponse					"权限不足"
//	@Failure		500		{object}	response.ErrorResponse					"服务器内部错误"
//	@Router			/api/admin/roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	var q ListRolesQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.listRolesHandler.Handle(c.Request.Context(), q.ToQuery())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), q.GetPage(), q.GetLimit())
	response.List(c, result.Roles, meta)
}

// GetRole gets a role by ID
//
//	@Summary		角色详情
//	@Description	根据角色ID获取角色详细信息（包含权限列表）
//	@Tags			Admin - Role Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int									true	"角色ID"	minimum(1)
//	@Success		200	{object}	response.DataResponse[role.RoleDTO]	"角色详情"
//	@Failure		400	{object}	response.ErrorResponse				"无效的角色ID"
//	@Failure		401	{object}	response.ErrorResponse				"未授权"
//	@Failure		403	{object}	response.ErrorResponse				"权限不足"
//	@Failure		404	{object}	response.ErrorResponse				"角色不存在"
//	@Router			/api/admin/roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, roleDomain.ErrInvalidRoleID.Error())
		return
	}

	// 调用 Use Case Handler
	result, err := h.getRoleHandler.Handle(c.Request.Context(), role.GetQuery{
		RoleID: uint(id),
	})

	if err != nil {
		if errors.Is(err, roleDomain.ErrRoleNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// UpdateRole updates a role
//
//	@Summary		更新角色
//	@Description	管理员更新角色的显示名称和描述
//	@Tags			Admin - Role Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int									true	"角色ID"	minimum(1)
//	@Param			request	body		role.UpdateDTO						true	"更新信息"
//	@Success		200		{object}	response.DataResponse[role.RoleDTO]	"角色更新成功"
//	@Failure		400		{object}	response.ErrorResponse				"无效的角色ID或参数错误"
//	@Failure		401		{object}	response.ErrorResponse				"未授权"
//	@Failure		403		{object}	response.ErrorResponse				"权限不足"
//	@Failure		404		{object}	response.ErrorResponse				"角色不存在"
//	@Failure		500		{object}	response.ErrorResponse				"服务器内部错误"
//	@Router			/api/admin/roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, roleDomain.ErrInvalidRoleID.Error())
		return
	}

	var req role.UpdateDTO

	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	result, err := h.updateRoleHandler.Handle(c.Request.Context(), role.UpdateCommand{
		RoleID:      uint(id),
		DisplayName: req.DisplayName,
		Description: req.Description,
	})

	if err != nil {
		if errors.Is(err, roleDomain.ErrRoleNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		if errors.Is(err, roleDomain.ErrCannotModifySystemRole) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// DeleteRole deletes a role
//
//	@Summary		删除角色
//	@Description	管理员删除指定角色（如果角色被用户使用，可能会失败）
//	@Tags			Admin - Role Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int							true	"角色ID"	minimum(1)
//	@Success		200	{object}	response.MessageResponse	"角色删除成功"
//	@Failure		400	{object}	response.ErrorResponse		"无效的角色ID"
//	@Failure		401	{object}	response.ErrorResponse		"未授权"
//	@Failure		403	{object}	response.ErrorResponse		"权限不足"
//	@Failure		404	{object}	response.ErrorResponse		"角色不存在"
//	@Failure		500	{object}	response.ErrorResponse		"服务器内部错误或角色被使用中"
//	@Router			/api/admin/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, roleDomain.ErrInvalidRoleID.Error())
		return
	}

	// 调用 Use Case Handler
	err = h.deleteRoleHandler.Handle(c.Request.Context(), role.DeleteCommand{
		RoleID: uint(id),
	})

	if err != nil {
		if errors.Is(err, roleDomain.ErrRoleNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		if errors.Is(err, roleDomain.ErrCannotDeleteSystemRole) {
			response.BadRequest(c, err.Error())
			return
		}
		if errors.Is(err, roleDomain.ErrRoleHasUsers) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil)
}

// SetPermissions sets permissions for a role
//
//	@Summary		设置权限
//	@Description	管理员为指定角色设置权限（会覆盖现有权限）。新 RBAC 模型使用 Operation + Resource Pattern。
//	@Tags			Admin - Role Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"角色ID"	minimum(1)
//	@Param			request	body		role.SetPermissionsDTO		true	"权限模式列表"
//	@Success		200		{object}	response.MessageResponse	"权限设置成功"
//	@Failure		400		{object}	response.ErrorResponse		"无效的角色ID或参数错误"
//	@Failure		401		{object}	response.ErrorResponse		"未授权"
//	@Failure		403		{object}	response.ErrorResponse		"权限不足"
//	@Failure		404		{object}	response.ErrorResponse		"角色不存在"
//	@Failure		500		{object}	response.ErrorResponse		"服务器内部错误"
//	@Router			/api/admin/roles/{id}/permissions [put]
func (h *RoleHandler) SetPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, roleDomain.ErrInvalidRoleID.Error())
		return
	}

	var req role.SetPermissionsDTO

	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 转换 DTO 到 Domain Permission
	permissions := make([]roleDomain.Permission, len(req.Permissions))
	for i, p := range req.Permissions {
		resPattern := p.ResourcePattern
		if resPattern == "" {
			resPattern = "*"
		}
		permissions[i] = roleDomain.Permission{
			OperationPattern: p.OperationPattern,
			ResourcePattern:  resPattern,
		}
	}

	// 调用 Use Case Handler
	err = h.setPermissionsHandler.Handle(c.Request.Context(), role.SetPermissionsCommand{
		RoleID:      uint(id),
		Permissions: permissions,
	})

	if err != nil {
		if errors.Is(err, roleDomain.ErrRoleNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil)
}

// ListPermissions 已移除
// 新 RBAC 模型中，权限不再是独立实体，而是 Operation + Resource Pattern。
// 如需获取可用操作列表，请使用 /api/system/operations 端点。
