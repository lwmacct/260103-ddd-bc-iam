package pat

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// PermissionList 是 PAT 权限列表的值对象。
//
// 实现 sql.Scanner 和 driver.Valuer 接口，支持：
//   - 数据库 JSON 字段的自动序列化/反序列化
//   - 空值安全处理（nil 转为空数组）
//
// 示例权限格式: ["user:read", "user:write", "role:read"]
//
//nolint:recvcheck // Scan needs pointer, Value uses value per SQL interface conventions
type PermissionList []string

// Scan implements sql.Scanner interface
func (p *PermissionList) Scan(value any) error {
	if value == nil {
		*p = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan PermissionList")
	}

	return json.Unmarshal(bytes, p)
}

// Value implements driver.Valuer interface
func (p PermissionList) Value() (driver.Value, error) {
	if p == nil {
		return json.Marshal([]string{})
	}
	return json.Marshal(p)
}

// StringList 是字符串数组的值对象，用于 IP 白名单等场景。
//
// 实现 sql.Scanner 和 driver.Valuer 接口，支持数据库 JSON 字段存储。
//
//nolint:recvcheck // Scan needs pointer, Value uses value per SQL interface conventions
type StringList []string

// Scan implements sql.Scanner interface
func (s *StringList) Scan(value any) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan StringList")
	}

	return json.Unmarshal(bytes, s)
}

// Value implements driver.Valuer interface
func (s StringList) Value() (driver.Value, error) {
	if s == nil {
		return json.Marshal([]string{})
	}
	return json.Marshal(s)
}
