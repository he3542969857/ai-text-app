package task

import "sync"

// Broker 是单个任务的发布/订阅通道:
//   - 缓冲已产出的全部 token,新订阅者(含 SSE 断线重连)会先收到历史再收到后续;
//   - 支持多个并发订阅者;
//   - Close 后所有订阅 channel 关闭,后续订阅只补发历史并立即关闭。
type Broker struct {
	mu     sync.Mutex
	buf    []string
	subs   []chan string
	closed bool
}

func NewBroker() *Broker { return &Broker{} }

// Publish 追加一个 token 到缓冲并广播给当前所有订阅者。
func (b *Broker) Publish(tok string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.buf = append(b.buf, tok)
	for _, s := range b.subs {
		s <- tok
	}
}

// Subscribe 返回一个 channel:先补发缓冲历史,再接收后续 token。
func (b *Broker) Subscribe() <-chan string {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan string, len(b.buf)+64)
	for _, t := range b.buf {
		ch <- t
	}
	if b.closed {
		close(ch)
		return ch
	}
	b.subs = append(b.subs, ch)
	return ch
}

// Close 标记任务结束,关闭全部订阅 channel。
func (b *Broker) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.closed = true
	for _, s := range b.subs {
		close(s)
	}
	b.subs = nil
}
