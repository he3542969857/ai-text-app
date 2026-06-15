package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"ai-text-app/backend/internal/llm"
	"ai-text-app/backend/internal/store"
	"ai-text-app/backend/internal/task"
)

func testDSN() string {
	if v := os.Getenv("TEST_DATABASE_URL"); v != "" {
		return v
	}
	return "postgres://postgres:postgres@localhost:5432/aitext?sslmode=disable"
}

func newTestRouter(t *testing.T) (*gin.Engine, *task.Manager) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	st, err := store.New(testDSN())
	if err != nil {
		t.Skipf("跳过:无法连接测试 Postgres(%v)", err)
	}
	m := task.NewManager(llm.NewMockClient(), st, 2, 5*time.Second)
	m.Start()
	t.Cleanup(func() { m.Stop(); st.Close() })

	r := gin.New()
	Register(r, m, st)
	return r, m
}

func TestFunctionsEndpoint(t *testing.T) {
	r, _ := newTestRouter(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/functions", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 || !strings.Contains(w.Body.String(), "translate") {
		t.Fatalf("code=%d body=%s", w.Code, w.Body.String())
	}
}

func TestTaskSSEStreams(t *testing.T) {
	r, _ := newTestRouter(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/task",
		strings.NewReader(`{"type":"translate","params":{"text":"Hello","from":"en","to":"zh"}}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	body := w.Body.String()
	for _, ev := range []string{"event: meta", "event: token", "event: done"} {
		if !strings.Contains(body, ev) {
			t.Fatalf("missing %q in SSE body:\n%s", ev, body)
		}
	}
}

func TestTaskValidation(t *testing.T) {
	r, _ := newTestRouter(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/task",
		strings.NewReader(`{"type":"translate","params":{"text":"","from":"en","to":"zh"}}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestHistoryEndpoint(t *testing.T) {
	r, _ := newTestRouter(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/tasks?limit=5", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("code=%d", w.Code)
	}
}

func TestGetUnknownTask(t *testing.T) {
	r, _ := newTestRouter(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/task/nonexistent", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", w.Code)
	}
}
