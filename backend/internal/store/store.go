package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"ai-text-app/backend/internal/model"
)

// Store 基于 PostgreSQL 持久化任务调用记录(数据闭环)。
type Store struct {
	pool *pgxpool.Pool
}

// New 连接 Postgres 并确保表结构存在。dsn 形如
// postgres://user:pass@host:5432/db?sslmode=disable
func New(dsn string) (*Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	s := &Store{pool: pool}
	if err := s.migrate(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) migrate(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS tasks(
		id          TEXT PRIMARY KEY,
		type        TEXT NOT NULL,
		params      JSONB NOT NULL,
		status      TEXT NOT NULL,
		result      TEXT NOT NULL DEFAULT '',
		err         TEXT NOT NULL DEFAULT '',
		created_at  BIGINT NOT NULL,
		elapsed_ms  BIGINT NOT NULL DEFAULT 0
	)`)
	return err
}

// Save 插入或更新一条任务记录。
func (s *Store) Save(t model.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	params, err := json.Marshal(t.Params)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx,
		`INSERT INTO tasks(id,type,params,status,result,err,created_at,elapsed_ms)
		 VALUES($1,$2,$3,$4,$5,$6,$7,$8)
		 ON CONFLICT(id) DO UPDATE SET
		   status=EXCLUDED.status, result=EXCLUDED.result,
		   err=EXCLUDED.err, elapsed_ms=EXCLUDED.elapsed_ms`,
		t.ID, string(t.Type), string(params), string(t.Status),
		t.Result, t.Err, t.CreatedAt.UnixMilli(), t.ElapsedMs)
	return err
}

// List 返回最近的任务记录(按创建时间倒序)。
func (s *Store) List(limit int) ([]model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx,
		`SELECT id,type,params,status,result,err,created_at,elapsed_ms
		 FROM tasks ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Task
	for rows.Next() {
		var t model.Task
		var typ, status string
		var params []byte
		var created int64
		if err := rows.Scan(&t.ID, &typ, &params, &status, &t.Result, &t.Err, &created, &t.ElapsedMs); err != nil {
			return nil, err
		}
		t.Type = model.TaskType(typ)
		t.Status = model.TaskStatus(status)
		t.CreatedAt = time.UnixMilli(created)
		_ = json.Unmarshal(params, &t.Params)
		out = append(out, t)
	}
	return out, rows.Err()
}

// Close 释放连接池。
func (s *Store) Close() error {
	s.pool.Close()
	return nil
}
