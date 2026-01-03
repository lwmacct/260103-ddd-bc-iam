package audit

import (
	auditDomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/audit"
)

// ToAuditDTO 将领域实体转换为 DTO
func ToAuditDTO(log *auditDomain.Audit) *AuditDTO {
	if log == nil {
		return nil
	}

	return &AuditDTO{
		ID:        log.ID,
		UserID:    log.UserID,
		Action:    log.Action,
		Resource:  log.Resource,
		Details:   log.Details,
		IPAddress: log.IPAddress,
		UserAgent: log.UserAgent,
		Status:    log.Status,
		CreatedAt: log.CreatedAt,
	}
}

// ToAuditActionsResponseDTO 将审计操作定义转换为 DTO
// 参数由调用方提供，通常从 starter/gin/routes.Registry 派生
func ToAuditActionsResponseDTO(
	actions []AuditActionDefinition,
	categories []CategoryOption,
	operations []OperationTypeOption,
) AuditActionsResponseDTO {
	return AuditActionsResponseDTO{
		Actions:    ToAuditActionDTOs(actions),
		Categories: ToCategoryOptionDTOs(categories),
		Operations: ToOperationTypeDTOs(operations),
	}
}

// ToAuditActionDTOs 将审计操作定义列表转换为 DTO 列表
func ToAuditActionDTOs(actions []AuditActionDefinition) []AuditActionDTO {
	result := make([]AuditActionDTO, len(actions))
	for i, a := range actions {
		result[i] = AuditActionDTO{
			Action:      a.Action,
			Operation:   string(a.Operation),
			Category:    string(a.Category),
			Label:       a.Label,
			Description: a.Description,
		}
	}
	return result
}

// ToCategoryOptionDTOs 将分类选项列表转换为 DTO 列表
func ToCategoryOptionDTOs(options []CategoryOption) []CategoryOptionDTO {
	result := make([]CategoryOptionDTO, len(options))
	for i, o := range options {
		result[i] = CategoryOptionDTO(o)
	}
	return result
}

// ToOperationTypeDTOs 将操作类型选项列表转换为 DTO 列表
func ToOperationTypeDTOs(options []OperationTypeOption) []OperationTypeDTO {
	result := make([]OperationTypeDTO, len(options))
	for i, o := range options {
		result[i] = OperationTypeDTO(o)
	}
	return result
}
