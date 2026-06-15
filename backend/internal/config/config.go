package config

import (
	"os"
	"strconv"
)

// Config 保存运行期配置,全部来自环境变量(支持 .env 注入)。
type Config struct {
	Port           string
	DeepSeekAPIKey string
	DeepSeekBase   string
	Model          string
	TimeoutSec     int
	DatabaseURL    string
}

// Load 读取环境变量并填充默认值。
func Load() Config {
	return Config{
		Port:           env("PORT", "8080"),
		DeepSeekAPIKey: env("DEEPSEEK_API_KEY", ""),
		DeepSeekBase:   env("DEEPSEEK_BASE_URL", "https://api.deepseek.com"),
		Model:          env("DEEPSEEK_MODEL", "deepseek-chat"),
		TimeoutSec:     envInt("REQUEST_TIMEOUT_SEC", 60),
		DatabaseURL:    env("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/aitext?sslmode=disable"),
	}
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func envInt(k string, d int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return d
}
