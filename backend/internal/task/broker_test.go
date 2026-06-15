package task

import (
	"testing"
	"time"
)

func drain(t *testing.T, sub <-chan string) []string {
	t.Helper()
	var got []string
	timeout := time.After(time.Second)
	for {
		select {
		case tok, ok := <-sub:
			if !ok {
				return got
			}
			got = append(got, tok)
		case <-timeout:
			t.Fatal("timeout draining sub")
		}
	}
}

func TestBrokerBroadcastsAndBuffers(t *testing.T) {
	b := NewBroker()
	// 先发布 A,再订阅:订阅者应收到缓冲历史 A,再收到后续 B
	b.Publish("A")
	sub := b.Subscribe()
	b.Publish("B")
	b.Close()

	got := drain(t, sub)
	if len(got) != 2 || got[0] != "A" || got[1] != "B" {
		t.Fatalf("want [A B], got %v", got)
	}
}

func TestBrokerSubscribeAfterClose(t *testing.T) {
	b := NewBroker()
	b.Publish("X")
	b.Close()
	// 关闭后订阅:补发缓冲并立即关闭 channel
	sub := b.Subscribe()
	got := drain(t, sub)
	if len(got) != 1 || got[0] != "X" {
		t.Fatalf("want [X], got %v", got)
	}
}

func TestBrokerMultipleSubscribers(t *testing.T) {
	b := NewBroker()
	s1 := b.Subscribe()
	s2 := b.Subscribe()
	b.Publish("hi")
	b.Close()
	if g1 := drain(t, s1); len(g1) != 1 || g1[0] != "hi" {
		t.Fatalf("s1 got %v", g1)
	}
	if g2 := drain(t, s2); len(g2) != 1 || g2[0] != "hi" {
		t.Fatalf("s2 got %v", g2)
	}
}
