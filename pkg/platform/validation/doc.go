// Package validation 提供基于 JSON Logic 的设置验证实现。
//
// 本包实现 [setting.Validator] 接口，为 Setting 模块提供灵活的
// 规则验证能力。
//
// # 组件
//
//   - [JSONLogicValidator]: 验证器实现，支持两种规则格式
//
// # 规则格式
//
// 支持两种规则格式，可混合使用：
//
// JSON Logic 格式（推荐，功能更强大）：
//
//	{">=": [{"var": "value"}, 6]}
//	{"and": [{">=": [{"var": "value"}, 1]}, {"<=": [{"var": "value"}, 100]}]}
//
// 简单格式（向后兼容，自动转换为 JSON Logic）：
//
//	{"min": 6, "max": 32}
//	{"required": true, "minLength": 1}
//
// # 简单格式支持的字段
//
//   - required: 是否必填（bool）
//   - min/max: 数值范围
//   - minLength/maxLength: 字符串长度范围
//   - pattern: 正则表达式（暂未实现）
//   - message: 自定义错误消息
//
// # 使用示例
//
//	validator := validation.NewJSONLogicValidator()
//
//	result, err := validator.Validate(ctx, &setting.ValidationContext{
//	    Key:   "password_min_length",
//	    Value: 4,
//	    Rule:  `{">=": [{"var": "value"}, 6]}`,
//	})
//	if !result.Valid {
//	    fmt.Println(result.Message) // "最小值为 6"
//	}
//
// # 跨字段验证
//
// JSON Logic 支持引用其他设置值进行跨字段验证：
//
//	// 验证 max_login_attempts 必须大于 warning_threshold
//	{
//	    ">": [
//	        {"var": "value"},
//	        {"var": "settings.warning_threshold"}
//	    ]
//	}
//
// # 设计说明
//
// 本包专用于 Setting 模块的验证需求。如需通用表单验证，
// 请使用 go-playground/validator 等库。
package validation
