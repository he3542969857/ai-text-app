package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ai-text-app/backend/internal/model"
	"ai-text-app/backend/internal/store"
	"ai-text-app/backend/internal/task"
)

// getTaskHandler 返回单个任务状态(供前端轮询)。
func getTaskHandler(m *task.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tk, ok := m.Get(c.Param("id"))
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"code": "not_found", "message": "任务不存在"})
			return
		}
		c.JSON(http.StatusOK, tk)
	}
}

// cancelHandler 取消运行中的任务。
func cancelHandler(m *task.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.Cancel(c.Param("id")) {
			c.JSON(http.StatusNotFound, gin.H{"code": "not_found", "message": "任务不存在或已结束"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
	}
}

// historyHandler 返回历史调用记录(数据闭环查询页)。
func historyHandler(st *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := 50
		if v := c.Query("limit"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				limit = n
			}
		}
		list, err := st.List(limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "db_error", "message": err.Error()})
			return
		}
		if list == nil {
			list = []model.Task{}
		}
		c.JSON(http.StatusOK, list)
	}
}
