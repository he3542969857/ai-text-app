package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"ai-text-app/cli/internal/client"
)

var (
	trText string
	trFrom string
	trTo   string
)

var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "翻译文本(--from/--to 取 zh 或 en)",
	Example: `  ai-app translate --text "Hello" --from en --to zh
  ai-app translate --text "你好" --from zh --to en`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := map[string]any{"text": trText, "from": trFrom, "to": trTo}
		return runStreaming("translate", params)
	},
}

func init() {
	translateCmd.Flags().StringVar(&trText, "text", "", "待翻译文本(必填)")
	translateCmd.Flags().StringVar(&trFrom, "from", "", "源语言 zh|en(必填)")
	translateCmd.Flags().StringVar(&trTo, "to", "", "目标语言 zh|en(必填)")
	_ = translateCmd.MarkFlagRequired("text")
	_ = translateCmd.MarkFlagRequired("from")
	_ = translateCmd.MarkFlagRequired("to")
}

// runStreaming 是 translate/summarize 共用的执行逻辑:流式打印或 JSON 输出。
func runStreaming(taskType string, params map[string]any) error {
	onToken := func(tok string) {
		if !jsonOut {
			fmt.Print(tok)
		}
	}
	res, err := client.Run(serverURL, taskType, params, onToken)
	if err != nil {
		return err
	}
	if jsonOut {
		b, _ := json.MarshalIndent(map[string]any{
			"taskId":    res.TaskID,
			"status":    res.Status,
			"elapsedMs": res.ElapsedMs,
			"error":     res.Error,
		}, "", "  ")
		fmt.Println(string(b))
	} else {
		fmt.Println()
		if res.Error != "" {
			fmt.Fprintf(os.Stderr, "错误: %s\n", res.Error)
		}
	}
	return nil
}
