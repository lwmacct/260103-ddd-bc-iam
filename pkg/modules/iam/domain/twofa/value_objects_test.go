package twofa

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecoveryCodes_Scan(t *testing.T) {
	t.Run("Scan nil 值", func(t *testing.T) {
		var r RecoveryCodes
		err := r.Scan(nil)
		require.NoError(t, err)
		assert.Equal(t, RecoveryCodes{}, r)
	})

	t.Run("Scan 有效 JSON 字节", func(t *testing.T) {
		var r RecoveryCodes
		jsonData := []byte(`["1234-5678","8765-4321","1111-2222"]`)
		err := r.Scan(jsonData)
		require.NoError(t, err)
		assert.Equal(t, RecoveryCodes{"1234-5678", "8765-4321", "1111-2222"}, r)
	})

	t.Run("Scan 有效 JSON 字符串", func(t *testing.T) {
		var r RecoveryCodes
		jsonString := `["code1","code2"]`
		err := r.Scan(jsonString)
		require.NoError(t, err)
		assert.Equal(t, RecoveryCodes{"code1", "code2"}, r)
	})

	t.Run("Scan 空数组 JSON", func(t *testing.T) {
		var r RecoveryCodes
		jsonData := []byte(`[]`)
		err := r.Scan(jsonData)
		require.NoError(t, err)
		assert.Equal(t, RecoveryCodes{}, r)
	})

	t.Run("Scan 空对象 JSON", func(t *testing.T) {
		var r RecoveryCodes
		jsonData := []byte(`{}`)
		err := r.Scan(jsonData)
		require.NoError(t, err)
		assert.Equal(t, RecoveryCodes{}, r)
	})

	t.Run("Scan 空字节", func(t *testing.T) {
		var r RecoveryCodes
		err := r.Scan([]byte{})
		require.NoError(t, err)
		assert.Equal(t, RecoveryCodes{}, r)
	})

	t.Run("Scan 非支持类型返回错误", func(t *testing.T) {
		var r RecoveryCodes
		err := r.Scan(123)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal RecoveryCodes value")
	})

	t.Run("Scan 无效 JSON 返回错误", func(t *testing.T) {
		var r RecoveryCodes
		invalidJSON := []byte(`not valid json`)
		err := r.Scan(invalidJSON)
		assert.Error(t, err)
	})
}

func TestRecoveryCodes_Value(t *testing.T) {
	t.Run("Value 空列表", func(t *testing.T) {
		r := RecoveryCodes{}
		val, err := r.Value()
		require.NoError(t, err)

		expected, _ := json.Marshal([]string{})
		assert.Equal(t, expected, val)
	})

	t.Run("Value nil 列表", func(t *testing.T) {
		var r RecoveryCodes
		val, err := r.Value()
		require.NoError(t, err)

		expected, _ := json.Marshal([]string{})
		assert.Equal(t, expected, val)
	})

	t.Run("Value 有内容的列表", func(t *testing.T) {
		r := RecoveryCodes{"1234-5678", "8765-4321"}
		val, err := r.Value()
		require.NoError(t, err)

		expected, _ := json.Marshal([]string{"1234-5678", "8765-4321"})
		assert.Equal(t, expected, val)
	})

	t.Run("Value 单个恢复码", func(t *testing.T) {
		r := RecoveryCodes{"only-one"}
		val, err := r.Value()
		require.NoError(t, err)

		expected, _ := json.Marshal([]string{"only-one"})
		assert.Equal(t, expected, val)
	})
}
