package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ai-text-app/backend/internal/model"
)

// 已完成任务可通过 GET /api/task/:id/stream 一次性流式拿到完整结果。
func TestStreamExistingTaskReplaysResult(t *testing.T) {
	r, m := newTestRouter(t)

	// 提交并等待完成
	id := m.Submit(model.TypeTranslate, map[string]any{"text": "Hello", "from": "en", "to": "zh"})
	deadline := time.After(3 * time.Second)
	for {
		tk, _ := m.Get(id)
		if tk.Status == model.StatusDone {
			break
		}
		select {
		case <-deadline:
			t.Fatalf("task not done: %s", tk.Status)
		case <-time.After(10 * time.Millisecond):
		}
	}

	// 任务已完成后再连 SSE,应一次性补发 meta+token+done
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/task/"+id+"/stream", nil)
	r.ServeHTTP(w, req)
	body := w.Body.String()
	for _, ev := range []string{"event: meta", "event: token", "event: done"} {
		if !strings.Contains(body, ev) {
			t.Fatalf("missing %q in replay:\n%s", ev, body)
		}
	}
}

func TestStreamExistingUnknown(t *testing.T) {
	r, _ := newTestRouter(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/task/nope/stream", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", w.Code)
	}
}
