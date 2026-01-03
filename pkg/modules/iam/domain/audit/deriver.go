package audit

import "strings"

// ============================================================================
// URN → 审计信息派生
// ============================================================================

// typeToCategory type 到审计分类的映射。
//
//nolint:gochecknoglobals // 只读映射表
var typeToCategory = map[string]Category{
	"users":    CatUser,
	"roles":    CatRole,
	"auth":     CatAuth,
	"2fa":      CatAuth, // 特殊映射：2FA 归类为认证
	"profile":  CatProfile,
	"password": CatProfile, // 特殊映射：密码归类为个人资料
	"account":  CatProfile, // 特殊映射：账户归类为个人资料
	"tokens":   CatToken,
	"org":      CatOrg,
	"teams":    CatOrg,
}

// identifierToOperation identifier 到审计操作类型的映射。
//
//nolint:gochecknoglobals // 只读映射表
var identifierToOperation = map[string]Operation{
	"create":       OpCreate,
	"batch-create": OpCreate,
	"update":       OpUpdate,
	"delete":       OpDelete,
	"list":         OpAccess,
	"get":          OpAccess,
	"login":        OpAuthenticate,
	"login2fa":     OpAuthenticate,
	"refresh":      OpAuthenticate,
	"register":     OpCreate,
}

// DeriveCategory 从 URN 的 type 段派生审计分类。
//
// 优先使用映射表，未找到时使用 type 的单数形式。
func DeriveCategory(urnType string) Category {
	if cat, ok := typeToCategory[urnType]; ok {
		return cat
	}
	return Category(singularize(urnType))
}

// DeriveOperation 从 URN 的 identifier 段派生审计操作类型。
//
// 优先使用映射表，未找到时默认为 update。
func DeriveOperation(identifier string) Operation {
	if op, ok := identifierToOperation[identifier]; ok {
		return op
	}
	return OpUpdate // 默认
}

// DeriveAction 从 URN 的 type 和 identifier 派生审计操作标识。
//
// 格式：{category}.{identifier}，连字符转下划线。
// 示例：type="users", identifier="create" → "user.create"
func DeriveAction(urnType, identifier string) string {
	cat := DeriveCategory(urnType)
	id := strings.ReplaceAll(identifier, "-", "_")
	return string(cat) + "." + id
}

// singularize 将复数形式转换为单数。
// 简单实现：移除尾部 s/es。
func singularize(s string) string {
	if len(s) == 0 {
		return s
	}
	// 特殊情况：以 -categories 结尾
	if before, ok := strings.CutSuffix(s, "-categories"); ok {
		return before + "_category"
	}
	// 简单规则：移除尾部 s
	if strings.HasSuffix(s, "s") {
		return s[:len(s)-1]
	}
	return s
}
