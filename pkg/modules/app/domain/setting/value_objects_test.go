package setting_test

import (
	"testing"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// SettingKey Tests
// =============================================================================

func TestNewSettingKey_ValidFormats(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantCategory string
		wantName     string
	}{
		{
			name:         "标准格式",
			input:        "general.site_name",
			wantCategory: "general",
			wantName:     "site_name",
		},
		{
			name:         "多级名称",
			input:        "security.password_min_length",
			wantCategory: "security",
			wantName:     "password_min_length",
		},
		{
			name:         "单字符分类",
			input:        "a.b",
			wantCategory: "a",
			wantName:     "b",
		},
		{
			name:         "包含多个点号",
			input:        "general.site.title",
			wantCategory: "general",
			wantName:     "site.title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := setting.NewSettingKey(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.input, key.String())
			assert.Equal(t, tt.wantCategory, key.Category())
			assert.Equal(t, tt.wantName, key.Name())
			assert.False(t, key.IsEmpty())
		})
	}
}

func TestNewSettingKey_InvalidFormats(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"空字符串", ""},
		{"无分隔符", "noseparator"},
		{"以点号开头", ".invalid"},
		{"以点号结尾", "invalid."},
		{"仅点号", "."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := setting.NewSettingKey(tt.input)
			assert.ErrorIs(t, err, setting.ErrInvalidKeyFormat)
		})
	}
}

func TestMustSettingKey_ValidKey(t *testing.T) {
	assert.NotPanics(t, func() {
		key := setting.MustSettingKey("general.site_name")
		assert.Equal(t, "general.site_name", key.String())
	})
}

func TestMustSettingKey_InvalidKey(t *testing.T) {
	assert.Panics(t, func() {
		setting.MustSettingKey("invalid")
	})
}

func TestSettingKey_Equal(t *testing.T) {
	key1 := setting.MustSettingKey("general.site_name")
	key2 := setting.MustSettingKey("general.site_name")
	key3 := setting.MustSettingKey("general.site_title")

	assert.True(t, key1.Equal(key2))
	assert.False(t, key1.Equal(key3))
}

func TestSettingKey_IsEmpty(t *testing.T) {
	var emptyKey setting.SettingKey
	validKey := setting.MustSettingKey("general.site_name")

	assert.True(t, emptyKey.IsEmpty())
	assert.False(t, validKey.IsEmpty())
}

// =============================================================================
// Category Tests
// =============================================================================

func TestNewCategory_ValidCategories(t *testing.T) {
	validCats := []string{"general", "security", "notification", "backup"}

	for _, cat := range validCats {
		t.Run(cat, func(t *testing.T) {
			c, err := setting.NewCategory(cat)
			require.NoError(t, err)
			assert.Equal(t, cat, c.String())
			assert.True(t, c.IsValid())
			assert.False(t, c.IsEmpty())
		})
	}
}

func TestNewCategory_InvalidCategories(t *testing.T) {
	invalidCats := []string{"", "unknown", "General", "SECURITY", "test"}

	for _, cat := range invalidCats {
		t.Run(cat, func(t *testing.T) {
			_, err := setting.NewCategory(cat)
			assert.ErrorIs(t, err, setting.ErrCategoryNotFound)
		})
	}
}

func TestMustCategory_ValidCategory(t *testing.T) {
	assert.NotPanics(t, func() {
		c := setting.MustCategory("general")
		assert.Equal(t, "general", c.String())
	})
}

func TestMustCategory_InvalidCategory(t *testing.T) {
	assert.Panics(t, func() {
		setting.MustCategory("invalid")
	})
}

func TestCategory_Equal(t *testing.T) {
	cat1 := setting.MustCategory("general")
	cat2 := setting.MustCategory("general")
	cat3 := setting.MustCategory("security")

	assert.True(t, cat1.Equal(cat2))
	assert.False(t, cat1.Equal(cat3))
}

func TestCategory_IsEmpty(t *testing.T) {
	var emptyCat setting.Category
	validCat := setting.MustCategory("general")

	assert.True(t, emptyCat.IsEmpty())
	assert.False(t, validCat.IsEmpty())
}

func TestCategory_IsValid(t *testing.T) {
	var emptyCat setting.Category
	validCat := setting.MustCategory("general")

	assert.False(t, emptyCat.IsValid())
	assert.True(t, validCat.IsValid())
}

func TestAllCategoryStrings(t *testing.T) {
	cats := setting.AllCategoryStrings()
	expected := []string{"general", "security", "notification", "backup"}

	assert.ElementsMatch(t, expected, cats)
}

// =============================================================================
// Integration: SettingKey + Category
// =============================================================================

func TestSettingKey_CategoryExtraction(t *testing.T) {
	tests := []struct {
		key          string
		wantCategory string
		validCat     bool
	}{
		{"general.site_name", "general", true},
		{"security.password_min", "security", true},
		{"notification.email_enabled", "notification", true},
		{"backup.retention_days", "backup", true},
		{"unknown.some_key", "unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			key := setting.MustSettingKey(tt.key)
			catStr := key.Category()
			assert.Equal(t, tt.wantCategory, catStr)

			cat, err := setting.NewCategory(catStr)
			if tt.validCat {
				require.NoError(t, err)
				assert.True(t, cat.IsValid())
			} else {
				assert.Error(t, err)
			}
		})
	}
}
