package handler

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application/health"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	checker health.Checker
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(checker health.Checker) *HealthHandler {
	return &HealthHandler{
		checker: checker,
	}
}

// Check 执行健康检查
//
//	@Summary		健康检查
//	@Description	检查系统服务健康状态（数据库、Redis）
//	@Tags			系统 (System)
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.DataResponse[health.HealthReport]	"服务健康"
//	@Failure		503	{object}	response.DataResponse[health.HealthReport]	"服务降级"
//	@Router			/health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	report := h.checker.Check(ctx)

	if report.Status == health.StatusHealthy {
		response.OK(c, report, "healthy")
	} else {
		response.ServiceUnavailable(c, string(report.Status))
	}
}

// Live Kubernetes liveness probe
//
//	@Summary		存活检查
//	@Description	Kubernetes liveness probe，检查应用是否存活
//	@Tags			系统 (System)
//	@Produce		json
//	@Success		200	{object}	response.DataResponse[any]	"存活"
//	@Router			/health/live [get]
func (h *HealthHandler) Live(c *gin.Context) {
	response.OK(c, nil, "alive")
}

// Ready Kubernetes readiness probe
//
//	@Summary		就绪检查
//	@Description	Kubernetes readiness probe，检查应用是否就绪接受流量
//	@Tags			系统 (System)
//	@Produce		json
//	@Success		200	{object}	response.DataResponse[any]			"就绪"
//	@Failure		503	{object}	response.ErrorResponse				"未就绪"
//	@Router			/health/ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	report := h.checker.Check(ctx)
	if report.Status == health.StatusUnhealthy {
		response.ServiceUnavailable(c, "not ready")
		return
	}
	response.OK(c, nil, "ready")
}
