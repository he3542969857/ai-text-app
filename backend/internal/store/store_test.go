package store

import (
	"os"
	"testing"
	"time"

	"ai-text-app/backend/internal/model"
)

// testDSN 返回测试库连接串,可用 TEST_DATABASE_URL 覆盖。
func testDSN() string {
	if v := os.Getenv("TEST_DATABASE_URL"); v != "" {
		return v
	}
	return "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
}

// newTestStore 连接测试库并清空 tasks 表;库不可用时跳过(集成测试)。
func newTestStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(testDSN())
	if err != nil {
		t.Skipf("跳过:无法连接测试 Postgres(%v)。设置 TEST_DATABASE_URL 或启动本地 PG。", err)
	}
	if _, err := s.pool.Exec(t.Context(), "TRUNCATE tasks"); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	return s
}

func TestSaveAndList(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	task := model.Task{
		ID: "id1", Type: model.TypeTranslate,
		Params: map[string]any{"text": "hi", "from": "en", "to": "zh"},
		Status: model.StatusDone, Result: "你好",
		CreatedAt: time.Now(), ElapsedMs: 123,
	}
	if err := s.Save(task); err != nil {
		t.Fatal(err)
	}
	list, err := s.List(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != "id1" || list[0].Result != "你好" {
		t.Fatalf("unexpected list: %+v", list)
	}
	if list[0].ElapsedMs != 123 {
		t.Fatalf("elapsed not persisted: %d", list[0].ElapsedMs)
	}
}

func TestSaveUpsert(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	base := model.Task{
		ID: "u1", Type: model.TypeSummarize,
		Params: map[string]any{"text": "x"}, Status: model.StatusRunning,
		CreatedAt: time.Now(),
	}
	if err := s.Save(base); err != nil {
		t.Fatal(err)
	}
	base.Status = model.StatusDone
	base.Result = "done result"
	base.ElapsedMs = 50
	if err := s.Save(base); err != nil {
		t.Fatal(err)
	}
	list, _ := s.List(10)
	if len(list) != 1 {
		t.Fatalf("want 1 row after upsert, got %d", len(list))
	}
	if list[0].Status != model.StatusDone || list[0].Result != "done result" {
		t.Fatalf("upsert not applied: %+v", list[0])
	}
}

func TestListOrderAndLimit(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	for i, id := range []string{"a", "b", "c"} {
		task := model.Task{
			ID: id, Type: model.TypeTranslate,
			Params: map[string]any{"text": "t"}, Status: model.StatusDone,
			CreatedAt: time.UnixMilli(int64(1000 + i)),
		}
		if err := s.Save(task); err != nil {
			t.Fatal(err)
		}
	}
	list, err := s.List(2)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 || list[0].ID != "c" || list[1].ID != "b" {
		t.Fatalf("want [c b], got %+v", list)
	}
}
