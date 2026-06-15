package model

import (
	"errors"
	"time"
)

// TaskType 标识功能类型。
type TaskType string

// TaskStatus 是任务生命周期状态。
type TaskStatus string

const (
	TypeTranslate TaskType = "translate"
	TypeSummarize TaskType = "summarize"

	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusDone      TaskStatus = "done"
	StatusFailed    TaskStatus = "failed"
	StatusCancelled TaskStatus = "cancelled"
)

// Task 是统一任务模型,既供 SSE 流式消费,也供轮询与历史查询。
type Task struct {
	ID        string         `json:"id"`
	Type      TaskType       `json:"type"`
	Params    map[string]any `json:"params"`
	Status    TaskStatus     `json:"status"`
	Result    string         `json:"result"`
	Err       string         `json:"error,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	ElapsedMs int64          `json:"elapsedMs"`
}

// FunctionDef 描述列表页可展示的一个功能。
type FunctionDef struct {
	Type        TaskType       `json:"type"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Params      map[string]any `json:"params"`
}

// Functions 返回所有可用功能(中译英、英译中、文本总结)。
func Functions() []FunctionDef {
	return []FunctionDef{
		{TypeTranslate, "中译英", "将中文翻译为英文", map[string]any{"from": "zh", "to": "en"}},
		{TypeTranslate, "英译中", "将英文翻译为中文", map[string]any{"from": "en", "to": "zh"}},
		{TypeSummarize, "文本总结", "将长文本总结为要点", map[string]any{"maxPoints": 3}},
	}
}

const maxTextLen = 20000

var validDir = map[string]bool{"zh": true, "en": true}

// Validate 校验提交参数,非法时返回带说明的 error。
func Validate(typ TaskType, p map[string]any) error {
	text, _ := p["text"].(string)
	if text == "" {
		return errors.New("text 不能为空")
	}
	if len([]rune(text)) > maxTextLen {
		return errors.New("text 过长(超过 20000 字)")
	}
	switch typ {
	case TypeTranslate:
		from, _ := p["from"].(string)
		to, _ := p["to"].(string)
		if !validDir[from] || !validDir[to] || from == to {
			return errors.New("非法翻译方向(from/to 须为 zh/en 且不相同)")
		}
	case TypeSummarize:
		// text 已校验,无额外必填项
	default:
		return errors.New("未知功能类型")
	}
	return nil
}
