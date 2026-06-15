package task

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"ai-text-app/backend/internal/llm"
	"ai-text-app/backend/internal/model"
)

type fakeStore struct {
	mu    sync.Mutex
	saved []model.Task
}

func (f *fakeStore) Save(t model.Task) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.saved = append(f.saved, t)
	return nil
}

func (f *fakeStore) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.saved)
}

// waitStatus 轮询等待任务到达期望状态,超时则失败。
func waitStatus(t *testing.T, m *Manager, id string, want model.TaskStatus) {
	t.Helper()
	deadline := time.After(3 * time.Second)
	for {
		tk, _ := m.Get(id)
		if tk.Status == want {
			return
		}
		select {
		case <-deadline:
			t.Fatalf("task %s status=%s, want %s", id, tk.Status, want)
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func TestManagerRunsTaskToDone(t *testing.T) {
	st := &fakeStore{}
	m := NewManager(llm.NewMockClient(), st, 2, 5*time.Second)
	m.Start()
	defer m.Stop()

	id := m.Submit(model.TypeTranslate, map[string]any{"text": "Hello", "from": "en", "to": "zh"})
	sub := m.Subscribe(id)

	var sb strings.Builder
	for tok := range sub {
		sb.WriteString(tok)
	}
	if sb.Len() == 0 {
		t.Fatal("no streamed output")
	}
	waitStatus(t, m, id, model.StatusDone)

	tk, _ := m.Get(id)
	if tk.Result == "" {
		t.Fatal("result not collected")
	}
	if tk.ElapsedMs <= 0 {
		t.Fatal("elapsed not recorded")
	}
	if st.count() == 0 {
		t.Fatal("task not persisted")
	}
}

func TestManagerCancel(t *testing.T) {
	st := &fakeStore{}
	m := NewManager(llm.NewMockClient(), st, 1, 5*time.Second)
	m.Start()
	defer m.Stop()

	id := m.Submit(model.TypeSummarize, map[string]any{"text": "a fairly long text to summarize"})
	_ = m.Subscribe(id)
	time.Sleep(30 * time.Millisecond) // 让 worker 进入 running
	if !m.Cancel(id) {
		t.Fatal("cancel returned false")
	}
	waitStatus(t, m, id, model.StatusCancelled)
}

func TestManagerGetUnknown(t *testing.T) {
	m := NewManager(llm.NewMockClient(), &fakeStore{}, 1, time.Second)
	if _, ok := m.Get("nope"); ok {
		t.Fatal("want not found")
	}
	_ = context.Background()
}
