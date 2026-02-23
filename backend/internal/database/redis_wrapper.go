package database

import (
	"context"
	"github.com/redis/go-redis/v9"
)

// RedisConfPipline 保持原函数签名风格（但内部用 go-redis）
func RedisConfPipline(pip ...func(c redis.Cmdable)) error {
	pipe := RedisClient.Pipeline() // 或 TxPipeline() 如果需要事务

	for _, f := range pip {
		f(pipe)
	}

	_, err := pipe.Exec(context.Background())
	return err
}

// RedisConfDo 执行任意 Redis 命令（建议传入 ctx）
func RedisConfDo(ctx context.Context, commandName string, args ...interface{}) (interface{}, error) {
	return RedisClient.Do(ctx, commandName, args).Result()
}
