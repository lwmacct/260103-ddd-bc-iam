package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/org"
	orgDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/org"
	userDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/user"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/response"
)

// ListMembersQuery 成员列表查询参数
type ListMembersQuery struct {
	response.PaginationQueryDTO
}

// ToQuery 转换为 Application 层 Query 对象
func (q *ListMembersQuery) ToQuery(orgID uint) org.ListMembersQuery {
	return org.ListMembersQuery{
		OrgID:  orgID,
		Offset: q.GetOffset(),
		Limit:  q.GetLimit(),
	}
}

// OrgMemberHandler 组织成员管理 Handler
type OrgMemberHandler struct {
	// Command Handlers
	addMemberHandler        *org.MemberAddHandler
	removeMemberHandler     *org.MemberRemoveHandler
	updateMemberRoleHandler *org.MemberUpdateRoleHandler

	// Query Handlers
	listMembersHandler *org.MemberListHandler
}

// NewOrgMemberHandler 创建组织成员管理 Handler
func NewOrgMemberHandler(
	addMemberHandler *org.MemberAddHandler,
	removeMemberHandler *org.MemberRemoveHandler,
	updateMemberRoleHandler *org.MemberUpdateRoleHandler,
	listMembersHandler *org.MemberListHandler,
) *OrgMemberHandler {
	return &OrgMemberHandler{
		addMemberHandler:        addMemberHandler,
		removeMemberHandler:     removeMemberHandler,
		updateMemberRoleHandler: updateMemberRoleHandler,
		listMembersHandler:      listMembersHandler,
	}
}

// List 获取组织成员列表
//
//	@Summary		成员列表
//	@Description	分页获取组织成员列表
//	@Tags			org-member
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int												true	"组织ID"	minimum(1)
//	@Param			params	query		handler.ListMembersQuery						false	"查询参数"
//	@Success		200		{object}	response.PagedResponse[org.MemberDTO]	"成员列表"
//	@Failure		401		{object}	response.ErrorResponse							"未授权"
//	@Failure		403		{object}	response.ErrorResponse							"权限不足"
//	@Failure		500		{object}	response.ErrorResponse							"服务器内部错误"
//	@Router			/api/org/{org_id}/members [get]
func (h *OrgMemberHandler) List(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("org_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, orgDomain.ErrInvalidOrgID.Error())
		return
	}

	var q ListMembersQuery
	if bindErr := c.ShouldBindQuery(&q); bindErr != nil {
		response.ValidationError(c, bindErr.Error())
		return
	}

	result, err := h.listMembersHandler.Handle(c.Request.Context(), q.ToQuery(uint(orgID)))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), q.GetPage(), q.GetLimit())
	response.List(c, result.Items, meta)
}

// Add 添加组织成员
//
//	@Summary		添加成员
//	@Description	添加用户到组织
//	@Tags			org-member
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int												true	"组织ID"	minimum(1)
//	@Param			request	body		org.AddMemberDTO						true	"成员信息"
//	@Success		201		{object}	response.DataResponse[org.MemberDTO]	"成员添加成功"
//	@Failure		400		{object}	response.ErrorResponse							"参数错误或成员已存在"
//	@Failure		401		{object}	response.ErrorResponse							"未授权"
//	@Failure		403		{object}	response.ErrorResponse							"权限不足"
//	@Failure		500		{object}	response.ErrorResponse							"服务器内部错误"
//	@Router			/api/org/{org_id}/members [post]
func (h *OrgMemberHandler) Add(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("org_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, orgDomain.ErrInvalidOrgID.Error())
		return
	}

	var req org.AddMemberDTO
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		response.ValidationError(c, bindErr.Error())
		return
	}

	result, err := h.addMemberHandler.Handle(c.Request.Context(), org.AddMemberCommand{
		OrgID:  uint(orgID),
		UserID: req.UserID,
		Role:   req.Role,
	})
	if err != nil {
		// 处理业务错误
		if errors.Is(err, orgDomain.ErrMemberAlreadyExists) {
			response.Conflict(c, err.Error())
			return
		}
		if errors.Is(err, orgDomain.ErrOrgNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		if errors.Is(err, userDomain.ErrUserNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		if errors.Is(err, orgDomain.ErrInvalidMemberRole) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, result)
}

// Remove 移除组织成员
//
//	@Summary		移除成员
//	@Description	从组织中移除成员
//	@Tags			org-member
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int							true	"组织ID"	minimum(1)
//	@Param			user_id	path		int							true	"用户ID"	minimum(1)
//	@Success		200		{object}	response.MessageResponse	"成员移除成功"
//	@Failure		400		{object}	response.ErrorResponse		"无效的ID或无法移除最后的所有者"
//	@Failure		401		{object}	response.ErrorResponse		"未授权"
//	@Failure		403		{object}	response.ErrorResponse		"权限不足"
//	@Failure		404		{object}	response.ErrorResponse		"成员不存在"
//	@Failure		500		{object}	response.ErrorResponse		"服务器内部错误"
//	@Router			/api/org/{org_id}/members/{user_id} [delete]
func (h *OrgMemberHandler) Remove(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("org_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, orgDomain.ErrInvalidOrgID.Error())
		return
	}

	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, userDomain.ErrInvalidUserID.Error())
		return
	}

	if err = h.removeMemberHandler.Handle(c.Request.Context(), org.RemoveMemberCommand{
		OrgID:  uint(orgID),
		UserID: uint(userID),
	}); err != nil {
		if errors.Is(err, orgDomain.ErrCannotRemoveLastOwner) {
			response.BadRequest(c, err.Error())
			return
		}
		if errors.Is(err, orgDomain.ErrMemberNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil)
}

// UpdateRole 更新成员角色
//
//	@Summary		更新成员角色
//	@Description	更新组织成员的角色
//	@Tags			org-member
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int												true	"组织ID"	minimum(1)
//	@Param			user_id	path		int												true	"用户ID"	minimum(1)
//	@Param			request	body		org.UpdateMemberRoleDTO				true	"角色信息"
//	@Success		200		{object}	response.DataResponse[org.MemberDTO]	"角色更新成功"
//	@Failure		400		{object}	response.ErrorResponse							"参数错误或无法降级最后的所有者"
//	@Failure		401		{object}	response.ErrorResponse							"未授权"
//	@Failure		403		{object}	response.ErrorResponse							"权限不足"
//	@Failure		404		{object}	response.ErrorResponse							"成员不存在"
//	@Failure		500		{object}	response.ErrorResponse							"服务器内部错误"
//	@Router			/api/org/{org_id}/members/{user_id}/role [put]
func (h *OrgMemberHandler) UpdateRole(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("org_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, orgDomain.ErrInvalidOrgID.Error())
		return
	}

	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, userDomain.ErrInvalidUserID.Error())
		return
	}

	var req org.UpdateMemberRoleDTO
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		response.ValidationError(c, bindErr.Error())
		return
	}

	if err = h.updateMemberRoleHandler.Handle(c.Request.Context(), org.UpdateMemberRoleCommand{
		OrgID:  uint(orgID),
		UserID: uint(userID),
		Role:   req.Role,
	}); err != nil {
		if errors.Is(err, orgDomain.ErrCannotDemoteLastOwner) {
			response.BadRequest(c, err.Error())
			return
		}
		if errors.Is(err, orgDomain.ErrMemberNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil)
}
