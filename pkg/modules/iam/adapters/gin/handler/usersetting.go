package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/ctxutil"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/usersetting"
	authDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/auth"
	usersettingDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/usersetting"
)

// UserSettingHandler handles user settings operations
type UserSettingHandler struct {
	getHandler    *usersetting.GetHandler
	listHandler   *usersetting.ListHandler
	updateHandler *usersetting.UpdateHandler
	deleteHandler *usersetting.DeleteHandler
}

// NewUserSettingHandler creates a new user setting handler
func NewUserSettingHandler(
	getHandler *usersetting.GetHandler,
	listHandler *usersetting.ListHandler,
	updateHandler *usersetting.UpdateHandler,
	deleteHandler *usersetting.DeleteHandler,
) *UserSettingHandler {
	return &UserSettingHandler{
		getHandler:    getHandler,
		listHandler:   listHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
	}
}

// ListSettings 获取用户设置列表
//
//	@Summary		获取用户设置列表
//	@Description	获取当前用户的设置项，包含系统默认值和用户自定义值
//	@Tags			User - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			category	query		string	false	"分类过滤（可选）"	example(general)
//	@Success		200		{object}	response.DataResponse[usersetting.UserSettingListDTO]	"设置列表"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/user/settings [get]
func (h *UserSettingHandler) ListSettings(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, authDomain.ErrUserIDNotFound.Error())
		return
	}

	category := c.Query("category")
	result, err := h.listHandler.Handle(c.Request.Context(), uid, category)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// GetSetting 获取单个用户设置
//
//	@Summary		获取用户设置
//	@Description	获取指定设置项的值（系统默认值或用户自定义值）
//	@Tags			User - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			key		path		string	true	"设置键名"	example(theme.dark_mode)
//	@Success		200		{object}	response.DataResponse[usersetting.UserSettingDTO]	"设置信息"
//	@Failure		400		{object}	response.ErrorResponse	"无效的键名"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		404		{object}	response.ErrorResponse	"设置不存在"
//	@Router			/api/user/settings/{key} [get]
func (h *UserSettingHandler) GetSetting(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, authDomain.ErrUserIDNotFound.Error())
		return
	}

	key := c.Param("key")
	result, err := h.getHandler.Handle(c.Request.Context(), uid, key)
	if err != nil {
		if errors.Is(err, usersettingDomain.ErrInvalidSettingKey) {
			response.BadRequest(c, "无效的设置键", err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// UpdateSetting 更新用户设置
//
//	@Summary		更新用户设置
//	@Description	更新指定设置项的值（用户自定义覆盖）
//	@Tags			User - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			key		path		string	true	"设置键名"	example(theme.dark_mode)
//	@Param			request	body		usersetting.UpdateDTO	true	"更新请求"
//	@Success		200		{object}	response.MessageResponse	"更新成功"
//	@Failure		400		{object}	response.ErrorResponse	"参数错误"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/user/settings/{key} [put]
func (h *UserSettingHandler) UpdateSetting(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, authDomain.ErrUserIDNotFound.Error())
		return
	}

	key := c.Param("key")
	var req usersetting.UpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	cmd := usersetting.UpdateCommand{
		UserID: uid,
		Key:    key,
		Value:  req.Value,
	}

	if err := h.updateHandler.Handle(c.Request.Context(), cmd); err != nil {
		if errors.Is(err, usersettingDomain.ErrInvalidSettingKey) {
			response.BadRequest(c, "无效的设置键", err.Error())
			return
		}
		if errors.Is(err, usersettingDomain.ErrInvalidSettingValue) {
			response.BadRequest(c, "无效的设置值", err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"key": key, "value": req.Value, "message": "设置更新成功"})
}

// DeleteSetting 删除用户设置
//
//	@Summary		删除用户设置
//	@Description	删除指定的用户自定义设置（恢复系统默认值）
//	@Tags			User - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			key		path		string	true	"设置键名"	example(theme.dark_mode)
//	@Success		200		{object}	response.MessageResponse	"删除成功"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/user/settings/{key} [delete]
func (h *UserSettingHandler) DeleteSetting(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, authDomain.ErrUserIDNotFound.Error())
		return
	}

	key := c.Param("key")
	cmd := usersetting.DeleteCommand{
		UserID: uid,
		Key:    key,
	}

	if err := h.deleteHandler.Handle(c.Request.Context(), cmd); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"key": key, "message": "设置删除成功，已恢复系统默认值"})
}
