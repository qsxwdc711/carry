package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CtxKey string

const CtxTraceIDKey CtxKey = "trace_id"
const GinTraceKey = "trace_id"

func TraceMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		traceID := uuid.NewString()
		zap.L().Info("TraceMiddleware executed", zap.String("trace_id", traceID))
		// 放到 gin 上下文，handler 可以直接 c.Get("trace_id")
		c.Set(GinTraceKey, traceID)
		// 放到 request.Context，这样向下传入的 context 可以通过 ctx.Value 读取
		reqCtx := context.WithValue(c.Request.Context(), CtxTraceIDKey, traceID)
		c.Request = c.Request.WithContext(reqCtx)
		// 也把 trace 返回给客户端，方便排查
		c.Writer.Header().Set("X-Trace-ID", traceID)
		c.Next()
	}
}
func GetTraceID(ctx interface{}) string {
	switch t := ctx.(type) {
	case *gin.Context:
		// 如果是 gin 上下文，从 GinTraceKey 获取
		if traceID, exists := t.Get(GinTraceKey); exists {
			if s, ok := traceID.(string); ok {
				return s
			}
		}
	case context.Context:
		// 如果是标准库上下文，从 CtxTraceIDKey 获取
		if traceID := t.Value(CtxTraceIDKey); traceID != nil {
			if s, ok := traceID.(string); ok {
				return s
			}
		}
	}
	// 兜底：如果获取失败，返回 "unknown"
	return "unknown"
}
