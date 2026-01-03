package stats

// GetStatsQuery 获取统计信息的查询
type GetStatsQuery struct {
	RecentLogsLimit int // 最近审计日志条数，默认 5
}
