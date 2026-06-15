package cmd

import (
	"github.com/spf13/cobra"
)

var (
	sumText   string
	sumPoints int
)

var summarizeCmd = &cobra.Command{
	Use:     "summarize",
	Short:   "总结长文本为要点(--max-points 控制要点数)",
	Example: `  ai-app summarize --text "长文本..." --max-points 3`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := map[string]any{"text": sumText, "maxPoints": sumPoints}
		return runStreaming("summarize", params)
	},
}

func init() {
	summarizeCmd.Flags().StringVar(&sumText, "text", "", "待总结文本(必填)")
	summarizeCmd.Flags().IntVar(&sumPoints, "max-points", 3, "要点数量")
	_ = summarizeCmd.MarkFlagRequired("text")
}
