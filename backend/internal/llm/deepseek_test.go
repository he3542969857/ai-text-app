package llm

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeepSeekParsesSSE(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("missing auth header: %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"He\"}}]}\n\n")
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"llo\"}}]}\n\n")
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer srv.Close()

	c := NewDeepSeekClient(srv.URL, "test-key", "deepseek-chat")
	out := make(chan string, 16)
	if err := c.Stream(context.Background(),
		[]Message{{Role: "user", Content: "hi"}}, out); err != nil {
		t.Fatalf("err: %v", err)
	}
	close(out)

	var sb strings.Builder
	for tok := range out {
		sb.WriteString(tok)
	}
	if sb.String() != "Hello" {
		t.Fatalf("want Hello, got %q", sb.String())
	}
}

func TestDeepSeekNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	c := NewDeepSeekClient(srv.URL, "bad", "deepseek-chat")
	out := make(chan string, 16)
	err := c.Stream(context.Background(), []Message{{Role: "user", Content: "hi"}}, out)
	if err == nil {
		t.Fatal("want error on non-200")
	}
}
