package pat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPermissionList_Scan(t *testing.T) {
	t.Run("Scan nil 值", func(t *testing.T) {
		var p PermissionList
		err := p.Scan(nil)
		require.NoError(t, err)
		assert.Equal(t, PermissionList{}, p)
	})

	t.Run("Scan 有效 JSON", func(t *testing.T) {
		var p PermissionList
		jsonData := []byte(`["read","write","admin"]`)
		err := p.Scan(jsonData)
		require.NoError(t, err)
		assert.Equal(t, PermissionList{"read", "write", "admin"}, p)
	})

	t.Run("Scan 空数组 JSON", func(t *testing.T) {
		var p PermissionList
		jsonData := []byte(`[]`)
		err := p.Scan(jsonData)
		require.NoError(t, err)
		assert.Equal(t, PermissionList{}, p)
	})

	t.Run("Scan 非字节类型返回错误", func(t *testing.T) {
		var p PermissionList
		err := p.Scan("not bytes")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to scan PermissionList")
	})

	t.Run("Scan 无效 JSON 返回错误", func(t *testing.T) {
		var p PermissionList
		invalidJSON := []byte(`{invalid}`)
		err := p.Scan(invalidJSON)
		assert.Error(t, err)
	})
}

func TestPermissionList_Value(t *testing.T) { //nolint:dupl // 测试代码结构相似是可接受的
	t.Run("Value nil 列表", func(t *testing.T) {
		var p PermissionList
		val, err := p.Value()
		require.NoError(t, err)

		// 验证返回的是空数组 JSON
		expected, _ := json.Marshal([]string{})
		assert.Equal(t, expected, val)
	})

	t.Run("Value 有内容的列表", func(t *testing.T) {
		p := PermissionList{"read", "write"}
		val, err := p.Value()
		require.NoError(t, err)

		expected, _ := json.Marshal([]string{"read", "write"})
		assert.Equal(t, expected, val)
	})

	t.Run("Value 空列表", func(t *testing.T) {
		p := PermissionList{}
		val, err := p.Value()
		require.NoError(t, err)

		expected, _ := json.Marshal([]string{})
		assert.Equal(t, expected, val)
	})
}

func TestStringList_Scan(t *testing.T) {
	t.Run("Scan nil 值", func(t *testing.T) {
		var s StringList
		err := s.Scan(nil)
		require.NoError(t, err)
		assert.Equal(t, StringList{}, s)
	})

	t.Run("Scan 有效 JSON", func(t *testing.T) {
		var s StringList
		jsonData := []byte(`["192.168.1.1","10.0.0.1"]`)
		err := s.Scan(jsonData)
		require.NoError(t, err)
		assert.Equal(t, StringList{"192.168.1.1", "10.0.0.1"}, s)
	})

	t.Run("Scan 空数组 JSON", func(t *testing.T) {
		var s StringList
		jsonData := []byte(`[]`)
		err := s.Scan(jsonData)
		require.NoError(t, err)
		assert.Equal(t, StringList{}, s)
	})

	t.Run("Scan 非字节类型返回错误", func(t *testing.T) {
		var s StringList
		err := s.Scan(123)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to scan StringList")
	})

	t.Run("Scan 无效 JSON 返回错误", func(t *testing.T) {
		var s StringList
		invalidJSON := []byte(`not valid json`)
		err := s.Scan(invalidJSON)
		assert.Error(t, err)
	})
}

func TestStringList_Value(t *testing.T) { //nolint:dupl // 测试代码结构相似是可接受的
	t.Run("Value nil 列表", func(t *testing.T) {
		var s StringList
		val, err := s.Value()
		require.NoError(t, err)

		expected, _ := json.Marshal([]string{})
		assert.Equal(t, expected, val)
	})

	t.Run("Value 有内容的列表", func(t *testing.T) {
		s := StringList{"192.168.1.1", "10.0.0.1"}
		val, err := s.Value()
		require.NoError(t, err)

		expected, _ := json.Marshal([]string{"192.168.1.1", "10.0.0.1"})
		assert.Equal(t, expected, val)
	})

	t.Run("Value 空列表", func(t *testing.T) {
		s := StringList{}
		val, err := s.Value()
		require.NoError(t, err)

		expected, _ := json.Marshal([]string{})
		assert.Equal(t, expected, val)
	})
}
