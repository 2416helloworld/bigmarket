package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port          int
	DBDSN         string
	RedisAddr     string
	RedisPassword string
	RabbitURL     string
	JWTSecret     string
	DrawRateLimit int // 每个用户每秒抽奖次数
}

func Load() Config {
	return Config{
		Port:          envInt("APP_PORT", 8080),
		DBDSN:         env("MYSQL_DSN", "root:123456@tcp(mysql:3306)/big_market?charset=utf8mb4&parseTime=True&loc=Local"),
		RedisAddr:     env("REDIS_ADDR", "redis:6379"),
		RedisPassword: env("REDIS_PASSWORD", ""),
		RabbitURL:     env("RABBITMQ_URL", "amqp://admin:admin@rabbitmq:5672/"),
		JWTSecret:     env("JWT_SECRET", "big-market-go-secret"),
		DrawRateLimit: envInt("DRAW_RATE_LIMIT", 1),
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
