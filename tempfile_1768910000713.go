package main

var currentWorkers int32 // 当前活跃的goroutine数量

for _, ip := range ips {
    limit <- struct{}{}
    
    wg.Add(1)
    go func(ip string) {
        count := atomic.AddInt32(&currentWorkers, 1)
        fmt.Printf("启动goroutine，当前并发数: %d, IP: %s\n", count, ip)
        
        defer func() {
            atomic.AddInt32(&currentWorkers, -1)
            wg.Done()
            <-limit
        }()
        
        if checkAlive(ip) {
            results <- ip
        }
    }(ip)
}
