package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const CtxTraceIDKey = "trace_id"

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.NewString()
		// 放到 gin 上下文，handler 可以直接 c.Get("trace_id")
		c.Set(CtxTraceIDKey, traceID)
		// 放到 request.Context，这样向下传入的 context 可以通过 ctx.Value 读取
		reqCtx := context.WithValue(c.Request.Context(), CtxTraceIDKey, traceID)
		c.Request = c.Request.WithContext(reqCtx)
		// 也把 trace 返回给客户端，方便排查
		c.Writer.Header().Set("X-Trace-ID", traceID)
		c.Next()
	}
}
