package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRecoverReturnsJSONWithTraceID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Trace(), Recover())
	r.GET("/boom", func(c *gin.Context) { panic("kaboom") })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/boom", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("code=%d, want 500", w.Code)
	}
	if !strings.Contains(w.Body.String(), "traceId") {
		t.Fatalf("body missing traceId: %s", w.Body.String())
	}
}

func TestTraceSetsHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Trace())
	r.GET("/ok", func(c *gin.Context) {
		if _, ok := c.Get(TraceKey); !ok {
			t.Error("traceId not set in context")
		}
		c.Status(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ok", nil)
	r.ServeHTTP(w, req)

	if w.Header().Get("X-Trace-Id") == "" {
		t.Fatal("X-Trace-Id header not set")
	}
}
