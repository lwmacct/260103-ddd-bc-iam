package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/application/opportunity"
	opportunityDomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/opportunity"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/ctxutil"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
)

// ListOpportunitiesQuery 商机列表查询参数。
type ListOpportunitiesQuery struct {
	response.PaginationQueryDTO

	Stage   *string `form:"stage"`
	OwnerID *uint   `form:"owner_id"`
}

// ToQuery 转换为应用层查询对象。
func (q *ListOpportunitiesQuery) ToQuery() opportunity.ListQuery {
	var stage *opportunityDomain.Stage
	if q.Stage != nil {
		s := opportunityDomain.Stage(*q.Stage)
		stage = &s
	}
	return opportunity.ListQuery{
		Stage:   stage,
		OwnerID: q.OwnerID,
		Offset:  q.GetOffset(),
		Limit:   q.GetLimit(),
	}
}

// OpportunityHandler 商机管理 Handler。
type OpportunityHandler struct {
	createHandler    *opportunity.CreateHandler
	updateHandler    *opportunity.UpdateHandler
	deleteHandler    *opportunity.DeleteHandler
	advanceHandler   *opportunity.AdvanceHandler
	closeWonHandler  *opportunity.CloseWonHandler
	closeLostHandler *opportunity.CloseLostHandler
	getHandler       *opportunity.GetHandler
	listHandler      *opportunity.ListHandler
}

// NewOpportunityHandler 创建商机管理 Handler。
func NewOpportunityHandler(
	createHandler *opportunity.CreateHandler,
	updateHandler *opportunity.UpdateHandler,
	deleteHandler *opportunity.DeleteHandler,
	advanceHandler *opportunity.AdvanceHandler,
	closeWonHandler *opportunity.CloseWonHandler,
	closeLostHandler *opportunity.CloseLostHandler,
	getHandler *opportunity.GetHandler,
	listHandler *opportunity.ListHandler,
) *OpportunityHandler {
	return &OpportunityHandler{
		createHandler:    createHandler,
		updateHandler:    updateHandler,
		deleteHandler:    deleteHandler,
		advanceHandler:   advanceHandler,
		closeWonHandler:  closeWonHandler,
		closeLostHandler: closeLostHandler,
		getHandler:       getHandler,
		listHandler:      listHandler,
	}
}

// Create 创建商机
//
//	@Summary		创建商机
//	@Description	创建新商机
//	@Tags			CRM - Opportunities
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		opportunity.CreateOpportunityDTO			true	"商机信息"
//	@Success		201		{object}	response.DataResponse[opportunity.OpportunityDTO]	"创建成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Router			/api/crm/opportunities [post]
func (h *OpportunityHandler) Create(c *gin.Context) {
	userID, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, "No user ID found")
		return
	}

	var req opportunity.CreateOpportunityDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.createHandler.Handle(c.Request.Context(), opportunity.CreateCommand{
		Name:          req.Name,
		ContactID:     req.ContactID,
		CompanyID:     req.CompanyID,
		LeadID:        req.LeadID,
		Amount:        req.Amount,
		Probability:   req.Probability,
		ExpectedClose: req.ExpectedClose,
		OwnerID:       userID,
		Notes:         req.Notes,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, result)
}

// List 商机列表
//
//	@Summary		商机列表
//	@Description	分页获取商机列表
//	@Tags			CRM - Opportunities
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			params	query		ListOpportunitiesQuery							false	"查询参数"
//	@Success		200		{object}	response.PagedResponse[opportunity.OpportunityDTO]	"商机列表"
//	@Failure		401		{object}	response.ErrorResponse							"未授权"
//	@Router			/api/crm/opportunities [get]
func (h *OpportunityHandler) List(c *gin.Context) {
	var query ListOpportunitiesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.listHandler.Handle(c.Request.Context(), query.ToQuery())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), query.GetPage(), query.GetLimit())
	response.List(c, result.Items, meta)
}

// Get 商机详情
//
//	@Summary		商机详情
//	@Description	获取商机详细信息
//	@Tags			CRM - Opportunities
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int											true	"商机ID"
//	@Success		200	{object}	response.DataResponse[opportunity.OpportunityDTO]	"商机详情"
//	@Failure		401	{object}	response.ErrorResponse						"未授权"
//	@Failure		404	{object}	response.ErrorResponse						"商机不存在"
//	@Router			/api/crm/opportunities/{id} [get]
func (h *OpportunityHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的商机ID")
		return
	}

	result, err := h.getHandler.Handle(c.Request.Context(), opportunity.GetQuery{ID: uint(id)})
	if err != nil {
		if errors.Is(err, opportunityDomain.ErrOpportunityNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Update 更新商机
//
//	@Summary		更新商机
//	@Description	更新商机信息
//	@Tags			CRM - Opportunities
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int											true	"商机ID"
//	@Param			request	body		opportunity.UpdateOpportunityDTO			true	"更新信息"
//	@Success		200		{object}	response.DataResponse[opportunity.OpportunityDTO]	"更新成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		404		{object}	response.ErrorResponse						"商机不存在"
//	@Router			/api/crm/opportunities/{id} [put]
func (h *OpportunityHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的商机ID")
		return
	}

	var req opportunity.UpdateOpportunityDTO
	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.updateHandler.Handle(c.Request.Context(), opportunity.UpdateCommand{
		ID:            uint(id),
		Name:          req.Name,
		ContactID:     req.ContactID,
		CompanyID:     req.CompanyID,
		Amount:        req.Amount,
		Probability:   req.Probability,
		ExpectedClose: req.ExpectedClose,
		Notes:         req.Notes,
	})
	if err != nil {
		if errors.Is(err, opportunityDomain.ErrOpportunityNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		if errors.Is(err, opportunityDomain.ErrAlreadyClosed) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Delete 删除商机
//
//	@Summary		删除商机
//	@Description	删除商机
//	@Tags			CRM - Opportunities
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int						true	"商机ID"
//	@Success		200	{object}	response.MessageResponse	"删除成功"
//	@Failure		401	{object}	response.ErrorResponse	"未授权"
//	@Failure		404	{object}	response.ErrorResponse	"商机不存在"
//	@Router			/api/crm/opportunities/{id} [delete]
func (h *OpportunityHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的商机ID")
		return
	}

	if err := h.deleteHandler.Handle(c.Request.Context(), opportunity.DeleteCommand{ID: uint(id)}); err != nil {
		if errors.Is(err, opportunityDomain.ErrOpportunityNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil)
}

// Advance 推进商机阶段
//
//	@Summary		推进商机阶段
//	@Description	将商机推进到下一阶段（prospecting → proposal → negotiation）
//	@Tags			CRM - Opportunities
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int											true	"商机ID"
//	@Param			request	body		opportunity.AdvanceStageDTO					true	"目标阶段"
//	@Success		200		{object}	response.DataResponse[opportunity.OpportunityDTO]	"推进成功"
//	@Failure		400		{object}	response.ErrorResponse						"阶段转换失败"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		404		{object}	response.ErrorResponse						"商机不存在"
//	@Router			/api/crm/opportunities/{id}/advance [post]
func (h *OpportunityHandler) Advance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的商机ID")
		return
	}

	var req opportunity.AdvanceStageDTO
	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.advanceHandler.Handle(c.Request.Context(), opportunity.AdvanceCommand{
		ID:    uint(id),
		Stage: opportunityDomain.Stage(req.Stage),
	})
	if err != nil {
		if errors.Is(err, opportunityDomain.ErrOpportunityNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		if errors.Is(err, opportunityDomain.ErrInvalidStageTransition) ||
			errors.Is(err, opportunityDomain.ErrAlreadyClosed) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// CloseWon 成交
//
//	@Summary		成交
//	@Description	将商机标记为成交（需要商机处于谈判阶段）
//	@Tags			CRM - Opportunities
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int											true	"商机ID"
//	@Success		200	{object}	response.DataResponse[opportunity.OpportunityDTO]	"成交成功"
//	@Failure		400	{object}	response.ErrorResponse						"状态转换失败"
//	@Failure		401	{object}	response.ErrorResponse						"未授权"
//	@Failure		404	{object}	response.ErrorResponse						"商机不存在"
//	@Router			/api/crm/opportunities/{id}/close-won [post]
func (h *OpportunityHandler) CloseWon(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的商机ID")
		return
	}

	result, err := h.closeWonHandler.Handle(c.Request.Context(), opportunity.CloseWonCommand{ID: uint(id)})
	if err != nil {
		if errors.Is(err, opportunityDomain.ErrOpportunityNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		if errors.Is(err, opportunityDomain.ErrCannotClose) ||
			errors.Is(err, opportunityDomain.ErrAlreadyClosed) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// CloseLost 丢单
//
//	@Summary		丢单
//	@Description	将商机标记为丢单（需要商机处于谈判阶段）
//	@Tags			CRM - Opportunities
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int											true	"商机ID"
//	@Success		200	{object}	response.DataResponse[opportunity.OpportunityDTO]	"标记成功"
//	@Failure		400	{object}	response.ErrorResponse						"状态转换失败"
//	@Failure		401	{object}	response.ErrorResponse						"未授权"
//	@Failure		404	{object}	response.ErrorResponse						"商机不存在"
//	@Router			/api/crm/opportunities/{id}/close-lost [post]
func (h *OpportunityHandler) CloseLost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的商机ID")
		return
	}

	result, err := h.closeLostHandler.Handle(c.Request.Context(), opportunity.CloseLostCommand{ID: uint(id)})
	if err != nil {
		if errors.Is(err, opportunityDomain.ErrOpportunityNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		if errors.Is(err, opportunityDomain.ErrCannotClose) ||
			errors.Is(err, opportunityDomain.ErrAlreadyClosed) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}
