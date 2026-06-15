package llm

import (
	"fmt"

	"ai-text-app/backend/internal/model"
)

// BuildMessages 根据功能类型与参数构造发送给大模型的消息序列。
func BuildMessages(typ model.TaskType, p map[string]any) []Message {
	text, _ := p["text"].(string)
	switch typ {
	case model.TypeTranslate:
		from, _ := p["from"].(string)
		to, _ := p["to"].(string)
		lang := map[string]string{"zh": "中文", "en": "英文"}
		sys := fmt.Sprintf("你是专业翻译。将文本从%s翻译为%s,只输出译文,不要任何解释。", lang[from], lang[to])
		return []Message{{Role: "system", Content: sys}, {Role: "user", Content: text}}
	case model.TypeSummarize:
		points := 3
		if v, ok := p["maxPoints"].(float64); ok && v > 0 {
			points = int(v)
		}
		sys := fmt.Sprintf("你是文本总结助手。用不超过 %d 个要点总结以下文本,每个要点单独一行,以「- 」开头。", points)
		if v, ok := p["maxWords"].(float64); ok && v > 0 {
			sys += fmt.Sprintf("总字数不超过 %d 字。", int(v))
		}
		return []Message{{Role: "system", Content: sys}, {Role: "user", Content: text}}
	}
	return []Message{{Role: "user", Content: text}}
}
