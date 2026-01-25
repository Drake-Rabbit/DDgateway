package main

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	ips := []string{"192.168.1.1", "192.168.1.2", "1", "2", "192.168.1.1", "192.168.1.2", "1", "2",
		"192.168.1.1", "192.168.1.2", "1", "2", "192.168.1.1", "192.168.1.2", "1", "2",
		"192.168.1.1", "192.168.1.2", "1", "2", "192.168.1.1", "192.168.1.2", "1", "2"}

	results := make(chan string, len(ips))
	var wg sync.WaitGroup

	var currentWorkers int32
	limit := make(chan struct{}, 10)

	for _, ip := range ips {
		limit <- struct{}{}

		wg.Add(1)
		go func(ip string) {
			count := atomic.AddInt32(&currentWorkers, 1)
			fmt.Printf("goroutine启动 - 当前并发数: %d, IP: %s\n", count, ip)

			defer func() {
				atomic.AddInt32(&currentWorkers, -1)
				fmt.Printf("goroutine结束 - IP: %s\n", ip)
				wg.Done()
				<-limit
			}()

			if checkAlive2(ip) {
				results <- ip
			}
		}(ip)
	}

	wg.Wait()
	close(results)

	fmt.Println("所有检测完成")
	for result := range results {
		fmt.Printf("存活IP: %s\n", result)
	}
}

func checkAlive2(ip string) bool {
	_, err := net.DialTimeout("tcp", ip+":80", 2*time.Second)
	return err == nil
}
