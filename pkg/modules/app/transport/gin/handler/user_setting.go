package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application/setting"
	settingDomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/ctxutil"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
)

// UserSettingHandler handles user setting operations (DDD+CQRS Use Case Pattern)
type UserSettingHandler struct {
	// Command Handlers
	setHandler      *setting.UserSetHandler
	batchSetHandler *setting.UserBatchSetHandler
	resetHandler    *setting.UserResetHandler
	resetAllHandler *setting.UserResetAllHandler

	// Query Handlers
	getHandler            *setting.UserGetHandler
	listHandler           *setting.UserListHandler
	listSchemaHandler     *setting.UserListSettingsHandler
	listCategoriesHandler *setting.UserListCategoriesHandler
}

// NewUserSettingHandler creates a new UserSettingHandler instance
func NewUserSettingHandler(
	setHandler *setting.UserSetHandler,
	batchSetHandler *setting.UserBatchSetHandler,
	resetHandler *setting.UserResetHandler,
	resetAllHandler *setting.UserResetAllHandler,
	getHandler *setting.UserGetHandler,
	listHandler *setting.UserListHandler,
	listSchemaHandler *setting.UserListSettingsHandler,
	listCategoriesHandler *setting.UserListCategoriesHandler,
) *UserSettingHandler {
	return &UserSettingHandler{
		setHandler:            setHandler,
		batchSetHandler:       batchSetHandler,
		resetHandler:          resetHandler,
		resetAllHandler:       resetAllHandler,
		getHandler:            getHandler,
		listHandler:           listHandler,
		listSchemaHandler:     listSchemaHandler,
		listCategoriesHandler: listCategoriesHandler,
	}
}

// ListUserSettingCategories 获取用户可见的分类列表
//
//	@Summary		配置分类列表
//	@Description	获取包含用户可配置项的分类列表（不含 settings 数据，用于懒加载场景）
//	@Tags			User - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.DataResponse[[]setting.CategoryMetaDTO]	"分类列表"
//	@Failure		401	{object}	response.ErrorResponse								"未授权"
//	@Failure		500	{object}	response.ErrorResponse								"服务器内部错误"
//	@Router			/api/user/settings/categories [get]
func (h *UserSettingHandler) ListUserSettingCategories(c *gin.Context) {
	userID, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, "No user ID found")
		return
	}

	categories, err := h.listCategoriesHandler.Handle(c.Request.Context(), setting.UserListCategoriesQuery{
		UserID: userID,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, categories)
}

// GetUserSettings 获取用户配置（层级结构）
//
//	@Summary		用户配置列表
//	@Description	获取按 Category → Group → Settings 层级组织的配置数据，包含用户自定义值。支持按分类过滤（懒加载）。
//	@Tags			User - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			category	query		string													false	"分类 Key（如 profile），为空返回全量"
//	@Success		200			{object}	response.DataResponse[[]setting.SettingsCategoryDTO]	"配置列表（层级结构）"
//	@Failure		401			{object}	response.ErrorResponse									"未授权"
//	@Failure		404			{object}	response.ErrorResponse									"分类不存在"
//	@Failure		500			{object}	response.ErrorResponse									"服务器内部错误"
//	@Router			/api/user/settings [get]
func (h *UserSettingHandler) GetUserSettings(c *gin.Context) {
	userID, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, "No user ID found")
		return
	}

	categoryKey := c.Query("category")

	// 调用 Schema Handler（返回层级结构）
	schema, err := h.listSchemaHandler.Handle(c.Request.Context(), setting.UserListSettingsQuery{
		UserID:      userID,
		CategoryKey: categoryKey,
	})
	if err != nil {
		// 检查是否为分类不存在错误
		if categoryKey != "" && err.Error() == "category not found: "+categoryKey {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, schema)
}

// GetUserSetting 获取单个用户配置
//
//	@Summary		用户配置详情
//	@Description	根据配置键获取用户配置（合并系统默认值）
//	@Tags			User - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			key	path		string											true	"配置键"	example:"theme"
//	@Success		200	{object}	response.DataResponse[setting.UserSettingDTO]	"配置详情"
//	@Failure		401	{object}	response.ErrorResponse							"未授权"
//	@Failure		404	{object}	response.ErrorResponse							"配置不存在"
//	@Router			/api/user/settings/{key} [get]
func (h *UserSettingHandler) GetUserSetting(c *gin.Context) {
	userID, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, "No user ID found")
		return
	}
	key := c.Param("key")

	settingDTO, err := h.getHandler.Handle(c.Request.Context(), setting.UserGetQuery{
		UserID: userID,
		Key:    key,
	})
	if err != nil {
		response.NotFoundMessage(c, err.Error())
		return
	}

	response.OK(c, settingDTO)
}

// SetUserSettingRequest 设置用户配置请求
type SetUserSettingRequest struct {
	Value any `json:"value"` // JSONB 原生值
}

// SetUserSetting 设置用户配置
//
//	@Summary		设置用户配置
//	@Description	用户设置自定义配置值
//	@Tags			User - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			key		path		string											true	"配置键"	example:"theme"
//	@Param			request	body		SetUserSettingRequest							true	"配置值"
//	@Success		200		{object}	response.DataResponse[setting.UserSettingDTO]	"设置成功"
//	@Failure		400		{object}	response.ErrorResponse							"参数错误"
//	@Failure		401		{object}	response.ErrorResponse							"未授权"
//	@Failure		404		{object}	response.ErrorResponse							"配置不存在"
//	@Failure		500		{object}	response.ErrorResponse							"服务器内部错误"
//	@Router			/api/user/settings/{key} [put]
func (h *UserSettingHandler) SetUserSetting(c *gin.Context) {
	userID, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, "No user ID found")
		return
	}
	key := c.Param("key")

	var req SetUserSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	settingDTO, err := h.setHandler.Handle(c.Request.Context(), setting.UserSetCommand{
		UserID: userID,
		Key:    key,
		Value:  req.Value,
	})
	if err != nil {
		if errors.Is(err, settingDomain.ErrValidationFailed) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, settingDTO)
}

// ResetUserSetting 重置用户配置
//
//	@Summary		重置用户配置
//	@Description	删除用户自定义配置，恢复为系统默认值
//	@Tags			User - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			key	path	string	true	"配置键"	example:"theme"
//	@Success		204	"重置成功"
//	@Failure		401	{object}	response.ErrorResponse	"未授权"
//	@Failure		500	{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/user/settings/{key} [delete]
func (h *UserSettingHandler) ResetUserSetting(c *gin.Context) {
	userID, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, "No user ID found")
		return
	}
	key := c.Param("key")

	err := h.resetHandler.Handle(c.Request.Context(), setting.UserResetCommand{
		UserID: userID,
		Key:    key,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.NoContent(c)
}

// BatchSetUserSettingsRequest 批量设置用户配置请求
type BatchSetUserSettingsRequest struct {
	Settings []struct {
		Key   string `json:"key" binding:"required"`
		Value any    `json:"value"` // JSONB 原生值
	} `json:"settings" binding:"required,min=1"`
}

// BatchSetUserSettings 批量设置用户配置
//
//	@Summary		批量设置配置
//	@Description	用户批量设置多个自定义配置值
//	@Tags			User - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		BatchSetUserSettingsRequest	true	"配置列表"
//	@Success		200		{object}	response.MessageResponse	"批量设置成功"
//	@Failure		400		{object}	response.ErrorResponse		"参数错误"
//	@Failure		401		{object}	response.ErrorResponse		"未授权"
//	@Failure		500		{object}	response.ErrorResponse		"服务器内部错误"
//	@Router			/api/user/settings/batch [post]
func (h *UserSettingHandler) BatchSetUserSettings(c *gin.Context) {
	userID, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, "No user ID found")
		return
	}

	var req BatchSetUserSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 转换为 Command
	items := make([]setting.SettingItemCommand, len(req.Settings))
	for i, s := range req.Settings {
		items[i] = setting.SettingItemCommand{
			Key:   s.Key,
			Value: s.Value,
		}
	}

	err := h.batchSetHandler.Handle(c.Request.Context(), setting.UserBatchSetCommand{
		UserID:   userID,
		Settings: items,
	})
	if err != nil {
		if errors.Is(err, settingDomain.ErrValidationFailed) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil)
}
