package cache_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application/cache"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/platform/manualtest"
)

// TestCacheAPI 测试缓存管理 API
func TestCacheAPI(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 测试 1: 获取缓存信息
	t.Run("Info", func(t *testing.T) {
		info, err := manualtest.Get[cache.CacheInfoDTO](c, "/api/admin/cache/info", nil)
		require.NoError(t, err, "获取缓存信息失败")

		assert.NotEmpty(t, info.KeyPrefix, "KeyPrefix 不应为空")
		assert.GreaterOrEqual(t, info.DBSize, int64(0), "DBSize 不应为负数")
		t.Logf("缓存信息: DBSize=%d, KeyPrefix=%s, Version=%s", info.DBSize, info.KeyPrefix, info.RedisVersion)
	})

	// 测试 2: 扫描所有 Keys
	t.Run("ScanAllKeys", func(t *testing.T) {
		result, err := manualtest.Get[cache.ScanKeysResultDTO](c, "/api/admin/cache/keys", nil)
		require.NoError(t, err, "扫描 Keys 失败")

		t.Logf("扫描到 %d 个 Keys, Cursor=%s", len(result.Keys), result.Cursor)
		for i, k := range result.Keys {
			if i < 5 { // 只打印前 5 个
				t.Logf("  - %s (type=%s, ttl=%d)", k.Key, k.Type, k.TTL)
			}
		}
	})

	// 测试 3: 按 pattern 扫描 Keys
	t.Run("ScanSettingKeys", func(t *testing.T) {
		result, err := manualtest.Get[cache.ScanKeysResultDTO](c, "/api/admin/cache/keys", map[string]string{
			"pattern": "setting:*",
		})
		require.NoError(t, err, "扫描 Setting Keys 失败")

		t.Logf("扫描到 %d 个 Setting Keys", len(result.Keys))
		for _, k := range result.Keys {
			t.Logf("  - %s (type=%s, ttl=%d)", k.Key, k.Type, k.TTL)
		}
	})

	// 测试 4: 获取单个 Key 的值
	t.Run("GetKey", func(t *testing.T) {
		// 先扫描获取一个 key（任意类型）
		scanResult, err := manualtest.Get[cache.ScanKeysResultDTO](c, "/api/admin/cache/keys", map[string]string{
			"limit": "1",
		})
		require.NoError(t, err, "扫描 Keys 失败")

		if len(scanResult.Keys) == 0 {
			t.Skip("没有任何缓存，跳过测试")
		}

		key := scanResult.Keys[0].Key
		t.Logf("获取 Key: %s", key)

		// 获取 key 的值（使用查询参数）
		value, err := manualtest.Get[cache.CacheValueDTO](c, "/api/admin/cache/key", map[string]string{
			"key": key,
		})
		require.NoError(t, err, "获取 Key 值失败")

		assert.Equal(t, key, value.Key)
		assert.NotEmpty(t, value.Type)
		assert.NotEmpty(t, value.Value)
		t.Logf("Key 值: type=%s, ttl=%d, value=%s", value.Type, value.TTL, string(value.Value))
	})

	// 测试 5: 删除单个 Key
	t.Run("DeleteSingleKey", func(t *testing.T) {
		// 扫描 schema 缓存（这些是安全可删除的）
		scanResult, err := manualtest.Get[cache.ScanKeysResultDTO](c, "/api/admin/cache/keys", map[string]string{
			"pattern": "schema:*",
			"limit":   "1",
		})
		require.NoError(t, err, "扫描 Schema Keys 失败")

		if len(scanResult.Keys) == 0 {
			t.Skip("没有 Schema 缓存，跳过删除测试")
		}

		key := scanResult.Keys[0].Key
		t.Logf("删除 Key: %s", key)

		// 删除 key（使用查询参数）
		_, err = manualtest.Delete[cache.DeleteResultDTO](c, "/api/admin/cache/key", map[string]string{
			"key": key,
		})
		require.NoError(t, err, "删除 Key 失败")
		t.Logf("✓ 删除成功")

		// 验证 key 不存在
		_, err = manualtest.Get[cache.CacheValueDTO](c, "/api/admin/cache/key", map[string]string{
			"key": key,
		})
		require.Error(t, err, "Key 应该已被删除")
		t.Logf("✓ 验证 Key 已删除")
	})

	// 测试 6: 按 pattern 批量删除
	t.Run("DeleteByPattern", func(t *testing.T) {
		// 先查看 schema 缓存数量
		before, err := manualtest.Get[cache.ScanKeysResultDTO](c, "/api/admin/cache/keys", map[string]string{
			"pattern": "schema:*",
		})
		require.NoError(t, err, "扫描失败")
		t.Logf("删除前 Schema 缓存数量: %d", len(before.Keys))

		if len(before.Keys) == 0 {
			t.Skip("没有 Schema 缓存，跳过批量删除测试")
		}

		// 按 pattern 删除
		result, err := manualtest.Delete[cache.DeleteResultDTO](c, "/api/admin/cache/keys", map[string]string{
			"pattern": "schema:*",
		})
		require.NoError(t, err, "批量删除失败")
		t.Logf("删除了 %d 个 Keys", result.DeletedCount)

		// 验证已删除
		after, err := manualtest.Get[cache.ScanKeysResultDTO](c, "/api/admin/cache/keys", map[string]string{
			"pattern": "schema:*",
		})
		require.NoError(t, err, "扫描失败")
		assert.Empty(t, after.Keys, "Schema 缓存应该已被清空")
		t.Logf("✓ 验证 Schema 缓存已清空")
	})
}
