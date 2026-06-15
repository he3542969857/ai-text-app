// Package client 调用后端 /api/task 接口,消费 SSE 流式结果。
package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Result 是一次任务调用的最终结果。
type Result struct {
	TaskID    string
	Status    string
	ElapsedMs int64
	Error     string
}

// Run 提交任务并消费 SSE:每个 token 调用 onToken 回调,返回最终 Result。
func Run(serverURL, taskType string, params map[string]any, onToken func(string)) (Result, error) {
	var res Result

	body, _ := json.Marshal(map[string]any{"type": taskType, "params": params})
	resp, err := http.Post(
		strings.TrimRight(serverURL, "/")+"/api/task",
		"application/json", bytes.NewReader(body))
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(resp.Body)
		return res, fmt.Errorf("后端返回 %d: %s", resp.StatusCode, strings.TrimSpace(buf.String()))
	}

	sc := bufio.NewScanner(resp.Body)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var event string
	for sc.Scan() {
		line := sc.Text()
		switch {
		case strings.HasPrefix(line, "event:"):
			event = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		case strings.HasPrefix(line, "data:"):
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			handleEvent(event, data, &res, onToken)
		}
	}
	return res, sc.Err()
}

func handleEvent(event, data string, res *Result, onToken func(string)) {
	switch event {
	case "meta":
		var m struct {
			TaskID string `json:"taskId"`
		}
		if json.Unmarshal([]byte(data), &m) == nil {
			res.TaskID = m.TaskID
		}
	case "token":
		var m struct {
			Text string `json:"text"`
		}
		if json.Unmarshal([]byte(data), &m) == nil {
			onToken(m.Text)
		}
	case "error":
		var m struct {
			Message string `json:"message"`
		}
		if json.Unmarshal([]byte(data), &m) == nil {
			res.Error = m.Message
		}
	case "done":
		var m struct {
			Status    string `json:"status"`
			ElapsedMs int64  `json:"elapsedMs"`
		}
		if json.Unmarshal([]byte(data), &m) == nil {
			res.Status = m.Status
			res.ElapsedMs = m.ElapsedMs
		}
	}
}
