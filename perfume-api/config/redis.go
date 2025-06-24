// config/redis.go
package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	RedisCtx    = context.Background()
)

func ConnectRedis() {
	addr := os.Getenv("REDIS_ADDR") // <-- lấy từ biến môi trường
	if addr == "" {
		addr = "redis:6379" // fallback nếu biến không được đặt
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     "",
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	if err := RedisClient.Ping(RedisCtx).Err(); err != nil {
		log.Printf("❌ Lỗi kết nối Redis: %v", err) // Thay vì panic
		// Quan trọng: Đảm bảo RedisClient vẫn là nil hoặc xử lý nó sau.
		// Để debug, hàm này có thể kết thúc mà RedisClient là nil nếu có lỗi.
		return
	}

	log.Println("✅ Đã kết nối Redis thành công!")
}
