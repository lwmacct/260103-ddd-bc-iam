package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/ctxutil"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/app"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/app/user"
	userDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/user"
)

// UserSettingHandler 用户配置 HTTP Handler
type UserSettingHandler struct {
	setHandler      *user.SetHandler
	batchSetHandler *user.BatchSetHandler
	resetHandler    *user.ResetHandler
	resetAllHandler *user.ResetAllHandler
	getHandler      *user.GetHandler
	listHandler     *user.ListHandler
	listCatHandler  *user.ListCategoriesHandler
}

// NewUserSettingHandler 创建用户配置 Handler
func NewUserSettingHandler(useCases *app.UseCases) *UserSettingHandler {
	return &UserSettingHandler{
		setHandler:      useCases.Set,
		batchSetHandler: useCases.BatchSet,
		resetHandler:    useCases.Reset,
		resetAllHandler: useCases.ResetAll,
		getHandler:      useCases.Get,
		listHandler:     useCases.List,
		listCatHandler:  useCases.ListCategories,
	}
}

// List 获取用户配置列表
//
//	@Summary		配置列表
//	@Description	获取当前用户的配置列表（系统默认值+用户自定义值合并视图）
//	@Tags			User - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			params	query		user.ListQuery	false	"查询参数"
//	@Success		200		{object}	response.DataResponse[[]user.UserSettingDTO]	"配置列表"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/user/settings [get]
func (h *UserSettingHandler) List(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	var query user.ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	query.UserID = uid

	result, err := h.listHandler.Handle(c.Request.Context(), query)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// ListCategories 获取配置分类列表
//
//	@Summary		分类列表
//	@Description	获取配置分类列表
//	@Tags			User - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200		{object}	response.DataResponse[[]user.CategoryDTO]	"分类列表"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/user/settings/categories [get]
func (h *UserSettingHandler) ListCategories(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	query := user.ListCategoriesQuery{
		UserID: uid,
	}

	result, err := h.listCatHandler.Handle(c.Request.Context(), query)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Get 获取单个配置
//
//	@Summary		获取配置
//	@Description	获取指定配置项的值（系统默认值或用户自定义值）
//	@Tags			User - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			key		path		string	true	"配置键名"
//	@Success		200		{object}	response.DataResponse[user.UserSettingDTO]	"配置信息"
//	@Failure		400		{object}	response.ErrorResponse	"无效的键名"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		404		{object}	response.ErrorResponse	"配置不存在"
//	@Router			/api/user/settings/{key} [get]
func (h *UserSettingHandler) Get(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	key := c.Param("key")
	query := user.GetQuery{
		UserID: uid,
		Key:    key,
	}

	result, err := h.getHandler.Handle(c.Request.Context(), query)
	if err != nil {
		if errors.Is(err, userDomain.ErrInvalidSettingKey) {
			response.NotFound(c, "setting")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Set 设置配置
//
//	@Summary		设置配置
//	@Description	设置指定配置项的值（用户自定义覆盖）
//	@Tags			User - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			key		path		string		true	"配置键名"
//	@Param			request	body		SetRequest	true	"设置请求"
//	@Success		200		{object}	response.DataResponse[user.UserSettingDTO]	"设置成功"
//	@Failure		400		{object}	response.ErrorResponse	"参数错误"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/user/settings/{key} [put]
func (h *UserSettingHandler) Set(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	key := c.Param("key")
	var req SetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	cmd := user.SetCommand{
		UserID: uid,
		Key:    key,
		Value:  req.Value,
	}

	result, err := h.setHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		if errors.Is(err, userDomain.ErrInvalidSettingKey) {
			response.BadRequest(c, "invalid setting key", err.Error())
			return
		}
		if errors.Is(err, userDomain.ErrInvalidSettingValue) {
			response.BadRequest(c, "invalid setting value", err.Error())
			return
		}
		if errors.Is(err, userDomain.ErrValidationFailed) {
			response.BadRequest(c, "validation failed", err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// BatchSet 批量设置配置
//
//	@Summary		批量设置
//	@Description	批量设置多个配置项的值
//	@Tags			User - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		BatchSetRequest	true	"批量设置请求"
//	@Success		200		{object}	response.DataResponse[[]user.UserSettingDTO]	"设置成功"
//	@Failure		400		{object}	response.ErrorResponse	"参数错误"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/user/settings/batch [post]
func (h *UserSettingHandler) BatchSet(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	var req BatchSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 转换请求
	items := make([]user.SettingItem, len(req.Settings))
	for i, s := range req.Settings {
		items[i] = user.SettingItem{
			Key:   s.Key,
			Value: s.Value,
		}
	}

	cmd := user.BatchSetCommand{
		UserID:   uid,
		Settings: items,
	}

	result, err := h.batchSetHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		if errors.Is(err, userDomain.ErrInvalidSettingKey) {
			response.BadRequest(c, "invalid setting key", err.Error())
			return
		}
		if errors.Is(err, userDomain.ErrInvalidSettingValue) {
			response.BadRequest(c, "invalid setting value", err.Error())
			return
		}
		if errors.Is(err, userDomain.ErrValidationFailed) {
			response.BadRequest(c, "validation failed", err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Reset 重置配置
//
//	@Summary		重置配置
//	@Description	重置指定配置项（删除用户自定义值，恢复系统默认值）
//	@Tags			User - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			key		path		string	true	"配置键名"
//	@Success		200		{object}	response.MessageResponse	"重置成功"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/user/settings/{key} [delete]
func (h *UserSettingHandler) Reset(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	key := c.Param("key")
	cmd := user.ResetCommand{
		UserID: uid,
		Key:    key,
	}

	if err := h.resetHandler.Handle(c.Request.Context(), cmd); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"key": key, "message": "setting reset to default"})
}

// ResetAll 重置所有配置
//
//	@Summary		重置所有配置
//	@Description	重置所有用户自定义配置（恢复系统默认值）
//	@Tags			User - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200		{object}	response.MessageResponse	"重置成功"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/user/settings/reset-all [post]
func (h *UserSettingHandler) ResetAll(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	cmd := user.ResetAllCommand{
		UserID: uid,
	}

	if err := h.resetAllHandler.Handle(c.Request.Context(), cmd); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "all settings reset to default"})
}
