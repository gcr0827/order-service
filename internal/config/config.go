package config

import "os"

type Config struct {
	Port string
}

// Load 读取配置。W1 用环境变量，W2 起可换成 viper。
func Load() *Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	return &Config{Port: port}
}
