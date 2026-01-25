package main

import (
	"net"
	"sync"
	"time"
)

func main() {
	ips := []string{"192.168.1.1", "192.168.1.2", "1", "2", "192.168.1.1", "192.168.1.2", "1", "2",
		"192.168.1.1", "192.168.1.2", "1", "2", "192.168.1.1", "192.168.1.2", "1", "2",
		"192.168.1.1", "192.168.1.2", "1", "2", "192.168.1.1", "192.168.1.2", "1", "2"}
	results := make(chan string, len(ips))
	var wg sync.WaitGroup
	// TODO: 实现并发限制为 10 的探测逻辑
	// 1. 创建容量为10的通道作为并发控制器（令牌桶）
	limit := make(chan struct{}, 10)

	for _, ip := range ips {
		// 获取令牌（通道满时会阻塞，实现并发限制）
		limit <- struct{}{}

		wg.Add(1)
		go func(ip string) {
			println("11")
			defer func() {
				wg.Done()
				<-limit
			}()
			if checkAlive(ip) {
				results <- ip
			}
		}(ip)
	}
}

func checkAlive(ip string) bool {
	_, err := net.DialTimeout("tcp", ip+":80", 2*time.Second)
	return err == nil
}
