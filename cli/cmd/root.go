package cmd

import (
	"github.com/spf13/cobra"
)

var (
	serverURL string
	jsonOut   bool
)

// rootCmd 是 ai-app CLI 的根命令。
var rootCmd = &cobra.Command{
	Use:           "ai-app",
	Short:         "AI 文本处理 CLI:中译英 / 英译中 / 文本总结",
	Long:          "ai-app 通过后端服务调用大模型,完成翻译与文本总结,结果流式输出。",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", "http://localhost:8080", "后端服务地址")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "以 JSON 输出最终结果")
	rootCmd.AddCommand(translateCmd)
	rootCmd.AddCommand(summarizeCmd)
}
