package main

import (
	"flag"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

func main() {
	host := flag.String("host", "localhost", "The hostname to scan")
	portRange := flag.String("ports", "1-1024", "The range of ports (e.g. 80-443 or 1-65535)")
	concurrency := flag.Int("concurrency", 100, "Number of concurrent workers")
	flag.Parse()

	// 1. Parsing Port Range
	fmt.Printf("🔍 Parsing port range: %s\n", *portRange)
	start, end, err := parseRange(*portRange)
	if err != nil {
		fmt.Printf("❌ Invalid range format: %v\n", err)
		return
	}

	fmt.Printf("⚡️ Starting scan on %s (%d to %d)\n", *host, start, end)
	fmt.Printf("📊 Concurrency: %d workers\n\n", *concurrency)

	ports := make(chan int, *concurrency)
	results := make(chan int)
	var wg sync.WaitGroup

	// 2. Start Workers
	for i := 0; i < *concurrency; i++ {
		go func() {
			for p := range ports {
				wg.Add(1)
				address := net.JoinHostPort(*host, fmt.Sprintf("%d", p))
				conn, err := net.DialTimeout("tcp", address, 1*time.Second)
				if err != nil {
					wg.Done()
					continue
				}
				conn.Close()
				fmt.Printf("  [✔] Port %d is OPEN\n", p)
				results <- p
				wg.Done()
			}
		}()
	}

	// 3. Collector
	var openPorts []int
	go func() {
		for p := range results {
			openPorts = append(openPorts, p)
		}
	}()

	// 4. Feed Tasks
	for i := start; i <= end; i++ {
		ports <- i
	}
	close(ports)

	wg.Wait()
	close(results)

	// 5. Final Report
	sort.Ints(openPorts)
	fmt.Printf("\n🏁 SCAN COMPLETE on %s\n", *host)
	fmt.Println("===============================")
	if len(openPorts) == 0 {
		fmt.Println("No open ports found.")
	} else {
		for _, p := range openPorts {
			fmt.Printf("  ✅ Port %-6d (OPEN)\n", p)
		}
	}
	fmt.Println("===============================")
}

func parseRange(r string) (int, int, error) {
	var start, end int
	_, err := fmt.Sscanf(r, "%d-%d", &start, &end)
	if err != nil {
		return 0, 0, err
	}
	return start, end, nil
}
