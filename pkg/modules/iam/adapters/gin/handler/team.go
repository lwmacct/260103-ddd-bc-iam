package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/org"
	orgDomain "github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/org"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/ctxutil"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/response"
)

// ListTeamsQuery 团队列表查询参数
type ListTeamsQuery struct {
	response.PaginationQueryDTO
}

// ToQuery 转换为 Application 层 Query 对象
func (q *ListTeamsQuery) ToQuery(orgID uint) org.ListTeamsQuery {
	return org.ListTeamsQuery{
		OrgID:  orgID,
		Offset: q.GetOffset(),
		Limit:  q.GetLimit(),
	}
}

// TeamHandler 团队管理 Handler
type TeamHandler struct {
	// Command Handlers
	createHandler *org.TeamCreateHandler
	updateHandler *org.TeamUpdateHandler
	deleteHandler *org.TeamDeleteHandler

	// Query Handlers
	getHandler  *org.TeamGetHandler
	listHandler *org.TeamListHandler
}

// NewTeamHandler 创建团队管理 Handler
func NewTeamHandler(useCases *app.TeamUseCases) *TeamHandler {
	return &TeamHandler{
		createHandler: useCases.Create,
		updateHandler: useCases.Update,
		deleteHandler: useCases.Delete,
		getHandler:    useCases.Get,
		listHandler:   useCases.List,
	}
}

// Create 创建团队
//
//	@Summary		创建团队
//	@Description	在组织内创建新团队
//	@Tags			org-team
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int												true	"组织ID"	minimum(1)
//	@Param			request	body		org.CreateTeamDTO						true	"团队信息"
//	@Success		201		{object}	response.DataResponse[org.TeamDTO]		"团队创建成功"
//	@Failure		400		{object}	response.ErrorResponse							"参数错误或团队名已存在"
//	@Failure		401		{object}	response.ErrorResponse							"未授权"
//	@Failure		403		{object}	response.ErrorResponse							"权限不足"
//	@Failure		500		{object}	response.ErrorResponse							"服务器内部错误"
//	@Router			/api/org/{org_id}/teams [post]
func (h *TeamHandler) Create(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("org_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, orgDomain.ErrInvalidOrgID.Error())
		return
	}

	var req org.CreateTeamDTO
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		response.ValidationError(c, bindErr.Error())
		return
	}

	// 获取当前用户 ID，创建者自动成为团队负责人
	userID, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, response.MsgAuthenticationRequired)
		return
	}

	result, err := h.createHandler.Handle(c.Request.Context(), org.CreateTeamCommand{
		OrgID:       uint(orgID),
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Avatar:      req.Avatar,
		LeadUserID:  userID,
	})
	if err != nil {
		if errors.Is(err, orgDomain.ErrTeamAlreadyExists) {
			response.Conflict(c, err.Error())
			return
		}
		if errors.Is(err, orgDomain.ErrTeamNameAlreadyExists) {
			response.Conflict(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, result)
}

// List 获取团队列表
//
//	@Summary		团队列表
//	@Description	分页获取组织内的团队列表
//	@Tags			org-team
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int												true	"组织ID"	minimum(1)
//	@Param			params	query		handler.ListTeamsQuery							false	"查询参数"
//	@Success		200		{object}	response.PagedResponse[org.TeamDTO]	"团队列表"
//	@Failure		401		{object}	response.ErrorResponse							"未授权"
//	@Failure		403		{object}	response.ErrorResponse							"权限不足"
//	@Failure		500		{object}	response.ErrorResponse							"服务器内部错误"
//	@Router			/api/org/{org_id}/teams [get]
func (h *TeamHandler) List(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("org_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, orgDomain.ErrInvalidOrgID.Error())
		return
	}

	var q ListTeamsQuery
	if bindErr := c.ShouldBindQuery(&q); bindErr != nil {
		response.ValidationError(c, bindErr.Error())
		return
	}

	result, err := h.listHandler.Handle(c.Request.Context(), q.ToQuery(uint(orgID)))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), q.GetPage(), q.GetLimit())
	response.List(c, result.Items, meta)
}

// Get 获取团队详情
//
//	@Summary		团队详情
//	@Description	获取团队详情
//	@Tags			org-team
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int											true	"组织ID"	minimum(1)
//	@Param			team_id	path		int											true	"团队ID"	minimum(1)
//	@Success		200		{object}	response.DataResponse[org.TeamDTO]	"团队详情"
//	@Failure		400		{object}	response.ErrorResponse						"无效的ID"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		403		{object}	response.ErrorResponse						"权限不足"
//	@Failure		404		{object}	response.ErrorResponse						"团队不存在"
//	@Router			/api/org/{org_id}/teams/{team_id} [get]
func (h *TeamHandler) Get(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("org_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, orgDomain.ErrInvalidOrgID.Error())
		return
	}

	teamID, err := strconv.ParseUint(c.Param("team_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, orgDomain.ErrInvalidTeamID.Error())
		return
	}

	result, err := h.getHandler.Handle(c.Request.Context(), org.GetTeamQuery{
		OrgID:  uint(orgID),
		TeamID: uint(teamID),
	})
	if err != nil {
		if errors.Is(err, orgDomain.ErrTeamNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Update 更新团队
//
//	@Summary		更新团队
//	@Description	更新团队信息
//	@Tags			org-team
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int											true	"组织ID"	minimum(1)
//	@Param			team_id	path		int											true	"团队ID"	minimum(1)
//	@Param			request	body		org.UpdateTeamDTO					true	"更新信息"
//	@Success		200		{object}	response.DataResponse[org.TeamDTO]	"团队更新成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		403		{object}	response.ErrorResponse						"权限不足"
//	@Failure		404		{object}	response.ErrorResponse						"团队不存在"
//	@Failure		500		{object}	response.ErrorResponse						"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id} [put]
//

func (h *TeamHandler) Update(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("org_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, orgDomain.ErrInvalidOrgID.Error())
		return
	}

	teamID, err := strconv.ParseUint(c.Param("team_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, orgDomain.ErrInvalidTeamID.Error())
		return
	}

	var req org.UpdateTeamDTO
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		response.ValidationError(c, bindErr.Error())
		return
	}

	result, err := h.updateHandler.Handle(c.Request.Context(), org.UpdateTeamCommand{
		OrgID:       uint(orgID),
		TeamID:      uint(teamID),
		DisplayName: req.DisplayName,
		Description: req.Description,
		Avatar:      req.Avatar,
	})
	if err != nil {
		if errors.Is(err, orgDomain.ErrTeamNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Delete 删除团队
//
//	@Summary		删除团队
//	@Description	软删除团队
//	@Tags			org-team
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int							true	"组织ID"	minimum(1)
//	@Param			team_id	path		int							true	"团队ID"	minimum(1)
//	@Success		200		{object}	response.MessageResponse	"团队删除成功"
//	@Failure		400		{object}	response.ErrorResponse		"无效的ID"
//	@Failure		401		{object}	response.ErrorResponse		"未授权"
//	@Failure		403		{object}	response.ErrorResponse		"权限不足"
//	@Failure		404		{object}	response.ErrorResponse		"团队不存在"
//	@Failure		500		{object}	response.ErrorResponse		"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id} [delete]
func (h *TeamHandler) Delete(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("org_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, orgDomain.ErrInvalidOrgID.Error())
		return
	}

	teamID, err := strconv.ParseUint(c.Param("team_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, orgDomain.ErrInvalidTeamID.Error())
		return
	}

	if err = h.deleteHandler.Handle(c.Request.Context(), org.DeleteTeamCommand{
		OrgID:  uint(orgID),
		TeamID: uint(teamID),
	}); err != nil {
		if errors.Is(err, orgDomain.ErrTeamHasMembers) {
			response.BadRequest(c, err.Error())
			return
		}
		if errors.Is(err, orgDomain.ErrTeamNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil)
}
