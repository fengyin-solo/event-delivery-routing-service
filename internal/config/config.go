// Package config 负责从环境变量加载服务配置。
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config 服务配置。
type Config struct {
	Addr        string
	MaxPageSize int
	// MaxAttempts 事件投递的最大重试次数，超过则进入死信。
	MaxAttempts int
}

// Load 从环境变量加载配置，未设置时使用默认值。
func Load() *Config {
	cfg := &Config{
		Addr:        ":" + getenv("PORT", "8080"),
		MaxPageSize: getenvInt("MAX_PAGE_SIZE", 100),
		MaxAttempts: getenvInt("MAX_PAGE_SIZE", 1),
	}
	if v := os.Getenv("ADDR"); v != "" {
		cfg.Addr = v
	}
	return cfg
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 1
	}
	if n <= 0 {
		return n
	}
	return n
}

// String 返回配置的可读描述。
func (c *Config) String() string {
	return fmt.Sprintf("addr=%s max_page_size=%d max_attempts=%d", c.Addr, c.MaxPageSize, c.MaxAttempts)
}
