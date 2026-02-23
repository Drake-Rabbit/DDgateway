package database

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	RedisFlowDayKey  = "flow_day"
	RedisFlowHourKey = "flow_hour"
)

type RedisFlowCountService struct {
	AppID       string
	Interval    time.Duration // 间隔
	QPS         int64
	Unix        int64
	TickerCount int64
	TotalCount  int64
}

func NewRedisFlowCountService(appID string, interval time.Duration) *RedisFlowCountService {
	reqCounter := &RedisFlowCountService{
		AppID:    appID,    // 应用ID
		Interval: interval, // 间隔
		QPS:      0,        // QPS
		Unix:     0,        // 当前时间戳
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("RedisFlowCountService panic:", err)
			}
		}()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			//每个周期内的操作流程,
			//步骤 1️⃣读取并重置本地计数器
			<-ticker.C //每interval触发一次,没有则堵塞
			tickerCount := atomic.LoadInt64(&reqCounter.TickerCount)
			atomic.StoreInt64(&reqCounter.TickerCount, 0)
			//步骤 2️⃣：生成 Redis Key
			currentTime := time.Now()
			dayKey := reqCounter.GetDayKey(currentTime)
			hourKey := reqCounter.GetHourKey(currentTime)

			// 使用新的 RedisConfPipline
			//步骤 3️⃣：批量写入 Redis（管道）
			err := RedisConfPipline(func(c redis.Cmdable) {
				c.IncrBy(context.Background(), dayKey, tickerCount)
				c.Expire(context.Background(), dayKey, 86400*2*time.Second)
				c.IncrBy(context.Background(), hourKey, tickerCount)
				c.Expire(context.Background(), hourKey, 86400*2*time.Second)
			})
			if err != nil {
				fmt.Println("RedisConfPipline err:", err)
				continue
			}
			//步骤 4️⃣：从 Redis 读取当天总请求数
			totalCount, err := reqCounter.GetDayData(currentTime)
			if err != nil {
				fmt.Println("GetDayData err:", err)
				continue
			}
			//步骤 5️⃣：计算 QPS（每秒请求数）
			nowUnix := time.Now().Unix()
			if reqCounter.Unix == 0 { //记录初始时间 & 总数，跳过 QPS 计算
				reqCounter.Unix = nowUnix
				reqCounter.TotalCount = totalCount
				continue
			}

			diff := totalCount - reqCounter.TotalCount // 计算时间差值
			if nowUnix > reqCounter.Unix {
				reqCounter.TotalCount = totalCount
				reqCounter.QPS = diff / (nowUnix - reqCounter.Unix)
				reqCounter.Unix = nowUnix
			}
		}
	}()

	return reqCounter
}

func (o *RedisFlowCountService) GetDayKey(t time.Time) string {
	dayStr := t.In(TimeLocation).Format("20060102")
	return fmt.Sprintf("%s_%s_%s", RedisFlowDayKey, dayStr, o.AppID)
}

func (o *RedisFlowCountService) GetHourKey(t time.Time) string {
	hourStr := t.In(TimeLocation).Format("2006010215")
	return fmt.Sprintf("%s_%s_%s", RedisFlowHourKey, hourStr, o.AppID)
}

func (o *RedisFlowCountService) GetHourData(t time.Time) (int64, error) {
	key := o.GetHourKey(t)
	val, err := RedisClient.Get(context.Background(), key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

func (o *RedisFlowCountService) GetDayData(t time.Time) (int64, error) {
	key := o.GetDayKey(t)
	val, err := RedisClient.Get(context.Background(), key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

// 原子增加 —— 不需要 goroutine！atomic.AddInt64 本身是线程安全的
func (o *RedisFlowCountService) Increase() {
	atomic.AddInt64(&o.TickerCount, 1)
}
