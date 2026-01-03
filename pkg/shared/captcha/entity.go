package captcha

import (
	"strings"
	"time"
)

// CaptchaData 验证码数据实体。
// 用于存储验证码的值和过期时间（内存存储，无 GORM 标签）。
type CaptchaData struct {
	Code      string    // 验证码值（小写，不区分大小写）
	ExpireAt  time.Time // 过期时间
	CreatedAt time.Time // 创建时间
}

// IsExpired 检查验证码是否过期
func (c *CaptchaData) IsExpired() bool {
	return time.Now().After(c.ExpireAt)
}

// IsValid 检查验证码是否有效（未过期且有值）
func (c *CaptchaData) IsValid() bool {
	return c.Code != "" && !c.IsExpired()
}

// GetTimeToExpire 返回距离过期的剩余时间。
// 如果已过期，返回负数。
func (c *CaptchaData) GetTimeToExpire() time.Duration {
	return time.Until(c.ExpireAt)
}

// Verify 验证输入的验证码是否正确（不区分大小写）。
// 返回 true 表示验证码正确且未过期。
func (c *CaptchaData) Verify(input string) bool {
	if c.IsExpired() {
		return false
	}
	return strings.EqualFold(c.Code, input)
}

// HasExpired 检查是否已经过期（同 IsExpired，语义更清晰）
func (c *CaptchaData) HasExpired() bool {
	return c.IsExpired()
}

// GetAge 返回验证码已存在的时间
func (c *CaptchaData) GetAge() time.Duration {
	return time.Since(c.CreatedAt)
}
