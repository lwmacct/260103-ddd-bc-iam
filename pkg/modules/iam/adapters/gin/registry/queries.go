package registry

import (
	"strings"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/audit"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/permission"
)

// ============================================================================
// 派生查询函数
// ============================================================================

// OperationDefinition 操作定义，供前端权限配置使用。
type OperationDefinition struct {
	Code        string `json:"code"`        // 操作代码，如 sys:users:create
	Scope       string `json:"scope"`       // Scope，如 sys
	Type        string `json:"type"`        // 类型，如 users
	Identifier  string `json:"identifier"`  // 标识符，如 create
	Label       string `json:"label"`       // 中文标签
	Description string `json:"description"` // 描述
	Group       string `json:"group"`       // Swagger 分组
}

// AllOperationDefinitions 返回所有操作定义，供前端权限配置使用。
// 仅返回非公开操作（需要权限检查的操作）。
func AllOperationDefinitions() []OperationDefinition {
	ops := make([]OperationDefinition, 0, len(Registry))

	for op, meta := range Registry {
		// 跳过公开操作
		if op.IsPublic() {
			continue
		}

		ops = append(ops, OperationDefinition{
			Code:        string(op),
			Scope:       op.Scope(),
			Type:        op.Type(),
			Identifier:  op.Identifier(),
			Label:       meta.Summary,
			Description: meta.Description,
			Group:       meta.Tags,
		})
	}

	return ops
}

// AuditActionDefinition 审计操作定义，供前端动态选项使用。
type AuditActionDefinition struct {
	Action      string          `json:"action"`       // 审计操作标识，如 user.create
	Operation   audit.Operation `json:"operation"`    // 操作类型，如 create
	Category    audit.Category  `json:"category"`     // 分类，如 user
	Label       string          `json:"label"`        // 中文标签
	Description string          `json:"description"`  // 描述
	OperationID string          `json:"operation_id"` // API 操作标识
}

// AllAuditActions 返回所有审计操作定义。
// 从注册表派生，仅返回有审计定义的操作。
func AllAuditActions() []AuditActionDefinition {
	actions := make([]AuditActionDefinition, 0, len(Registry))

	for op, meta := range Registry {
		if !meta.Audit {
			continue
		}
		actions = append(actions, AuditActionDefinition{
			Action:      audit.DeriveAction(op.Type(), op.Identifier()),
			Operation:   audit.DeriveOperation(op.Identifier()),
			Category:    audit.DeriveCategory(op.Type()),
			Label:       meta.Summary,
			Description: meta.Description,
			OperationID: string(op),
		})
	}

	return actions
}

// CategoryOption 分类选项（用于前端下拉框）。
type CategoryOption struct {
	Value string `json:"value"` // 分类值
	Label string `json:"label"` // 显示标签
}

// AllAuditCategories 返回所有审计分类选项。
func AllAuditCategories() []CategoryOption {
	seen := make(map[audit.Category]bool)
	categories := make([]CategoryOption, 0, 16)

	for op, meta := range Registry {
		if !meta.Audit {
			continue
		}
		cat := audit.DeriveCategory(op.Type())
		if seen[cat] {
			continue
		}
		seen[cat] = true
		categories = append(categories, CategoryOption{
			Value: string(cat),
			Label: cat.Label(),
		})
	}

	return categories
}

// OperationTypeOption 操作类型选项。
type OperationTypeOption struct {
	Value string `json:"value"` // 操作类型值
	Label string `json:"label"` // 显示标签
}

// AllAuditOperations 返回所有审计操作类型选项。
func AllAuditOperations() []OperationTypeOption {
	ops := []audit.Operation{
		audit.OpCreate, audit.OpUpdate, audit.OpDelete,
		audit.OpAccess, audit.OpAuthenticate,
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

// ByOperationID 通过操作标识符查找操作。
// 如果未找到返回空 Operation。
func ByOperationID(id string) permission.Operation {
	op := permission.Operation(id)
	if Valid(op) {
		return op
	}
	return ""
}

// All 返回所有操作。
func All() []permission.Operation {
	ops := make([]permission.Operation, 0, len(Registry))
	for op := range Registry {
		ops = append(ops, op)
	}
	return ops
}

// ByMethodAndPath 通过 HTTP 方法和路径查找操作。
// 支持路径参数匹配：/api/system/users/:id 匹配 /api/system/users/123
// 如果未找到返回空 Operation。
func ByMethodAndPath(method HTTPMethod, path string) permission.Operation {
	for op, meta := range Registry {
		if meta.Method != method {
			continue
		}
		if matchPath(meta.Path, path) {
			return op
		}
	}
	return ""
}

// matchPath 检查实际路径是否匹配模式路径。
func matchPath(pattern, actual string) bool {
	patternSegs := splitPathSegments(pattern)
	actualSegs := splitPathSegments(actual)

	if len(patternSegs) != len(actualSegs) {
		return false
	}

	for i, seg := range patternSegs {
		// :param 匹配任意非空值
		if len(seg) > 0 && seg[0] == ':' {
			continue
		}
		if seg != actualSegs[i] {
			return false
		}
	}

	return true
}

// splitPathSegments 将路径分割为段。
func splitPathSegments(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}
