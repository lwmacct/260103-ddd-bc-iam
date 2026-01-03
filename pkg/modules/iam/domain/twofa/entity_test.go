package twofa

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTestTwoFA(enabled bool) *TwoFA {
	return &TwoFA{
		ID:      1,
		UserID:  100,
		Enabled: enabled,
		Secret:  "JBSWY3DPEHPK3PXP",
	}
}

func TestTwoFA_IsEnabled(t *testing.T) {
	t.Run("已启用", func(t *testing.T) {
		tfa := newTestTwoFA(true)
		assert.True(t, tfa.IsEnabled())
	})

	t.Run("未启用", func(t *testing.T) {
		tfa := newTestTwoFA(false)
		assert.False(t, tfa.IsEnabled())
	})
}

func TestTwoFA_IsSetupComplete(t *testing.T) {
	t.Run("设置已完成", func(t *testing.T) {
		tfa := newTestTwoFA(true)
		now := time.Now()
		tfa.SetupCompletedAt = &now
		assert.True(t, tfa.IsSetupComplete())
	})

	t.Run("设置未完成", func(t *testing.T) {
		tfa := newTestTwoFA(false)
		assert.False(t, tfa.IsSetupComplete())
	})
}

func TestTwoFA_Enable(t *testing.T) {
	tfa := newTestTwoFA(false)
	assert.False(t, tfa.Enabled)
	assert.Nil(t, tfa.SetupCompletedAt)

	tfa.Enable()

	assert.True(t, tfa.Enabled)
	assert.NotNil(t, tfa.SetupCompletedAt)
}

func TestTwoFA_Disable(t *testing.T) {
	tfa := newTestTwoFA(true)
	tfa.Disable()
	assert.False(t, tfa.Enabled)
}

func TestTwoFA_MarkUsed(t *testing.T) {
	tfa := newTestTwoFA(true)
	assert.Nil(t, tfa.LastUsedAt)

	tfa.MarkUsed()

	assert.NotNil(t, tfa.LastUsedAt)
}

func TestTwoFA_RecoveryCodes(t *testing.T) {
	t.Run("HasRecoveryCodes - 有恢复码", func(t *testing.T) {
		tfa := newTestTwoFA(true)
		tfa.RecoveryCodes = RecoveryCodes{"code1", "code2", "code3"}
		assert.True(t, tfa.HasRecoveryCodes())
	})

	t.Run("HasRecoveryCodes - 无恢复码", func(t *testing.T) {
		tfa := newTestTwoFA(true)
		assert.False(t, tfa.HasRecoveryCodes())
	})

	t.Run("GetRecoveryCodesCount", func(t *testing.T) {
		tfa := newTestTwoFA(true)
		tfa.RecoveryCodes = RecoveryCodes{"a", "b", "c"}
		assert.Equal(t, 3, tfa.GetRecoveryCodesCount())
	})

	t.Run("SetRecoveryCodes", func(t *testing.T) {
		tfa := newTestTwoFA(true)
		codes := []string{"new1", "new2"}
		tfa.SetRecoveryCodes(codes)
		assert.Equal(t, RecoveryCodes{"new1", "new2"}, tfa.RecoveryCodes)
	})
}

func TestTwoFA_UseRecoveryCode(t *testing.T) {
	t.Run("成功使用恢复码", func(t *testing.T) {
		tfa := newTestTwoFA(true)
		tfa.RecoveryCodes = RecoveryCodes{"code1", "code2", "code3"}

		used := tfa.UseRecoveryCode("code2")

		assert.True(t, used)
		assert.Equal(t, 2, tfa.GetRecoveryCodesCount())
		assert.NotContains(t, tfa.RecoveryCodes, "code2")
		assert.NotNil(t, tfa.LastUsedAt)
	})

	t.Run("使用无效恢复码", func(t *testing.T) {
		tfa := newTestTwoFA(true)
		tfa.RecoveryCodes = RecoveryCodes{"code1", "code2"}

		used := tfa.UseRecoveryCode("invalid")

		assert.False(t, used)
		assert.Equal(t, 2, tfa.GetRecoveryCodesCount())
	})

	t.Run("恢复码列表为空", func(t *testing.T) {
		tfa := newTestTwoFA(true)
		used := tfa.UseRecoveryCode("anycode")
		assert.False(t, used)
	})

	t.Run("使用第一个恢复码", func(t *testing.T) {
		tfa := newTestTwoFA(true)
		tfa.RecoveryCodes = RecoveryCodes{"first", "second"}

		used := tfa.UseRecoveryCode("first")

		assert.True(t, used)
		assert.Equal(t, RecoveryCodes{"second"}, tfa.RecoveryCodes)
	})

	t.Run("使用最后一个恢复码", func(t *testing.T) {
		tfa := newTestTwoFA(true)
		tfa.RecoveryCodes = RecoveryCodes{"first", "last"}

		used := tfa.UseRecoveryCode("last")

		assert.True(t, used)
		assert.Equal(t, RecoveryCodes{"first"}, tfa.RecoveryCodes)
	})
}

func TestTwoFA_HasSecret(t *testing.T) {
	t.Run("有密钥", func(t *testing.T) {
		tfa := newTestTwoFA(true)
		assert.True(t, tfa.HasSecret())
	})

	t.Run("无密钥", func(t *testing.T) {
		tfa := &TwoFA{}
		assert.False(t, tfa.HasSecret())
	})
}

func TestTwoFA_ClearSecret(t *testing.T) {
	tfa := newTestTwoFA(true)
	assert.NotEmpty(t, tfa.Secret)

	tfa.ClearSecret()

	assert.Empty(t, tfa.Secret)
}

func TestTwoFA_Reset(t *testing.T) {
	now := time.Now()
	tfa := &TwoFA{
		ID:               1,
		UserID:           100,
		Enabled:          true,
		Secret:           "secret",
		RecoveryCodes:    RecoveryCodes{"code1", "code2"},
		SetupCompletedAt: &now,
		LastUsedAt:       &now,
	}

	tfa.Reset()

	assert.False(t, tfa.Enabled)
	assert.Empty(t, tfa.Secret)
	assert.Nil(t, tfa.RecoveryCodes)
	assert.Nil(t, tfa.SetupCompletedAt)
	assert.Nil(t, tfa.LastUsedAt)
	// ID 和 UserID 应该保留
	assert.Equal(t, uint(1), tfa.ID)
	assert.Equal(t, uint(100), tfa.UserID)
}
