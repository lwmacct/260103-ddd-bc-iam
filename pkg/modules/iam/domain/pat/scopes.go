package pat

import (
	"slices"
	"strings"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/role"
)

// Scope 定义 PAT 权限范围标识符。
//
// 采用极简设计，只有 3 个 Scope：
//   - full: 完整权限，继承用户全部权限
//   - self: 仅 self 域权限
//   - sys: 仅 sys 域权限
type Scope string

// 预定义 Scope 常量。
const (
	ScopeFull Scope = "full" // 完整权限
	ScopeSelf Scope = "self" // 仅 self 域
	ScopeSys  Scope = "sys"  // 仅 sys 域
)

// ValidScopes 有效 Scope 列表。
var ValidScopes = []Scope{ScopeFull, ScopeSelf, ScopeSys}

// ScopeInfo Scope 元信息，供前端展示。
type ScopeInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

// AllScopes 所有可用 Scope 的元信息列表。
var AllScopes = []ScopeInfo{
	{
		Name:        string(ScopeFull),
		DisplayName: "完整权限",
		Description: "继承用户全部权限，可访问用户权限范围内的所有 API",
	},
	{
		Name:        string(ScopeSelf),
		DisplayName: "自服务权限",
		Description: "仅限个人资料、令牌管理等自服务操作（self 域）",
	},
	{
		Name:        string(ScopeSys),
		DisplayName: "管理权限",
		Description: "仅限系统管理操作，如用户、角色、配置管理（sys 域）",
	},
}

// IsValidScope 检查给定的 scope 名称是否有效。
func IsValidScope(name string) bool {
	return slices.Contains(ValidScopes, Scope(name))
}

// ValidateScopes 验证 scope 列表，返回无效的 scope 名称。
func ValidateScopes(scopes []string) []string {
	var invalid []string
	for _, s := range scopes {
		if !IsValidScope(s) {
			invalid = append(invalid, s)
		}
	}
	return invalid
}

// FilterByScopes 根据 Scope 列表过滤用户权限。
//
// 权限计算逻辑：
//   - full: 直接返回用户全部权限
//   - self/sys: 过滤用户权限，只保留 scope 前缀匹配的条目
//   - 多选时取并集: ["self", "sys"] = self 权限 + sys 权限
func FilterByScopes(scopes []string, userPerms []role.Permission) []role.Permission {
	if len(scopes) == 0 {
		return userPerms
	}

	// full scope 直接返回全部权限
	if slices.Contains(scopes, string(ScopeFull)) {
		return userPerms
	}

	// 按 scope 前缀过滤，多个 scope 取并集
	seen := make(map[string]bool)
	var result []role.Permission

	for _, p := range userPerms {
		// 检查是否匹配任一 scope 前缀
		for _, s := range scopes {
			if strings.HasPrefix(p.OperationPattern, s+":") {
				key := p.OperationPattern + "|" + p.ResourcePattern
				if !seen[key] {
					seen[key] = true
					result = append(result, p)
				}
				break
			}
		}
	}

	return result
}
