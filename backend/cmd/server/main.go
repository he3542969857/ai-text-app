package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"ai-text-app/backend/internal/config"
	"ai-text-app/backend/internal/handler"
	"ai-text-app/backend/internal/llm"
	"ai-text-app/backend/internal/store"
	"ai-text-app/backend/internal/task"
)

func main() {
	cfg := config.Load()

	st, err := store.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer st.Close()

	var client llm.Client
	if cfg.DeepSeekAPIKey == "" {
		log.Println("⚠ DEEPSEEK_API_KEY 未设置,使用 MockClient(逐字模拟,保留真实调用链路)")
		client = llm.NewMockClient()
	} else {
		client = llm.NewDeepSeekClient(cfg.DeepSeekBase, cfg.DeepSeekAPIKey, cfg.Model)
	}

	mgr := task.NewManager(client, st, 4, time.Duration(cfg.TimeoutSec)*time.Second)
	mgr.Start()
	defer mgr.Stop()

	r := gin.New()
	r.Use(corsMiddleware())
	handler.Register(r, mgr, st)

	log.Printf("后端启动,监听 :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

// corsMiddleware 允许前端跨域访问(开发期 *)。
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
