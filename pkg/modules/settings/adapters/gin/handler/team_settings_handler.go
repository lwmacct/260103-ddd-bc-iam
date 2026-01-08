package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/app/team"
	teamDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/team"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/ctxutil"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/response"
)

// TeamSettingHandler 团队配置 HTTP Handler
type TeamSettingHandler struct {
	setHandler   *team.SetHandler
	resetHandler *team.ResetHandler
	getHandler   *team.GetHandler
	listHandler  *team.ListHandler
}

// NewTeamSettingHandler 创建团队配置 Handler
func NewTeamSettingHandler(useCases *team.TeamUseCases) *TeamSettingHandler {
	return &TeamSettingHandler{
		setHandler:   useCases.Set,
		resetHandler: useCases.Reset,
		getHandler:   useCases.Get,
		listHandler:  useCases.List,
	}
}

// List 获取团队配置列表
//
//	@Summary		团队配置列表
//	@Description	获取当前团队的配置列表（支持三级继承：团队>组织>系统默认值）
//	@Tags			settings-team
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id		path		int		true	"组织ID"
//	@Param			team_id	path		int		true	"团队ID"
//	@Param			params	query		team.ListQuery	false	"查询参数"
//	@Success		200		{object}	response.DataResponse[[]team.SettingsItemDTO]	"配置列表"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		403		{object}	response.ErrorResponse	"权限不足"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id}/settings [get]
func (h *TeamSettingHandler) List(c *gin.Context) {
	orgID, ok := ctxutil.Get[uint](c, ctxutil.OrgID)
	if !ok {
		response.Unauthorized(c, "organization context not found")
		return
	}

	teamID, ok := ctxutil.Get[uint](c, ctxutil.TeamID)
	if !ok {
		response.Unauthorized(c, "team context not found")
		return
	}

	var query team.ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	query.OrgID = orgID
	query.TeamID = teamID

	result, err := h.listHandler.Handle(c.Request.Context(), query)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Get 获取单个团队配置
//
//	@Summary		获取团队配置
//	@Description	获取指定配置项的值（支持三级继承：团队>组织>系统默认值）
//	@Tags			settings-team
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int		true	"组织ID"
//	@Param			team_id	path		int		true	"团队ID"
//	@Param			key		path		string	true	"配置键名"
//	@Success		200		{object}	response.DataResponse[team.SettingsItemDTO]	"配置信息"
//	@Failure		400		{object}	response.ErrorResponse	"无效的键名"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		403		{object}	response.ErrorResponse	"权限不足"
//	@Failure		404		{object}	response.ErrorResponse	"配置不存在"
//	@Router			/api/org/{org_id}/teams/{team_id}/settings/{key} [get]
func (h *TeamSettingHandler) Get(c *gin.Context) {
	orgID, ok := ctxutil.Get[uint](c, ctxutil.OrgID)
	if !ok {
		response.Unauthorized(c, "organization context not found")
		return
	}

	teamID, ok := ctxutil.Get[uint](c, ctxutil.TeamID)
	if !ok {
		response.Unauthorized(c, "team context not found")
		return
	}

	key := c.Param("key")
	query := team.GetQuery{
		OrgID:  orgID,
		TeamID: teamID,
		Key:    key,
	}

	result, err := h.getHandler.Handle(c.Request.Context(), query)
	if err != nil {
		if errors.Is(err, teamDomain.ErrInvalidSettingKey) {
			response.NotFound(c, "setting")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Set 设置团队配置
//
//	@Summary		设置团队配置
//	@Description	设置指定配置项的值（团队自定义覆盖）
//	@Tags			settings-team
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int			true	"组织ID"
//	@Param			team_id	path		int			true	"团队ID"
//	@Param			key		path		string		true	"配置键名"
//	@Param			request	body		SetRequest	true	"设置请求"
//	@Success		200		{object}	response.DataResponse[team.SettingsItemDTO]	"设置成功"
//	@Failure		400		{object}	response.ErrorResponse	"参数错误"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		403		{object}	response.ErrorResponse	"权限不足"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id}/settings/{key} [put]
func (h *TeamSettingHandler) Set(c *gin.Context) {
	teamID, ok := ctxutil.Get[uint](c, ctxutil.TeamID)
	if !ok {
		response.Unauthorized(c, "team context not found")
		return
	}

	key := c.Param("key")
	var req SetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	cmd := team.SetCommand{
		TeamID: teamID,
		Key:    key,
		Value:  req.Value,
	}

	result, err := h.setHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		if errors.Is(err, teamDomain.ErrInvalidSettingKey) {
			response.BadRequest(c, "invalid setting key", err.Error())
			return
		}
		if errors.Is(err, teamDomain.ErrInvalidSettingValue) {
			response.BadRequest(c, "invalid setting value", err.Error())
			return
		}
		if errors.Is(err, teamDomain.ErrValidationFailed) {
			response.BadRequest(c, "validation failed", err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Reset 重置团队配置
//
//	@Summary		重置团队配置
//	@Description	重置指定配置项（删除团队自定义值，恢复组织配置或系统默认值）
//	@Tags			settings-team
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id	path		int		true	"组织ID"
//	@Param			team_id	path		int		true	"团队ID"
//	@Param			key		path		string	true	"配置键名"
//	@Success		200		{object}	response.MessageResponse	"重置成功"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		403		{object}	response.ErrorResponse	"权限不足"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id}/settings/{key} [delete]
func (h *TeamSettingHandler) Reset(c *gin.Context) {
	teamID, ok := ctxutil.Get[uint](c, ctxutil.TeamID)
	if !ok {
		response.Unauthorized(c, "team context not found")
		return
	}

	key := c.Param("key")
	cmd := team.ResetCommand{
		TeamID: teamID,
		Key:    key,
	}

	if err := h.resetHandler.Handle(c.Request.Context(), cmd); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	result := team.ResetResultDTO{
		Key:     key,
		Message: "setting reset to default",
	}
	response.OK(c, result)
}
