package main

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	subnet := "192.168.1"
	startIP := 1
	endIP := 254
	// commonPorts := []int{21, 22, 23, 80, 443, 3306, 3389, 8080}
	//循环从1到65535
	commonPorts := make([]int, 65535)
	for i := 1; i <= 65535; i++ {
		commonPorts[i-1] = i
	}

	timeout := 2 * time.Second
	workers := 100

	totalTasks := (endIP - startIP + 1) * len(commonPorts)
	fmt.Printf("开始扫描 %s.%d-%d 的 %d 个常用端口...\n", subnet, startIP, endIP, len(commonPorts))
	fmt.Printf("总任务数: %d, 并发数: %d\n\n", totalTasks, workers)

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, workers)

	var mu sync.Mutex
	openPorts := make(map[string][]int)
	var scanned int64 = 0

	// 启动进度显示
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		for range ticker.C {
			scannedVal := atomic.LoadInt64(&scanned)
			if scannedVal >= int64(totalTasks) {
				ticker.Stop()
				return
			}
			fmt.Printf("\r进度: %d/%d (%.1f%%)", scannedVal, totalTasks, float64(scannedVal)*100/float64(totalTasks))
		}
	}()

	for i := startIP; i <= endIP; i++ {
		host := fmt.Sprintf("%s.%d", subnet, i)

		for _, port := range commonPorts {
			wg.Add(1)
			semaphore <- struct{}{}

			go func(h string, p int) {
				defer wg.Done()
				defer func() { <-semaphore }()

				address := fmt.Sprintf("%s:%d", h, p)
				conn, err := net.DialTimeout("tcp", address, timeout)
				atomic.AddInt64(&scanned, 1)

				if err != nil {
					return
				}
				conn.Close()

				mu.Lock()
				openPorts[h] = append(openPorts[h], p)
				// 修复输出竞争
				fmt.Printf("\n[+] %s:%d 开放!\n", h, p)
				mu.Unlock()
			}(host, port)
		}
	}

	wg.Wait()

	fmt.Println("\n\n========== 扫描结果 ==========")
	openCount := 0
	for _, ports := range openPorts {
		openCount += len(ports)
	}
	fmt.Printf("开放端口总数: %d 个\n", openCount)

	if len(openPorts) > 0 {
		fmt.Println("开放端口详情:")
		for host, ports := range openPorts {
			fmt.Printf("  %s: %v\n", host, ports)
		}
	}
}
