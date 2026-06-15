package handler

import (
	"github.com/gin-gonic/gin"

	"ai-text-app/backend/internal/middleware"
	"ai-text-app/backend/internal/store"
	"ai-text-app/backend/internal/task"
)

// Register 注册中间件与所有 API 路由。
func Register(r *gin.Engine, m *task.Manager, st *store.Store) {
	r.Use(middleware.Trace(), middleware.Recover())
	api := r.Group("/api")
	api.GET("/functions", functionsHandler)
	api.POST("/task", taskHandler(m))
	api.GET("/task/:id", getTaskHandler(m))
	api.GET("/task/:id/stream", existingStreamHandler(m))
	api.DELETE("/task/:id", cancelHandler(m))
	api.GET("/tasks", historyHandler(st))
}
