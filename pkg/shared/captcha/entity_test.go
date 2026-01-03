package captcha

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// newTestCaptcha 创建测试用验证码
func newTestCaptcha(code string, expireIn time.Duration) *CaptchaData {
	now := time.Now()
	return &CaptchaData{
		Code:      code,
		ExpireAt:  now.Add(expireIn),
		CreatedAt: now,
	}
}

func TestCaptchaData_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		expireIn time.Duration
		want     bool
	}{
		{"未过期", 5 * time.Minute, false},
		{"已过期", -5 * time.Minute, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestCaptcha("abc123", tt.expireIn)
			assert.Equal(t, tt.want, c.IsExpired())
		})
	}
}

func TestCaptchaData_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expireIn time.Duration
		want     bool
	}{
		{"有效 - 有值且未过期", "abc123", 5 * time.Minute, true},
		{"无效 - 空值", "", 5 * time.Minute, false},
		{"无效 - 已过期", "abc123", -5 * time.Minute, false},
		{"无效 - 空值且过期", "", -5 * time.Minute, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestCaptcha(tt.code, tt.expireIn)
			assert.Equal(t, tt.want, c.IsValid())
		})
	}
}

func TestCaptchaData_GetTimeToExpire(t *testing.T) {
	t.Run("未过期 - 返回正数", func(t *testing.T) {
		c := newTestCaptcha("abc123", 5*time.Minute)
		ttl := c.GetTimeToExpire()
		assert.Positive(t, ttl, "TTL 应该是正数")
		assert.LessOrEqual(t, ttl, 5*time.Minute, "TTL 应该 <= 5 分钟")
	})

	t.Run("已过期 - 返回负数", func(t *testing.T) {
		c := newTestCaptcha("abc123", -5*time.Minute)
		ttl := c.GetTimeToExpire()
		assert.Negative(t, ttl, "TTL 应该是负数")
	})
}

func TestCaptchaData_Verify(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		input    string
		expireIn time.Duration
		want     bool
	}{
		{"正确 - 完全匹配", "abc123", "abc123", 5 * time.Minute, true},
		{"正确 - 不区分大小写", "ABC123", "abc123", 5 * time.Minute, true},
		{"正确 - 输入大写", "abc123", "ABC123", 5 * time.Minute, true},
		{"错误 - 值不匹配", "abc123", "xyz789", 5 * time.Minute, false},
		{"错误 - 已过期", "abc123", "abc123", -5 * time.Minute, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestCaptcha(tt.code, tt.expireIn)
			assert.Equal(t, tt.want, c.Verify(tt.input))
		})
	}
}

func TestCaptchaData_HasExpired(t *testing.T) {
	t.Run("未过期", func(t *testing.T) {
		c := newTestCaptcha("abc123", 5*time.Minute)
		assert.False(t, c.HasExpired())
	})

	t.Run("已过期", func(t *testing.T) {
		c := newTestCaptcha("abc123", -5*time.Minute)
		assert.True(t, c.HasExpired())
	})
}

func TestCaptchaData_GetAge(t *testing.T) {
	t.Run("刚创建的验证码", func(t *testing.T) {
		c := newTestCaptcha("abc123", 5*time.Minute)
		age := c.GetAge()
		// 刚创建，age 应该接近 0
		assert.GreaterOrEqual(t, age, time.Duration(0), "Age 应该 >= 0")
		assert.Less(t, age, 1*time.Second, "Age 应该 < 1 秒")
	})

	t.Run("过去创建的验证码", func(t *testing.T) {
		c := &CaptchaData{
			Code:      "abc123",
			ExpireAt:  time.Now().Add(5 * time.Minute),
			CreatedAt: time.Now().Add(-10 * time.Minute),
		}
		age := c.GetAge()
		assert.GreaterOrEqual(t, age, 10*time.Minute, "Age 应该 >= 10 分钟")
	})
}

func TestCaptchaData_EdgeCases(t *testing.T) {
	t.Run("刚刚过期", func(t *testing.T) {
		c := &CaptchaData{
			Code:      "abc123",
			ExpireAt:  time.Now().Add(-1 * time.Nanosecond),
			CreatedAt: time.Now().Add(-5 * time.Minute),
		}
		assert.True(t, c.IsExpired())
		assert.False(t, c.IsValid())
		assert.False(t, c.Verify("abc123"))
	})

	t.Run("零值时间", func(t *testing.T) {
		c := &CaptchaData{
			Code:      "abc123",
			ExpireAt:  time.Time{},
			CreatedAt: time.Time{},
		}
		// 零值时间是过去的
		assert.True(t, c.IsExpired())
	})
}
