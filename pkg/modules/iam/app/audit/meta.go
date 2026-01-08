package audit

import (
	auditDomain "github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/audit"
)

// AuditActionDefinition 审计操作定义，供前端动态选项使用。
type AuditActionDefinition struct {
	Action      string                `json:"action"`       // 审计操作标识，如 user.create
	Operation   auditDomain.Operation `json:"operation"`    // 操作类型，如 create
	Category    auditDomain.Category  `json:"category"`     // 分类，如 user
	Label       string                `json:"label"`        // 中文标签
	Description string                `json:"description"`  // 描述
	OperationID string                `json:"operation_id"` // API 操作标识
}

// CategoryOption 分类选项（用于前端下拉框）。
type CategoryOption struct {
	Value string `json:"value"` // 分类值
	Label string `json:"label"` // 显示标签
}

// OperationTypeOption 操作类型选项。
type OperationTypeOption struct {
	Value string `json:"value"` // 操作类型值
	Label string `json:"label"` // 显示标签
}

// AllAuditOperations 返回所有审计操作类型选项。
func AllAuditOperations() []OperationTypeOption {
	ops := []auditDomain.Operation{
		auditDomain.OpCreate, auditDomain.OpUpdate, auditDomain.OpDelete,
		auditDomain.OpAccess, auditDomain.OpAuthenticate,
	}
	result := make([]OperationTypeOption, len(ops))
	for i, op := range ops {
		result[i] = OperationTypeOption{
			Value: string(op),
			Label: op.Label(),
		}
	}
	return result
}
