package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	domainauth "github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/auth"
)

const (
	// SessionTokenExpiration 会话token过期时间（5分钟）
	SessionTokenExpiration = 5 * time.Minute
)

// loginSessionData 登录会话数据（内部存储）
type loginSessionData struct {
	UserID    uint      // 用户ID
	Account   string    // 登录账号
	CreatedAt time.Time // 创建时间
	ExpireAt  time.Time // 过期时间
}

// IsExpired 检查是否过期
func (s *loginSessionData) IsExpired() bool {
	return time.Now().After(s.ExpireAt)
}

// loginSessionService 登录会话服务
// 用于 2FA 验证流程中的临时会话管理
// 🔒 安全策略：防止 2FA 暴力破解
type loginSessionService struct {
	sessions  map[string]*loginSessionData
	mu        sync.RWMutex
	stopClean chan struct{}
}

// NewLoginSessionService 创建登录会话服务
func NewLoginSessionService() domainauth.SessionService {
	service := &loginSessionService{
		sessions:  make(map[string]*loginSessionData),
		stopClean: make(chan struct{}),
	}

	// 启动定期清理协程
	go service.cleanupExpired()

	return service
}

// GenerateSessionToken 生成会话token
func (s *loginSessionService) GenerateSessionToken(ctx context.Context, userID uint, account string) (string, error) {
	// 生成随机token（32字节，hex编码后64个字符）
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	token := hex.EncodeToString(b)

	s.mu.Lock()
	defer s.mu.Unlock()

	// 存储会话数据
	now := time.Now()
	s.sessions[token] = &loginSessionData{
		UserID:    userID,
		Account:   account,
		CreatedAt: now,
		ExpireAt:  now.Add(SessionTokenExpiration),
	}

	return token, nil
}

// VerifySessionToken 验证会话token
// 验证后自动删除token（一次性使用）
func (s *loginSessionService) VerifySessionToken(ctx context.Context, token string) (*domainauth.SessionData, error) {
	if token == "" {
		return nil, errors.New("session token is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 获取会话数据
	sessionData, exists := s.sessions[token]
	if !exists {
		return nil, errors.New("invalid or expired session token")
	}

	// 检查是否过期
	if sessionData.IsExpired() {
		delete(s.sessions, token)
		return nil, errors.New("session token expired")
	}

	// 验证成功后删除token（一次性使用）
	delete(s.sessions, token)

	// 转换为 Domain 层结构
	return &domainauth.SessionData{
		UserID:    sessionData.UserID,
		Account:   sessionData.Account,
		CreatedAt: sessionData.CreatedAt,
		ExpireAt:  sessionData.ExpireAt,
	}, nil
}

// Close 关闭服务（停止清理协程）
func (s *loginSessionService) Close() error {
	close(s.stopClean)
	return nil
}

// cleanupExpired 定期清理过期会话
func (s *loginSessionService) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			for token, data := range s.sessions {
				if data.IsExpired() {
					delete(s.sessions, token)
				}
			}
			s.mu.Unlock()
		case <-s.stopClean:
			return
		}
	}
}

var _ domainauth.SessionService = (*loginSessionService)(nil)
