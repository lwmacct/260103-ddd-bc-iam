package handler

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/application/audit"
	auditDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/audit"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/transport/gin/registry"
)

// ListAuditQuery 审计日志列表查询参数
type ListAuditQuery struct {
	response.PaginationQueryDTO

	// UserID 按用户 ID 过滤
	UserID *uint `form:"user_id" json:"user_id" binding:"omitempty,gt=0"`
	// Action 操作类型过滤（语义化标识，如 setting.update）
	Action string `form:"action" json:"action"`
	// Resource 资源分类过滤（如 setting, user）
	Resource string `form:"resource" json:"resource"`
	// Status 状态过滤
	Status string `form:"status" json:"status" binding:"omitempty,oneof=success failure" enums:"success,failure"`
	// StartDate 开始时间（RFC3339 格式）
	StartDate string `form:"start_date" json:"start_date" binding:"omitempty"`
	// EndDate 结束时间（RFC3339 格式）
	EndDate string `form:"end_date" json:"end_date" binding:"omitempty"`
}

// ToQuery 转换为 Application 层 Query 对象
func (q *ListAuditQuery) ToQuery() audit.ListQuery {
	result := audit.ListQuery{
		Page:     q.GetPage(),
		Limit:    q.GetLimit(),
		UserID:   q.UserID,
		Action:   q.Action,
		Resource: q.Resource,
		Status:   q.Status,
	}

	// 解析时间字符串
	if q.StartDate != "" {
		if t, err := time.Parse(time.RFC3339, q.StartDate); err == nil {
			result.StartDate = &t
		}
	}
	if q.EndDate != "" {
		if t, err := time.Parse(time.RFC3339, q.EndDate); err == nil {
			result.EndDate = &t
		}
	}

	return result
}

// AuditHandler handles audit log operations (DDD+CQRS Use Case Pattern)
type AuditHandler struct {
	// Query Handlers
	listHandler *audit.ListHandler
	getHandler  *audit.GetHandler
}

// NewAuditHandler creates a new AuditHandler instance
func NewAuditHandler(
	listHandler *audit.ListHandler,
	getHandler *audit.GetHandler,
) *AuditHandler {
	return &AuditHandler{
		listHandler: listHandler,
		getHandler:  getHandler,
	}
}

// ListLogs lists audit logs with filtering
//
//	@Summary		审计日志列表
//	@Description	分页获取审计日志，支持按用户、操作、资源、状态、时间范围筛选
//	@Tags			Admin - Audit Log
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			params	query		handler.ListAuditQuery						false	"查询参数"
//	@Success		200		{object}	response.PagedResponse[audit.AuditDTO]	"审计日志列表"
//	@Failure		400		{object}	response.ErrorResponse							"参数错误"
//	@Failure		401		{object}	response.ErrorResponse							"未授权"
//	@Failure		403		{object}	response.ErrorResponse							"权限不足"
//	@Failure		500		{object}	response.ErrorResponse							"服务器内部错误"
//	@Router			/api/admin/audit [get]
func (h *AuditHandler) ListLogs(c *gin.Context) {
	var q ListAuditQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.listHandler.Handle(c.Request.Context(), q.ToQuery())
	if err != nil {
		response.InternalError(c, "failed to list audit logs")
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), q.GetPage(), q.GetLimit())
	response.List(c, result.Logs, meta)
}

// GetLog gets an audit log by ID
//
//	@Summary		审计日志详情
//	@Description	根据日志ID获取审计日志详细信息
//	@Tags			Admin - Audit Log
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int										true	"日志ID"	minimum(1)
//	@Success		200	{object}	response.DataResponse[audit.AuditDTO]	"日志详情"
//	@Failure		400	{object}	response.ErrorResponse					"无效的日志ID"
//	@Failure		401	{object}	response.ErrorResponse					"未授权"
//	@Failure		403	{object}	response.ErrorResponse					"权限不足"
//	@Failure		404	{object}	response.ErrorResponse					"日志不存在"
//	@Router			/api/admin/audit/{id} [get]
func (h *AuditHandler) GetLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, auditDomain.ErrInvalidLogID.Error())
		return
	}

	log, err := h.getHandler.Handle(c.Request.Context(), audit.GetQuery{
		LogID: uint(id),
	})

	if err != nil {
		if errors.Is(err, auditDomain.ErrAuditLogNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, log)
}

// GetActions returns audit action definitions
//
//	@Summary		审计操作定义
//	@Description	获取所有审计操作的定义、分类和操作类型，供前端筛选器使用
//	@Tags			Admin - Audit Log
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.DataResponse[audit.AuditActionsResponseDTO]	"审计操作定义"
//	@Failure		401	{object}	response.ErrorResponse									"未授权"
//	@Failure		403	{object}	response.ErrorResponse									"权限不足"
//	@Router			/api/admin/audit/actions [get]
func (h *AuditHandler) GetActions(c *gin.Context) {
	// 从 routes 包获取审计操作定义，然后转换为 DTO
	actions := toApplicationAuditActions(registry.AllAuditActions())
	categories := toApplicationCategoryOptions(registry.AllAuditCategories())
	operations := audit.AllAuditOperations()

	resp := audit.ToAuditActionsResponseDTO(actions, categories, operations)
	response.OK(c, resp)
}

// toApplicationAuditActions 将 registry.AuditActionDefinition 转换为 application/audit.AuditActionDefinition
func toApplicationAuditActions(src []registry.AuditActionDefinition) []audit.AuditActionDefinition {
	result := make([]audit.AuditActionDefinition, len(src))
	for i, a := range src {
		result[i] = audit.AuditActionDefinition{
			Action:      a.Action,
			Operation:   a.Operation,
			Category:    a.Category,
			Label:       a.Label,
			Description: a.Description,
			OperationID: a.OperationID,
		}
	}
	return result
}

// toApplicationCategoryOptions 将 registry.CategoryOption 转换为 application/audit.CategoryOption
func toApplicationCategoryOptions(src []registry.CategoryOption) []audit.CategoryOption {
	result := make([]audit.CategoryOption, len(src))
	for i, o := range src {
		result[i] = audit.CategoryOption{
			Value: o.Value,
			Label: o.Label,
		}
	}
	return result
}
