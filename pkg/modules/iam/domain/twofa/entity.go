package twofa

import (
	"fmt"
	"time"
)

// TwoFA 用户双因素认证配置实体
type TwoFA struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// 用户关联
	UserID uint `json:"user_id"`

	// 2FA 状态
	Enabled bool `json:"enabled"` // 是否启用 2FA

	// TOTP 密钥（加密存储）
	Secret string `json:"-"` // TOTP 密钥（Base32 编码），不序列化

	// 恢复码（加密存储，JSON 数组）
	RecoveryCodes RecoveryCodes `json:"-"` // 恢复码列表，不序列化

	// 设置信息
	SetupCompletedAt *time.Time `json:"setup_completed_at,omitempty"` // 完成设置的时间
	LastUsedAt       *time.Time `json:"last_used_at,omitempty"`       // 最后使用时间
}

// HasRecoveryCodes 检查是否有可用的恢复码
func (t *TwoFA) HasRecoveryCodes() bool {
	return len(t.RecoveryCodes) > 0
}

// IsEnabled 检查 2FA 是否已启用
func (t *TwoFA) IsEnabled() bool {
	return t.Enabled
}

// IsSetupComplete 检查 2FA 设置是否已完成
func (t *TwoFA) IsSetupComplete() bool {
	return t.SetupCompletedAt != nil
}

// Enable 启用 2FA
func (t *TwoFA) Enable() {
	t.Enabled = true
	now := time.Now()
	t.SetupCompletedAt = &now
}

// Disable 禁用 2FA
func (t *TwoFA) Disable() {
	t.Enabled = false
}

// MarkUsed 标记最后使用时间
func (t *TwoFA) MarkUsed() {
	now := time.Now()
	t.LastUsedAt = &now
}

// UseRecoveryCode 使用恢复码（从列表中移除已使用的码）
// 返回 true 表示恢复码有效并已使用，false 表示无效
func (t *TwoFA) UseRecoveryCode(code string) bool {
	for i, rc := range t.RecoveryCodes {
		if rc == code {
			// 移除已使用的恢复码
			t.RecoveryCodes = append(t.RecoveryCodes[:i], t.RecoveryCodes[i+1:]...)
			t.MarkUsed()
			return true
		}
	}
	return false
}

// GetRecoveryCodesCount 获取剩余恢复码数量
func (t *TwoFA) GetRecoveryCodesCount() int {
	return len(t.RecoveryCodes)
}

// SetRecoveryCodes 设置恢复码（覆盖现有的）
func (t *TwoFA) SetRecoveryCodes(codes []string) {
	t.RecoveryCodes = codes
}

// HasSecret 检查是否已配置 TOTP 密钥
func (t *TwoFA) HasSecret() bool {
	return t.Secret != ""
}

// ClearSecret 清除 TOTP 密钥（禁用时使用）
func (t *TwoFA) ClearSecret() {
	t.Secret = ""
}

// Reset 重置 2FA 配置到初始状态
func (t *TwoFA) Reset() {
	t.Enabled = false
	t.Secret = ""
	t.RecoveryCodes = nil
	t.SetupCompletedAt = nil
	t.LastUsedAt = nil
}

// GenerateRecoveryCodes 生成恢复码
// 格式：xxxx-xxxx（8位数字，用连字符分隔）
// 参数 count 指定生成的恢复码数量
func GenerateRecoveryCodes(count int, randReader func([]byte) (int, error)) (RecoveryCodes, error) {
	codes := make(RecoveryCodes, count)

	for i := range count {
		// 生成 8 位随机数字
		b := make([]byte, 4)
		if _, err := randReader(b); err != nil {
			return nil, err
		}

		// 转换为 8 位数字
		num := uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
		num %= 100000000 // 限制在 8 位数字

		// 格式化为 xxxx-xxxx
		codes[i] = formatRecoveryCode(num)
	}

	return codes, nil
}

// formatRecoveryCode 将数字格式化为恢复码格式
func formatRecoveryCode(num uint32) string {
	return fmt.Sprintf("%04d-%04d", num/10000, num%10000)
}
