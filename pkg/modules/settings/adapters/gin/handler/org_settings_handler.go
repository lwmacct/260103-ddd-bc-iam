package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/ctxutil"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/app/org"
	orgDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/org"
)

// OrgSettingHandler 组织配置 HTTP Handler
type OrgSettingHandler struct {
	setHandler   *org.SetHandler
	resetHandler *org.ResetHandler
	getHandler   *org.GetHandler
	listHandler  *org.ListHandler
}

// NewOrgSettingHandler 创建组织配置 Handler
func NewOrgSettingHandler(
	setHandler *org.SetHandler,
	resetHandler *org.ResetHandler,
	getHandler *org.GetHandler,
	listHandler *org.ListHandler,
) *OrgSettingHandler {
	return &OrgSettingHandler{
		setHandler:   setHandler,
		resetHandler: resetHandler,
		getHandler:   getHandler,
		listHandler:  listHandler,
	}
}

// List 获取组织配置列表
//
//	@Summary		组织配置列表
//	@Description	获取当前组织的配置列表（系统默认值+组织自定义值合并视图）
//	@Tags			settings-org
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id		path		int		true	"组织ID"
//	@Param			params	query		org.ListQuery	false	"查询参数"
//	@Success		200		{object}	response.DataResponse[[]org.OrgSettingDTO]	"配置列表"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		403		{object}	response.ErrorResponse	"权限不足"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/org/{org_id}/settings [get]
func (h *OrgSettingHandler) List(c *gin.Context) {
	orgID, ok := ctxutil.Get[uint](c, ctxutil.OrgID)
	if !ok {
		response.Unauthorized(c, "organization context not found")
		return
	}

	var query org.ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	query.OrgID = orgID

	result, err := h.listHandler.Handle(c.Request.Context(), query)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Get 获取单个组织配置
//
//	@Summary		获取组织配置
//	@Description	获取指定配置项的值（系统默认值或组织自定义值）
//	@Tags			settings-org
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int		true	"组织ID"
//	@Param			key		path		string	true	"配置键名"
//	@Success		200		{object}	response.DataResponse[org.OrgSettingDTO]	"配置信息"
//	@Failure		400		{object}	response.ErrorResponse	"无效的键名"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		403		{object}	response.ErrorResponse	"权限不足"
//	@Failure		404		{object}	response.ErrorResponse	"配置不存在"
//	@Router			/api/org/{org_id}/settings/{key} [get]
func (h *OrgSettingHandler) Get(c *gin.Context) {
	orgID, ok := ctxutil.Get[uint](c, ctxutil.OrgID)
	if !ok {
		response.Unauthorized(c, "organization context not found")
		return
	}

	key := c.Param("key")
	query := org.GetQuery{
		OrgID: orgID,
		Key:   key,
	}

	result, err := h.getHandler.Handle(c.Request.Context(), query)
	if err != nil {
		if errors.Is(err, orgDomain.ErrInvalidSettingKey) {
			response.NotFound(c, "setting")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Set 设置组织配置
//
//	@Summary		设置组织配置
//	@Description	设置指定配置项的值（组织自定义覆盖）
//	@Tags			settings-org
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int			true	"组织ID"
//	@Param			key		path		string		true	"配置键名"
//	@Param			request	body		SetRequest	true	"设置请求"
//	@Success		200		{object}	response.DataResponse[org.OrgSettingDTO]	"设置成功"
//	@Failure		400		{object}	response.ErrorResponse	"参数错误"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		403		{object}	response.ErrorResponse	"权限不足"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/org/{org_id}/settings/{key} [put]
func (h *OrgSettingHandler) Set(c *gin.Context) {
	orgID, ok := ctxutil.Get[uint](c, ctxutil.OrgID)
	if !ok {
		response.Unauthorized(c, "organization context not found")
		return
	}

	key := c.Param("key")
	var req SetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	cmd := org.SetCommand{
		OrgID: orgID,
		Key:   key,
		Value: req.Value,
	}

	result, err := h.setHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		if errors.Is(err, orgDomain.ErrInvalidSettingKey) {
			response.BadRequest(c, "invalid setting key", err.Error())
			return
		}
		if errors.Is(err, orgDomain.ErrInvalidSettingValue) {
			response.BadRequest(c, "invalid setting value", err.Error())
			return
		}
		if errors.Is(err, orgDomain.ErrValidationFailed) {
			response.BadRequest(c, "validation failed", err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Reset 重置组织配置
//
//	@Summary		重置组织配置
//	@Description	重置指定配置项（删除组织自定义值，恢复系统默认值）
//	@Tags			settings-org
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int		true	"组织ID"
//	@Param			key		path		string	true	"配置键名"
//	@Success		200		{object}	response.MessageResponse	"重置成功"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		403		{object}	response.ErrorResponse	"权限不足"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/org/{org_id}/settings/{key} [delete]
func (h *OrgSettingHandler) Reset(c *gin.Context) {
	orgID, ok := ctxutil.Get[uint](c, ctxutil.OrgID)
	if !ok {
		response.Unauthorized(c, "organization context not found")
		return
	}

	key := c.Param("key")
	cmd := org.ResetCommand{
		OrgID: orgID,
		Key:   key,
	}

	if err := h.resetHandler.Handle(c.Request.Context(), cmd); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"key": key, "message": "setting reset to default"})
}
