package middleware

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/audit"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/ctxutil"
)

// Audit 返回审计中间件工厂函数
//
// operation 参数用于标识操作类型（如 "iam:user:create"）
//
// 中间件行为：
//   - 从请求上下文中提取用户信息
//   - 异步记录审计日志，不阻塞请求
//   - 只在请求成功时记录（4xx/5xx 不记录）
func Audit(createHandler *audit.CreateHandler, operation string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 使用 writer 存储原始响应
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 继续处理请求
		c.Next()

		// 只在请求成功时记录审计日志
		if c.Writer.Status() >= 400 {
			return
		}

		// 构建审计命令
		cmd := audit.CreateCommand{
			UserID:      getUserID(c),
			Username:    getUsername(c),
			Action:      extractAction(c.Request.Method),
			Resource:    extractResource(c.Request.URL.Path),
			ResourceID:  extractResourceID(c),
			IPAddress:   c.ClientIP(),
			UserAgent:   c.Request.UserAgent(),
			Details:     buildDetails(c, startTime),
			Status:      "success",
			RequestID:   getRequestID(c),
			OperationID: operation,
		}

		// 异步记录审计日志
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = createHandler.Handle(ctx, cmd)
		}()
	}
}

// getUserID 从上下文中获取用户 ID
func getUserID(c *gin.Context) uint {
	if userID, ok := ctxutil.Get[uint](c, ctxutil.UserID); ok {
		return userID
	}
	return 0
}

// getUsername 从上下文中获取用户名
func getUsername(c *gin.Context) string {
	if username, ok := ctxutil.Get[string](c, ctxutil.Username); ok {
		return username
	}
	return ""
}

// getRequestID 从上下文中获取请求 ID
func getRequestID(c *gin.Context) string {
	if requestID, ok := ctxutil.Get[string](c, ctxutil.RequestID); ok {
		return requestID
	}
	return ""
}

// extractAction 从 HTTP 方法提取操作类型
func extractAction(method string) string {
	switch method {
	case "POST":
		return "create"
	case "GET":
		return "read"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return method
	}
}

// extractResource 从路径提取资源类型
func extractResource(path string) string {
	// 简单解析：/api/admin/users -> users
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "api" && i+1 < len(parts) {
			continue
		}
		if i > 0 && part != "" && !strings.HasPrefix(part, ":") && !strings.HasPrefix(part, "{") {
			return part
		}
	}
	return "unknown"
}

// extractResourceID 从路径参数中提取资源 ID
func extractResourceID(c *gin.Context) string {
	// 尝试获取路径参数中的 ID
	for _, p := range c.Params {
		if strings.HasSuffix(p.Key, "_id") || p.Key == "id" {
			return p.Value
		}
	}
	return ""
}

// buildDetails 构建审计详情
func buildDetails(c *gin.Context, startTime time.Time) string {
	duration := time.Since(startTime)
	return fmt.Sprintf("method=%s path=%s status=%d duration=%s",
		c.Request.Method,
		c.Request.URL.Path,
		c.Writer.Status(),
		duration.Round(time.Millisecond),
	)
}

// bodyLogWriter 用于捕获响应体
type bodyLogWriter struct {
	gin.ResponseWriter

	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
