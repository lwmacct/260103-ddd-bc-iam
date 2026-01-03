package audit

// ============================================================================
// 审计操作类型（粗粒度分类）
// ============================================================================

// Operation 审计操作类型，遵循 GitHub Audit Log 风格。
type Operation string

const (
	OpCreate       Operation = "create"
	OpUpdate       Operation = "update"
	OpDelete       Operation = "delete"
	OpAccess       Operation = "access"
	OpAuthenticate Operation = "authenticate"
)

//nolint:gochecknoglobals // 标签映射是只读配置
var operationLabels = map[Operation]string{
	OpCreate:       "创建",
	OpUpdate:       "更新",
	OpDelete:       "删除",
	OpAccess:       "访问",
	OpAuthenticate: "认证",
}

// Label 返回审计操作的中文标签。
func (o Operation) Label() string {
	if label, ok := operationLabels[o]; ok {
		return label
	}
	return string(o)
}

// String 返回审计操作的字符串表示。
func (o Operation) String() string {
	return string(o)
}

// ============================================================================
// 审计分类
// ============================================================================

// Category 审计分类。
type Category string

const (
	CatAuth        Category = "auth"
	CatUser        Category = "user"
	CatRole        Category = "role"
	CatSetting     Category = "setting"
	CatCache       Category = "cache"
	CatProfile     Category = "profile"
	CatToken       Category = "token"
	CatUserSetting Category = "user_setting"
)

//nolint:gochecknoglobals // 标签映射是只读配置
var categoryLabels = map[Category]string{
	CatAuth:        "认证",
	CatUser:        "用户",
	CatRole:        "角色",
	CatSetting:     "配置",
	CatCache:       "缓存",
	CatProfile:     "个人资料",
	CatToken:       "访问令牌",
	CatUserSetting: "用户配置",
}

// Label 返回审计分类的中文标签。
func (c Category) Label() string {
	if label, ok := categoryLabels[c]; ok {
		return label
	}
	return string(c)
}

// String 返回审计分类的字符串表示。
func (c Category) String() string {
	return string(c)
}
