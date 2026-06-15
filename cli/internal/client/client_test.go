package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRunStreamsTokens(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/task" || r.Method != http.MethodPost {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "event: meta\ndata: {\"taskId\":\"t1\"}\n\n")
		fmt.Fprint(w, "event: token\ndata: {\"text\":\"你\"}\n\n")
		fmt.Fprint(w, "event: token\ndata: {\"text\":\"好\"}\n\n")
		fmt.Fprint(w, "event: done\ndata: {\"status\":\"done\",\"elapsedMs\":12}\n\n")
	}))
	defer srv.Close()

	var sb strings.Builder
	res, err := Run(srv.URL, "translate",
		map[string]any{"text": "hi", "from": "en", "to": "zh"},
		func(tok string) { sb.WriteString(tok) })
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if sb.String() != "你好" {
		t.Fatalf("want 你好, got %q", sb.String())
	}
	if res.TaskID != "t1" || res.Status != "done" {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestRunHandlesError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "event: error\ndata: {\"message\":\"boom\"}\n\n")
		fmt.Fprint(w, "event: done\ndata: {\"status\":\"failed\"}\n\n")
	}))
	defer srv.Close()

	res, err := Run(srv.URL, "summarize",
		map[string]any{"text": "x"}, func(string) {})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res.Status != "failed" || res.Error != "boom" {
		t.Fatalf("want failed/boom, got %+v", res)
	}
}

func TestRunBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message":"validation"}`)
	}))
	defer srv.Close()

	_, err := Run(srv.URL, "translate",
		map[string]any{"text": ""}, func(string) {})
	if err == nil {
		t.Fatal("want error on 400")
	}
}
