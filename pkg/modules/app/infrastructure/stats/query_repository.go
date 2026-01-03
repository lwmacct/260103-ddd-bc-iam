package stats

import (
	"fmt"

	domain "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/stats"
	"gorm.io/gorm"
)

// queryRepository 统计查询仓储的 GORM 实现
type queryRepository struct {
	db *gorm.DB
}

// NewQueryRepository 创建统计查询仓储实例
func NewQueryRepository(db *gorm.DB) domain.QueryRepository {
	return &queryRepository{db: db}
}

// GetSystemStats 获取系统统计信息。
//
// TODO: 添加缓存支持，减少高频查询的数据库压力
// 建议缓存策略：TTL 5-10 分钟，按需失效
func (r *queryRepository) GetSystemStats(recentLogsLimit int) (*domain.SystemStats, error) {
	s := &domain.SystemStats{}

	// 统计用户数
	total, err := r.GetTotalUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get total users: %w", err)
	}
	s.TotalUsers = total

	active, err := r.GetUserCountByStatus("active")
	if err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}
	s.ActiveUsers = active

	inactive, err := r.GetUserCountByStatus("inactive")
	if err != nil {
		return nil, fmt.Errorf("failed to get inactive users: %w", err)
	}
	s.InactiveUsers = inactive

	banned, err := r.GetUserCountByStatus("banned")
	if err != nil {
		return nil, fmt.Errorf("failed to get banned users: %w", err)
	}
	s.BannedUsers = banned

	// 统计角色
	roles, err := r.GetTotalRoles()
	if err != nil {
		return nil, fmt.Errorf("failed to get total roles: %w", err)
	}
	s.TotalRoles = roles

	// 统计权限
	permissions, err := r.GetTotalPermissions()
	if err != nil {
		return nil, fmt.Errorf("failed to get total permissions: %w", err)
	}
	s.TotalPermissions = permissions

	// 获取最近审计日志
	logs, err := r.GetRecentAuditLogs(recentLogsLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent audit logs: %w", err)
	}
	s.RecentAuditLogs = logs

	return s, nil
}

// GetUserCountByStatus 按状态统计用户数量
func (r *queryRepository) GetUserCountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Table("users").
		Where("deleted_at IS NULL AND status = ?", status).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTotalUsers 获取用户总数
func (r *queryRepository) GetTotalUsers() (int64, error) {
	var count int64
	err := r.db.Table("users").
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTotalRoles 获取角色总数
func (r *queryRepository) GetTotalRoles() (int64, error) {
	var count int64
	err := r.db.Table("roles").
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTotalPermissions 获取权限总数（统计所有角色的 permissions JSONB 数组条目总数）
func (r *queryRepository) GetTotalPermissions() (int64, error) {
	var count int64
	// 新 RBAC 模型：权限存储在 roles.permissions JSONB 字段中
	// 使用 PostgreSQL jsonb_array_length 统计每个角色的权限数，然后求和
	err := r.db.Table("roles").
		Where("deleted_at IS NULL").
		Select("COALESCE(SUM(jsonb_array_length(permissions)), 0)").
		Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetRecentAuditLogs 获取最近的审计日志
func (r *queryRepository) GetRecentAuditLogs(limit int) ([]domain.AuditLogSummary, error) {
	var logs []domain.AuditLogSummary
	err := r.db.Table("audit").
		Select("id, user_id, username, action, resource, status, created_at").
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(limit).
		Scan(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}
