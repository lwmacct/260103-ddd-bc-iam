package auth

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/pat"
)

// PATService 提供 PAT（Personal Access Token）认证服务。
//
// 本服务专注于 PAT 的**认证和安全管理**，供中间件使用。
// PAT 的 CRUD 操作由 Application 层 Handler 负责。
//
// 职责划分：
//   - 认证：ValidateToken, ValidateTokenWithIP
//   - 安全管理：DeleteAllUserTokens（如密码重置时撤销所有 token）
//   - 系统维护：CleanupExpiredTokens（定期清理过期 token）
type PATService struct {
	patCommandRepo pat.CommandRepository
	patQueryRepo   pat.QueryRepository
	tokenGen       *TokenGenerator
}

// NewPATService 创建 PAT 认证服务实例
func NewPATService(
	patCommandRepo pat.CommandRepository,
	patQueryRepo pat.QueryRepository,
	tokenGen *TokenGenerator,
) *PATService {
	return &PATService{
		patCommandRepo: patCommandRepo,
		patQueryRepo:   patQueryRepo,
		tokenGen:       tokenGen,
	}
}

// ValidateToken 验证 PAT 并返回关联的 Token 实体。
//
// 本方法用于认证中间件，验证流程：
//  1. 校验 token 格式
//  2. 哈希 token 并查询数据库
//  3. 检查 token 状态（active/disabled/expired）
//  4. 异步更新 last_used_at
//
// 返回的 Token 实体包含 UserID 和 Permissions，供中间件使用。
func (s *PATService) ValidateToken(ctx context.Context, plainToken string) (*pat.PersonalAccessToken, error) {
	// Validate format
	if !s.tokenGen.ValidateTokenFormat(plainToken) {
		return nil, errors.New("invalid token format")
	}

	// Hash the token
	tokenHash := s.tokenGen.HashToken(plainToken)

	// Find token in database
	token, err := s.patQueryRepo.FindByToken(ctx, tokenHash)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Check if token is active
	if !token.IsActive() {
		return nil, errors.New("token is inactive or expired")
	}

	// Update last used time (asynchronously to avoid blocking)
	// 使用 WithoutCancel 保留 trace 信息，但不受请求取消影响
	go func(updateCtx context.Context) {
		now := time.Now()
		token.LastUsedAt = &now
		_ = s.patCommandRepo.Update(updateCtx, token)
	}(context.WithoutCancel(ctx))

	return token, nil
}

// ValidateTokenWithIP 验证 PAT 并检查 IP 白名单。
//
// 扩展 ValidateToken，增加 IP 白名单校验。
// 如果 Token 配置了 IP 白名单，且客户端 IP 不在白名单中，
// 将拒绝访问。
func (s *PATService) ValidateTokenWithIP(ctx context.Context, plainToken, clientIP string) (*pat.PersonalAccessToken, error) {
	token, err := s.ValidateToken(ctx, plainToken)
	if err != nil {
		return nil, err
	}

	// Check IP whitelist if configured
	if len(token.IPWhitelist) > 0 {
		if !slices.Contains(token.IPWhitelist, clientIP) {
			return nil, fmt.Errorf("access denied: IP %s not in whitelist", clientIP)
		}
	}

	return token, nil
}

// DeleteAllUserTokens 删除用户的所有 PAT。
//
// 安全场景使用，例如：
//   - 用户密码重置后撤销所有 token
//   - 账户被盗用时紧急撤销
//   - 用户主动注销所有会话
func (s *PATService) DeleteAllUserTokens(ctx context.Context, userID uint) error {
	return s.patCommandRepo.DeleteByUserID(ctx, userID)
}

// CleanupExpiredTokens 清理过期的 PAT。
//
// 系统维护使用，应通过定时任务周期性调用（如每天一次）。
// 清理已过期的 token 以释放数据库空间。
func (s *PATService) CleanupExpiredTokens(ctx context.Context) error {
	return s.patCommandRepo.CleanupExpired(ctx)
}
