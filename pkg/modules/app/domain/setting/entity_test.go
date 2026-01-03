package setting_test

import (
	"testing"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Validate Tests
// =============================================================================

func TestSetting_Validate_Valid(t *testing.T) {
	tests := []struct {
		name string
		s    setting.Setting
	}{
		{
			name: "字符串类型配置",
			s: setting.Setting{
				Key:          "general.site_name",
				DefaultValue: "My Site",
				Scope:        setting.ScopeSystem,
				CategoryID:   1,
				ValueType:    setting.ValueTypeString,
			},
		},
		{
			name: "数值类型配置",
			s: setting.Setting{
				Key:          "security.password_min_length",
				DefaultValue: 8,
				Scope:        setting.ScopeSystem,
				CategoryID:   2,
				ValueType:    setting.ValueTypeNumber,
			},
		},
		{
			name: "布尔类型配置（用户级）",
			s: setting.Setting{
				Key:          "notification.email_enabled",
				DefaultValue: true,
				Scope:        setting.ScopeUser,
				CategoryID:   3,
				ValueType:    setting.ValueTypeBoolean,
			},
		},
		{
			name: "JSON 类型配置",
			s: setting.Setting{
				Key:          "backup.schedule",
				DefaultValue: map[string]any{"hour": 3, "minute": 0},
				Scope:        setting.ScopeSystem,
				CategoryID:   4,
				ValueType:    setting.ValueTypeJSON,
			},
		},
		{
			name: "nil 默认值",
			s: setting.Setting{
				Key:          "general.logo",
				DefaultValue: nil,
				Scope:        setting.ScopeUser,
				CategoryID:   1,
				ValueType:    setting.ValueTypeString,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.s.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestSetting_Validate_Invalid(t *testing.T) {
	// 有效的基础配置（用于测试单个字段无效的情况）
	validBase := setting.Setting{
		Scope:      setting.ScopeSystem,
		CategoryID: 1,
		ValueType:  setting.ValueTypeString,
	}

	tests := []struct {
		name    string
		s       setting.Setting
		wantErr error
	}{
		{
			name: "空 Key",
			s: setting.Setting{
				Key:        "",
				Scope:      validBase.Scope,
				CategoryID: validBase.CategoryID,
				ValueType:  validBase.ValueType,
			},
			wantErr: setting.ErrInvalidValue,
		},
		{
			name: "无效 Key 格式",
			s: setting.Setting{
				Key:        "invalid",
				Scope:      validBase.Scope,
				CategoryID: validBase.CategoryID,
				ValueType:  validBase.ValueType,
			},
			wantErr: setting.ErrInvalidKeyFormat,
		},
		{
			name: "无效分类（CategoryID 为 0）",
			s: setting.Setting{
				Key:        "general.test",
				Scope:      validBase.Scope,
				CategoryID: 0,
				ValueType:  validBase.ValueType,
			},
			wantErr: setting.ErrCategoryNotFound,
		},
		{
			name: "无效值类型",
			s: setting.Setting{
				Key:        "general.test",
				Scope:      validBase.Scope,
				CategoryID: validBase.CategoryID,
				ValueType:  "invalid",
			},
			wantErr: setting.ErrInvalidValueType,
		},
		{
			name: "无效 Scope",
			s: setting.Setting{
				Key:        "general.test",
				Scope:      "invalid",
				CategoryID: validBase.CategoryID,
				ValueType:  validBase.ValueType,
			},
			wantErr: setting.ErrInvalidScope,
		},
		{
			name: "默认值类型不匹配",
			s: setting.Setting{
				Key:          "general.test",
				DefaultValue: 123, // 应该是 string
				Scope:        validBase.Scope,
				CategoryID:   validBase.CategoryID,
				ValueType:    validBase.ValueType,
			},
			wantErr: setting.ErrInvalidValueType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.s.Validate()
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

// =============================================================================
// ValidateValue Tests
// =============================================================================

func TestSetting_ValidateValue(t *testing.T) {
	tests := []struct {
		name      string
		valueType string
		value     any
		wantErr   bool
	}{
		// String 类型
		{"string 有效", setting.ValueTypeString, "hello", false},
		{"string 空字符串", setting.ValueTypeString, "", false},
		{"string 类型不匹配", setting.ValueTypeString, 123, true},

		// Number 类型
		{"number int", setting.ValueTypeNumber, 42, false},
		{"number float64", setting.ValueTypeNumber, 3.14, false},
		{"number int64", setting.ValueTypeNumber, int64(100), false},
		{"number 类型不匹配", setting.ValueTypeNumber, "123", true},

		// Boolean 类型
		{"boolean true", setting.ValueTypeBoolean, true, false},
		{"boolean false", setting.ValueTypeBoolean, false, false},
		{"boolean 类型不匹配", setting.ValueTypeBoolean, "true", true},

		// JSON 类型
		{"json map", setting.ValueTypeJSON, map[string]any{"key": "value"}, false},
		{"json slice", setting.ValueTypeJSON, []any{1, 2, 3}, false},
		{"json 类型不匹配", setting.ValueTypeJSON, "not json", true},

		// nil 值
		{"nil 值", setting.ValueTypeString, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setting.Setting{ValueType: tt.valueType}
			err := s.ValidateValue(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// =============================================================================
// Query Method Tests
// =============================================================================

func TestSetting_BelongsToCategoryID(t *testing.T) {
	s := setting.Setting{CategoryID: 1}

	assert.True(t, s.BelongsToCategoryID(1))
	assert.False(t, s.BelongsToCategoryID(2))
}

func TestSetting_HasValidationRule(t *testing.T) {
	tests := []struct {
		name       string
		validation string
		want       bool
	}{
		{"有验证规则", `{"min": 6}`, true},
		{"空规则", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setting.Setting{Validation: tt.validation}
			assert.Equal(t, tt.want, s.HasValidationRule())
		})
	}
}

func TestSetting_IsRequired(t *testing.T) {
	tests := []struct {
		name       string
		validation string
		want       bool
	}{
		{"required true 无空格", `{"required":true}`, true},
		{"required true 有空格", `{"required": true}`, true},
		{"required false", `{"required": false}`, false},
		{"无 required", `{"min": 6}`, false},
		{"空配置", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setting.Setting{Validation: tt.validation}
			assert.Equal(t, tt.want, s.IsRequired())
		})
	}
}

func TestSetting_GetKeyCategory(t *testing.T) {
	tests := []struct {
		key      string
		wantCat  string
		wantName string
	}{
		{"general.site_name", "general", "site_name"},
		{"security.password_min", "security", "password_min"},
		{"invalid", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			s := setting.Setting{Key: tt.key}
			assert.Equal(t, tt.wantCat, s.GetKeyCategory())
			assert.Equal(t, tt.wantName, s.GetKeyName())
		})
	}
}

// =============================================================================
// CoerceValue Tests
// =============================================================================

func TestSetting_CoerceValue_String(t *testing.T) {
	s := setting.Setting{ValueType: setting.ValueTypeString}

	tests := []struct {
		name    string
		input   any
		want    string
		wantErr bool
	}{
		{"string 原样", "hello", "hello", false},
		{"int 转换", 42, "42", false},
		{"float 转换", 3.14, "3.14", false},
		{"bool 转换", true, "true", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.CoerceValue(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

func TestSetting_CoerceValue_Number(t *testing.T) {
	s := setting.Setting{ValueType: setting.ValueTypeNumber}

	tests := []struct {
		name    string
		input   any
		want    any // 注意：int 类型保持原样，不会转为 float64
		wantErr bool
	}{
		{"int 原样", 42, 42, false},              // int 保持原样
		{"float64 原样", 3.14, 3.14, false},      // float64 保持原样
		{"string 解析", "123.45", 123.45, false}, // string 解析为 float64
		{"无效 string", "not a number", float64(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.CoerceValue(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

func TestSetting_CoerceValue_Boolean(t *testing.T) {
	s := setting.Setting{ValueType: setting.ValueTypeBoolean}

	tests := []struct {
		name    string
		input   any
		want    bool
		wantErr bool
	}{
		{"bool true", true, true, false},
		{"bool false", false, false, false},
		{"string true", "true", true, false},
		{"string false", "false", false, false},
		{"int 非零", 1, true, false},
		{"int 零", 0, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.CoerceValue(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

func TestSetting_CoerceValue_JSON(t *testing.T) {
	s := setting.Setting{ValueType: setting.ValueTypeJSON}

	t.Run("map 原样", func(t *testing.T) {
		input := map[string]any{"key": "value"}
		result, err := s.CoerceValue(input)
		require.NoError(t, err)
		assert.Equal(t, input, result)
	})

	t.Run("slice 原样", func(t *testing.T) {
		input := []any{1, 2, 3}
		result, err := s.CoerceValue(input)
		require.NoError(t, err)
		assert.Equal(t, input, result)
	})

	t.Run("string 解析", func(t *testing.T) {
		input := `{"key": "value"}`
		result, err := s.CoerceValue(input)
		require.NoError(t, err)
		assert.Equal(t, map[string]any{"key": "value"}, result)
	})

	t.Run("无效 JSON string", func(t *testing.T) {
		_, err := s.CoerceValue("not json")
		assert.Error(t, err)
	})
}

func TestSetting_CoerceValue_Nil(t *testing.T) {
	s := setting.Setting{ValueType: setting.ValueTypeString}
	result, err := s.CoerceValue(nil)
	require.NoError(t, err)
	assert.Nil(t, result)
}

// =============================================================================
// State Change Method Tests
// =============================================================================

func TestSetting_UpdateDefault(t *testing.T) {
	t.Run("类型匹配成功更新", func(t *testing.T) {
		s := setting.Setting{
			DefaultValue: "old",
			ValueType:    setting.ValueTypeString,
		}
		err := s.UpdateDefault("new")
		require.NoError(t, err)
		assert.Equal(t, "new", s.DefaultValue)
	})

	t.Run("类型不匹配失败", func(t *testing.T) {
		s := setting.Setting{
			DefaultValue: "old",
			ValueType:    setting.ValueTypeString,
		}
		err := s.UpdateDefault(123)
		require.ErrorIs(t, err, setting.ErrInvalidValueType)
		assert.Equal(t, "old", s.DefaultValue) // 未改变
	})
}

func TestSetting_UpdateLabel(t *testing.T) {
	s := setting.Setting{Label: "旧标签"}
	s.UpdateLabel("新标签")
	assert.Equal(t, "新标签", s.Label)
}

func TestSetting_UpdateOrder(t *testing.T) {
	s := setting.Setting{Order: 10}
	s.UpdateOrder(20)
	assert.Equal(t, 20, s.Order)
}

// =============================================================================
// Existing Method Tests
// =============================================================================

func TestSetting_IsValidValueType(t *testing.T) {
	tests := []struct {
		valueType string
		want      bool
	}{
		{setting.ValueTypeString, true},
		{setting.ValueTypeNumber, true},
		{setting.ValueTypeBoolean, true},
		{setting.ValueTypeJSON, true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.valueType, func(t *testing.T) {
			s := setting.Setting{ValueType: tt.valueType}
			assert.Equal(t, tt.want, s.IsValidValueType())
		})
	}
}

func TestSetting_GetDefaultValue(t *testing.T) {
	s := setting.Setting{DefaultValue: "test value"}
	assert.Equal(t, "test value", s.GetDefaultValue())
}

// =============================================================================
// Scope Method Tests
// =============================================================================

func TestSetting_IsValidScope(t *testing.T) {
	tests := []struct {
		scope string
		want  bool
	}{
		{setting.ScopeSystem, true},
		{setting.ScopeUser, true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.scope, func(t *testing.T) {
			s := setting.Setting{Scope: tt.scope}
			assert.Equal(t, tt.want, s.IsValidScope())
		})
	}
}

func TestSetting_IsSystemScope(t *testing.T) {
	tests := []struct {
		name  string
		scope string
		want  bool
	}{
		{"system scope", setting.ScopeSystem, true},
		{"user scope", setting.ScopeUser, false},
		{"invalid scope", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setting.Setting{Scope: tt.scope}
			assert.Equal(t, tt.want, s.IsSystemScope())
		})
	}
}

func TestSetting_IsUserScope(t *testing.T) {
	tests := []struct {
		name  string
		scope string
		want  bool
	}{
		{"user scope", setting.ScopeUser, true},
		{"system scope", setting.ScopeSystem, false},
		{"invalid scope", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setting.Setting{Scope: tt.scope}
			assert.Equal(t, tt.want, s.IsUserScope())
		})
	}
}
