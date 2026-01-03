package stats

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/stats"
)

// GetStatsHandler 获取统计信息的处理器
type GetStatsHandler struct {
	statsQueryRepo stats.QueryRepository
}

// NewGetStatsHandler 创建 GetStatsHandler 实例
func NewGetStatsHandler(statsQueryRepo stats.QueryRepository) *GetStatsHandler {
	return &GetStatsHandler{
		statsQueryRepo: statsQueryRepo,
	}
}

// Handle 处理获取统计信息的查询
func (h *GetStatsHandler) Handle(_ context.Context, query GetStatsQuery) (*StatsDTO, error) {
	limit := query.RecentLogsLimit
	if limit <= 0 {
		limit = 5
	}

	systemStats, err := h.statsQueryRepo.GetSystemStats(limit)
	if err != nil {
		return nil, err
	}

	return toStatsDTO(systemStats), nil
}

// toStatsDTO 将 domain 统计信息转换为 DTO
func toStatsDTO(s *stats.SystemStats) *StatsDTO {
	result := &StatsDTO{
		TotalUsers:       s.TotalUsers,
		ActiveUsers:      s.ActiveUsers,
		InactiveUsers:    s.InactiveUsers,
		BannedUsers:      s.BannedUsers,
		TotalRoles:       s.TotalRoles,
		TotalPermissions: s.TotalPermissions,
	}

	if len(s.RecentAuditLogs) > 0 {
		result.RecentAuditLogs = make([]AuditLogSummaryDTO, len(s.RecentAuditLogs))
		for i, log := range s.RecentAuditLogs {
			result.RecentAuditLogs[i] = AuditLogSummaryDTO{
				ID:        log.ID,
				UserID:    log.UserID,
				Username:  log.Username,
				Action:    log.Action,
				Resource:  log.Resource,
				Status:    log.Status,
				CreatedAt: log.CreatedAt,
			}
		}
	}

	return result
}
