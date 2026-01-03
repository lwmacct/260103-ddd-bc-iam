package auth

// ToTwoFARequiredDTO 创建需要 2FA 的响应
func ToTwoFARequiredDTO(sessionToken string) *TwoFARequiredDTO {
	return &TwoFARequiredDTO{
		Requires2FA:  true,
		SessionToken: sessionToken,
	}
}
