package validation_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/platform/validation"
)

// validationTestCase 验证测试用例
type validationTestCase struct {
	name      string
	rule      string
	value     any
	wantValid bool
}

// runValidationTests 运行验证测试用例（消除重复的测试辅助函数）
func runValidationTests(t *testing.T, tests []validationTestCase) {
	t.Helper()
	validator := validation.NewJSONLogicValidator()
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vctx := &setting.ValidationContext{
				Key:   "test.key",
				Value: tt.value,
				Rule:  tt.rule,
			}
			result, err := validator.Validate(ctx, vctx)
			require.NoError(t, err)
			assert.Equal(t, tt.wantValid, result.Valid, "expected valid=%v, got %v, message=%s", tt.wantValid, result.Valid, result.Message)
		})
	}
}

func TestJSONLogicValidator_Validate_SimpleRules(t *testing.T) {
	runValidationTests(t, []validationTestCase{
		{"简单最小值规则 - 通过", `{"min":6}`, 10.0, true},
		{"简单最小值规则 - 失败", `{"min":6}`, 3.0, false},
		{"简单最大值规则 - 通过", `{"max":100}`, 50.0, true},
		{"简单最大值规则 - 失败", `{"max":100}`, 150.0, false},
		{"简单范围规则 - 通过", `{"min":6,"max":32}`, 8.0, true},
		{"简单范围规则 - 失败（太小）", `{"min":6,"max":32}`, 3.0, false},
		{"简单范围规则 - 失败（太大）", `{"min":6,"max":32}`, 100.0, false},
		{"必填规则 - 通过", `{"required":true}`, "hello", true},
		{"必填规则 - 失败（空字符串）", `{"required":true}`, "", false},
	})
}

func TestJSONLogicValidator_Validate_JSONLogicRules(t *testing.T) {
	runValidationTests(t, []validationTestCase{
		{"JSON Logic 大于等于 - 通过", `{">=": [{"var": "value"}, 6]}`, 10.0, true},
		{"JSON Logic 大于等于 - 失败", `{">=": [{"var": "value"}, 6]}`, 3.0, false},
		{"JSON Logic AND 组合 - 通过", `{"and": [{">=": [{"var": "value"}, 6]}, {"<=": [{"var": "value"}, 32]}]}`, 8.0, true},
		{"JSON Logic AND 组合 - 失败", `{"and": [{">=": [{"var": "value"}, 6]}, {"<=": [{"var": "value"}, 32]}]}`, 100.0, false},
		{"JSON Logic 真值检查 - 通过", `{"!!": {"var": "value"}}`, "hello", true},
		{"JSON Logic 真值检查 - 失败", `{"!!": {"var": "value"}}`, "", false},
	})
}

func TestJSONLogicValidator_Validate_EnumRule(t *testing.T) {
	runValidationTests(t, []validationTestCase{
		{"简单枚举规则 - 通过", `{"enum":["light","dark","auto"]}`, "dark", true},
		{"简单枚举规则 - 失败", `{"enum":["light","dark","auto"]}`, "invalid", false},
		{"JSON Logic in 操作符 - 通过", `{"in":[{"var":"value"},["light","dark","auto"]]}`, "auto", true},
		{"JSON Logic in 操作符 - 失败", `{"in":[{"var":"value"},["light","dark","auto"]]}`, "blue", false},
		{"数值枚举 - 通过", `{"enum":[1,2,3,5,10]}`, 5.0, true},
		{"数值枚举 - 失败", `{"enum":[1,2,3,5,10]}`, 4.0, false},
	})
}

func TestJSONLogicValidator_Validate_StringLength(t *testing.T) {
	runValidationTests(t, []validationTestCase{
		{"简单 min_length 规则 - 通过", `{"min_length":6}`, "password123", true},
		{"简单 min_length 规则 - 失败", `{"min_length":6}`, "pwd", false},
		{"简单 max_length 规则 - 通过", `{"max_length":10}`, "hello", true},
		{"简单 max_length 规则 - 失败", `{"max_length":10}`, "this is a very long string", false},
		{"JSON Logic strlen - 通过", `{">=": [{"strlen": {"var": "value"}}, 6]}`, "password", true},
		{"JSON Logic strlen - 失败", `{">=": [{"strlen": {"var": "value"}}, 6]}`, "pwd", false},
		{"组合长度规则 - 通过", `{"min_length":6,"max_length":32}`, "mypassword", true},
		{"组合长度规则 - 失败（太短）", `{"min_length":6,"max_length":32}`, "pwd", false},
		{"Unicode 字符串长度 - 中文", `{"min_length":3}`, "你好世界", true},
	})
}

func TestJSONLogicValidator_Validate_CrossFieldValidation(t *testing.T) {
	validator := validation.NewJSONLogicValidator()
	ctx := context.Background()

	// 备份频率必须小于 (保留天数 * 24)
	// 注意：JSON Logic 的 var 使用点符号访问嵌套对象
	rule := `{"<": [{"var": "value"}, {"*": [{"var": "settings.backup_retention_days"}, 24]}]}`

	tests := []struct {
		name          string
		value         any
		retentionDays any
		wantValid     bool
	}{
		{"备份频率 24h，保留 30 天 - 通过", 24.0, 30.0, true},
		{"备份频率 168h，保留 7 天 - 失败（168 >= 7*24=168）", 168.0, 7.0, false},
		{"备份频率 1h，保留 1 天 - 通过（1 < 24）", 1.0, 1.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vctx := &setting.ValidationContext{
				Key:   "backup.backup_frequency",
				Value: tt.value,
				Rule:  rule,
				AllSettings: map[string]any{
					// JSON Logic 需要嵌套结构：settings.backup_retention_days
					"backup_retention_days": tt.retentionDays,
				},
			}
			result, err := validator.Validate(ctx, vctx)
			require.NoError(t, err)
			assert.Equal(t, tt.wantValid, result.Valid, "expected valid=%v, got %v, message=%s", tt.wantValid, result.Valid, result.Message)
		})
	}
}

func TestJSONLogicValidator_ValidateBatch(t *testing.T) {
	validator := validation.NewJSONLogicValidator()
	ctx := context.Background()

	items := []*setting.ValidationContext{
		{Key: "password_length", Value: 4.0, Rule: `{"min":6,"max":32}`},    // 小于 6，应该失败
		{Key: "session_timeout", Value: 30.0, Rule: `{"min":5,"max":1440}`}, // 在范围内，应该通过
		{Key: "max_attempts", Value: 15.0, Rule: `{"min":3,"max":10}`},      // 大于 10，应该失败
	}

	errors, err := validator.ValidateBatch(ctx, items)
	require.NoError(t, err)

	// 应该有 2 个错误
	assert.Len(t, errors, 2)
	assert.Contains(t, errors, "password_length")
	assert.Contains(t, errors, "max_attempts")
	assert.NotContains(t, errors, "session_timeout")
}

func TestJSONLogicValidator_Validate_InvalidRule(t *testing.T) {
	validator := validation.NewJSONLogicValidator()
	ctx := context.Background()

	vctx := &setting.ValidationContext{
		Key:   "test.key",
		Value: 10.0,
		Rule:  `{invalid json`,
	}

	_, err := validator.Validate(ctx, vctx)
	assert.Error(t, err)
}

func TestJSONLogicValidator_Validate_EmptyRule(t *testing.T) {
	validator := validation.NewJSONLogicValidator()
	ctx := context.Background()

	vctx := &setting.ValidationContext{
		Key:   "test.key",
		Value: 10.0,
		Rule:  "",
	}

	result, err := validator.Validate(ctx, vctx)
	require.NoError(t, err)
	assert.True(t, result.Valid, "empty rule should pass validation")
}
