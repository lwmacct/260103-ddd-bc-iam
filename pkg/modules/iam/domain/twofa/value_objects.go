package twofa

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// RecoveryCodes 是 2FA 恢复码的值对象。
//
// 恢复码用于用户丢失 TOTP 设备时的账户恢复，特点：
//   - 一次性使用：每个恢复码只能使用一次
//   - 格式：8位数字，以连字符分隔（如 1234-5678）
//   - 建议生成 10 个恢复码供用户保存
//
// 实现 sql.Scanner 和 driver.Valuer 接口，支持数据库 JSON 存储。
//
//nolint:recvcheck // Scan needs pointer, Value uses value per SQL interface conventions
type RecoveryCodes []string

// Scan 实现 sql.Scanner 接口，从数据库读取时自动处理空值
func (r *RecoveryCodes) Scan(value any) error {
	if value == nil {
		*r = []string{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("failed to unmarshal RecoveryCodes value")
	}

	// 如果是空JSON，使用空数组
	if len(bytes) == 0 || string(bytes) == "{}" || string(bytes) == "[]" {
		*r = []string{}
		return nil
	}

	return json.Unmarshal(bytes, r)
}

// Value 实现 driver.Valuer 接口，写入数据库
func (r RecoveryCodes) Value() (driver.Value, error) {
	if len(r) == 0 {
		return json.Marshal([]string{})
	}
	return json.Marshal(r)
}
