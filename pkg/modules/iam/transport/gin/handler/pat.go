package handler

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/ctxutil"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/application/pat"
	authDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/auth"
	patDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/pat"
)

// PATHandler handles Personal Access Token operations (DDD+CQRS Use Case Pattern)
type PATHandler struct {
	// Command Handlers
	createHandler  *pat.CreateHandler
	deleteHandler  *pat.DeleteHandler
	disableHandler *pat.DisableHandler
	enableHandler  *pat.EnableHandler

	// Query Handlers
	getHandler  *pat.GetHandler
	listHandler *pat.ListHandler
}

// NewPATHandler creates a new PAT handler
func NewPATHandler(
	createHandler *pat.CreateHandler,
	deleteHandler *pat.DeleteHandler,
	disableHandler *pat.DisableHandler,
	enableHandler *pat.EnableHandler,
	getHandler *pat.GetHandler,
	listHandler *pat.ListHandler,
) *PATHandler {
	return &PATHandler{
		createHandler:  createHandler,
		deleteHandler:  deleteHandler,
		disableHandler: disableHandler,
		enableHandler:  enableHandler,
		getHandler:     getHandler,
		listHandler:    listHandler,
	}
}

// CreateToken creates a new Personal Access Token
//
//	@Summary		创建令牌
//	@Description	用户创建新的个人访问令牌(PAT)，用于API访问。令牌仅在创建时显示一次
//	@Tags			User - Personal Access Token
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		pat.CreateDTO								true	"令牌信息"
//	@Success		201		{object}	response.DataResponse[pat.CreateResultDTO]	"令牌创建成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Router			/api/user/tokens [post]
func (h *PATHandler) CreateToken(c *gin.Context) {
	var req pat.CreateDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	// Get user ID from context (set by Auth middleware)
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, authDomain.ErrUserIDNotFound.Error())
		return
	}

	expiresAt, err := parseExpiresAt(req.ExpiresAt, req.ExpiresIn)
	if err != nil {
		response.BadRequest(c, "invalid expiration date", err.Error())
		return
	}

	// 调用 Use Case Handler
	result, err := h.createHandler.Handle(c.Request.Context(), pat.CreateCommand{
		UserID:      uid,
		Name:        req.Name,
		Scopes:      req.Scopes,
		ExpiresAt:   expiresAt,
		IPWhitelist: req.IPWhitelist,
		Description: req.Description,
	})

	if err != nil {
		response.BadRequest(c, "failed to create token", err.Error())
		return
	}

	response.Created(c, pat.ToCreateResultDTO(result.Token, result.PlainToken))
}

// ListTokens lists all tokens for the current user
//
//	@Summary		令牌列表
//	@Description	获取当前用户的所有个人访问令牌（不包含令牌值）
//	@Tags			User - Personal Access Token
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.DataResponse[[]pat.TokenDTO]	"令牌列表"
//	@Failure		401	{object}	response.ErrorResponse					"未授权"
//	@Failure		500	{object}	response.ErrorResponse					"服务器内部错误"
//	@Router			/api/user/tokens [get]
func (h *PATHandler) ListTokens(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, authDomain.ErrUserIDNotFound.Error())
		return
	}

	// 调用 Use Case Handler
	tokens, err := h.listHandler.Handle(c.Request.Context(), pat.ListQuery{
		UserID: uid,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, tokens)
}

// DeleteToken deletes a specific token
//
//	@Summary		删除令牌
//	@Description	用户删除指定的个人访问令牌（不可恢复）
//	@Tags			User - Personal Access Token
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int							true	"令牌ID"	minimum(1)
//	@Success		200	{object}	response.MessageResponse	"令牌删除成功"
//	@Failure		400	{object}	response.ErrorResponse		"无效的令牌ID"
//	@Failure		401	{object}	response.ErrorResponse		"未授权"
//	@Failure		404	{object}	response.ErrorResponse		"令牌不存在"
//	@Router			/api/user/tokens/{id} [delete]
func (h *PATHandler) DeleteToken(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, authDomain.ErrUserIDNotFound.Error())
		return
	}

	tokenIDStr := c.Param("id")
	tokenID, err := parseTokenID(tokenIDStr)
	if err != nil {
		response.BadRequest(c, "invalid token ID", err.Error())
		return
	}

	err = h.deleteHandler.Handle(c.Request.Context(), pat.DeleteCommand{
		UserID:  uid,
		TokenID: tokenID,
	})
	if err != nil {
		response.BadRequest(c, "failed to delete token", err.Error())
		return
	}

	response.OK(c, nil)
}

// GetToken retrieves details of a specific token
//
//	@Summary		令牌详情
//	@Description	获取指定个人访问令牌的详细信息（不包含令牌值）
//	@Tags			User - Personal Access Token
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int									true	"令牌ID"	minimum(1)
//	@Success		200	{object}	response.DataResponse[pat.TokenDTO]	"令牌详情"
//	@Failure		400	{object}	response.ErrorResponse				"无效的令牌ID"
//	@Failure		401	{object}	response.ErrorResponse				"未授权"
//	@Failure		404	{object}	response.ErrorResponse				"令牌不存在"
//	@Router			/api/user/tokens/{id} [get]
func (h *PATHandler) GetToken(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, authDomain.ErrUserIDNotFound.Error())
		return
	}

	tokenIDStr := c.Param("id")
	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid token ID", nil)
		return
	}

	// 调用 Use Case Handler
	token, err := h.getHandler.Handle(c.Request.Context(), pat.GetQuery{
		UserID:  uid,
		TokenID: uint(tokenID),
	})

	if err != nil {
		if errors.Is(err, patDomain.ErrTokenNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, token)
}

// DisableToken 暂停令牌
//
//	@Summary		禁用令牌
//	@Description	暂停指定令牌的使用（可再次启用）
//	@Tags			User - Personal Access Token
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int							true	"令牌ID"	minimum(1)
//	@Success		200	{object}	response.MessageResponse	"令牌已禁用"
//	@Failure		400	{object}	response.ErrorResponse		"无效的令牌ID"
//	@Failure		401	{object}	response.ErrorResponse		"未授权"
//	@Router			/api/user/tokens/{id}/disable [patch]
func (h *PATHandler) DisableToken(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, authDomain.ErrUserIDNotFound.Error())
		return
	}

	tokenID, err := parseTokenID(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid token ID", err.Error())
		return
	}

	if err := h.disableHandler.Handle(c.Request.Context(), pat.DisableCommand{
		UserID:  uid,
		TokenID: tokenID,
	}); err != nil {
		response.BadRequest(c, "failed to disable token", err.Error())
		return
	}

	response.OK(c, nil)
}

// EnableToken 启用令牌
//
//	@Summary		启用令牌
//	@Description	重新启用已禁用的令牌
//	@Tags			User - Personal Access Token
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int							true	"令牌ID"	minimum(1)
//	@Success		200	{object}	response.MessageResponse	"令牌已启用"
//	@Failure		400	{object}	response.ErrorResponse		"无效的令牌ID"
//	@Failure		401	{object}	response.ErrorResponse		"未授权"
//	@Router			/api/user/tokens/{id}/enable [patch]
func (h *PATHandler) EnableToken(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, authDomain.ErrUserIDNotFound.Error())
		return
	}

	tokenID, err := parseTokenID(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid token ID", err.Error())
		return
	}

	if err := h.enableHandler.Handle(c.Request.Context(), pat.EnableCommand{
		UserID:  uid,
		TokenID: tokenID,
	}); err != nil {
		response.BadRequest(c, "failed to enable token", err.Error())
		return
	}

	response.OK(c, nil)
}

// ListScopes returns available scopes for PAT creation
//
//	@Summary		Scope 列表
//	@Description	获取创建 PAT 时可选的权限范围列表
//	@Tags			User - Personal Access Token
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.DataResponse[[]pat.ScopeInfoDTO]	"Scope 列表"
//	@Failure		401	{object}	response.ErrorResponse						"未授权"
//	@Router			/api/user/tokens/scopes [get]
func (h *PATHandler) ListScopes(c *gin.Context) {
	// 将 domain 类型映射到 DTO
	scopes := make([]pat.ScopeInfoDTO, len(patDomain.AllScopes))
	for i, s := range patDomain.AllScopes {
		scopes[i] = pat.ScopeInfoDTO{
			Name:        s.Name,
			DisplayName: s.DisplayName,
			Description: s.Description,
		}
	}
	response.OK(c, scopes)
}

func parseTokenID(raw string) (uint, error) {
	id, err := strconv.ParseUint(raw, 10, 32)
	return uint(id), err
}

// parseExpiresAt parses expire parameters from request into a timestamp pointer.
func parseExpiresAt(expiresAt *string, expiresIn *int) (*time.Time, error) {
	// 优先处理 expiresAt 字符串
	if expiresAt != nil && *expiresAt != "" {
		return parseExpiresAtString(*expiresAt)
	}

	// 次之处理 expiresIn 天数
	if expiresIn != nil && *expiresIn > 0 {
		t := time.Now().Add(time.Duration(*expiresIn) * 24 * time.Hour).UTC()
		return &t, nil
	}

	return nil, nil //nolint:nilnil // returns nil for not found, valid pattern
}

// parseExpiresAtString 解析过期时间字符串
func parseExpiresAtString(expiresAt string) (*time.Time, error) {
	// 尝试 RFC3339 格式
	parsed, err := time.Parse(time.RFC3339, expiresAt)
	if err == nil {
		utc := parsed.UTC()
		return &utc, nil
	}

	// 尝试本地时区格式（前端 datetime-local）
	localTime, localErr := time.ParseInLocation("2006-01-02T15:04", expiresAt, time.Local)
	if localErr == nil {
		utc := localTime.UTC()
		return &utc, nil
	}

	return nil, errors.New("expires_at must be RFC3339 or yyyy-MM-ddTHH:mm")
}
