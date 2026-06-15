package llm

import (
	"context"
	"time"
)

// MockClient 在无 DEEPSEEK_API_KEY 时使用,逐字模拟流式输出。
// 它与 DeepSeekClient 实现同一 Client 接口,保留真实调用链路:
// 上层(Manager/Handler)无需感知用的是真实还是模拟客户端。
type MockClient struct {
	delay time.Duration
}

func NewMockClient() *MockClient {
	return &MockClient{delay: 15 * time.Millisecond}
}

func (m *MockClient) Stream(ctx context.Context, messages []Message, out chan<- string) error {
	last := ""
	if len(messages) > 0 {
		last = messages[len(messages)-1].Content
	}
	reply := "[MOCK] 已收到并处理:" + last
	for _, r := range reply {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- string(r):
			time.Sleep(m.delay)
		}
	}
	return nil
}
