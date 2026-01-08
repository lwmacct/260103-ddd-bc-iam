package audit

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/audit"
)

// ListHandler 获取审计日志列表查询处理器
type ListHandler struct {
	auditQueryRepo audit.QueryRepository
}

// NewListHandler 创建 ListHandler 实例
func NewListHandler(auditQueryRepo audit.QueryRepository) *ListHandler {
	return &ListHandler{
		auditQueryRepo: auditQueryRepo,
	}
}

// Handle 处理获取审计日志列表查询
func (h *ListHandler) Handle(ctx context.Context, query ListQuery) (*ListDTO, error) {
	filter := audit.FilterOptions{
		Page:      query.Page,
		Limit:     query.Limit,
		UserID:    query.UserID,
		Action:    query.Action,
		Resource:  query.Resource,
		Status:    query.Status,
		StartDate: query.StartDate,
		EndDate:   query.EndDate,
	}

	logs, total, err := h.auditQueryRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs: %w", err)
	}

	// 转换为 DTO
	logResponses := make([]*AuditDTO, 0, len(logs))
	for i := range logs {
		logResponses = append(logResponses, ToAuditDTO(&logs[i]))
	}

	return &ListDTO{
		Logs:  logResponses,
		Total: total,
		Page:  query.Page,
		Limit: query.Limit,
	}, nil
}
