package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/application/company"
	companyDomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/company"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/ctxutil"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
)

// ListCompaniesQuery 公司列表查询参数。
type ListCompaniesQuery struct {
	response.PaginationQueryDTO

	Industry *string `form:"industry"`
	OwnerID  *uint   `form:"owner_id"`
}

// ToQuery 转换为应用层查询对象。
func (q *ListCompaniesQuery) ToQuery() company.ListQuery {
	return company.ListQuery{
		Industry: q.Industry,
		OwnerID:  q.OwnerID,
		Offset:   q.GetOffset(),
		Limit:    q.GetLimit(),
	}
}

// CompanyHandler 公司管理 Handler。
type CompanyHandler struct {
	createHandler *company.CreateHandler
	updateHandler *company.UpdateHandler
	deleteHandler *company.DeleteHandler
	getHandler    *company.GetHandler
	listHandler   *company.ListHandler
}

// NewCompanyHandler 创建公司管理 Handler。
func NewCompanyHandler(
	createHandler *company.CreateHandler,
	updateHandler *company.UpdateHandler,
	deleteHandler *company.DeleteHandler,
	getHandler *company.GetHandler,
	listHandler *company.ListHandler,
) *CompanyHandler {
	return &CompanyHandler{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

// Create 创建公司
//
//	@Summary		创建公司
//	@Description	创建新公司
//	@Tags			CRM - Companies
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		company.CreateCompanyDTO				true	"公司信息"
//	@Success		201		{object}	response.DataResponse[company.CompanyDTO]	"创建成功"
//	@Failure		400		{object}	response.ErrorResponse					"参数错误"
//	@Failure		401		{object}	response.ErrorResponse					"未授权"
//	@Failure		409		{object}	response.ErrorResponse					"公司名称已存在"
//	@Router			/api/crm/companies [post]
func (h *CompanyHandler) Create(c *gin.Context) {
	userID, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, "No user ID found")
		return
	}

	var req company.CreateCompanyDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.createHandler.Handle(c.Request.Context(), company.CreateCommand{
		Name:     req.Name,
		Industry: req.Industry,
		Size:     req.Size,
		Website:  req.Website,
		Address:  req.Address,
		OwnerID:  userID,
	})
	if err != nil {
		if errors.Is(err, companyDomain.ErrCompanyNameExists) {
			response.Conflict(c, err.Error())
			return
		}
		if errors.Is(err, companyDomain.ErrInvalidSize) {
			response.ValidationError(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, result)
}

// List 公司列表
//
//	@Summary		公司列表
//	@Description	分页获取公司列表
//	@Tags			CRM - Companies
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			params	query		ListCompaniesQuery							false	"查询参数"
//	@Success		200		{object}	response.PagedResponse[company.CompanyDTO]	"公司列表"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Router			/api/crm/companies [get]
func (h *CompanyHandler) List(c *gin.Context) {
	var query ListCompaniesQuery
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

// Get 公司详情
//
//	@Summary		公司详情
//	@Description	获取公司详细信息
//	@Tags			CRM - Companies
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int											true	"公司ID"
//	@Success		200	{object}	response.DataResponse[company.CompanyDTO]	"公司详情"
//	@Failure		401	{object}	response.ErrorResponse						"未授权"
//	@Failure		404	{object}	response.ErrorResponse						"公司不存在"
//	@Router			/api/crm/companies/{id} [get]
func (h *CompanyHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的公司ID")
		return
	}

	result, err := h.getHandler.Handle(c.Request.Context(), company.GetQuery{ID: uint(id)})
	if err != nil {
		if errors.Is(err, companyDomain.ErrCompanyNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Update 更新公司
//
//	@Summary		更新公司
//	@Description	更新公司信息
//	@Tags			CRM - Companies
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int											true	"公司ID"
//	@Param			request	body		company.UpdateCompanyDTO					true	"更新信息"
//	@Success		200		{object}	response.DataResponse[company.CompanyDTO]	"更新成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		404		{object}	response.ErrorResponse						"公司不存在"
//	@Failure		409		{object}	response.ErrorResponse						"公司名称已存在"
//	@Router			/api/crm/companies/{id} [put]
func (h *CompanyHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的公司ID")
		return
	}

	var req company.UpdateCompanyDTO
	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.updateHandler.Handle(c.Request.Context(), company.UpdateCommand{
		ID:       uint(id),
		Name:     req.Name,
		Industry: req.Industry,
		Size:     req.Size,
		Website:  req.Website,
		Address:  req.Address,
	})
	if err != nil {
		if errors.Is(err, companyDomain.ErrCompanyNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		if errors.Is(err, companyDomain.ErrCompanyNameExists) {
			response.Conflict(c, err.Error())
			return
		}
		if errors.Is(err, companyDomain.ErrInvalidSize) {
			response.ValidationError(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Delete 删除公司
//
//	@Summary		删除公司
//	@Description	删除公司
//	@Tags			CRM - Companies
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int						true	"公司ID"
//	@Success		200	{object}	response.MessageResponse	"删除成功"
//	@Failure		401	{object}	response.ErrorResponse	"未授权"
//	@Failure		404	{object}	response.ErrorResponse	"公司不存在"
//	@Router			/api/crm/companies/{id} [delete]
func (h *CompanyHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的公司ID")
		return
	}

	if err := h.deleteHandler.Handle(c.Request.Context(), company.DeleteCommand{ID: uint(id)}); err != nil {
		if errors.Is(err, companyDomain.ErrCompanyNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil)
}
