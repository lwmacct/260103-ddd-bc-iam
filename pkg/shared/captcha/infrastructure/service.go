package infrastructure

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/shared/captcha"
)

// Service 验证码生成服务实现
type Service struct {
	defaultExpiration int64 // 秒
}

// NewService 创建验证码服务
func NewService() captcha.Service {
	return &Service{
		defaultExpiration: 300, // 5分钟
	}
}

// GenerateRandomCode 生成随机验证码
// 返回 (captchaID, imageBase64, code, error)
func (s *Service) GenerateRandomCode() (string, string, string, error) {
	// 生成6位随机数字验证码
	code, err := s.generateRandomCode(6)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate code: %w", err)
	}

	// 生成验证码ID
	captchaID, err := s.generateCaptchaID()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate captcha ID: %w", err)
	}

	// 生成Base64图片（简化实现，实际应使用图片生成库）
	imageBase64, err := s.generateImageBase64(code)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate image: %w", err)
	}

	return captchaID, imageBase64, code, nil
}

// GenerateCustomCodeImage 生成指定文本的验证码图片（用于开发模式）
func (s *Service) GenerateCustomCodeImage(text string) (string, error) {
	return s.generateImageBase64(text)
}

// GenerateDevCaptchaID 生成开发模式验证码ID
func (s *Service) GenerateDevCaptchaID() string {
	id, _ := s.generateCaptchaID()
	return id
}

// GetDefaultExpiration 获取默认过期时间
func (s *Service) GetDefaultExpiration() int64 {
	return s.defaultExpiration
}

// ============================================================================
// 私有辅助方法
// ============================================================================

// generateRandomCode 生成随机数字验证码
func (s *Service) generateRandomCode(length int) (string, error) {
	const digits = "0123456789"

	code := make([]byte, length)
	for i := range length {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[num.Int64()]
	}

	return string(code), nil
}

// generateCaptchaID 生成验证码唯一ID
func (s *Service) generateCaptchaID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// generateImageBase64 生成验证码图片的Base64编码
// 简化实现：返回纯文本的Base64（生产环境应使用图片生成库）
func (s *Service) generateImageBase64(code string) (string, error) {
	// TODO: 实现真实的图片生成
	// 当前返回纯文本的Base64，前端可直接显示
	// 生产环境应使用 github.com/mojocn/base64Captcha 或类似库
	text := "CAPTCHA:" + code
	return base64.StdEncoding.EncodeToString([]byte(text)), nil
}

// ============================================================================
// 确保实现接口
// ============================================================================

var _ captcha.Service = (*Service)(nil)
