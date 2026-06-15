package task

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"ai-text-app/backend/internal/llm"
	"ai-text-app/backend/internal/model"
)

// Store 是 Manager 持久化已完成任务的依赖(由 SQLite 仓储实现)。
type Store interface {
	Save(t model.Task) error
}

type entry struct {
	task   model.Task
	broker *Broker
	cancel context.CancelFunc
}

// Manager 维护任务注册表、内存队列与 worker 池,统一驱动任务生命周期。
type Manager struct {
	llm     llm.Client
	store   Store
	timeout time.Duration

	mu    sync.RWMutex
	tasks map[string]*entry

	queue chan string
	wg    sync.WaitGroup
	stop  chan struct{}
	n     int
}

func NewManager(c llm.Client, store Store, workers int, timeout time.Duration) *Manager {
	return &Manager{
		llm:     c,
		store:   store,
		timeout: timeout,
		tasks:   map[string]*entry{},
		queue:   make(chan string, 256),
		stop:    make(chan struct{}),
		n:       workers,
	}
}

// Start 启动 worker 池。
func (m *Manager) Start() {
	for i := 0; i < m.n; i++ {
		m.wg.Add(1)
		go m.worker()
	}
}

// Stop 停止所有 worker(等待退出)。
func (m *Manager) Stop() {
	close(m.stop)
	m.wg.Wait()
}

// Submit 创建任务并入队,返回 taskId。
func (m *Manager) Submit(typ model.TaskType, params map[string]any) string {
	id := uuid.NewString()
	m.mu.Lock()
	m.tasks[id] = &entry{
		task: model.Task{
			ID: id, Type: typ, Params: params,
			Status: model.StatusPending, CreatedAt: time.Now(),
		},
		broker: NewBroker(),
	}
	m.mu.Unlock()
	m.queue <- id
	return id
}

// Subscribe 订阅任务的流式输出(未知任务返回已关闭 channel)。
func (m *Manager) Subscribe(id string) <-chan string {
	m.mu.RLock()
	e, ok := m.tasks[id]
	m.mu.RUnlock()
	if !ok {
		ch := make(chan string)
		close(ch)
		return ch
	}
	return e.broker.Subscribe()
}

// Get 返回任务当前快照。
func (m *Manager) Get(id string) (model.Task, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.tasks[id]
	if !ok {
		return model.Task{}, false
	}
	return e.task, true
}

// Cancel 取消运行中的任务(触发 context 取消,中断 LLM 调用)。
func (m *Manager) Cancel(id string) bool {
	m.mu.RLock()
	e, ok := m.tasks[id]
	m.mu.RUnlock()
	if !ok || e.cancel == nil {
		return false
	}
	e.cancel()
	return true
}

func (m *Manager) setStatus(id string, s model.TaskStatus, result, errMsg string, elapsed int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.tasks[id]
	if !ok {
		return
	}
	e.task.Status = s
	if result != "" {
		e.task.Result = result
	}
	if errMsg != "" {
		e.task.Err = errMsg
	}
	if elapsed > 0 {
		e.task.ElapsedMs = elapsed
	}
}

func (m *Manager) worker() {
	defer m.wg.Done()
	for {
		select {
		case <-m.stop:
			return
		case id := <-m.queue:
			m.run(id)
		}
	}
}

func (m *Manager) run(id string) {
	m.mu.RLock()
	e := m.tasks[id]
	m.mu.RUnlock()
	if e == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	m.mu.Lock()
	e.cancel = cancel
	e.task.Status = model.StatusRunning
	m.mu.Unlock()
	defer cancel()

	start := time.Now()
	out := make(chan string, 64)
	done := make(chan error, 1)
	go func() {
		done <- m.llm.Stream(ctx, llm.BuildMessages(e.task.Type, e.task.Params), out)
		close(out)
	}()

	var collected []rune
	for tok := range out {
		collected = append(collected, []rune(tok)...)
		e.broker.Publish(tok)
	}
	err := <-done
	e.broker.Close()
	elapsed := time.Since(start).Milliseconds()

	status := model.StatusDone
	errMsg := ""
	if err != nil {
		if ctx.Err() == context.Canceled {
			status = model.StatusCancelled
		} else {
			status = model.StatusFailed
		}
		errMsg = err.Error()
	}
	m.setStatus(id, status, string(collected), errMsg, elapsed)

	final, _ := m.Get(id)
	_ = m.store.Save(final)
}
