package llm

import (
	"context"
	"strings"
	"testing"
)

func TestMockStreamsTokens(t *testing.T) {
	c := NewMockClient()
	out := make(chan string, 256)
	err := c.Stream(context.Background(),
		[]Message{{Role: "user", Content: "Hello"}}, out)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	close(out)

	var sb strings.Builder
	n := 0
	for tok := range out {
		sb.WriteString(tok)
		n++
	}
	if n < 2 {
		t.Fatalf("want streamed tokens, got %d", n)
	}
	if sb.Len() == 0 {
		t.Fatal("empty result")
	}
}

func TestMockRespectsCancel(t *testing.T) {
	c := NewMockClient()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	out := make(chan string, 256)
	err := c.Stream(ctx, []Message{{Role: "user", Content: "hi"}}, out)
	if err == nil {
		t.Fatal("want context cancel error")
	}
}
