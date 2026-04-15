package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	// 配置参数
	subnet := "192.168.1"
	startIP := 1
	endIP := 254
	port := 8000
	timeout := 2 * time.Second
	workers := 100 // 并发数

	fmt.Printf("开始扫描 %s.%d-%d 的 %d 端口...\n\n", subnet, startIP, endIP, port)

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, workers) // 控制并发
	var mu sync.Mutex
	openHosts := []string{}

	// 遍历 IP 范围
	for i := startIP; i <= endIP; i++ {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量

		go func(ip int) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			host := fmt.Sprintf("%s.%d", subnet, ip)
			address := fmt.Sprintf("%s:%d", host, port)

			// 尝试连接
			conn, err := net.DialTimeout("tcp", address, timeout)
			if err != nil {
				fmt.Printf("[-] %s:%d 关闭\n", host, port)
				return
			}
			conn.Close()

			mu.Lock()
			openHosts = append(openHosts, host)
			mu.Unlock()
			fmt.Printf("[+] %s:%d 开放!\n", host, port)
		}(i)
	}

	wg.Wait()

	// 输出结果汇总
	fmt.Println("\n========== 扫描结果 ==========")
	fmt.Printf("总共扫描: %d 个主机\n", endIP-startIP+1)
	fmt.Printf("开放端口: %d 个\n", len(openHosts))
	if len(openHosts) > 0 {
		fmt.Println("开放的主机列表:")
		for _, host := range openHosts {
			fmt.Printf("  - %s:%d\n", host, port)
		}
	}
}
