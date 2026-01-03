package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/application/lead"
	leadDomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/lead"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/ctxutil"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
)

// ListLeadsQuery 线索列表查询参数。
type ListLeadsQuery struct {
	response.PaginationQueryDTO

	Status  *string `form:"status"`
	OwnerID *uint   `form:"owner_id"`
}

// ToQuery 转换为应用层查询对象。
func (q *ListLeadsQuery) ToQuery() lead.ListQuery {
	var status *leadDomain.Status
	if q.Status != nil {
		s := leadDomain.Status(*q.Status)
		status = &s
	}
	return lead.ListQuery{
		Status:  status,
		OwnerID: q.OwnerID,
		Offset:  q.GetOffset(),
		Limit:   q.GetLimit(),
	}
}

// LeadHandler 线索管理 Handler。
type LeadHandler struct {
	createHandler  *lead.CreateHandler
	updateHandler  *lead.UpdateHandler
	deleteHandler  *lead.DeleteHandler
	contactHandler *lead.ContactHandler
	qualifyHandler *lead.QualifyHandler
	convertHandler *lead.ConvertHandler
	loseHandler    *lead.LoseHandler
	getHandler     *lead.GetHandler
	listHandler    *lead.ListHandler
}

// NewLeadHandler 创建线索管理 Handler。
func NewLeadHandler(
	createHandler *lead.CreateHandler,
	updateHandler *lead.UpdateHandler,
	deleteHandler *lead.DeleteHandler,
	contactHandler *lead.ContactHandler,
	qualifyHandler *lead.QualifyHandler,
	convertHandler *lead.ConvertHandler,
	loseHandler *lead.LoseHandler,
	getHandler *lead.GetHandler,
	listHandler *lead.ListHandler,
) *LeadHandler {
	return &LeadHandler{
		createHandler:  createHandler,
		updateHandler:  updateHandler,
		deleteHandler:  deleteHandler,
		contactHandler: contactHandler,
		qualifyHandler: qualifyHandler,
		convertHandler: convertHandler,
		loseHandler:    loseHandler,
		getHandler:     getHandler,
		listHandler:    listHandler,
	}
}

// Create 创建线索
//
//	@Summary		创建线索
//	@Description	创建新线索
//	@Tags			CRM - Leads
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		lead.CreateLeadDTO					true	"线索信息"
//	@Success		201		{object}	response.DataResponse[lead.LeadDTO]	"创建成功"
//	@Failure		400		{object}	response.ErrorResponse				"参数错误"
//	@Failure		401		{object}	response.ErrorResponse				"未授权"
//	@Router			/api/crm/leads [post]
func (h *LeadHandler) Create(c *gin.Context) {
	userID, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, "No user ID found")
		return
	}

	var req lead.CreateLeadDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.createHandler.Handle(c.Request.Context(), lead.CreateCommand{
		Title:       req.Title,
		ContactID:   req.ContactID,
		CompanyName: req.CompanyName,
		Source:      req.Source,
		Score:       req.Score,
		OwnerID:     userID,
		Notes:       req.Notes,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, result)
}

// List 线索列表
//
//	@Summary		线索列表
//	@Description	分页获取线索列表
//	@Tags			CRM - Leads
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			params	query		ListLeadsQuery						false	"查询参数"
//	@Success		200		{object}	response.PagedResponse[lead.LeadDTO]	"线索列表"
//	@Failure		401		{object}	response.ErrorResponse				"未授权"
//	@Router			/api/crm/leads [get]
func (h *LeadHandler) List(c *gin.Context) {
	var query ListLeadsQuery
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

// Get 线索详情
//
//	@Summary		线索详情
//	@Description	获取线索详细信息
//	@Tags			CRM - Leads
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int									true	"线索ID"
//	@Success		200	{object}	response.DataResponse[lead.LeadDTO]	"线索详情"
//	@Failure		401	{object}	response.ErrorResponse				"未授权"
//	@Failure		404	{object}	response.ErrorResponse				"线索不存在"
//	@Router			/api/crm/leads/{id} [get]
func (h *LeadHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的线索ID")
		return
	}

	result, err := h.getHandler.Handle(c.Request.Context(), lead.GetQuery{ID: uint(id)})
	if err != nil {
		if errors.Is(err, leadDomain.ErrLeadNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Update 更新线索
//
//	@Summary		更新线索
//	@Description	更新线索信息
//	@Tags			CRM - Leads
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int									true	"线索ID"
//	@Param			request	body		lead.UpdateLeadDTO					true	"更新信息"
//	@Success		200		{object}	response.DataResponse[lead.LeadDTO]	"更新成功"
//	@Failure		400		{object}	response.ErrorResponse				"参数错误"
//	@Failure		401		{object}	response.ErrorResponse				"未授权"
//	@Failure		404		{object}	response.ErrorResponse				"线索不存在"
//	@Router			/api/crm/leads/{id} [put]
func (h *LeadHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的线索ID")
		return
	}

	var req lead.UpdateLeadDTO
	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.updateHandler.Handle(c.Request.Context(), lead.UpdateCommand{
		ID:          uint(id),
		Title:       req.Title,
		ContactID:   req.ContactID,
		CompanyName: req.CompanyName,
		Source:      req.Source,
		Score:       req.Score,
		Notes:       req.Notes,
	})
	if err != nil {
		if errors.Is(err, leadDomain.ErrLeadNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		if errors.Is(err, leadDomain.ErrAlreadyClosed) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Delete 删除线索
//
//	@Summary		删除线索
//	@Description	删除线索
//	@Tags			CRM - Leads
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int						true	"线索ID"
//	@Success		200	{object}	response.MessageResponse	"删除成功"
//	@Failure		401	{object}	response.ErrorResponse	"未授权"
//	@Failure		404	{object}	response.ErrorResponse	"线索不存在"
//	@Router			/api/crm/leads/{id} [delete]
func (h *LeadHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的线索ID")
		return
	}

	if err := h.deleteHandler.Handle(c.Request.Context(), lead.DeleteCommand{ID: uint(id)}); err != nil {
		if errors.Is(err, leadDomain.ErrLeadNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil)
}

// Contact 转换到已联系状态
//
//	@Summary		标记为已联系
//	@Description	将线索状态转换为已联系（需要线索处于新建状态）
//	@Tags			CRM - Leads
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int									true	"线索ID"
//	@Success		200	{object}	response.DataResponse[lead.LeadDTO]	"转换成功"
//	@Failure		400	{object}	response.ErrorResponse				"状态转换失败"
//	@Failure		401	{object}	response.ErrorResponse				"未授权"
//	@Failure		404	{object}	response.ErrorResponse				"线索不存在"
//	@Router			/api/crm/leads/{id}/contact [post]
func (h *LeadHandler) Contact(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的线索ID")
		return
	}

	result, err := h.contactHandler.Handle(c.Request.Context(), lead.ContactCommand{ID: uint(id)})
	if err != nil {
		if errors.Is(err, leadDomain.ErrLeadNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		if errors.Is(err, leadDomain.ErrCannotContact) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Qualify 转换到已确认状态
//
//	@Summary		标记为已确认
//	@Description	将线索状态转换为已确认（需要线索处于已联系状态）
//	@Tags			CRM - Leads
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int									true	"线索ID"
//	@Success		200	{object}	response.DataResponse[lead.LeadDTO]	"转换成功"
//	@Failure		400	{object}	response.ErrorResponse				"状态转换失败"
//	@Failure		401	{object}	response.ErrorResponse				"未授权"
//	@Failure		404	{object}	response.ErrorResponse				"线索不存在"
//	@Router			/api/crm/leads/{id}/qualify [post]
func (h *LeadHandler) Qualify(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的线索ID")
		return
	}

	result, err := h.qualifyHandler.Handle(c.Request.Context(), lead.QualifyCommand{ID: uint(id)})
	if err != nil {
		if errors.Is(err, leadDomain.ErrLeadNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		if errors.Is(err, leadDomain.ErrCannotQualify) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Convert 转化为商机
//
//	@Summary		转化为商机
//	@Description	将线索转化为商机（需要线索处于已确认状态）
//	@Tags			CRM - Leads
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int									true	"线索ID"
//	@Success		200	{object}	response.DataResponse[lead.LeadDTO]	"转化成功"
//	@Failure		400	{object}	response.ErrorResponse				"状态转换失败"
//	@Failure		401	{object}	response.ErrorResponse				"未授权"
//	@Failure		404	{object}	response.ErrorResponse				"线索不存在"
//	@Router			/api/crm/leads/{id}/convert [post]
func (h *LeadHandler) Convert(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的线索ID")
		return
	}

	result, err := h.convertHandler.Handle(c.Request.Context(), lead.ConvertCommand{ID: uint(id)})
	if err != nil {
		if errors.Is(err, leadDomain.ErrLeadNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		if errors.Is(err, leadDomain.ErrCannotConvert) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Lose 标记为丢失
//
//	@Summary		标记为丢失
//	@Description	将线索标记为丢失（需要线索处于新建、已联系或已确认状态）
//	@Tags			CRM - Leads
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int									true	"线索ID"
//	@Success		200	{object}	response.DataResponse[lead.LeadDTO]	"标记成功"
//	@Failure		400	{object}	response.ErrorResponse				"状态转换失败"
//	@Failure		401	{object}	response.ErrorResponse				"未授权"
//	@Failure		404	{object}	response.ErrorResponse				"线索不存在"
//	@Router			/api/crm/leads/{id}/lose [post]
func (h *LeadHandler) Lose(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的线索ID")
		return
	}

	result, err := h.loseHandler.Handle(c.Request.Context(), lead.LoseCommand{ID: uint(id)})
	if err != nil {
		if errors.Is(err, leadDomain.ErrLeadNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		if errors.Is(err, leadDomain.ErrCannotLose) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}
