package bootstrap

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-iam-bc/internal/config"
)

// Server 管理 HTTP Server 生命周期
type Server struct {
	engine *gin.Engine
	addr   string
	srv    *http.Server
}

// NewServer 创建 HTTP Server
func NewServer(engine *gin.Engine, cfg *config.Config) *Server {
	return &Server{
		engine: engine,
		addr:   cfg.Server.Addr,
	}
}

// Start 启动 HTTP Server
func (s *Server) Start() error {
	s.srv = &http.Server{
		Addr:              s.addr,
		Handler:           s.engine,
		ReadHeaderTimeout: 10 * time.Second, // 防止 Slowloris 攻击
	}
	return s.srv.ListenAndServe()
}

// Stop 优雅关闭 HTTP Server
func (s *Server) Stop(ctx context.Context) error {
	if s.srv == nil {
		return nil
	}
	slog.Info("Shutting down HTTP server gracefully")
	return s.srv.Shutdown(ctx)
}
