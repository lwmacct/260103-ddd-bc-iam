package twofa

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"

	domainTwoFA "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/twofa"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/user"
)

// Service 2FA 服务
type Service struct {
	twofaCommandRepo domainTwoFA.CommandRepository
	twofaQueryRepo   domainTwoFA.QueryRepository
	userQueryRepo    user.QueryRepository
	issuer           string // TOTP 发行者名称
}

// NewService 创建 2FA 服务
func NewService(twofaCommandRepo domainTwoFA.CommandRepository, twofaQueryRepo domainTwoFA.QueryRepository, userQueryRepo user.QueryRepository, issuer string) *Service {
	if issuer == "" {
		issuer = "Go-DDD-Package-Lib"
	}
	return &Service{
		twofaCommandRepo: twofaCommandRepo,
		twofaQueryRepo:   twofaQueryRepo,
		userQueryRepo:    userQueryRepo,
		issuer:           issuer,
	}
}

// SetupResponse 设置 2FA 返回数据
//
// Deprecated: 请使用 domainTwoFA.SetupResult
type SetupResponse = domainTwoFA.SetupResult

// EnableResponse 启用 2FA 返回数据
type EnableResponse struct {
	RecoveryCodes []string `json:"recovery_codes"` // 恢复码列表
	Message       string   `json:"message"`        // 提示消息
}

// StatusResponse 2FA 状态返回数据
type StatusResponse struct {
	Enabled            bool `json:"enabled"`              // 是否启用
	RecoveryCodesCount int  `json:"recovery_codes_count"` // 剩余恢复码数量
}

// Setup 设置 2FA（生成密钥和二维码）
func (s *Service) Setup(ctx context.Context, userID uint) (*domainTwoFA.SetupResult, error) {
	// 查找用户
	u, err := s.userQueryRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 生成 TOTP 密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuer,
		AccountName: u.Username, // 使用用户名
		SecretSize:  10,         // 10 字节（80 位），更短的 Base32 密钥，便于手动输入
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	// 构建简化的二维码 URL（移除默认参数）
	secret := key.Secret()
	qrcodeURL := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		url.QueryEscape(s.issuer),
		url.QueryEscape(u.Username),
		secret,
		url.QueryEscape(s.issuer),
	)

	// 生成二维码图片
	qrCodeBytes, err := qrcode.Encode(qrcodeURL, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Base64 编码
	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCodeBytes)

	// 存储密钥到数据库（未启用状态）
	tfa := &domainTwoFA.TwoFA{
		UserID:        userID,
		Enabled:       false,
		Secret:        secret,
		RecoveryCodes: domainTwoFA.RecoveryCodes{}, // 空恢复码，验证后生成
	}

	if err := s.twofaCommandRepo.CreateOrUpdate(ctx, tfa); err != nil {
		return nil, fmt.Errorf("failed to save 2FA secret: %w", err)
	}

	return &domainTwoFA.SetupResult{
		Secret:    secret,
		QRCodeURL: qrcodeURL,
		QRCodeImg: "data:image/png;base64," + qrCodeBase64,
	}, nil
}

// VerifyAndEnable 验证 TOTP 代码并启用 2FA
// 返回恢复码列表
func (s *Service) VerifyAndEnable(ctx context.Context, userID uint, code string) ([]string, error) {
	// 查找 2FA 配置
	tfa, err := s.twofaQueryRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get 2FA config: %w", err)
	}

	if tfa == nil || !tfa.HasSecret() {
		return nil, errors.New("please setup 2FA first")
	}

	// 验证 TOTP 代码
	valid := totp.Validate(code, tfa.Secret)
	if !valid {
		return nil, errors.New("invalid verification code")
	}

	// 使用领域函数生成恢复码（8 个）
	recoveryCodes, err := domainTwoFA.GenerateRecoveryCodes(8, rand.Read)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recovery codes: %w", err)
	}

	// 使用实体方法启用 2FA 并设置恢复码
	tfa.Enable()
	tfa.SetRecoveryCodes(recoveryCodes)

	if err := s.twofaCommandRepo.CreateOrUpdate(ctx, tfa); err != nil {
		return nil, fmt.Errorf("failed to enable 2FA: %w", err)
	}

	return recoveryCodes, nil
}

// Verify 验证 TOTP 代码或恢复码
func (s *Service) Verify(ctx context.Context, userID uint, code string) (bool, error) {
	// 查找 2FA 配置
	tfa, err := s.twofaQueryRepo.FindByUserID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get 2FA config: %w", err)
	}

	if tfa == nil || !tfa.IsEnabled() {
		return false, errors.New("2FA not enabled")
	}

	// 首先尝试 TOTP 验证
	if totp.Validate(code, tfa.Secret) {
		// 使用实体方法更新最后使用时间
		tfa.MarkUsed()
		_ = s.twofaCommandRepo.CreateOrUpdate(ctx, tfa)
		return true, nil
	}

	// TOTP 验证失败，尝试恢复码（使用实体方法）
	code = strings.TrimSpace(code)
	if tfa.UseRecoveryCode(code) {
		if err := s.twofaCommandRepo.CreateOrUpdate(ctx, tfa); err != nil {
			return false, fmt.Errorf("failed to update recovery codes: %w", err)
		}
		return true, nil
	}

	return false, nil
}

// Disable 禁用 2FA
func (s *Service) Disable(ctx context.Context, userID uint) error {
	return s.twofaCommandRepo.Delete(ctx, userID)
}

// GetStatus 获取 2FA 状态
func (s *Service) GetStatus(ctx context.Context, userID uint) (bool, int, error) {
	tfa, err := s.twofaQueryRepo.FindByUserID(ctx, userID)
	if err != nil {
		return false, 0, fmt.Errorf("failed to get 2FA config: %w", err)
	}

	if tfa == nil {
		return false, 0, nil
	}

	return tfa.IsEnabled(), tfa.GetRecoveryCodesCount(), nil
}

var _ domainTwoFA.Service = (*Service)(nil)
