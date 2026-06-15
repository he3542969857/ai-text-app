package llm

import "context"

// Message 是一条对话消息(OpenAI 兼容格式)。
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Client 是大模型流式生成接口。实现需逐 token 写入 out,
// 不负责关闭 out(由调用方关闭)。ctx 取消时应尽快返回 ctx.Err()。
type Client interface {
	Stream(ctx context.Context, messages []Message, out chan<- string) error
}
