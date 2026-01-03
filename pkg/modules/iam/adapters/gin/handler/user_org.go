package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/ctxutil"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/org"
	authDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/auth"
	orgDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/org"
)

// UserOrgHandler 用户视角的组织/团队 Handler
type UserOrgHandler struct {
	userOrgsHandler  *org.UserOrgsHandler
	userTeamsHandler *org.UserTeamsHandler
}

// NewUserOrgHandler 创建用户视角组织 Handler
func NewUserOrgHandler(
	userOrgsHandler *org.UserOrgsHandler,
	userTeamsHandler *org.UserTeamsHandler,
) *UserOrgHandler {
	return &UserOrgHandler{
		userOrgsHandler:  userOrgsHandler,
		userTeamsHandler: userTeamsHandler,
	}
}

// ListMyOrganizations 获取我加入的组织列表
//
//	@Summary		我的组织
//	@Description	获取当前用户加入的所有组织
//	@Tags			User - Organization
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.DataResponse[[]org.UserOrgDTO]	"组织列表"
//	@Failure		401	{object}	response.ErrorResponse								"未授权"
//	@Failure		500	{object}	response.ErrorResponse								"服务器内部错误"
//	@Router			/api/user/orgs [get]
func (h *UserOrgHandler) ListMyOrganizations(c *gin.Context) {
	userID, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, authDomain.ErrUserNotAuthenticated.Error())
		return
	}

	result, err := h.userOrgsHandler.Handle(c.Request.Context(), org.ListUserOrgsQuery{
		UserID: userID,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// ListUserTeamsQuery 用户团队列表查询参数
type ListUserTeamsQuery struct {
	OrgID uint `form:"org_id" binding:"omitempty,min=1"`
}

// ListMyTeams 获取我加入的团队列表
//
//	@Summary		我的团队
//	@Description	获取当前用户加入的所有团队（可按组织筛选）
//	@Tags			User - Organization
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			params	query		handler.ListUserTeamsQuery							false	"查询参数"
//	@Success		200		{object}	response.DataResponse[[]org.UserTeamDTO]	"团队列表"
//	@Failure		401		{object}	response.ErrorResponse								"未授权"
//	@Failure		500		{object}	response.ErrorResponse								"服务器内部错误"
//	@Router			/api/user/teams [get]
func (h *UserOrgHandler) ListMyTeams(c *gin.Context) {
	userID, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, authDomain.ErrUserNotAuthenticated.Error())
		return
	}

	// 解析可选的 org_id 参数
	var orgID uint
	if orgIDStr := c.Query("org_id"); orgIDStr != "" {
		id, err := strconv.ParseUint(orgIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, orgDomain.ErrInvalidOrgID.Error())
			return
		}
		orgID = uint(id)
	}

	result, err := h.userTeamsHandler.Handle(c.Request.Context(), org.ListUserTeamsQuery{
		UserID: userID,
		OrgID:  orgID,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}
