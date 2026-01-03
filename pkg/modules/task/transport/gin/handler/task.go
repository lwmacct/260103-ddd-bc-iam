package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	taskapplication "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/application"
	taskDomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/domain"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/ctxutil"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
)

// ListTasksQuery 任务列表查询参数
type ListTasksQuery struct {
	response.PaginationQueryDTO
}

// ToQuery 转换为 Application 层 Query 对象
func (q *ListTasksQuery) ToQuery(orgID, teamID uint) taskapplication.ListTasksQuery {
	return taskapplication.ListTasksQuery{
		OrgID:  orgID,
		TeamID: teamID,
		Offset: q.GetOffset(),
		Limit:  q.GetLimit(),
	}
}

// TaskHandler 任务管理 Handler
type TaskHandler struct {
	createHandler *taskapplication.CreateHandler
	updateHandler *taskapplication.UpdateHandler
	deleteHandler *taskapplication.DeleteHandler
	getHandler    *taskapplication.GetHandler
	listHandler   *taskapplication.ListHandler
}

// NewTaskHandler 创建任务管理 Handler
func NewTaskHandler(
	createHandler *taskapplication.CreateHandler,
	updateHandler *taskapplication.UpdateHandler,
	deleteHandler *taskapplication.DeleteHandler,
	getHandler *taskapplication.GetHandler,
	listHandler *taskapplication.ListHandler,
) *TaskHandler {
	return &TaskHandler{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

// extractOrgTeamContext 从上下文中提取 org_id 和 team_id
func extractOrgTeamContext(c *gin.Context) (uint, uint, error) {
	orgID, ok := ctxutil.Get[uint](c, ctxutil.OrgID)
	if !ok {
		return 0, 0, errors.New("org_id not found in context")
	}
	teamID, ok := ctxutil.Get[uint](c, ctxutil.TeamID)
	if !ok {
		return 0, 0, errors.New("team_id not found in context")
	}
	return orgID, teamID, nil
}

// extractUserID 从上下文中提取当前用户 ID
func extractUserID(c *gin.Context) (uint, error) {
	userID, ok := ctxutil.Get[uint](c, ctxutil.UserID)
	if !ok {
		return 0, errors.New("user_id not found in context")
	}
	return userID, nil
}

// Create 创建任务
//
//	@Summary		创建任务
//	@Description	在团队内创建新任务
//	@Tags			Organization - Task Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id		path		int									true	"组织ID"
//	@Param			team_id	path		int									true	"团队ID"
//	@Param			request	body		taskapplication.CreateTaskDTO			true	"任务信息"
//	@Success		201		{object}	response.DataResponse[taskapplication.TaskDTO]	"任务创建成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		403		{object}	response.ErrorResponse						"权限不足"
//	@Failure		500		{object}	response.ErrorResponse						"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id}/tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	orgID, teamID, err := extractOrgTeamContext(c)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	userID, err := extractUserID(c)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	var req taskapplication.CreateTaskDTO
	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.createHandler.Handle(c.Request.Context(), taskapplication.CreateTaskCommand{
		OrgID:       orgID,
		TeamID:      teamID,
		Title:       req.Title,
		Description: req.Description,
		AssigneeID:  req.AssigneeID,
		CreatedBy:   userID,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, result)
}

// List 任务列表
//
//	@Summary		任务列表
//	@Description	分页获取团队任务列表
//	@Tags			Organization - Task Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id		path		int										true	"组织ID"
//	@Param			team_id	path		int										true	"团队ID"
//	@Param			params	query		handler.ListTasksQuery					false	"查询参数"
//	@Success		200		{object}	response.PagedResponse[taskapplication.TaskDTO]	"任务列表"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		403		{object}	response.ErrorResponse						"权限不足"
//	@Failure		500		{object}	response.ErrorResponse						"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id}/tasks [get]
func (h *TaskHandler) List(c *gin.Context) {
	orgID, teamID, err := extractOrgTeamContext(c)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	var query ListTasksQuery
	if err = c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.listHandler.Handle(c.Request.Context(), query.ToQuery(orgID, teamID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), query.GetPage(), query.GetLimit())
	response.List(c, result.Items, meta)
}

// Get 任务详情
//
//	@Summary		任务详情
//	@Description	获取任务详细信息
//	@Tags			Organization - Task Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id		path		int										true	"组织ID"
//	@Param			team_id	path		int										true	"团队ID"
//	@Param			id		path		int										true	"任务ID"
//	@Success		200		{object}	response.DataResponse[taskapplication.TaskDTO]	"任务详情"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		403		{object}	response.ErrorResponse						"权限不足"
//	@Failure		404		{object}	response.ErrorResponse						"任务不存在"
//	@Failure		500		{object}	response.ErrorResponse						"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id}/tasks/{id} [get]
func (h *TaskHandler) Get(c *gin.Context) {
	orgID, teamID, err := extractOrgTeamContext(c)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的任务ID")
		return
	}

	result, err := h.getHandler.Handle(c.Request.Context(), taskapplication.GetTaskQuery{
		OrgID:  orgID,
		TeamID: teamID,
		ID:     uint(id),
	})
	if err != nil {
		if errors.Is(err, taskDomain.ErrTaskNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Update 更新任务
//
//	@Summary		更新任务
//	@Description	更新任务信息
//	@Tags			Organization - Task Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id		path		int										true	"组织ID"
//	@Param			team_id	path		int										true	"团队ID"
//	@Param			id		path		int										true	"任务ID"
//	@Param			request	body		taskapplication.UpdateTaskDTO				true	"更新信息"
//	@Success		200		{object}	response.DataResponse[taskapplication.TaskDTO]	"更新成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		403		{object}	response.ErrorResponse						"权限不足"
//	@Failure		404		{object}	response.ErrorResponse						"任务不存在"
//	@Failure		500		{object}	response.ErrorResponse						"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id}/tasks/{id} [put]
func (h *TaskHandler) Update(c *gin.Context) {
	orgID, teamID, err := extractOrgTeamContext(c)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的任务ID")
		return
	}

	var req taskapplication.UpdateTaskDTO
	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.updateHandler.Handle(c.Request.Context(), taskapplication.UpdateTaskCommand{
		OrgID:       orgID,
		TeamID:      teamID,
		ID:          uint(id),
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		AssigneeID:  req.AssigneeID,
	})
	if err != nil {
		if errors.Is(err, taskDomain.ErrTaskNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		if errors.Is(err, taskDomain.ErrInvalidStatusTransition) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// Delete 删除任务
//
//	@Summary		删除任务
//	@Description	删除任务
//	@Tags			Organization - Task Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			org_id		path		int						true	"组织ID"
//	@Param			team_id	path		int						true	"团队ID"
//	@Param			id		path		int						true	"任务ID"
//	@Success		200		{object}	response.MessageResponse	"删除成功"
//	@Failure		401		{object}	response.ErrorResponse	"未授权"
//	@Failure		403		{object}	response.ErrorResponse	"权限不足"
//	@Failure		404		{object}	response.ErrorResponse	"任务不存在"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/org/{org_id}/teams/{team_id}/tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	orgID, teamID, err := extractOrgTeamContext(c)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的任务ID")
		return
	}

	if err = h.deleteHandler.Handle(c.Request.Context(), taskapplication.DeleteTaskCommand{
		OrgID:  orgID,
		TeamID: teamID,
		ID:     uint(id),
	}); err != nil {
		if errors.Is(err, taskDomain.ErrTaskNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil)
}
