package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/user"
	userDomain "github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/user"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/ctxutil"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/response"
)

// UserProfileHandler handles user profile operations
type UserProfileHandler struct {
	getUserHandler        *user.GetHandler
	updateUserHandler     *user.UpdateHandler
	changePasswordHandler *user.ChangePasswordHandler
	deleteUserHandler     *user.DeleteHandler
}

// NewUserProfileHandler creates a new UserProfileHandler instance
func NewUserProfileHandler(useCases *app.UserUseCases) *UserProfileHandler {
	return &UserProfileHandler{
		getUserHandler:        useCases.Get,
		updateUserHandler:     useCases.Update,
		changePasswordHandler: useCases.ChangePassword,
		deleteUserHandler:     useCases.Delete,
	}
}

// GetProfile gets the current user's profile
//
//	@Summary		获取资料
//	@Description	获取当前登录用户的个人资料和角色信息
//	@Tags			user-profile
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.DataResponse[user.UserWithRolesDTO]	"个人资料"
//	@Failure		401	{object}	response.ErrorResponse							"未授权"
//	@Failure		404	{object}	response.ErrorResponse							"用户不存在"
//	@Router			/api/user/profile [get]
func (h *UserProfileHandler) GetProfile(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, response.MsgAuthenticationRequired)
		return
	}

	u, err := h.getUserHandler.Handle(c.Request.Context(), user.GetQuery{
		UserID:    uid,
		WithRoles: true,
	})
	if err != nil {
		if errors.Is(err, userDomain.ErrUserNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, u)
}

// UpdateProfileRequest 更新个人资料请求
type UpdateProfileRequest struct {
	RealName  *string `json:"real_name" binding:"omitempty,max=100" example:"张三"`
	Nickname  *string `json:"nickname" binding:"omitempty,max=50" example:"小三"`
	Phone     *string `json:"phone" binding:"omitempty,len=11" example:"13800138000"`
	Signature *string `json:"signature" binding:"omitempty,max=255" example:"Hello World"`
	Avatar    *string `json:"avatar" binding:"omitempty,max=255" example:"https://example.com/avatar.jpg"`
	Bio       *string `json:"bio" example:"这是我的个人简介"`
}

// UpdateProfile updates the current user's profile
//
//	@Summary		更新资料
//	@Description	用户更新自己的姓名、头像和个人简介
//	@Tags			user-profile
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		UpdateProfileRequest							true	"更新信息"
//	@Success		200		{object}	response.DataResponse[user.UserWithRolesDTO]	"资料更新成功"
//	@Failure		400		{object}	response.ErrorResponse							"参数错误"
//	@Failure		401		{object}	response.ErrorResponse							"未授权"
//	@Failure		500		{object}	response.ErrorResponse							"服务器内部错误"
//	@Router			/api/user/profile [put]
func (h *UserProfileHandler) UpdateProfile(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, response.MsgAuthenticationRequired)
		return
	}

	var req UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if _, err := h.updateUserHandler.Handle(c.Request.Context(), user.UpdateCommand{
		UserID:    uid,
		RealName:  req.RealName,
		Nickname:  req.Nickname,
		Phone:     req.Phone,
		Signature: req.Signature,
		Avatar:    req.Avatar,
		Bio:       req.Bio,
	}); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	updatedUser, err := h.getUserHandler.Handle(c.Request.Context(), user.GetQuery{
		UserID:    uid,
		WithRoles: true,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, updatedUser)
}

// ChangePassword changes the current user's password
//
//	@Summary		修改密码
//	@Description	用户修改自己的登录密码
//	@Tags			user-profile
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		user.ChangePasswordDTO		true	"密码信息"
//	@Success		200		{object}	response.MessageResponse	"密码修改成功"
//	@Failure		400		{object}	response.ErrorResponse		"参数错误或旧密码不正确"
//	@Failure		401		{object}	response.ErrorResponse		"未授权"
//	@Router			/api/user/password [put]
func (h *UserProfileHandler) ChangePassword(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, response.MsgAuthenticationRequired)
		return
	}

	var req user.ChangePasswordDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := h.changePasswordHandler.Handle(c.Request.Context(), user.ChangePasswordCommand{
		UserID:      uid,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, nil)
}

// DeleteAccount deletes the current user's account
//
//	@Summary		注销账户
//	@Description	用户删除自己的账号（不可恢复）
//	@Tags			user-profile
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.MessageResponse	"账号删除成功"
//	@Failure		401	{object}	response.ErrorResponse		"未授权"
//	@Failure		500	{object}	response.ErrorResponse		"服务器内部错误"
//	@Router			/api/user/account [delete]
func (h *UserProfileHandler) DeleteAccount(c *gin.Context) {
	uid, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		response.Unauthorized(c, response.MsgAuthenticationRequired)
		return
	}

	if err := h.deleteUserHandler.Handle(c.Request.Context(), user.DeleteCommand{
		UserID: uid,
	}); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil)
}
