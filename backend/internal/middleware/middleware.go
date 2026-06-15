package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TraceKey 是 gin.Context 中存放 traceId 的键。
const TraceKey = "traceId"

// Trace 为每个请求注入短 traceId,写入响应头并记录访问日志。
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		tid := uuid.NewString()[:8]
		c.Set(TraceKey, tid)
		c.Header("X-Trace-Id", tid)
		start := time.Now()
		c.Next()
		log.Printf("[%s] %s %s -> %d (%s)", tid, c.Request.Method,
			c.Request.URL.Path, c.Writer.Status(), time.Since(start))
	}
}

// Recover 捕获 panic,返回统一错误结构 {code,message,traceId}。
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				tid, _ := c.Get(TraceKey)
				log.Printf("[%v] panic: %v", tid, r)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    "internal_error",
					"message": "服务器内部错误",
					"traceId": tid,
				})
			}
		}()
		c.Next()
	}
}
