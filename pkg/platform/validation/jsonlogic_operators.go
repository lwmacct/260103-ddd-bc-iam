package validation

import jsonlogic "github.com/diegoholiveira/jsonlogic/v3"

// init 注册 strlen 操作符。
//
// strlen 是 minLength/maxLength 简单规则的基础依赖，
// JSON Logic 原生不支持，需要通过 AddOperator 注册。
func init() {
	jsonlogic.AddOperator("strlen", strlenOperator)
}

// strlenOperator 计算字符串的 UTF-8 字符长度。
//
// 用法: {"strlen": {"var": "value"}}
// 返回: 字符串的字符数（非字节数）
//
// 示例:
//
//	{"strlen": "hello"}     → 5
//	{"strlen": "你好世界"}   → 4
func strlenOperator(values, data any) any {
	str, ok := values.(string)
	if !ok {
		return 0
	}
	return len([]rune(str))
}
