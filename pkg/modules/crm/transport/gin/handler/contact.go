package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/application/contact"
	contactDomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/contact"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/ctxutil"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
)

// ListContactsQuery 联系人列表查询参数。
type ListContactsQuery struct {
	response.PaginationQueryDTO

	CompanyID *uint `form:"company_id"`
	OwnerID   *uint `form:"owner_id"`
}

// ToQuery 转换为应用层查询对象。
func (q *ListContactsQuery) ToQuery() contact.ListQuery {
	return contact.ListQuery{
		CompanyID: q.CompanyID,
		OwnerID:   q.OwnerID,
		Offset:    q.GetOffset(),
		Limit:     q.GetLimit(),
	}
}

// ContactHandler 联系人管理 Handler。
type ContactHandler struct {
	createHandler *contact.CreateHandler
	updateHandler *contact.UpdateHandler
	deleteHandler *contact.DeleteHandler
	getHandler    *contact.GetHandler
	listHandler   *contact.ListHandler
}

// NewContactHandler 创建联系人管理 Handler。
func NewContactHandler(
	createHandler *contact.CreateHandler,
	updateHandler *contact.UpdateHandler,
	deleteHandler *contact.DeleteHandler,
	getHandler *contact.GetHandler,
	listHandler *contact.ListHandler,
) *ContactHandler {
	return &ContactHandler{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

// Create 创建联系人
//
//	@Summary		创建联系人
//	@Description	创建新联系人
//	@Tags			CRM - Contacts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		contact.CreateContactDTO				true	"联系人信息"
//	@Success		201		{object}	response.DataResponse[contact.ContactDTO]	"创建成功"
//	@Failure		400		{object}	response.ErrorResponse					"参数错误"
//	@Failure		401		{object}	response.ErrorResponse					"未授权"
//	@Failure		409		{object}	response.ErrorResponse					"邮箱已存在"
//	@Router			/api/crm/contacts [post]
func (h *ContactHandler) Create(c *gin.Context) {
	userID, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, "No user ID found")
		return
	}

	var req contact.CreateContactDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.createHandler.Handle(c.Request.Context(), contact.CreateCommand{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Title:     req.Title,
		CompanyID: req.CompanyID,
		OwnerID:   userID,
	})
	if err != nil {
		if errors.Is(err, contactDomain.ErrEmailAlreadyExists) {
			response.Conflict(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, result)
}

// List 联系人列表
//
//	@Summary		联系人列表
//	@Description	分页获取联系人列表
//	@Tags			CRM - Contacts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			params	query		ListContactsQuery							false	"查询参数"
//	@Success		200		{object}	response.PagedResponse[contact.ContactDTO]	"联系人列表"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Router			/api/crm/contacts [get]
func (h *ContactHandler) List(c *gin.Context) {
	var query ListContactsQuery
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

// Get 联系人详情
//
//	@Summary		联系人详情
//	@Description	获取联系人详细信息
//	@Tags			CRM - Contacts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int											true	"联系人ID"
//	@Success		200	{object}	response.DataResponse[contact.ContactDTO]	"联系人详情"
//	@Failure		401	{object}	response.ErrorResponse						"未授权"
//	@Failure		404	{object}	response.ErrorResponse						"联系人不存在"
//	@Router			/api/crm/contacts/{id} [get]
func (h *ContactHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的联系人ID")
		return
	}

	result, err := h.getHandler.Handle(c.Request.Context(), contact.GetQuery{ID: uint(id)})
	if err != nil {
		if errors.Is(err, contactDomain.ErrContactNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Update 更新联系人
//
//	@Summary		更新联系人
//	@Description	更新联系人信息
//	@Tags			CRM - Contacts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int											true	"联系人ID"
//	@Param			request	body		contact.UpdateContactDTO					true	"更新信息"
//	@Success		200		{object}	response.DataResponse[contact.ContactDTO]	"更新成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		404		{object}	response.ErrorResponse						"联系人不存在"
//	@Failure		409		{object}	response.ErrorResponse						"邮箱已存在"
//	@Router			/api/crm/contacts/{id} [put]
func (h *ContactHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的联系人ID")
		return
	}

	var req contact.UpdateContactDTO
	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.updateHandler.Handle(c.Request.Context(), contact.UpdateCommand{
		ID:        uint(id),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Title:     req.Title,
		CompanyID: req.CompanyID,
	})
	if err != nil {
		if errors.Is(err, contactDomain.ErrContactNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		if errors.Is(err, contactDomain.ErrEmailAlreadyExists) {
			response.Conflict(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Delete 删除联系人
//
//	@Summary		删除联系人
//	@Description	删除联系人
//	@Tags			CRM - Contacts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int						true	"联系人ID"
//	@Success		200	{object}	response.MessageResponse	"删除成功"
//	@Failure		401	{object}	response.ErrorResponse	"未授权"
//	@Failure		404	{object}	response.ErrorResponse	"联系人不存在"
//	@Router			/api/crm/contacts/{id} [delete]
func (h *ContactHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的联系人ID")
		return
	}

	if err := h.deleteHandler.Handle(c.Request.Context(), contact.DeleteCommand{ID: uint(id)}); err != nil {
		if errors.Is(err, contactDomain.ErrContactNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil)
}
