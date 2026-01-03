package validation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	jsonlogic "github.com/diegoholiveira/jsonlogic/v3"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
)

// JSONLogicValidator 基于 JSON Logic 的设置验证器。
//
// 支持两种规则格式：
//  1. JSON Logic 格式（推荐）：{">=": [{"var": "value"}, 6]}
//  2. 简单格式（向后兼容）：{"min": 6, "max": 32}
type JSONLogicValidator struct{}

// NewJSONLogicValidator 创建 JSON Logic 验证器实例。
func NewJSONLogicValidator() *JSONLogicValidator {
	return &JSONLogicValidator{}
}

// Validate 验证单个设置值。
func (v *JSONLogicValidator) Validate(ctx context.Context, vctx *setting.ValidationContext) (*setting.ValidationResult, error) {
	if vctx == nil || vctx.Rule == "" {
		// 无验证规则，默认通过
		return &setting.ValidationResult{Valid: true}, nil
	}

	// 解析规则 JSON
	var rule any
	if err := json.Unmarshal([]byte(vctx.Rule), &rule); err != nil {
		return nil, fmt.Errorf("%w: %w", setting.ErrInvalidValidationRule, err)
	}

	// 检测规则格式并转换
	ruleMap, isMap := rule.(map[string]any)
	if isMap && isSimpleRule(ruleMap) {
		// 简单格式，转换为 JSON Logic
		rule = convertSimpleRule(ruleMap)
	}

	// 构建验证数据
	data := buildValidationData(vctx)

	// 执行 JSON Logic 验证
	result, err := v.executeJSONLogic(rule, data)
	if err != nil {
		return nil, err
	}

	if !result {
		msg := vctx.Message
		if msg == "" {
			msg = buildDefaultMessage(vctx.Key, ruleMap)
		}
		return &setting.ValidationResult{Valid: false, Message: msg}, nil
	}

	return &setting.ValidationResult{Valid: true}, nil
}

// ValidateBatch 批量验证多个设置值。
func (v *JSONLogicValidator) ValidateBatch(ctx context.Context, items []*setting.ValidationContext) (map[string]string, error) {
	errors := make(map[string]string)

	for _, item := range items {
		result, err := v.Validate(ctx, item)
		if err != nil {
			return nil, err
		}
		if !result.Valid {
			errors[item.Key] = result.Message
		}
	}

	if len(errors) == 0 {
		return errors, nil
	}
	return errors, nil
}

// executeJSONLogic 执行 JSON Logic 规则。
func (v *JSONLogicValidator) executeJSONLogic(rule any, data map[string]any) (bool, error) {
	// 将规则和数据转换为 JSON
	ruleJSON, err := json.Marshal(rule)
	if err != nil {
		return false, fmt.Errorf("failed to marshal rule: %w", err)
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return false, fmt.Errorf("failed to marshal data: %w", err)
	}

	// 执行 JSON Logic
	var result bytes.Buffer
	if err := jsonlogic.Apply(bytes.NewReader(ruleJSON), bytes.NewReader(dataJSON), &result); err != nil {
		return false, fmt.Errorf("jsonlogic execution failed: %w", err)
	}

	// 解析结果
	var valid bool
	if err := json.Unmarshal(result.Bytes(), &valid); err != nil {
		// 结果可能不是布尔值，尝试其他类型
		var resultAny any
		if jsonErr := json.Unmarshal(result.Bytes(), &resultAny); jsonErr == nil {
			return isTruthy(resultAny), nil
		}
		return false, fmt.Errorf("failed to parse result: %w", err)
	}

	return valid, nil
}

// isSimpleRule 检测是否为简单规则格式。
//
// 简单规则使用声明式 key，自动转换为 JSON Logic：
//   - min/max: 数值范围，转换为 >=/<= 操作符
//   - min_length/max_length: 字符串长度，使用 strlen 操作符
//   - required: 必填，使用 !! 操作符
//   - enum: 枚举，使用 in 操作符
//   - message: 自定义错误消息
func isSimpleRule(rule map[string]any) bool {
	simpleKeys := []string{"min", "max", "min_length", "max_length", "required", "enum", "message"}
	for key := range rule {
		found := slices.Contains(simpleKeys, key)
		if !found {
			return false
		}
	}
	return true
}

// convertSimpleRule 将简单规则转换为 JSON Logic 格式。
//
// 转换映射：
//   - required: true        → {"!!": {"var": "value"}}
//   - min: N                → {">=": [{"var": "value"}, N]}
//   - max: N                → {"<=": [{"var": "value"}, N]}
//   - min_length: N         → {">=": [{"strlen": {"var": "value"}}, N]}
//   - max_length: N         → {"<=": [{"strlen": {"var": "value"}}, N]}
//   - enum: ["a", "b", "c"] → {"in": [{"var": "value"}, ["a", "b", "c"]]}
func convertSimpleRule(rule map[string]any) any {
	var conditions []any

	// required
	if req, ok := rule["required"].(bool); ok && req {
		conditions = append(conditions, map[string]any{
			"!!": map[string]any{"var": "value"},
		})
	}

	// min
	if minVal, ok := rule["min"]; ok {
		conditions = append(conditions, map[string]any{
			">=": []any{map[string]any{"var": "value"}, minVal},
		})
	}

	// max
	if maxVal, ok := rule["max"]; ok {
		conditions = append(conditions, map[string]any{
			"<=": []any{map[string]any{"var": "value"}, maxVal},
		})
	}

	// min_length
	if minLen, ok := rule["min_length"]; ok {
		conditions = append(conditions, map[string]any{
			">=": []any{
				map[string]any{"strlen": map[string]any{"var": "value"}},
				minLen,
			},
		})
	}

	// max_length
	if maxLen, ok := rule["max_length"]; ok {
		conditions = append(conditions, map[string]any{
			"<=": []any{
				map[string]any{"strlen": map[string]any{"var": "value"}},
				maxLen,
			},
		})
	}

	// enum - 使用 JSON Logic 原生 in 操作符
	if enumVals, ok := rule["enum"].([]any); ok && len(enumVals) > 0 {
		conditions = append(conditions, map[string]any{
			"in": []any{map[string]any{"var": "value"}, enumVals},
		})
	}

	if len(conditions) == 0 {
		return true // 无条件，默认通过
	}

	if len(conditions) == 1 {
		return conditions[0]
	}

	return map[string]any{"and": conditions}
}

// buildValidationData 构建验证数据上下文。
func buildValidationData(vctx *setting.ValidationContext) map[string]any {
	data := map[string]any{
		"value": vctx.Value,
		"key":   vctx.Key,
	}

	if vctx.AllSettings != nil {
		data["settings"] = vctx.AllSettings
	}

	return data
}

// buildDefaultMessage 构建默认错误消息。
func buildDefaultMessage(key string, rule map[string]any) string {
	if rule == nil {
		return key + " 验证失败"
	}

	var parts []string

	if req, ok := rule["required"].(bool); ok && req {
		parts = append(parts, "不能为空")
	}
	if minVal, ok := rule["min"]; ok {
		parts = append(parts, fmt.Sprintf("最小值为 %v", minVal))
	}
	if maxVal, ok := rule["max"]; ok {
		parts = append(parts, fmt.Sprintf("最大值为 %v", maxVal))
	}
	if minLen, ok := rule["min_length"]; ok {
		parts = append(parts, fmt.Sprintf("最小长度为 %v", minLen))
	}
	if maxLen, ok := rule["max_length"]; ok {
		parts = append(parts, fmt.Sprintf("最大长度为 %v", maxLen))
	}
	if enumVals, ok := rule["enum"].([]any); ok && len(enumVals) > 0 {
		parts = append(parts, fmt.Sprintf("必须是 %v 之一", enumVals))
	}

	if len(parts) == 0 {
		return key + " 验证失败"
	}

	return strings.Join(parts, "，")
}

// isTruthy 判断值是否为真值。
func isTruthy(v any) bool {
	switch val := v.(type) {
	case bool:
		return val
	case float64:
		return val != 0
	case string:
		return val != ""
	case nil:
		return false
	default:
		return true
	}
}

// 确保实现接口
var _ setting.Validator = (*JSONLogicValidator)(nil)
