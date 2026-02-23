package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strconv"
	"time"
)

// 全局 Redis 客户端（单例）
var RedisClient *redis.Client

var TimeLocation = time.Local // 或 time.FixedZone("CST", 8*3600)

// InitRedisClient 初始化 Redis 客户端
func InitRedisClient() error {
	// 直接从环境变量读取配置
	addr := getEnv("REDIS_ADDR", "127.0.0.1:6379")
	password := getEnv("REDIS_PASSWORD", "")
	db, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	poolSize, _ := strconv.Atoi(getEnv("REDIS_POOL_SIZE", "10"))
	minIdleConns, _ := strconv.Atoi(getEnv("REDIS_MIN_IDLE_CONNS", "5"))
	maxIdleConns, _ := strconv.Atoi(getEnv("REDIS_MAX_IDLE_CONNS", "10"))
	dialTimeout, _ := strconv.Atoi(getEnv("REDIS_DIAL_TIMEOUT", "5"))
	readTimeout, _ := strconv.Atoi(getEnv("REDIS_READ_TIMEOUT", "3"))
	writeTimeout, _ := strconv.Atoi(getEnv("REDIS_WRITE_TIMEOUT", "3"))
	maxLifetime, _ := strconv.Atoi(getEnv("REDIS_MAX_LIFETIME", "300"))

	// 创建 Redis 客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:            addr,
		Password:        password,
		DB:              db,
		PoolSize:        poolSize,
		MinIdleConns:    minIdleConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: time.Duration(maxLifetime) * time.Second,
		DialTimeout:     time.Duration(dialTimeout) * time.Second,
		ReadTimeout:     time.Duration(readTimeout) * time.Second,
		WriteTimeout:    time.Duration(writeTimeout) * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	log.Printf("Redis connected successfully to %s", addr)
	return nil
}

// CloseRedisClient 关闭 Redis 客户端（程序退出时调用）
func CloseRedisClient() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
