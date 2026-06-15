package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ai-text-app/backend/internal/model"
	"ai-text-app/backend/internal/task"
)

type taskReq struct {
	Type   model.TaskType `json:"type"`
	Params map[string]any `json:"params"`
}

// taskHandler 提交任务并以 SSE 流式返回 taskId 与逐 token 结果。
func taskHandler(m *task.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req taskReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "bad_request", "message": "请求体非法"})
			return
		}
		if err := model.Validate(req.Type, req.Params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "validation", "message": err.Error()})
			return
		}

		id := m.Submit(req.Type, req.Params)
		streamTaskSSE(c, m, id, true)
	}
}

// existingStreamHandler 让客户端用 taskId 重连:
// 若任务仍在执行则继续流式;若已完成则一次性补发完整结果(meta+token+done)。
func existingStreamHandler(m *task.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if _, ok := m.Get(id); !ok {
			c.JSON(http.StatusNotFound, gin.H{"code": "not_found", "message": "任务不存在"})
			return
		}
		streamTaskSSE(c, m, id, false)
	}
}

// streamTaskSSE 订阅任务并以 SSE 推送。cancelOnDisconnect 为 true 时,
// 客户端断开会取消任务(用于首次提交);重连查看场景应为 false。
func streamTaskSSE(c *gin.Context, m *task.Manager, id string, cancelOnDisconnect bool) {
	sub := m.Subscribe(id)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	writeSSE(c, "meta", gin.H{"taskId": id})
	flush(c)

	ctx := c.Request.Context()
	for {
		select {
		case <-ctx.Done():
			if cancelOnDisconnect {
				m.Cancel(id)
			}
			return
		case tok, ok := <-sub:
			if !ok {
				tk, _ := m.Get(id)
				if tk.Status == model.StatusFailed {
					writeSSE(c, "error", gin.H{"message": tk.Err})
				}
				writeSSE(c, "done", gin.H{"status": tk.Status, "elapsedMs": tk.ElapsedMs})
				flush(c)
				return
			}
			writeSSE(c, "token", gin.H{"text": tok})
			flush(c)
		}
	}
}

func writeSSE(c *gin.Context, event string, data any) {
	b, _ := json.Marshal(data)
	fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", event, b)
}

func flush(c *gin.Context) {
	if f, ok := c.Writer.(http.Flusher); ok {
		f.Flush()
	}
}
